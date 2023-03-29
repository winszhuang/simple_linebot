package main

import (
	_ "embed"
	"fmt"
	c "linebot/constants"
	dbService "linebot/db"
	o "linebot/enum"
	. "linebot/handler"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v7/linebot/httphandler"
)

type Merchant struct {
	Name  string
	Phone string
}

type TmpInfo struct {
	action   o.Operate
	question o.Keyword
	data     Merchant
}

var (
	//go:embed richmenu.png
	richMenuImg                    []byte
	richMenuImgFileNameInBuildTime string
	userTmpInfo                    = make(map[string]TmpInfo)
	inputMode                      = false
)

func main() {
	// check is dev
	if os.Getenv("ISPROD") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	if err := initRichMenuImgPath(); err != nil {
		log.Fatal(err)
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

	if err = c.GenerateRichMenu(bot, richMenuImgFileNameInBuildTime); err != nil {
		log.Fatal(err)
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		for _, event := range events {
			userId := event.Source.UserID
			fmt.Println(userId)

			eventHandler := &EventHandler{Event: event, Bot: bot, UserId: userId}

			if err := initUserInDb(userId, bot); err != nil {
				log.Fatal()
			}

			switch event.Type {
			case linebot.EventTypePostback:
				handlePostback(eventHandler)
			case linebot.EventTypeMessage:
				handleMessage(eventHandler)
			}
		}
	})
	http.Handle("/callback", handler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func handleMessage(eh *EventHandler) {
	switch message := eh.Event.Message.(type) {
	case *linebot.TextMessage:
		// 額外處理有條件限制的情況
		if tmpInfo, ok := userTmpInfo[eh.UserId]; ok {
			if tmpInfo.action == o.Add {
				if tmpInfo.question == o.Name {
					if dbService.IsRestaurantSaved(eh.UserId, message.Text) {
						if err := eh.SendText("店家已存在，請重新其他店家"); err != nil {
							log.Fatal(err)
						}
						return
					}
					tmpInfo.data = Merchant{Name: message.Text}
					tmpInfo.question = o.Phone
					if err := eh.SendText("請輸入商家電話"); err != nil {
						log.Fatal(err)
					}
					userTmpInfo[eh.UserId] = tmpInfo
				} else if tmpInfo.question == o.Phone {
					tmpInfo.data.Phone = message.Text
					userTmpInfo[eh.UserId] = tmpInfo
					restaurant, err := dbService.CreateRestaurant(tmpInfo.data.Name, tmpInfo.data.Phone, "")
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
					delete(userTmpInfo, eh.UserId)
				}
			}
			if tmpInfo.action == o.Remove {
				userTmpInfo[eh.UserId] = tmpInfo
				if !dbService.IsRestaurantSaved(eh.UserId, message.Text) {
					if err := eh.SendText("找不到該店家名稱，請重新輸入"); err != nil {
						log.Fatal(err)
					}
					return
				}
				if err := dbService.RemoveRestaurantFromUser(eh.UserId, message.Text); err != nil {
					log.Fatalf("RemoveRestaurantFromUser error!!, %v", err)
				}
				if err := eh.SendText("刪除店家成功"); err != nil {
					log.Fatal(err)
				}
				delete(userTmpInfo, eh.UserId)
			}
			return
		}
		fmt.Println("回傳過來的文字是: ", message.Text)
		if err := eh.SendText("請遵照以下菜單來做功能選擇並輸入對應內容\nps目前尚未開放電腦版輸入指令"); err != nil {
			log.Fatal(err)
		}
	}
}

func handlePostback(eh *EventHandler) {
	switch eh.Event.Postback.Data {
	case o.Add:
		userTmpInfo[eh.UserId] = TmpInfo{action: o.Add, question: o.Name}
		if err := eh.SendText("請輸入商家名稱"); err != nil {
			log.Fatal()
		}
	case o.List:
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
	case string(o.Pick):
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
	case o.Remove:
		userTmpInfo[eh.UserId] = TmpInfo{action: o.Remove, question: o.Name}
		if err := eh.SendText("請輸入商家名稱"); err != nil {
			log.Fatal()
		}
	}
}

func initRichMenuImgPath() error {
	f, err := os.Create("menu.png")
	if err != nil {
		return err
	}
	_, err = f.WriteAt(richMenuImg, 0)
	if err != nil {
		return err
	}
	richMenuImgFileNameInBuildTime = f.Name()
	return nil
}

func initUserInDb(userId string, bot *linebot.Client) error {
	if dbService.IsUserExists(userId) {
		return nil
	}

	userData, err := bot.GetProfile(userId).Do()
	if err != nil {
		return err
	}

	return dbService.CreateUser(
		userData.DisplayName,
		userData.Language,
		userData.PictureURL,
		userData.UserID,
	)
}
