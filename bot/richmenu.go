package bot

import (
	"log"
)

func (bc *BotClient) SetupRichMenu() {
	// #TODO
	resp, err := bc.GetDefaultRichMenu().Do()
	if err != nil {
		panic(err)
	}

	log.Printf("resp.RichMenuID: %s", resp.RichMenuID)
	if resp.RichMenuID == "" {
		log.Fatal("richmenu not found!!")
	}
}
