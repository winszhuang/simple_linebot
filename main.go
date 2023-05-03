package main

import (
	_ "embed"
	"fmt"
	c "linebot/constants"
	dbService "linebot/db"
	. "linebot/handler"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v7/linebot/httphandler"
)

var (
	//go:embed richmenu.png
	richMenuImg []byte
	//go:embed richmenu.json
	richMenuJson                   []byte
	richMenuImgFileNameInBuildTime string
)

func main() {
	// check is dev
	if os.Getenv("ISPROD") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	menuImgPath, err := newFileInBuildTime("menu.png", richMenuImg)
	if err != nil {
		log.Fatal("build menu png error!!")
	}
	menuJsonPath, err := newFileInBuildTime("menu.json", richMenuJson)
	if err != nil {
		log.Fatal("build menu json error!!")
	}

	// connect db
	if err := dbService.InitDB(); err != nil {
		log.Fatalf("Error occurred: %v", err)
	}

	handler, err := httphandler.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// #NOTICE 我猜這裏可能有問題，可能一個客戶對應一個client?
	bot, err := handler.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	if err = c.SetupRichMenu(bot, menuImgPath, menuJsonPath); err != nil {
		log.Fatal(err)
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		for _, event := range events {
			userId := event.Source.UserID
			fmt.Println(userId)

			eventHandler := &EventHandler{Event: event, Bot: bot, UserId: userId}
			userInputData := LoadUserInputData(userId)

			if err := dbService.InitUserInDb(userId, bot); err != nil {
				log.Fatal()
			}

			switch event.Type {
			case linebot.EventTypePostback:
				ResetUserInputData(userId)
				handlePostback(eventHandler, userInputData)
			case linebot.EventTypeMessage:
				handleMessage(eventHandler, userInputData)
			}
		}
	})
	http.Handle("/callback", handler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func handleMessage(eh *EventHandler, userInputInfo *UserInputInfo) {
	switch messageData := eh.Event.Message.(type) {
	case *linebot.TextMessage:
		message := strings.TrimSpace(messageData.Text)
		if c.IsDirective(message) {
			ResetUserInputData(eh.UserId)
			switch message {
			case c.Add:
				userInputInfo.SetMode(c.Add).SetQuestion(c.Name)
				fmt.Println("當前userInputInfo: ", userInputInfo)
				if err := eh.SendText("請輸入商家名稱"); err != nil {
					log.Fatal(err)
				}
			case c.Remove:
				userInputInfo.SetMode(c.Remove).SetQuestion(c.Name)
				if err := eh.SendText("請輸入商家名稱"); err != nil {
					log.Fatal(err)
				}
			case c.List:
				showRestaurantList(eh)
			case string(c.Pick):
				showRandomRestaurant(eh)
			case c.Near:
				showNearByRestaurants(eh)
			}
		} else {
			isUserInSomeMode := userInputInfo.IsInMode()
			if isUserInSomeMode {
				switch userInputInfo.GetMode() {
				case c.Add:
					handleAddMode(eh, userInputInfo, message)
				case c.Remove:
					handleRemoveMode(eh, userInputInfo, message)
				}
			} else {
				showRule(eh)
			}
		}
	}
}

func handlePostback(eh *EventHandler, userInputInfo *UserInputInfo) {
	switch eh.Event.Postback.Data {
	case c.Add:
		userInputInfo.SetMode(c.Add).SetQuestion(c.Name)
		if err := eh.SendText("請輸入商家名稱"); err != nil {
			log.Fatal()
		}
	case c.List:
		showRestaurantList(eh)
	case string(c.Pick):
		showRandomRestaurant(eh)
	case c.Remove:
		userInputInfo.SetMode(c.Remove).SetQuestion(c.Name)
		if err := eh.SendText("請輸入商家名稱"); err != nil {
			log.Fatal()
		}
	}
}

func handleAddMode(eh *EventHandler, userInputInfo *UserInputInfo, text string) {
	switch userInputInfo.GetQuestion() {
	case c.Name:
		if dbService.IsRestaurantSaved(eh.UserId, text) {
			if err := eh.SendText("店家已存在，請重新其他店家"); err != nil {
				log.Fatal(err)
			}
			return
		}
		userInputInfo.
			SetQuestion(c.Phone).
			SetData(func(m *Merchant) *Merchant {
				m.Name = text
				return m
			})
		if err := eh.SendText("請輸入商家電話"); err != nil {
			log.Fatal(err)
		}
	case c.Phone:
		userInputInfo.SetData(func(m *Merchant) *Merchant {
			m.Phone = text
			return m
		})
		data := userInputInfo.GetData()
		restaurant, err := dbService.CreateRestaurant(data.Name, data.Phone, "")
		if err != nil {
			log.Fatalf("CreateRestaurant error!!, %v", err)
		}
		err = dbService.AddRestaurantToUser(eh.UserId, restaurant.ID)
		if err != nil {
			log.Fatalf("AddRestaurantToUser error!!, %v", err)
		}
		msg := fmt.Sprintf("店家[%v]成功新增!!", restaurant.Name)
		if err := eh.SendText(msg); err != nil {
			log.Fatal(err)
		}
		ResetUserInputData(eh.UserId)
	}
}

func handleRemoveMode(eh *EventHandler, userInputInfo *UserInputInfo, text string) {
	switch userInputInfo.GetQuestion() {
	case c.Name:
		if !dbService.IsRestaurantSaved(eh.UserId, text) {
			if err := eh.SendText("找不到該店家名稱，請重新輸入"); err != nil {
				log.Fatal(err)
			}
			return
		}
		if err := dbService.RemoveRestaurantFromUser(eh.UserId, text); err != nil {
			log.Fatalf("RemoveRestaurantFromUser error!!, %v", err)
		}
		if err := eh.SendText("刪除店家成功"); err != nil {
			log.Fatal(err)
		}
		ResetUserInputData(eh.UserId)
	}
}

func newFileInBuildTime(newFilePathName string, goEmbedFile []byte) (string, error) {
	f, err := os.Create(newFilePathName)
	if err != nil {
		return "", err
	}
	_, err = f.WriteAt(goEmbedFile, 0)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func showRestaurantList(eh *EventHandler) {
	restaurants, err := dbService.GetRestaurantListByUser(eh.UserId)
	if err != nil {
		log.Fatalf("GetRestaurantsByUser error!! : %v", err)
	}

	str := "店家列表如下\n"
	for _, restaurant := range restaurants {
		str += "---" + "\n"
		str += restaurant.Name + "\n"
		str += restaurant.Phone + "\n"
	}

	if err := eh.SendText(str); err != nil {
		log.Fatal()
	}
}

func showRandomRestaurant(eh *EventHandler) {
	if dbService.IsUserRestaurantEmpty(eh.UserId) {
		if err := eh.SendText("尚未有店家，請先加入店家再做隨機選店!!"); err != nil {
			log.Fatal()
		}
		return
	}

	restaurant, err := dbService.PickRestaurantFromUser(eh.UserId)
	if err != nil {
		log.Fatalf("PickRestaurantFromUser error!! : %v", err)
	}

	restaurantInfo := restaurant.Name + "\n" + restaurant.Phone
	if err := eh.SendText(restaurantInfo); err != nil {
		log.Fatal()
	}
}

func showNearByRestaurants(eh *EventHandler) {
	// temp data
	list := []c.RestaurantInfo{
		{
			Name:             "吉利蛋餅",
			Rating:           4.5,
			UserRatingsTotal: 294,
			Vicinity:         "No. 68號, Section 1, Dalian Road, Beitun District",
			BusinessStatus:   "OPERATIONAL",
			Lat:              24.1762394,
			Lng:              120.6734827,
			Photo:            "https://lh3.googleusercontent.com/places/AJQcZqKccUzcZKW3Fc0jtggYqrjhd0nZLGJXmJQ3FxFBW0sFiY6apX88XX_2qa3jxqa353wL6tUxwn0mdjVVrh727Foj9u5jSdIbLYk=s1600-w400",
		},
		{
			Name:             "真北方早餐店/湯包、蛋餅",
			Rating:           3.9,
			UserRatingsTotal: 286,
			Vicinity:         "No. 52, Section 2, Beiping Road, North District",
			BusinessStatus:   "OPERATIONAL",
			Lat:              24.1715908,
			Lng:              120.6734085,
			Photo:            "https://lh3.googleusercontent.com/places/AJQcZqK4G5DEgyEoWx1oSbtpd66n0aohRlSU-aKHTMesqNKjpxdqzVa8vpPI2udIwD-1GU13lwH-bEMf2DtA0kUyKdr29pzclRyPoVc=s1600-w400",
		},
	}
	flexContainer := c.CreateBubble(list)

	_, err := eh.Bot.ReplyMessage(
		eh.Event.ReplyToken,
		linebot.NewFlexMessage("Restaurant List", flexContainer),
	).Do()
	if err != nil {
		log.Fatal(err)
	}
}

func showRule(eh *EventHandler) {
	if err := eh.SendText("使用說明如下\n---\n電腦版指令:\n/list     -> 查看所有儲存過的店家\n/pick   -> 隨機挑選儲存過的某個店家\n/add    -> 新增店家資訊\n/rm     -> 刪除店家資訊\n---\n手機板請點擊下方選單按鈕操作"); err != nil {
		log.Fatal(err)
	}
}
