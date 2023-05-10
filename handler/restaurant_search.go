package handler

import (
	"linebot/component"
	"linebot/store"
)

const MAX_BUBBLE_COUNT = 10

func ShowRestaurantSearch(replyToken string, eventContext *EventContext, pageIndex int, lat, lng float64) {
	bc := eventContext.botClient

	locationStore := store.GetLocationStore()
	restaurantList, err := locationStore.List(store.ListParams{
		Lat:       lat,
		Lng:       lng,
		PageIndex: pageIndex,
		PageSize:  MAX_BUBBLE_COUNT,
	})

	if err != nil {
		bc.SendText(replyToken, "取得附近店家失敗!")
		return
	}

	if len(restaurantList) == 0 {
		bc.SendText(replyToken, "附近沒有店家!")
		return
	}

	var nextPageIndex int
	if len(restaurantList) < MAX_BUBBLE_COUNT {
		nextPageIndex = 0
	} else {
		nextPageIndex = pageIndex + 1
	}

	carouselContainer := component.CreateCarouselWithNext(
		restaurantList,
		nextPageIndex,
		lat,
		lng,
	)
	bc.SendFlex(replyToken, "Restaurant List", carouselContainer)
}
