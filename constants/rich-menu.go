package constants

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func GenerateRichMenu(bot *linebot.Client, imgPath string) error {
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
					Data: string(Pick),
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 0, Y: 843, Width: 833, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: Add,
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 833, Y: 843, Width: 833, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: List,
				},
			},
			{
				Bounds: linebot.RichMenuBounds{X: 1666, Y: 843, Width: 833, Height: 843},
				Action: linebot.RichMenuAction{
					Type: linebot.RichMenuActionTypePostback,
					Data: Remove,
				},
			},
		},
	}

	res, err := bot.CreateRichMenu(richMenu).Do()
	if err != nil {
		return err
	}

	if _, err := bot.UploadRichMenuImage(res.RichMenuID, imgPath).Do(); err != nil {
		return err
	}

	if _, err := bot.SetDefaultRichMenu(res.RichMenuID).Do(); err != nil {
		return err
	}

	fmt.Println("menu is created success")
	return nil
}
