package commands

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
	"github.com/sirupsen/logrus"
)

type cryptoResponse struct {
	Symbol        string `json:"symbol"`
	ChangePercent string `json:"priceChangePercent"`
	LastPrice     string `json:"lastPrice"`
}

// Функция вывода курса криптовалюты SHIB
func SendCryptoInfo(botUrl string, chatId int) {

	// Отправка запроса
	resp, err := http.Get("https://api2.binance.com/api/v3/ticker/24hr?symbol=SHIBBUSD")
	if err != nil {
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		logrus.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var response = new(cryptoResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Printf("in SendCryptoInfo: json.Unmarshal err: %v", err)
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Формирование и отправка результата
	if response.ChangePercent[0] == '-' {
		telegram.SendMsg(botUrl, chatId, "За сегодняшний день "+response.Symbol+" упал на "+response.ChangePercent[1:]+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		telegram.SendRandomShibaSticker(botUrl, chatId, true)
	} else {
		telegram.SendMsg(botUrl, chatId, "За сегодняшний день "+response.Symbol+" вырос на "+response.ChangePercent+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		telegram.SendRandomShibaSticker(botUrl, chatId, false)
	}

}
