package main

import (
	_ "embed"
	"fmt"
	c "linebot/constants"
	dbService "linebot/db"
	o "linebot/enum"
	merchant "linebot/handler"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v7/linebot/httphandler"
)

type TmpInfo struct {
	action   o.Operate
	question o.Keyword
	data     merchant.Merchant
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

			if err := initUserInDb(userId, bot); err != nil {
				log.Fatal()
			}

			switch event.Type {
			case linebot.EventTypePostback:
				handlePostback(bot, event, userId)
			case linebot.EventTypeMessage:
				handleMessage(bot, event, userId)
			}
		}
	})
	http.Handle("/callback", handler)

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func handleMessage(bot *linebot.Client, event *linebot.Event, userId string) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		// 額外處理有條件限制的情況
		if tmpInfo, ok := userTmpInfo[userId]; ok {
			if tmpInfo.action == o.Add {
				if tmpInfo.question == o.Name {
					if dbService.IsRestaurantSaved(userId, message.Text) {
						_, err := bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("店家已存在，請重新其他店家"),
						).Do()
						if err != nil {
							log.Fatal(err)
						}
						return
					}
					tmpInfo.data = merchant.Merchant{Name: message.Text}
					tmpInfo.question = o.Phone
					_, err := bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("請輸入商家電話"),
					).Do()
					if err != nil {
						log.Fatal(err)
					}
					userTmpInfo[userId] = tmpInfo
				} else if tmpInfo.question == o.Phone {
					tmpInfo.data.Phone = message.Text
					userTmpInfo[userId] = tmpInfo
					restaurant, err := dbService.CreateRestaurant(tmpInfo.data.Name, tmpInfo.data.Phone, "")
					if err != nil {
						log.Fatalf("CreateRestaurant error!!, %v", err)
					}
					err = dbService.AddRestaurantToUser(userId, restaurant.ID)
					if err != nil {
						log.Fatalf("AddRestaurantToUser error!!, %v", err)
					}
					msg := fmt.Sprintf("店家[%v]成功新增!!", restaurant.Name)
					_, err = bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage(msg),
					).Do()
					if err != nil {
						log.Fatal(err)
					}
					delete(userTmpInfo, userId)
				}
			}
			if tmpInfo.action == o.Remove {
				userTmpInfo[userId] = tmpInfo
				if !dbService.IsRestaurantSaved(userId, message.Text) {
					_, err := bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("找不到該店家名稱，請重新輸入"),
					).Do()
					if err != nil {
						log.Fatal(err)
					}
					return
				}
				if err := dbService.RemoveRestaurantFromUser(userId, message.Text); err != nil {
					log.Fatalf("RemoveRestaurantFromUser error!!, %v", err)
				}
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewTextMessage("刪除店家成功"),
				).Do()
				if err != nil {
					log.Fatal(err)
				}
				delete(userTmpInfo, userId)
			}
			return
		}
		fmt.Println("回傳過來的文字是: ", message.Text)
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage("請遵照以下菜單來做功能選擇並輸入對應內容\nps目前尚未開放電腦版輸入指令"),
		).Do()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handlePostback(bot *linebot.Client, event *linebot.Event, userId string) {
	switch event.Postback.Data {
	case o.Add:
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage("請輸入商家名稱"),
		).Do()
		userTmpInfo[userId] = TmpInfo{action: o.Add, question: o.Name}
		if err != nil {
			log.Fatal(err)
		}
	case o.List:
		str := "店家列表如下\n"
		restaurants, err := dbService.GetRestaurantsByUser(userId)
		if err != nil {
			log.Fatalf("GetRestaurantsByUser error!! : %v", err)
		}

		for _, restaurant := range restaurants {
			str += "---" + "\n"
			str += restaurant.Name + "\n"
			str += restaurant.Phone + "\n"
		}
		_, err = bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(str),
		).Do()
		if err != nil {
			log.Fatal(err)
		}
	case string(o.Pick):
		if dbService.IsUserRestaurantEmpty(userId) {
			_, err := bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage("尚未有店家，請先加入店家再做隨機選店!!"),
			).Do()
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		restaurant, err := dbService.PickRestaurantFromUser(userId)
		if err != nil {
			log.Fatalf("PickRestaurantFromUser error!! : %v", err)
		}

		restaurantInfo := restaurant.Name + "\n" + restaurant.Phone
		_, err = bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(restaurantInfo),
		).Do()
		if err != nil {
			log.Fatal(err)
		}
	case o.Remove:
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage("請輸入商家名稱"),
		).Do()
		userTmpInfo[userId] = TmpInfo{action: o.Remove, question: o.Name}
		if err != nil {
			log.Fatal(err)
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
