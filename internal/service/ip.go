package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

// Структура респонса ip-api
type ipApiResponse struct {
	Status      string `json:"status"`
	CountryName string `json:"country"`
	Region      string `json:"regionName"`
	Zip         string `json:"zip"`
}

// Функция поиска местоположения по IP
func (s *BotService) SendIPInfo(chatId int64, IP string) {

	// Проверка на пустой IP
	if IP == "" {
		s.SendMessage(chatId, "Чтобы узнать страну по ip, отправьте:\n\n/ip 67.77.77.7")
		return
	}

	// Проверка на localhost
	if IP == "127.0.0.1" {
		s.SendMessage(chatId, "Айпишник локалхоста")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIYLGGzR7310Hqf8K2sljgcQF8kgOpYAAJTAAOtZbwUo9c59oswVBQjBA")
		return
	}

	// Проверка корректности ввода
	ipArray := strings.Split(IP, ".")
	if len(ipArray) != 4 {
		s.SendMessage(chatId, "Не могу обработать этот IP")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}
	for _, ipPart := range ipArray {
		num, err := strconv.Atoi(ipPart)
		if err != nil || num < 0 || num > 255 || (ipPart != fmt.Sprint(num)) {
			s.SendMessage(chatId, "Неправильно набран IP")
			s.SendSticker(chatId, "CAACAgIAAxkBAAIY4GG13SepKZJisWVrMrzQ9JyRpWFrAAJKAAOtZbwUiXsNXgiPepIjBA")
			return
		}
	}

	// Формирование url
	apiUrl := "http://ip-api.com/json/" + IP

	// Отправка запроса API
	resp, err := http.Get(apiUrl)
	if err != nil {
		s.SendMessage(chatId, "Внутренняя ошибка")
		slog.Error("http.Get error", "error", err, "request", apiUrl)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var response = new(ipApiResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		slog.Error("json.Unmarshal err", "error", err, "request", apiUrl)
		s.SendMessage(chatId, "Внутренняя ошибка")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод сообщения для респонса без страны
	if response.Status != "success" {
		s.SendMessage(chatId, "Не могу найти этот IP")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод результатов поиска
	s.SendMessage(chatId, "Нашёл! Страна происхождения - "+response.CountryName+" "+"\n"+
		"Регион - "+response.Region+"\n"+
		"Индекс - "+response.Zip+"\n\n"+
		"Мы не храним IP, которые просят проверить пользователи, весь код можно найти на гитхабе.")
	s.SendSticker(chatId, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")

}
