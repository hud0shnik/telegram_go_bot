package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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
func (s *BotService) SendGithubInfo(chatId int64, username string) {

	// Проверка параметра
	if username == "" {
		s.SendMessage(chatId, "Синтаксис команды:\n\n/github <b>[id]</b>\n\nПример:\n/github <b>hud0shnik</b>")
		return
	}

	apiUrl := "https://githubstatsapi.vercel.app/api/v2/user?type=string&id=" + username

	// Отправка запроса
	resp, err := http.Get(apiUrl)
	if err != nil {
		slog.Error("http.Get error", "error", err, "request", apiUrl)
		s.SendMessage(chatId, "Внутренняя ошибка")
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// Продолжение выполнения кода
	case 404:
		s.SendMessage(chatId, "Пользователь не найден")
		return
	case 400:
		s.SendMessage(chatId, "Плохой реквест")
		return
	default:
		s.SendMessage(chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var user = new(infoResponse)
	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error("in SendGithubInfo: json.Unmarshal error", "error", err, "request", apiUrl)
		s.SendMessage(chatId, "Внутренняя ошибка")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Отправка данных пользователю
	s.SendPhoto(chatId, user.Avatar,
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

// TODO: Добавить функцию получения количества коммитов за определенную дату

// Функция вывода количества коммитов пользователя GitHub
func (s *BotService) SendCommits(chatId int64, username string) {

	// Проверка параметра
	if username == "" {
		s.SendMessage(chatId, "Синтаксис команды:\n\n/commits <b>[id]</b>\n\nПример:\n/commits <b>hud0shnik</b>")
		return
	}

	apiUrl := "https://githubstatsapi.vercel.app/api/v2/commits?id=" + username

	// Отправка запроса моему API
	resp, err := http.Get(apiUrl)

	// Проверка на ошибку
	if err != nil {
		s.SendMessage(chatId, "Внутренняя ошибка")
		slog.Error("http.Get error", "error", err, "request", apiUrl)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// Продолжение выполнения кода
	case 404:
		s.SendMessage(chatId, "Пользователь не найден")
		return
	case 400:
		s.SendMessage(chatId, "Плохой реквест")
		return
	default:
		s.SendMessage(chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var user = new(commitsResponse)
	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error("in SendCommits: json.Unmarshal error", "error", err, "request", apiUrl)
		s.SendMessage(chatId, "Внутренняя ошибка")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод данных пользователю
	switch user.Color {
	case 1:
		s.SendMessage(chatId, fmt.Sprintf("Коммитов за <i>сегодня</i> <b>%d</b>", user.Commits))
		s.SendSticker(chatId, "CAACAgIAAxkBAAIYwmG11bAfndI1wciswTEVJUEdgB2jAAI5AAOtZbwUdHz8lasybOojBA")
	case 2:
		s.SendMessage(chatId, fmt.Sprintf("Коммитов за <i>сегодня</i> <b>%d</b>, неплохо!", user.Commits))
		s.SendSticker(chatId, "CAACAgIAAxkBAAIXWmGyDE1aVXGUY6lcjKxx9bOn0JA1AAJOAAOtZbwUIWzOXysr2zwjBA")
	case 3:
		s.SendMessage(chatId, fmt.Sprintf("Коммитов за <i>сегодня</i> <b>%d</b>, отлично!!", user.Commits))
		s.SendSticker(chatId, "CAACAgIAAxkBAAIYymG11mMdODUQUZGsQO97V9O0ZLJCAAJeAAOtZbwUvL_TIkzK-MsjBA")
	case 4:
		s.SendMessage(chatId, fmt.Sprintf("Коммитов за <i>сегодня</i> <b>%d</b>, прекрасно!!!", user.Commits))
		s.SendSticker(chatId, "CAACAgIAAxkBAAIXXGGyDFClr69PKZXJo9dlYMbyilXLAAI1AAOtZbwU9aVxXMUw5eAjBA")
	default:
		s.SendMessage(chatId, "Коммитов нет...")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIYG2GzRVNm_d_mVDIOaiLXkGukArlTAAJDAAOtZbwU_-iXZG7hfLsjBA")
	}

}
