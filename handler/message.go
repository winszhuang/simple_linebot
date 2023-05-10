package handler

import (
	"linebot/constant"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func HandleMessage(event *linebot.Event, eventContext *EventContext) {
	userInputData := eventContext.userInputData
	bc := eventContext.botClient

	switch messageData := event.Message.(type) {
	case *linebot.LocationMessage:
		ShowRestaurantSearch(event.ReplyToken, eventContext, 1, messageData.Latitude, messageData.Longitude)
	case *linebot.TextMessage:
		message := strings.TrimSpace(messageData.Text)
		if constant.IsDirective(message) {
			userInputData.Reset()
			directive := constant.Directive(message)
			switch directive {
			case constant.Add:
				userInputData.SetMode(constant.Add).SetQuestion(constant.Name)
				bc.SendText(event.ReplyToken, "請輸入商家名稱")
			case constant.Remove:
				userInputData.SetMode(constant.Remove).SetQuestion(constant.Name)
				bc.SendText(event.ReplyToken, "請輸入商家名稱")
			case constant.List:
				ShowRestaurantList(event, eventContext)
			case constant.Pick:
				ShowRandomRestaurant(event, eventContext)
			case constant.Near:
				bc.SendText(event.ReplyToken, "請先傳送位置資訊給我")
			}
		} else {
			if userInputData.IsInMode() {
				switch userInputData.Mode {
				case constant.Add:
					HandleAddMode(event, eventContext)
				case constant.Remove:
					HandleRemoveMode(event, eventContext)
				}
			} else {
				ShowRule(event, bc)
			}
		}
	}
}
