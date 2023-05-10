package bot

import (
	"linebot/config"
	"linebot/model"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type BotClient struct {
	*linebot.Client
}

var bc *BotClient

func Init() {
	bot, err := linebot.New(config.C.ChannelSecret, config.C.ChannelToken)
	if err != nil {
		panic(err)
	}

	bc = &BotClient{bot}
	bc.SetWebHookUrl()
	bc.SetupRichMenu()
}

func GetBotClient() *BotClient {
	return bc
}

func (bc *BotClient) SetWebHookUrl() {
	_, err := bc.Client.SetWebhookEndpointURL(config.C.WebHookUrl).Do()
	if err != nil {
		panic(err)
	}
}

func (bc *BotClient) SendText(replyToken, text string) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}

func (bc *BotClient) SendFlex(replyToken string, altText string, flexContainer linebot.FlexContainer) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewFlexMessage(altText, flexContainer),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}

func (bc *BotClient) SendRestaurantList(replyToken string, list []model.RestaurantInfo) {
	str := "列表如下\n"
	for _, item := range list {
		str += "---" + "\n"
		str += item.Name + "\n"
		str += item.FormattedPhoneNumber + "\n"
	}

	bc.SendText(replyToken, str)
}
