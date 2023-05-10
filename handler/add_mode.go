package handler

import (
	"fmt"
	"linebot/constant"
	"linebot/store"
	"log"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func HandleAddMode(event *linebot.Event, eventContext *EventContext) {
	userId := event.Source.UserID
	userInputData := eventContext.userInputData
	dbStore := eventContext.dbStore
	bc := eventContext.botClient
	text := event.Message.(*linebot.TextMessage).Text

	switch userInputData.Question {
	case constant.Name:
		if dbStore.IsRestaurantSaved(userId, text) {
			bc.SendText(event.ReplyToken, "店家已存在，請重新其他店家")
			return
		}
		userInputData.
			SetQuestion(constant.Phone).
			SetData(func(m *store.Merchant) *store.Merchant {
				m.Name = text
				return m
			})
		bc.SendText(event.ReplyToken, "請輸入商家電話")
	case constant.Phone:
		userInputData.SetData(func(m *store.Merchant) *store.Merchant {
			m.Phone = text
			return m
		})
		data := userInputData.Data
		restaurant, err := dbStore.CreateRestaurant(data.Name, data.Phone, "")
		if err != nil {
			log.Fatalf("CreateRestaurant error!!, %v", err)
		}
		err = dbStore.AddRestaurantToUser(userId, restaurant.ID)
		if err != nil {
			log.Fatalf("AddRestaurantToUser error!!, %v", err)
		}
		bc.SendText(event.ReplyToken, fmt.Sprintf("店家[%v]成功新增!!", restaurant.Name))
		userInputData.Reset()
	}
}
