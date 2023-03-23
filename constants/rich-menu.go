package constants

import (
	"fmt"
	o "linebot/enum"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func GenerateRichMenu(bot *linebot.Client, imgPath string) {
	richMenu := linebot.RichMenu{
		Size:        linebot.RichMenuSize{Width: 2500, Height: 1686},
		Selected:    true,
		Name:        "Menu",
		ChatBarText: "點我收合選單",
		Areas: []linebot.AreaDetail{
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 0, Width: 2500, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: string(o.Pick),
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 843, Width: 833, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: o.Add,
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 833, Y: 843, Width: 833, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: o.List,
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1666, Y: 843, Width: 833, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: o.Remove,
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

	_, err1 := bot.UploadRichMenuImage(res.RichMenuID, imgPath).Do()
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
