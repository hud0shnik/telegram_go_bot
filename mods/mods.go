package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Структуры для работы API

type cryptoResponse struct {
	Symbol        string `json:"symbol"`
	ChangePercent string `json:"priceChangePercent"`
	LastPrice     string `json:"lastPrice"`
}

type ipApiResponse struct {
	Status      string `json:"status"`
	CountryName string `json:"country"`
	Region      string `json:"regionName"`
	Zip         string `json:"zip"`
}

type dogResponse struct {
	DogUrl string `json:"message"`
}

type infoResponse struct {
	Success       bool   `json:"success"`
	Error         string `json:"error"`
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

type commitsResponse struct {
	Success  bool   `json:"success"`
	Error    string `json:"error"`
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

type osuUserInfo struct {
	Success        bool     `json:"success"`
	Error          string   `json:"error"`
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
func SendOsuInfo(botUrl string, chatId int, username string) {

	// Проверка параметра
	if username == "" {
		SendMsg(botUrl, chatId, "Синтаксис команды:\n\n/osu <b>[id]</b>\n\nПример:\n/osu <b>hud0shnik</b>")
		return
	}

	// Отправка запроса OsuStatsApi
	resp, err := http.Get("https://osustatsapi.vercel.app/api/v2/user?type=string&id=" + username)
	if err != nil {
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		log.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// При хорошем статусе респонса, продолжение выполнения кода
	case 404:
		SendMsg(botUrl, chatId, "Пользователь не найден")
		return
	case 400:
		SendMsg(botUrl, chatId, "Плохой реквест")
		return
	default:
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(osuUserInfo)
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("json.Unmarshal err: %v", err)
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
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
	SendPict(botUrl, chatId, user.AvatarUrl, responseText)

}

// Функция вывода количества коммитов пользователя GitHub
func SendCommits(botUrl string, chatId int, username, date string) {

	// Проверка параметра
	if username == "" {
		SendMsg(botUrl, chatId, "Синтаксис команды:\n\n/commits <b>[id]</b>\n\nПример:\n/commits <b>hud0shnik</b>")
		return
	}

	// Отправка запроса моему API
	resp, err := http.Get("https://githubstatsapi.vercel.app/api/v2/commits?id=" + username + "&date=" + date)

	// Проверка на ошибку
	if err != nil {
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		log.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// При хорошем статусе респонса, продолжение выполнения кода
	case 404:
		SendMsg(botUrl, chatId, "Пользователь не найден")
		return
	case 400:
		SendMsg(botUrl, chatId, "Плохой реквест")
		return
	default:
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(commitsResponse)
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("in SendCommits: json.Unmarshal err: %v", err)
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Если поле пустое, меняет date на "сегодня"
	if date == "" {
		date = "сегодня"
	}

	// Вывод данных пользователю
	switch user.Color {
	case 1:
		SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>", date, user.Commits))
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYwmG11bAfndI1wciswTEVJUEdgB2jAAI5AAOtZbwUdHz8lasybOojBA")
	case 2:
		SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>, неплохо!", date, user.Commits))
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIXWmGyDE1aVXGUY6lcjKxx9bOn0JA1AAJOAAOtZbwUIWzOXysr2zwjBA")
	case 3:
		SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>, отлично!!", date, user.Commits))
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYymG11mMdODUQUZGsQO97V9O0ZLJCAAJeAAOtZbwUvL_TIkzK-MsjBA")
	case 4:
		SendMsg(botUrl, chatId, fmt.Sprintf("Коммитов за <i>%s</i> <b>%d</b>, прекрасно!!!", date, user.Commits))
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIXXGGyDFClr69PKZXJo9dlYMbyilXLAAI1AAOtZbwU9aVxXMUw5eAjBA")
	default:
		SendMsg(botUrl, chatId, "Коммитов нет...")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYG2GzRVNm_d_mVDIOaiLXkGukArlTAAJDAAOtZbwU_-iXZG7hfLsjBA")
	}

}

// Функция вывода информации о пользователе GitHub
func SendGithubInfo(botUrl string, chatId int, username string) {

	// Проверка параметра
	if username == "" {
		SendMsg(botUrl, chatId, "Синтаксис команды:\n\n/github <b>[id]</b>\n\nПример:\n/github <b>hud0shnik</b>")
		return
	}

	// Отправка запроса
	resp, err := http.Get("https://githubstatsapi.vercel.app/api/v2/user?id=" + username)
	if err != nil {
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		log.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Проверка респонса
	switch resp.StatusCode {
	case 200:
		// При хорошем статусе респонса, продолжение выполнения кода
	case 404:
		SendMsg(botUrl, chatId, "Пользователь не найден")
		return
	case 400:
		SendMsg(botUrl, chatId, "Плохой реквест")
		return
	default:
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		return
	}

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(infoResponse)
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("in SendGithubInfo: json.Unmarshal err: %v", err)
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Отправка данных пользователю
	SendPict(botUrl, chatId, user.Avatar,
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

// Функция вывода курса криптовалюты SHIB
func SendCryptoInfo(botUrl string, chatId int) {

	// Отправка запроса
	resp, err := http.Get("https://api2.binance.com/api/v3/ticker/24hr?symbol=SHIBBUSD")
	if err != nil {
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		log.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)
	var response = new(cryptoResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("in SendCryptoInfo: json.Unmarshal err: %v", err)
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Формирование и отправка результата
	if response.ChangePercent[0] == '-' {
		SendMsg(botUrl, chatId, "За сегодняшний день "+response.Symbol+" упал на "+response.ChangePercent[1:]+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, chatId, true)
	} else {
		SendMsg(botUrl, chatId, "За сегодняшний день "+response.Symbol+" вырос на "+response.ChangePercent+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, chatId, false)
	}

}

// Функция нахождения местоположения по IP адресу
func SendIPInfo(botUrl string, chatId int, IP string) {

	// Проверка на пустой IP
	if IP == "" {
		SendMsg(botUrl, chatId, "Чтобы узнать страну по ip, отправьте:\n\n/ip 67.77.77.7")
		return
	}

	// Проверка на localhost
	if IP == "127.0.0.1" {
		SendMsg(botUrl, chatId, "Айпишник локалхоста")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIYLGGzR7310Hqf8K2sljgcQF8kgOpYAAJTAAOtZbwUo9c59oswVBQjBA")
		return
	}

	// Проверка корректности ввода
	ipArray := strings.Split(IP, ".")
	if len(ipArray) != 4 {
		SendMsg(botUrl, chatId, "Не могу обработать этот IP")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}
	for _, ipPart := range ipArray {
		num, err := strconv.Atoi(ipPart)
		if err != nil || num < 0 || num > 255 || (ipPart != fmt.Sprint(num)) {
			SendMsg(botUrl, chatId, "Неправильно набран IP")
			SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4GG13SepKZJisWVrMrzQ9JyRpWFrAAJKAAOtZbwUiXsNXgiPepIjBA")
			return
		}
	}

	// Отправка запроса API
	resp, err := http.Get("http://ip-api.com/json/" + IP)
	if err != nil {
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		log.Printf("http.Get error: %s", err)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := ioutil.ReadAll(resp.Body)
	var response = new(ipApiResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("in SendIPInfo: json.Unmarshal err: %v", err)
		SendMsg(botUrl, chatId, "Внутренняя ошибка")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод сообщения для респонса без страны
	if response.Status != "success" {
		SendMsg(botUrl, chatId, "Не могу найти этот IP")
		SendStck(botUrl, chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод результатов поиска
	SendMsg(botUrl, chatId, "Нашёл! Страна происхождения - "+response.CountryName+" "+"\n"+
		"Регион - "+response.Region+"\n"+
		"Индекс - "+response.Zip+"\n\n"+
		"Мы не храним IP, которые просят проверить пользователи, весь код можно найти на гитхабе.")
	SendStck(botUrl, chatId, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")

}

// Функция генерации псевдослучайных чисел
func random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)

}

// Функция броска монетки
func FlipCoin(botUrl string, chatId int) {
	if random(2) == 0 {
		SendMsg(botUrl, chatId, "Орёл")
		return
	}
	SendMsg(botUrl, chatId, "Решка")
}

// Функция вывода списка всех команд
func Help(botUrl string, chatId int) {
	SendMsg(botUrl, chatId, "Привет👋🏻, вот список команд:"+"\n\n"+
		"/commits <u>username</u> <u>date</u> - коммиты пользователя за день"+"\n\n"+
		"/github <u>username</u> - информация о пользователе GitHub"+"\n\n"+
		"/osu <u>username</u> - информация о пользователе Osu"+"\n\n"+
		"/ip <u>ip_address</u> - узнать страну по ip"+"\n\n"+
		"/crypto - узнать текущий курс криптовалюты SHIB"+"\n\n"+
		"/d <b>n</b> - кинуть <b>n</b>-гранную кость"+"\n\n"+
		"/coin - бросить монетку"+"\n\n"+
		"Также можешь позадавать вопросы, я на них отвечу 🙃")

}

// Функция броска n-гранного кубика
func RollDice(botUrl string, chatId int, parameter string) {

	// Проверка параметра
	if parameter == "" {
		SendMsg(botUrl, chatId, "Пожалуйста укажи количество граней\nНапример /d <b>20</b>")
		return
	}

	// Считывание числа граней
	num, err := strconv.Atoi(parameter)
	if err != nil || num < 1 {
		SendMsg(botUrl, chatId, "Это вообще кубик?🤨")
		return
	}

	// Проверка на d10 (единственный кубик, который имеет грань со значением "0")
	if num == 10 {
		SendMsg(botUrl, chatId, strconv.Itoa(random(10)))
		return
	}

	// Бросок
	SendMsg(botUrl, chatId, strconv.Itoa(1+random(num)))

}

// Функция генерации случайных ответов
func Ball8(botUrl string, chatId int) {

	// Массив ответов
	answers := [10]string{
		"Да, конечно!",
		"100%",
		"Да.",
		"100000000%",
		"Точно да!",
		"Нет, пфф",
		"Нееееееееееет",
		"Точно нет!",
		"Нет, нет, нет",
		"Нет.",
	}

	// Выбор случайного ответа
	SendMsg(botUrl, chatId, answers[random(10)])

}

// Функция проверки всех команд
func Check(botUrl string, chatId int) {

	// Проверка на мой id
	if chatId == viper.GetInt("DanyaChatId") {

		// Вызов функций для тестирования
		SendOsuInfo(botUrl, chatId, "")
		SendCommits(botUrl, chatId, "", "")
		SendGithubInfo(botUrl, chatId, "")
		SendCryptoInfo(botUrl, chatId)
		SendIPInfo(botUrl, chatId, "67.77.77.7")
		SendRandomSticker(botUrl, chatId)

	} else {

		// Вывод для других пользователей
		SendMsg(botUrl, chatId, "Error 403! Beep Boop... Forbidden! Access denied 🤖")

	}

}
