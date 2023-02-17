package mods

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Структуры для отправки сообщений, стикеров и картинок
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

// Функция отправки сообщения
func SendMsg(botUrl string, update Update, msg string) error {

	// Формирование сообщения
	botMessage := SendMessage{
		ChatId: update.Message.Chat.ChatId,
		Text:   msg,
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}

	// Отправка сообщения
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("sendMessage error: %s", err)
		return err
	}
	return nil
}

// Функция отправки стикера
func SendStck(botUrl string, update Update, url string) error {

	// Формирование стикера
	botStickerMessage := SendSticker{
		ChatId:     update.Message.Chat.ChatId,
		StickerUrl: url,
	}
	buf, err := json.Marshal(botStickerMessage)
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}
	// Отправка стикера
	_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("sendSticker error: %s", err)
		return err
	}
	return nil
}

// Функция отправки картинки
func SendPict(botUrl string, update Update, pic SendPhoto) error {

	// Формирование картинки
	buf, err := json.Marshal(pic)
	if err != nil {
		log.Printf("json.Marshal error: %s", err)
		return err
	}
	// Отправка картинки
	_, err = http.Post(botUrl+"/sendPhoto", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("sendPhoto error: %s", err)
		return err
	}
	return nil
}

// Функция отправки случайного стикера с собакой
func SendRandomShibaSticker(botUrl string, update Update, sadFlag bool) {

	// Массив стикеров
	var stickers [5]string

	// Проверка на тип стикеров
	if sadFlag {

		// Запись стикеров в массив
		stickers = [5]string{
			"CAACAgIAAxkBAAIWzmGvey9t1OC7aV0860j69WsT9G-DAAJ-AQACK15TC4qyw0Zen8nxIwQ",
			"CAACAgIAAxkBAAIWz2GvezDv5uKkBgRlqhAW3oK1dzFlAAKAAQACK15TC6DmST8rBLf3IwQ",
			"CAACAgIAAxkBAAIW0mGvezKCd2-1xoEwA2hKLGsN1-izAAKIAQACK15TCzT4pMalZQrlIwQ",
			"CAACAgIAAxkBAAIW1GGvezXS4RnzDeu0Lw_L2Sw4YA94AAKDAQACK15TCwzud-biO4E7IwQ",
			"CAACAgIAAxkBAAIW1mGvezmO36icAAH_ayJKj0ybA-yDVgAChAEAAiteUwtgPKr0UyWrYyME",
		}
	} else {

		// Запись стикеров в массив
		stickers = [5]string{
			"CAACAgIAAxkBAAIM7mF7830wgmsiYJ5xHTEZjHgJ_YphAAKRAQACK15TC92mC_kqIE5PIQQ",
			"CAACAgIAAxkBAAIM8mF785AXsxybm8IbstiOBA8vc7ujAAKHAQACK15TC3gn1k2Gf2lgIQQ",
			"CAACAgIAAxkBAAIM8GF784o9uWLTWhdCbaiH3xebHlDpAAKKAQACK15TCxtDbMsAAT60RCEE",
			"CAACAgIAAxkBAAITiGGOKl7peNxJRfBRLWvZDikLRMrxAAKMAQACK15TCzSpEXiTiKA5IgQ",
			"CAACAgIAAxkBAAITimGOKmYIQWpBWdEvs-J-RS4RWJZwAAKBAQACK15TC14KbD5sAAF4tCIE",
		}
	}

	// Отправка случайного стикера
	SendStck(botUrl, update, stickers[Random(len(stickers))])
}

// Отправка случайного стикера
func SendRandomSticker(botUrl string, update Update) error {

	// Открытие json файла со стикерами
	fileU, err := os.Open("mods/stickers.json")
	if err != nil {
		log.Fatalf("os.Open error: %s", err)
	}
	defer fileU.Close()

	// Запись стикеров в массив
	bodyU, _ := ioutil.ReadAll(fileU)
	stickers := [359]string{}
	json.Unmarshal(bodyU, &stickers)

	// Отправка случайного стикера
	SendStck(botUrl, update, stickers[Random(len(stickers))])
	return nil
}
