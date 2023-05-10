package main

import (
	"linebot/bot"
	"linebot/config"
	"linebot/handler"
	"linebot/store"
	"net/http"
)

func init() {
	config.Init()
	bot.Init()
	store.InitDB()
	store.InitLocationStore()
}

func main() {
	http.HandleFunc("/callback", handler.Entry)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	if err := http.ListenAndServe(":"+config.C.Port, nil); err != nil {
		panic(err)
	}
}
