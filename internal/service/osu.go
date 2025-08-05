package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Структура респонса osustatsapi
type osuUserInfo struct {
	Username       string   `json:"username"`
	Names          []string `json:"previous_usernames"`
	AvatarUrl      string   `json:"avatar_url"`
	CountryCode    string   `json:"country_code"`
	GlobalRank     string   `json:"global_rank"`
	CountryRank    string   `json:"country_rank"`
	PP             string   `json:"pp"`
	PlayTime       string   `json:"play_time"`
	SSH            string   `json:"ssh"`
	SS             string   `json:"ss"`
	SH             string   `json:"sh"`
	S              string   `json:"s"`
	A              string   `json:"a"`
	RankedScore    string   `json:"ranked_score"`
	Accuracy       string   `json:"accuracy"`
	PlayCount      string   `json:"play_count"`
	TotalScore     string   `json:"total_score"`
	TotalHits      string   `json:"total_hits"`
	MaximumCombo   string   `json:"maximum_combo"`
	Replays        string   `json:"replays"`
	Level          string   `json:"level"`
	SupportLvl     string   `json:"support_level"`
	DefaultGroup   string   `json:"default_group"`
	IsOnline       string   `json:"is_online"`
	IsActive       string   `json:"is_active"`
	IsDeleted      string   `json:"is_deleted"`
	IsNat          string   `json:"is_nat"`
	IsModerator    string   `json:"is_moderator"`
	IsBot          string   `json:"is_bot"`
	IsSilenced     string   `json:"is_silenced"`
	IsRestricted   string   `json:"is_restricted"`
	IsLimitedBn    string   `json:"is_limited_bn"`
	IsSupporter    string   `json:"is_supporter"`
	ProfileColor   string   `json:"profile_color"`
	PmFriendsOnly  string   `json:"pm_friends_only"`
	PostCount      string   `json:"post_count"`
	FollowersCount string   `json:"follower_count"`
	Medals         string   `json:"medals"`
}

// Функция вывода информации о пользователе Osu
func (s *BotService) SendOsuInfo(chatId int64, username string) {

	// Проверка параметра
	if username == "" {
		s.api.Send(tgbotapi.NewMessage(chatId, "Синтаксис команды:\n\n/osu <b>[id]</b>\n\nПример:\n/osu <b>hud0shnik</b>"))
		return
	}

	// Формирование url
	apiUrl := "https://osustatsapi.vercel.app/api/user?type=string&id=" + username

	// Отправка запроса OsuStatsApi
	resp, err := http.Get(apiUrl)
	if err != nil {
		s.api.Send(tgbotapi.NewMessage(chatId, "Внутренняя ошибка"))
		slog.Error("http.Get error", "error", err, "request", apiUrl)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// Продолжение выполнения кода
	case 404:
		s.api.Send(tgbotapi.NewMessage(chatId, "Пользователь не найден"))
		return
	case 400:
		s.api.Send(tgbotapi.NewMessage(chatId, "Плохой реквест"))
		return
	default:
		s.api.Send(tgbotapi.NewMessage(chatId, "Внутренняя ошибка"))
		return
	}

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var user = new(osuUserInfo)
	err = json.Unmarshal(body, &user)
	if err != nil {
		slog.Error("json.Unmarshal err", "error", err, "request", apiUrl)
		s.api.Send(tgbotapi.NewMessage(chatId, "Внутренняя ошибка"))
		s.api.Send(tgbotapi.NewSticker(chatId, tgbotapi.FileID("CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")))
		return
	}

	// Формирование текста респонса

	responseText := "Информация о <b>" + user.Username + "</b>\n"

	if len(user.Names) != 0 {
		responseText += "Aka " + user.Names[0] + "\n"
	}

	responseText += "Код страны " + user.CountryCode + "\n" +
		"Рейтинг в мире <b>" + user.GlobalRank + "</b>\n" +
		"Рейтинг в стране <b>" + user.CountryRank + "</b>\n" +
		"Точность попаданий <b>" + user.Accuracy + "%</b>\n" +
		"PP <b>" + user.PP + "</b>\n" +
		"-------карты---------\n" +
		"SSH: <b>" + user.SSH + "</b>\n" +
		"SH: <b>" + user.SH + "</b>\n" +
		"SS: <b>" + user.SS + "</b>\n" +
		"S: <b>" + user.S + "</b>\n" +
		"A: <b>" + user.A + "</b>\n" +
		"---------------------------\n" +
		"Рейтинговые очки <b>" + user.RankedScore + "</b>\n" +
		"Количество игр <b>" + user.PlayCount + "</b>\n" +
		"Всего очков <b>" + user.TotalScore + "</b>\n" +
		"Всего попаданий <b>" + user.TotalHits + "</b>\n" +
		"Максимальное комбо <b>" + user.MaximumCombo + "</b>\n" +
		"Реплеев просмотрено другими <b>" + user.Replays + "</b>\n" +
		"Уровень <b>" + user.Level + "</b>\n" +
		"---------------------------\n" +
		"Время в игре <i>" + user.PlayTime + "</i>\n" +
		"Достижений <i>" + user.Medals + "</i>\n"

	if user.SupportLvl != "0" {
		responseText += "Уровень подписки " + user.SupportLvl + "\n"
	}

	if user.PostCount != "0" {
		responseText += "Постов на форуме " + user.PostCount + "\n"
	}

	if user.FollowersCount != "0" {
		responseText += "Подписчиков " + user.FollowersCount + "\n"
	}

	if user.IsOnline == "true" {
		responseText += "Сейчас онлайн\n"
	} else {
		responseText += "Сейчас не в сети\n"
	}

	if user.IsActive == "true" {
		responseText += "Аккаунт активен\n"
	} else {
		responseText += "Аккаунт не активен\n"
	}

	if user.IsDeleted == "true" {
		responseText += "Аккаунт удалён\n"
	}

	if user.IsBot == "true" {
		responseText += "Это аккаунт бота\n"
	}

	if user.IsNat == "true" {
		responseText += "Это аккаунт члена команды оценки номинаций\n"
	}

	if user.IsModerator == "true" {
		responseText += "Это аккаунт модератора\n"
	}

	if user.ProfileColor != "" {
		responseText += "Цвет профиля <b>" + user.ProfileColor + "<b>\n"
	}

	// Отправка данных пользователю
	photo := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(user.AvatarUrl))
	photo.Caption = responseText
	photo.ParseMode = "HTML"

	// Отправка фото
	s.SendPhoto(chatId, user.AvatarUrl, responseText)

}
