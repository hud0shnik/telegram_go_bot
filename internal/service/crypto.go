package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type cryptoResponse struct {
	Symbol        string `json:"symbol"`
	ChangePercent string `json:"priceChangePercent"`
	LastPrice     string `json:"lastPrice"`
}

// Функция вывода курса криптовалюты SHIB
func (s *BotService) SendCryptoInfo(chatId int64) {

	// Отправка запроса
	resp, err := http.Get("https://api2.binance.com/api/v3/ticker/24hr?symbol=SHIBUSDT")
	if err != nil {
		s.SendMessage(chatId, "Внутренняя ошибка")
		slog.Error("http.Get error", "error", err)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var response = new(cryptoResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		slog.Error("in SendCryptoInfo: json.Unmarshal error", "error", err)
		s.SendMessage(chatId, "Внутренняя ошибка")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Формирование и отправка результата
	if response.ChangePercent[0] == '-' {
		s.SendMessage(chatId, "За сегодняшний день SHIB упал на "+response.ChangePercent[1:]+"%\n"+
			"до отметки в "+response.LastPrice+"USDT\n\n")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
	} else {
		s.SendMessage(chatId, "За сегодняшний день SHIB вырос на "+response.ChangePercent+"%\n"+
			"до отметки в "+response.LastPrice+"USDT\n\n")
		s.SendSticker(chatId, "CAACAgIAAxkBAAICd2iVNDnS_14zSxDh_l9Wf6Vb2vGkAAJeAAOtZbwUvL_TIkzK-Ms2BA")
	}

}
