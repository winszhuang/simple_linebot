package handler

import "github.com/line/line-bot-sdk-go/v7/linebot"

func ShowRestaurantList(event *linebot.Event, eventContext *EventContext) {
	userId := event.Source.UserID
	dbStore := eventContext.dbStore
	bc := eventContext.botClient

	list, err := dbStore.GetRestaurantListByUser(userId)
	if err != nil {
		bc.SendText(event.ReplyToken, "取得儲存的商家失敗!!")
		return
	}
	bc.SendRestaurantList(event.ReplyToken, list)
}
