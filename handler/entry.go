package handler

import (
	"linebot/bot"
	"linebot/store"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type EventContext struct {
	userInputData *store.InputCache
	dbStore       *store.DBStore
	botClient     *bot.BotClient
	// 後續可繼續擴充
}

func Entry(w http.ResponseWriter, r *http.Request) {
	bc := bot.GetBotClient()
	events, err := bc.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		userInputStore := store.GetUserInputStore()
		userInputData := userInputStore.LoadUserInputData(event.Source.UserID)

		dbStore := store.GetDBStore()
		eventContext := &EventContext{userInputData, dbStore, bc}

		switch event.Type {
		case linebot.EventTypePostback:
			HandlePostback(event, eventContext)
		case linebot.EventTypeMessage:
			HandleMessage(event, eventContext)
		}
	}
}
