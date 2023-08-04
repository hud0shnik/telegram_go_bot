package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
	"github.com/sirupsen/logrus"
)

// Структура статистики пользователя
type infoResponse struct {
	Username      string `json:"username"`
	Name          string `json:"name"`
	Followers     string `json:"followers"`
	Following     string `json:"following"`
	Repositories  string `json:"repositories"`
	Packages      string `json:"packages"`
	Stars         string `json:"stars"`
	Contributions string `json:"contributions"`
	Avatar        string `json:"avatar"`
}

// Структура количества коммитов
type commitsResponse struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функция вывода информации о пользователе GitHub
func SendGithubInfo(botUrl string, chatId int, username string) {

	// Проверка параметра
	if username == "" {
		telegram.SendMsg(botUrl, chatId, "Синтаксис команды:\n\n/github <b>[id]</b>\n\nПример:\n/github <b>hud0shnik</b>")
		return
	}

	// Отправка запроса
	resp, err := http.Get("https://githubstatsapi.vercel.app/api/v2/user?type=string&id=" + username)
	if err != nil {
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		logrus.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// Продолжение выполнения кода
	case 404:
		telegram.SendMsg(botUrl, chatId, "Пользователь не найден")
		return
	case 400:
		telegram.SendMsg(botUrl, chatId, "Плохой реквест")
		return
	default:
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var user = new(infoResponse)
	err = json.Unmarshal(body, &user)
	if err != nil {
		logrus.Printf("in SendGithubInfo: json.Unmarshal err: %v", err)
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Отправка данных пользователю
	telegram.SendPict(botUrl, chatId, user.Avatar,
		"Информация о "+user.Username+":\n"+
			"Имя "+user.Name+"\n"+
			"Поставленных звезд "+user.Stars+" ⭐\n"+
			"Подписчиков "+user.Followers+" 🤩\n"+
			"Подписок "+user.Following+" 🕵️\n"+
			"Репозиториев "+user.Repositories+" 📘\n"+
			"Пакетов "+user.Packages+" 📦\n"+
			"Контрибуций за год "+user.Contributions+" 🟩\n"+
			"Ссылка на аватар:\n "+user.Avatar)

}

// Функция вывода количества коммитов пользователя GitHub
func SendCommits(botUrl string, chatId int, username, date string) {

	// Проверка параметра
	if username == "" {
		telegram.SendMsg(botUrl, chatId, "Синтаксис команды:\n\n/commits <b>[id]</b>\n\nПример:\n/commits <b>hud0shnik</b>")
		return
	}

	// Отправка запроса моему API
	resp, err := http.Get("https://githubstatsapi.vercel.app/api/v2/commits?id=" + username + "&date=" + date)

	// Проверка на ошибку
	if err != nil {
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		logrus.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// Продолжение выполнения кода
	case 404:
		telegram.SendMsg(botUrl, chatId, "Пользователь не найден")
		return
	case 400:
		telegram.SendMsg(botUrl, chatId, "Плохой реквест")
		return
	default:
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var user = new(commitsResponse)
	err = json.Unmarshal(body, &user)
	if err != nil {
		logrus.Printf("in SendCommits: json.Unmarshal err: %v", err)
		telegram.SendMsg(botUrl, chatId, "Внутренняя ошибка")
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Если поле пустое, меняет date на "сегодня"
	if date == "" {
		date = "сегодня"
	}

	// Вывод данных пользователю
	switch user.Color {
	case 1:
		telegram.SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>", date, user.Commits))
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYwmG11bAfndI1wciswTEVJUEdgB2jAAI5AAOtZbwUdHz8lasybOojBA")
	case 2:
		telegram.SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>, неплохо!", date, user.Commits))
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIXWmGyDE1aVXGUY6lcjKxx9bOn0JA1AAJOAAOtZbwUIWzOXysr2zwjBA")
	case 3:
		telegram.SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>, отлично!!", date, user.Commits))
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYymG11mMdODUQUZGsQO97V9O0ZLJCAAJeAAOtZbwUvL_TIkzK-MsjBA")
	case 4:
		telegram.SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>, прекрасно!!!", date, user.Commits))
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIXXGGyDFClr69PKZXJo9dlYMbyilXLAAI1AAOtZbwU9aVxXMUw5eAjBA")
	default:
		telegram.SendMsg(botUrl, chatId, "Коммитов нет...")
		telegram.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYG2GzRVNm_d_mVDIOaiLXkGukArlTAAJDAAOtZbwU_-iXZG7hfLsjBA")
	}

}
