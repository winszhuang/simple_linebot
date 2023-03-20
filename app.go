package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/line/line-bot-sdk-go/v7/linebot/httphandler"
)

func main() {
	fmt.Println("----------------")
	fmt.Println("近來瞜")
	fmt.Println("CHANNEL_SECRET")
	fmt.Println(os.Getenv("CHANNEL_SECRET"))

	// ----------
	// 本地測試再開啟這段
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// ----------

	handler, err := httphandler.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("success")
	}

	handler.HandleEvents(func(events []*linebot.Event, r *http.Request) {
		bot, err := handler.NewClient()
		if err != nil {
			log.Fatal(err)
			return
		}
		generateRichMenu(bot)
		handleMessage(bot, events, r)
	})
	http.Handle("/callback", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["message"] = "test test ~~"
		jsonRes, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(jsonRes)
	})

	// This is just a sample code.
	// For actually use, you must support HTTPS by using `ListenAndServeTLS`, reverse proxy or etc.
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("OK")
	}
}

func handleMessage(bot *linebot.Client, events []*linebot.Event, r *http.Request) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				fmt.Println(message)
				bytes, errr := os.ReadFile("constants/bubble-container.json")
				if errr != nil {
					fmt.Println("read bubble-container.json error")
				}

				container, jsonErr := linebot.UnmarshalFlexMessageJSON(bytes)
				if jsonErr != nil {
					fmt.Println("UnmarshalFlexMessageJSON error")
				}
				_, err := bot.ReplyMessage(
					event.ReplyToken,
					linebot.NewFlexMessage("測試", container),
				).Do()
				// _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("8888888888888888")).Do()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

func generateRichMenu(bot *linebot.Client) {
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    true,
		Name:        "Menu1",
		ChatBarText: "點我收合選單",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 1250, Height: 212},
				Action: linebot.RichMenuAction{
					Type:            linebot.RichMenuActionTypeRichMenuSwitch,
					RichMenuAliasID: "richmenu-alias-a",
					Data:            "action=richmenu-changed-to-a",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1250, Y: 0, Width: 1250, Height: 212},
				Action: linebot.RichMenuAction{
					Type:            linebot.RichMenuActionTypeRichMenuSwitch,
					RichMenuAliasID: "richmenu-alias-b",
					Data:            "action=richmenu-changed-to-b",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 212, Width: 1250, Height: 737},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: "action=buy&itemid=123",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1250, Y: 212, Width: 1250, Height: 737},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeURI,
					URI:  "https://developers.line.me/",
					Text: "click me",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 949, Width: 1250, Height: 737},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeMessage,
					Text: "hello world!",
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1250, Y: 949, Width: 1250, Height: 737},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypeDatetimePicker,
					Data: "datetime picker!",
					Mode: "datetime",
				},
			},
		},
	}
	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("menu is created success")
	}

	_, err1 := bot.UploadRichMenuImage(res.RichMenuID, "test.png").Do()
	if err1 != nil {
		fmt.Println("UploadRichMenuImage fails!!")
		log.Fatal(err)
	} else {
		fmt.Println("UploadRichMenuImage success!!")
	}

	_, err2 := bot.SetDefaultRichMenu(res.RichMenuID).Do()
	if err2 != nil {
		fmt.Println("SetDefaultRichMenu fails!!")
		log.Fatal(err)
	} else {
		fmt.Println("SetDefaultRichMenu success~~")
	}
}
