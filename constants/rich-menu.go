package constants

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type CreateRichMenuResponse struct {
	RichMenuId string `json:"richMenuId"`
}

func SetupRichMenu(bot *linebot.Client, imgPath string, jsonPath string) error {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	menuId, err := createRichMenu(bot, data)
	if err != nil {
		return err
	}

	if _, err := bot.UploadRichMenuImage(menuId, imgPath).Do(); err != nil {
		return err
	}

	if _, err := bot.SetDefaultRichMenu(menuId).Do(); err != nil {
		return err
	}

	return err
}

// 官方go-sdk沒有支援input-option... 自己造輪子
func createRichMenu(bot *linebot.Client, data []byte) (string, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.line.me/v2/bot/richmenu",
		bytes.NewReader(data),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHANNEL_TOKEN"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	source, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	response := CreateRichMenuResponse{}
	if err := json.Unmarshal(source, &response); err != nil {
		return "", fmt.Errorf("Unmarshal response body failed:", err)
	}

	fmt.Println("create menu success!!")

	return response.RichMenuId, nil
}
