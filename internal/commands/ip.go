package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/hud0shnik/telegram_go_bot/internal/send"
	"github.com/sirupsen/logrus"
)

// Структура респонса ip-api
type ipApiResponse struct {
	Status      string `json:"status"`
	CountryName string `json:"country"`
	Region      string `json:"regionName"`
	Zip         string `json:"zip"`
}

// Функция поиска местоположения по IP
func SendIPInfo(botUrl string, chatId int, IP string) {

	// Проверка на пустой IP
	if IP == "" {
		send.SendMsg(botUrl, chatId, "Чтобы узнать страну по ip, отправьте:\n\n/ip 67.77.77.7")
		return
	}

	// Проверка на localhost
	if IP == "127.0.0.1" {
		send.SendMsg(botUrl, chatId, "Айпишник локалхоста")
		send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYLGGzR7310Hqf8K2sljgcQF8kgOpYAAJTAAOtZbwUo9c59oswVBQjBA")
		return
	}

	// Проверка корректности ввода
	ipArray := strings.Split(IP, ".")
	if len(ipArray) != 4 {
		send.SendMsg(botUrl, chatId, "Не могу обработать этот IP")
		send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}
	for _, ipPart := range ipArray {
		num, err := strconv.Atoi(ipPart)
		if err != nil || num < 0 || num > 255 || (ipPart != fmt.Sprint(num)) {
			send.SendMsg(botUrl, chatId, "Неправильно набран IP")
			send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4GG13SepKZJisWVrMrzQ9JyRpWFrAAJKAAOtZbwUiXsNXgiPepIjBA")
			return
		}
	}

	// Отправка запроса API
	resp, err := http.Get("http://ip-api.com/json/" + IP)
	if err != nil {
		send.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		logrus.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)
	var response = new(ipApiResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Printf("in SendIPInfo: json.Unmarshal err: %v", err)
		send.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод сообщения для респонса без страны
	if response.Status != "success" {
		send.SendMsg(botUrl, chatId, "Не могу найти этот IP")
		send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод результатов поиска
	send.SendMsg(botUrl, chatId, "Нашёл! Страна происхождения - "+response.CountryName+" "+"\n"+
		"Регион - "+response.Region+"\n"+
		"Индекс - "+response.Zip+"\n\n"+
		"Мы не храним IP, которые просят проверить пользователи, весь код можно найти на гитхабе.")
	send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")

}
