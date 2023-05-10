package handler

import (
	"linebot/constant"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func HandleRemoveMode(event *linebot.Event, eventContext *EventContext) {
	userId := event.Source.UserID
	userInputData := eventContext.userInputData
	dbStore := eventContext.dbStore
	bc := eventContext.botClient
	text := event.Message.(*linebot.TextMessage).Text

	switch userInputData.Question {
	case constant.Name:
		if !dbStore.IsRestaurantSaved(userId, text) {
			bc.SendText(event.ReplyToken, "找不到該店家名稱，請重新輸入")
			return
		}
		if err := dbStore.RemoveRestaurantFromUser(userId, text); err != nil {
			log.Fatalf("RemoveRestaurantFromUser error!!, %v", err)
		}
		bc.SendText(event.ReplyToken, "刪除店家成功")
		userInputData.Reset()
	}
}
