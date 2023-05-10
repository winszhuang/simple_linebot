package handler

import (
	"linebot/constant"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func HandlePostback(event *linebot.Event, eventContext *EventContext) {
	userInputData := eventContext.userInputData
	bc := eventContext.botClient
	data := event.Postback.Data

	if constant.IsDirective(data) {
		directive := constant.Directive(data)
		switch directive {
		case constant.Add:
			userInputData.SetMode(constant.Add).SetQuestion(constant.Name)
			bc.SendText(event.ReplyToken, "請輸入商家名稱")
		case constant.List:
			ShowRestaurantList(event, eventContext)
		case constant.Pick:
			ShowRandomRestaurant(event, eventContext)
		case constant.Remove:
			userInputData.SetMode(constant.Remove).SetQuestion(constant.Name)
			bc.SendText(event.ReplyToken, "請輸入商家名稱")
		}
	}

	if strings.Contains(event.Postback.Data, "pageIndex") {
		pageIndex, lat, lng := parseQueryToLocationData(event.Postback.Data)
		ShowRestaurantSearch(event.ReplyToken, eventContext, pageIndex, lat, lng)
		return
	}
}

func parseQueryToLocationData(query string) (int, float64, float64) {
	parts := strings.Split(query, ",")
	var lat, lng float64
	var pageIndex int
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			continue
		}
		key := kv[0]
		value := kv[1]
		switch key {
		case "lat":
			lat, _ = strconv.ParseFloat(value, 64)
		case "lng":
			lng, _ = strconv.ParseFloat(value, 64)
		case "pageIndex":
			pageIndex, _ = strconv.Atoi(value)
		default:
		}
	}
	return pageIndex, lat, lng
}
