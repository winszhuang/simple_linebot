package handler

import (
	"linebot/bot"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const RULE = "使用說明如下\n---\n電腦版指令:\n/list     -> 查看所有儲存過的店家\n/pick   -> 隨機挑選儲存過的某個店家\n/add    -> 新增店家資訊\n/rm     -> 刪除店家資訊\n---\n手機板請點擊下方選單按鈕操作"

func ShowRule(event *linebot.Event, bc *bot.BotClient) {
	bc.SendText(event.ReplyToken, RULE)
}
