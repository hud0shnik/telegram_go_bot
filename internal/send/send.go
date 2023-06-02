package send

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"tgBot/internal/utils"
)

// Структуры для отправки сообщений, стикеров и картинок

type sendMessage struct {
	ChatId    int    `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type sendSticker struct {
	ChatId     int    `json:"chat_id"`
	StickerUrl string `json:"sticker"`
}

type sendPhoto struct {
	ChatId    int    `json:"chat_id"`
	PhotoUrl  string `json:"photo"`
	Caption   string `json:"caption"`
	ParseMode string `json:"parse_mode"`
}

// Функция отправки сообщения
func SendMsg(botUrl string, chatId int, msg string) error {

	// Формирование сообщения
	buf, err := json.Marshal(sendMessage{
		ChatId:    chatId,
		Text:      msg,
		ParseMode: "HTML",
	})
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
func SendStck(botUrl string, chatId int, url string) error {

	// Формирование стикера
	buf, err := json.Marshal(sendSticker{
		ChatId:     chatId,
		StickerUrl: url,
	})
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
func SendPict(botUrl string, chatId int, photoUrl, caption string) error {

	// Формирование картинки
	buf, err := json.Marshal(sendPhoto{
		ChatId:    chatId,
		PhotoUrl:  photoUrl,
		Caption:   caption,
		ParseMode: "HTML",
	})
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
func SendRandomShibaSticker(botUrl string, chatId int, sadFlag bool) {

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
	SendStck(botUrl, chatId, stickers[utils.Random(len(stickers))])
}

// Отправка случайного стикера
func SendRandomSticker(botUrl string, chatId int) error {

	// Открытие json файла со стикерами
	file, err := os.Open("internal/send/stickers.json")
	if err != nil {
		log.Fatalf("os.Open error: %s", err)
	}
	defer file.Close()

	// Запись стикеров в массив
	stickers := [359]string{}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("ioutil.ReadAll error: %s", err)
	}
	err = json.Unmarshal(body, &stickers)
	if err != nil {
		log.Printf("json.Unmarshal error: %s", err)
	}

	// Отправка случайного стикера
	SendStck(botUrl, chatId, stickers[utils.Random(len(stickers))])

	return nil
}
