package mods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type SendMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type SendSticker struct {
	ChatId     int    `json:"chat_id"`
	StickerUrl string `json:"sticker"`
}
type SendPhoto struct {
	ChatId   int    `json:"chat_id"`
	PhotoUrl string `json:"photo"`
	Caption  string `json:"caption"`
}

func SendMsg(botUrl string, update Update, msg string) error {
	botMessage := SendMessage{
		ChatId: update.Message.Chat.ChatId,
		Text:   msg,
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		SendErrorMessage(botUrl, update, 2)
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendMessage method error: ", err)
		SendErrorMessage(botUrl, update, 5)
		return err
	}
	return nil
}

func SendStck(botUrl string, update Update, url string) error {
	botStickerMessage := SendSticker{
		ChatId:     update.Message.Chat.ChatId,
		StickerUrl: url,
	}
	buf, err := json.Marshal(botStickerMessage)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		SendErrorMessage(botUrl, update, 2)
		return err
	}
	_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendSticker method error: ", err)
		SendErrorMessage(botUrl, update, 3)
		return err
	}
	return nil
}

func SendPict(botUrl string, update Update, pic SendPhoto) error {
	buf, err := json.Marshal(pic)
	if err != nil {
		fmt.Println("Marshal json error: ", err)
		SendErrorMessage(botUrl, update, 2)
		return err
	}
	_, err = http.Post(botUrl+"/sendPhoto", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("SendPhoto method error: ", err)
		SendErrorMessage(botUrl, update, 4)
		return err
	}
	return nil
}

func SendRandomShibaSticker(botUrl string, update Update, sadFlag bool) {
	var stickers [5]string
	if sadFlag {
		stickers = [5]string{
			"CAACAgIAAxkBAAIWzmGvey9t1OC7aV0860j69WsT9G-DAAJ-AQACK15TC4qyw0Zen8nxIwQ",
			"CAACAgIAAxkBAAIWz2GvezDv5uKkBgRlqhAW3oK1dzFlAAKAAQACK15TC6DmST8rBLf3IwQ",
			"CAACAgIAAxkBAAIW0mGvezKCd2-1xoEwA2hKLGsN1-izAAKIAQACK15TCzT4pMalZQrlIwQ",
			"CAACAgIAAxkBAAIW1GGvezXS4RnzDeu0Lw_L2Sw4YA94AAKDAQACK15TCwzud-biO4E7IwQ",
			"CAACAgIAAxkBAAIW1mGvezmO36icAAH_ayJKj0ybA-yDVgAChAEAAiteUwtgPKr0UyWrYyME",
		}
	} else {
		stickers = [5]string{
			"CAACAgIAAxkBAAIM7mF7830wgmsiYJ5xHTEZjHgJ_YphAAKRAQACK15TC92mC_kqIE5PIQQ",
			"CAACAgIAAxkBAAIM8mF785AXsxybm8IbstiOBA8vc7ujAAKHAQACK15TC3gn1k2Gf2lgIQQ",
			"CAACAgIAAxkBAAIM8GF784o9uWLTWhdCbaiH3xebHlDpAAKKAQACK15TCxtDbMsAAT60RCEE",
			"CAACAgIAAxkBAAITiGGOKl7peNxJRfBRLWvZDikLRMrxAAKMAQACK15TCzSpEXiTiKA5IgQ",
			"CAACAgIAAxkBAAITimGOKmYIQWpBWdEvs-J-RS4RWJZwAAKBAQACK15TC14KbD5sAAF4tCIE",
		}
	}
	SendStck(botUrl, update, stickers[Random(len(stickers))])
}

func SendRandomSticker(botUrl string, update Update) error {
	fileU, err := os.Open("mods/stickers.json")
	if err != nil {
		fmt.Println(err)
		SendErrorMessage(botUrl, update, 6)
		os.Exit(1)
	}
	defer fileU.Close()

	bodyU, _ := ioutil.ReadAll(fileU)
	stickers := [359]string{}

	json.Unmarshal(bodyU, &stickers)

	SendStck(botUrl, update, stickers[Random(len(stickers))])
	return nil
}
