package telegram

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Структуры для работы с Telegram API

type TelegramResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  message `json:"message"`
}

type message struct {
	Chat    chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker sticker `json:"sticker"`
}

type chat struct {
	ChatId int `json:"id"`
}

type sticker struct {
	File_id string `json:"file_id"`
}

// Функция получения апдейтов
func GetUpdates(botUrl string, offset int) ([]Update, error) {

	// Rest запрос для получения апдейтов
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Запись и обработка полученных данных
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse TelegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}
