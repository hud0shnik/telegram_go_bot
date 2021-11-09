package mods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendMsg(botUrl string, update Update, msg string) error {
	botMessage := SendMessage{
		ChatId: update.Message.Chat.ChatId,
		Text:   msg,
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func SendStck(botUrl string, update Update, url string) error {
	botStickerMessage := SendSticker{
		ChatId:  update.Message.Chat.ChatId,
		Sticker: url,
	}
	buf, err := json.Marshal(botStickerMessage)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		return err
	}
	_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendSticker method error: ", err)
		return err
	}
	return nil
}

func SendPict(botUrl string, update Update, pic SendPhoto) error {
	buf, err := json.Marshal(pic)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		return err
	}
	_, err = http.Post(botUrl+"/sendPhoto", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendPhoto method error: ", err)
		return err
	}
	return nil
}