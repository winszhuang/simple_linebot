package handler

import (
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func ShowRandomRestaurant(event *linebot.Event, eventContext *EventContext) {
	userId := event.Source.UserID
	dbStore := eventContext.dbStore
	bc := eventContext.botClient

	if dbStore.IsUserRestaurantEmpty(userId) {
		bc.SendText(event.ReplyToken, "尚未有店家，請先加入店家再做隨機選店!!")
		return
	}

	restaurant, err := dbStore.PickRestaurantFromUser(userId)
	if err != nil {
		log.Fatalf("PickRestaurantFromUser error!! : %v", err)
	}

	restaurantInfo := restaurant.Name + "\n" + restaurant.Phone
	bc.SendText(event.ReplyToken, restaurantInfo)
}
