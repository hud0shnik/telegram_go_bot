package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Структуры для работы с Telegram API
type TelegramResponse struct {
	Result []Update `json:"result"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat    Chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker Sticker `json:"sticker"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type Sticker struct {
	File_id string `json:"file_id"`
}

// Структуры для работы с другими API
type RedditResponse struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Nsfw    bool   `json:"nsfw"`
	Spoiler bool   `json:"spoiler"`
}

type CryptoResponse struct {
	Symbol        string `json:"symbol"`
	ChangePercent string `json:"priceChangePercent"`
	LastPrice     string `json:"lastPrice"`
}

type IP2CountryResponse struct {
	CountryName  string `json:"countryName"`
	CountryEmoji string `json:"countryEmoji"`
}

type DogResponse struct {
	DogUrl string `json:"message"`
}

type InfoResponse struct {
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

type CommitsResponse struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

type OsuUserInfo struct {
	Username                 string     `json:"username"`
	Names                    string     `json:"previous_usernames"`
	AvatarUrl                string     `json:"avatar_url"`
	UserID                   string     `json:"id"`
	CountryCode              string     `json:"country_code"`
	GlobalRank               string     `json:"global_rank"`
	CountryRank              string     `json:"country_rank"`
	PP                       string     `json:"pp"`
	PlayTime                 string     `json:"play_time"`
	SSH                      string     `json:"ssh"`
	SS                       string     `json:"ss"`
	SH                       string     `json:"sh"`
	S                        string     `json:"s"`
	A                        string     `json:"a"`
	RankedScore              string     `json:"ranked_score"`
	Accuracy                 string     `json:"accuracy"`
	PlayCount                string     `json:"play_count"`
	TotalScore               string     `json:"total_score"`
	TotalHits                string     `json:"total_hits"`
	MaximumCombo             string     `json:"maximum_combo"`
	Replays                  string     `json:"replays"`
	Level                    string     `json:"level"`
	SupportLvl               string     `json:"support_level"`
	DefaultGroup             string     `json:"default_group"`
	IsOnline                 string     `json:"is_online"`
	IsActive                 string     `json:"is_active"`
	IsDeleted                string     `json:"is_deleted"`
	IsBot                    string     `json:"is_bot"`
	IsSupporter              string     `json:"is_supporter"`
	LastVisit                string     `json:"last_visit"`
	ProfileColor             string     `json:"profile_color"`
	RankedBeatmapsetCount    string     `json:"ranked_beatmapset_count"`
	PendingBeatmapsetCount   string     `json:"pending_beatmapset_count"`
	PmFriendsOnly            string     `json:"pm_friends_only"`
	GraveyardBeatmapsetCount string     `json:"graveyard_beatmapset_count"`
	BeatmapPlaycountsCount   string     `json:"beatmap_playcounts_count"`
	CommentsCount            string     `json:"comments_count"`
	FavoriteBeatmapsetCount  string     `json:"favorite_beatmapset_count"`
	GuestBeatmapsetCount     string     `json:"guest_beatmapset_count"`
	FollowerCount            string     `json:"follower_count"`
	BestBeatMap              OsuBeatMap `json:"best_beat_map"`
}

type OsuBeatMap struct {
	DifficultyRating string `json:"difficulty_rating"`
	Id               string `json:"id"`
	Rank             string `json:"rank"`
	EndedAt          string `json:"ended_at"`
	Accuracy         string `json:"accuracy"`
	MaximumCombo     string `json:"maximum_combo"`
	PP               string `json:"pp"`
	Passed           string `json:"passed"`
	TotalScore       string `json:"total_score"`
	LegacyPerfect    string `json:"legacy_perfect"`
	Replay           string `json:"replay"`
	Mode             string `json:"mode"`
	Status           string `json:"status"`
	TotalLength      string `json:"total_length"`
	Ar               string `json:"ar"`
	Bpm              string `json:"bpm"`
	Convert          string `json:"convert"`
	CountCircles     string `json:"count_circles"`
	CountSliders     string `json:"count_sliders"`
	CountSpinners    string `json:"count_spinners"`
	Cs               string `json:"cs"`
	DeletedAt        string `json:"deleted_at"`
	Drain            string `json:"drain"`
	HitLength        string `json:"hit_length"`
	IsScoreable      string `json:"is_scoreable"`
	LastUpdated      string `json:"last_updated"`
	ModeInt          string `json:"mode_int"`
	PassCount        string `json:"pass_count"`
	PlayCount        string `json:"play_count"`
	Ranked           string `json:"ranked"`
	Url              string `json:"url"`
	Checksum         string `json:"checksum"`
	Creator          string `json:"creator"`
	FavoriteCount    string `json:"favorite_count"`
	Hype             string `json:"hype"`
	Nsfw             string `json:"nsfw"`
	Offset           string `json:"offset"`
	Spotlight        string `json:"spotlight"`
}

// Функция вывода списка всех команд
func Help(botUrl string, update Update) {
	SendMsg(botUrl, update, "Привет👋🏻, вот список команд:"+
		"\n\n/commits username date - коммиты пользователя за день"+
		"\n\n/github username - информация о пользователе GitHub"+
		"\n\n/osu username - информация о пользователе Osu"+
		"\n\n/ip 67.77.77.7 - узнать страну по ip"+
		"\n\n/crypto - узнать текущий курс криптовалюты SHIB"+
		"\n\n/d20 - кинуть д20, вместо 20 можно поставить любое число"+
		"\n\n/coin - подбросить монетку"+
		"\n\n/meme - мем с Reddit"+
		"\n\n/cat и /parrot - картинка кота или попугая "+
		"\n\nТакже можешь позадавать вопросы, я на них отвечу 🙃")
}

// Функция генерации псевдослучайных чисел
func Random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

// Функция броска n-гранного кубика
func Dice(msg string) string {

	// Считывание числа граней
	num, err := strconv.Atoi(msg[2:])

	// Проверки на невозможное количество граней
	if err != nil {
		return "Это вообще кубик?🤨"
	}
	if num < 1 {
		return "как я по твоему кину такой кубик? Через четвёртое пространство?🤨"
	}

	// Проверка на d10 (единственный кубик, который имеет грань со значением "0")
	if num == 10 {
		return strconv.Itoa(Random(10))
	}

	// Бросок
	return strconv.Itoa(1 + Random(num))

}

// Функция генерации случайных ответов
func Ball8(botUrl string, update Update) {

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
	SendMsg(botUrl, update, answers[Random(len(answers))])

}

// Функция броска монетки
func Coin(botUrl string, update Update) {
	if Random(2) == 0 {
		SendMsg(botUrl, update, "Орёл")
	} else {
		SendMsg(botUrl, update, "Решка")
	}
}

// Функция отправки случайного поста с Reddit
func SendFromReddit(botUrl string, update Update, board string) error {

	// Отправка реквеста
	url := "https://meme-api.herokuapp.com/gimme/" + board
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)

	// Проверка на ошибку
	if err != nil {
		fmt.Println("Meme API error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return err
	}

	// Запись респонса
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response = new(RedditResponse)
	json.Unmarshal(body, &response)

	// Проверка на запрещёнку
	if response.Nsfw || response.Spoiler {
		response.Url = "https://belikebill.ga/billgen-API.php?default=1"
		response.Title = "Картинка оказалась со спойлером или nsfw-контентом, поэтому я заменил её на это"
	}

	// Формирование сообщения
	botImageMessage := SendPhoto{
		ChatId:   update.Message.Chat.ChatId,
		PhotoUrl: response.Url,
		Caption:  response.Title,
	}

	// Отправка результата
	SendPict(botUrl, update, botImageMessage)
	return nil

}

// Функция вывода курса криптовалюты SHIB
func SendCryptoData(botUrl string, update Update) {

	// Отправка реквеста
	req, _ := http.NewRequest("GET", "https://api2.binance.com/api/v3/ticker/24hr?symbol=SHIBBUSD", nil)
	res, err := http.DefaultClient.Do(req)

	// Проверка на ошибку
	if err != nil {
		fmt.Println("Binance API error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}

	// Запись респонса
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response = new(CryptoResponse)
	json.Unmarshal(body, &response)

	// Формирование и отправка результата
	if response.ChangePercent[0] == '-' {
		SendMsg(botUrl, update, "За сегодняшний день "+response.Symbol+" упал на "+response.ChangePercent[1:]+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, update, true)
	} else {
		SendMsg(botUrl, update, "За сегодняшний день "+response.Symbol+" вырос на "+response.ChangePercent+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, update, false)
	}

}

// Функция проверки всех команд
func Check(botUrl string, update Update) {

	// Проверка на мой id
	if update.Message.Chat.ChatId == viper.GetInt("DanyaChatId") {

		// Временная метка начала проверки
		start := time.Now()

		// Вызов всех команд
		SendCryptoData(botUrl, update)
		SendFromReddit(botUrl, update, "")
		Coin(botUrl, update)
		Help(botUrl, update)
		SendCommits(botUrl, update, "hud0shnik")
		SendMsg(botUrl, update, Dice("/d20"))
		Ball8(botUrl, update)
		SendRandomSticker(botUrl, update)
		SendFromReddit(botUrl, update, "parrots")

		// Отправка ошибок
		for i := 1; i < 7; i++ {
			SendErrorMessage(botUrl, update, i)
		}

		// Отправка результата
		SendMsg(botUrl, update, "Проверка заняла "+time.Since(start).String())
		return
	}

	// Вывод для других пользователей
	SendMsg(botUrl, update, "Error 403! Beep Boop... Forbidden! Access denied 🤖")

}

// Функция отправки сообщений об ошибках
func SendErrorMessage(botUrl string, update Update, errorCode int) {

	// Генерация текста ошибки по коду
	var result string
	switch errorCode {
	case 1:
		result = "Ошибка работы API"
	case 2:
		result = "Ошибка работы json.Marshal()"
	case 3:
		result = "Ошибка работы метода SendSticker"
	case 4:
		result = "Ошибка работы метода SendPhoto"
	case 5:
		result = "Ошибка работы метода SendMessage"
	case 6:
		result = "Ошибка работы stickers.json"
	default:
		result = "Неизвестная ошибка"
	}

	// Анонимное оповещение меня
	var updateDanya Update
	updateDanya.Message.Chat.ChatId = viper.GetInt("DanyaChatId")
	SendMsg(botUrl, updateDanya, "Дань, тут у одного из пользователей "+result+", надеюсь он скоро тебе о ней напишет.")

	// Вывод ошибки пользователю с просьбой связаться со мной для её устранения
	result += ", пожалуйста свяжитесь с моим создателем для устранения проблемы \n\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik"
	SendMsg(botUrl, update, result)
}

// Функция вывода информации о пользователе GitHub
func SendInfo(botUrl string, update Update, parameters string) {

	// Отправка запроса моему API
	resp, err := http.Get("https://hud0shnikgitapi.herokuapp.com/user/" + parameters)

	// Проверка на ошибку
	if err != nil {
		fmt.Println("GithubGoAPI error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(InfoResponse)
	json.Unmarshal(body, &user)

	// Отправка данных пользователю
	SendPict(botUrl, update, SendPhoto{
		PhotoUrl: user.Avatar,
		ChatId:   update.Message.Chat.ChatId,
		Caption: "Информация о " + user.Username + ":\n" +
			"Имя " + user.Name + "\n" +
			"Поставленных звезд " + user.Stars + " ⭐\n" +
			"Подписчиков " + user.Followers + " 🤩\n" +
			"Подписок " + user.Following + " 🕵️\n" +
			"Репозиториев " + user.Repositories + " 📘\n" +
			"Пакетов " + user.Packages + " 📦\n" +
			"Контрибуций за год " + user.Contributions + " 🟩\n" +
			"Ссылка на аватар:\n " + user.Avatar,
	})

}

// Функция вывода количества коммитов пользователя GitHub
func SendCommits(botUrl string, update Update, parameters string) {

	// Индекс пробела и дата
	i, date := 0, ""

	// Поиск конца юзернейма и начала даты
	for ; i < len(parameters); i++ {
		if parameters[i] == ' ' {
			break
		}
	}

	// Если дата задана, записывает её
	if i != len(parameters) {
		date = parameters[i+1:]
	}

	// Отправка запроса моему API
	resp, err := http.Get("https://hud0shnikgitapi.herokuapp.com/commits/" + parameters[:i] + "/" + date)

	// Проверка на ошибку
	if err != nil {
		fmt.Println("GithubGoAPI error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(CommitsResponse)
	json.Unmarshal(body, &user)

	// Если поле пустое, меняет date на "сегодня"
	if date == "" {
		date = "сегодня"
	}

	// Вывод данных пользователю
	switch user.Color {
	case 1:
		SendMsg(botUrl, update, "Коммитов за "+date+" "+strconv.Itoa(user.Commits))
		SendStck(botUrl, update, "CAACAgIAAxkBAAIYwmG11bAfndI1wciswTEVJUEdgB2jAAI5AAOtZbwUdHz8lasybOojBA")
	case 2:
		SendMsg(botUrl, update, "Коммитов за "+date+" "+strconv.Itoa(user.Commits)+", неплохо!")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIXWmGyDE1aVXGUY6lcjKxx9bOn0JA1AAJOAAOtZbwUIWzOXysr2zwjBA")
	case 3:
		SendMsg(botUrl, update, "Коммитов за "+date+" "+strconv.Itoa(user.Commits)+", отлично!!")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIYymG11mMdODUQUZGsQO97V9O0ZLJCAAJeAAOtZbwUvL_TIkzK-MsjBA")
	case 4:
		SendMsg(botUrl, update, "Коммитов за "+date+" "+strconv.Itoa(user.Commits)+", прекрасно!!!")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIXXGGyDFClr69PKZXJo9dlYMbyilXLAAI1AAOtZbwU9aVxXMUw5eAjBA")
	default:
		SendMsg(botUrl, update, "Коммитов нет")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIYG2GzRVNm_d_mVDIOaiLXkGukArlTAAJDAAOtZbwU_-iXZG7hfLsjBA")
	}

}

// Функция нахождения местоположения по IP адресу
func CheckIPAdress(botUrl string, update Update, IP string) {

	// Проверка на localhost
	if IP == "127.0.0.1" {
		SendMsg(botUrl, update, "Айпишник локалхоста, ага")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIYLGGzR7310Hqf8K2sljgcQF8kgOpYAAJTAAOtZbwUo9c59oswVBQjBA")
		return
	}

	// Проверка корректности ввода
	ipArray := strings.Split(IP, ".")
	if len(ipArray) != 4 {
		SendMsg(botUrl, update, "Не могу обработать этот IP")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}
	for _, ipPart := range ipArray {
		num, err := strconv.Atoi(ipPart)
		if err != nil || num < 0 || num > 255 || (ipPart != fmt.Sprint(num)) {
			SendMsg(botUrl, update, "Неправильно набран IP")
			SendStck(botUrl, update, "CAACAgIAAxkBAAIY4GG13SepKZJisWVrMrzQ9JyRpWFrAAJKAAOtZbwUiXsNXgiPepIjBA")
			return
		}
	}

	// Отправка запроса
	SendMsg(botUrl, update, "Ищу...")
	url := "https://api.ip2country.info/ip?" + IP
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("IP2Country API error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}

	// Запись респонса
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response = new(IP2CountryResponse)
	json.Unmarshal(body, &response)

	// Вывод сообщения для спец айпишников
	if response.CountryName == "" {
		SendMsg(botUrl, update, "Не могу найти этот IP")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод результатов поиска
	SendMsg(botUrl, update, "Нашёл! Страна происхождения - "+response.CountryName+" "+response.CountryEmoji+
		"\n\nМы не храним IP, которые просят проверить пользователи, весь код можно найти на гитхабе.")
	SendStck(botUrl, update, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")

}

// Функция инициализации конфига (всех токенов)
func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

// Функция вывода информации о пользователе Osu!
func SendOsuInfo(botUrl string, update Update, parameters string) {

	// Отправка запроса моему API
	resp, err := http.Get("https://osu-stats-api.herokuapp.com/user/" + parameters)

	// Проверка на ошибку
	if err != nil {
		fmt.Println("OsuStatsAPI error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}

	// Запись респонса
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(OsuUserInfo)
	json.Unmarshal(body, &user)

	// Формирование текста респонса
	responseText := "Информация о " + user.Username + "\n"

	if user.Names != "" {
		responseText += "Aka " + user.Names + "\n"
	}

	responseText += "Код страны " + user.CountryCode + "\n" +
		"Рейтинг в мире " + user.GlobalRank + "\n" +
		"Рейтинг в стране " + user.CountryRank + "\n" +
		"Точность попаданий " + user.Accuracy + "%\n" +
		"PP " + user.PP + "\n" +
		"-------карты---------\n" +
		"SSH: " + user.SSH + "\n" +
		"SH: " + user.SH + "\n" +
		"SS: " + user.SS + "\n" +
		"S: " + user.S + "\n" +
		"A: " + user.A + "\n" +
		"---------------------------\n" +
		"Рейтинговые очки " + user.RankedScore + "\n" +
		"Количество игр " + user.PlayCount + "\n" +
		"Всего очков " + user.TotalScore + "\n" +
		"Всего попаданий " + user.TotalHits + "\n" +
		"Максимальное комбо " + user.MaximumCombo + "\n" +
		"Реплеев просмотрено другими " + user.Replays + "\n" +
		"Уровень " + user.Level + "\n" +
		"---------------------------\n" +
		"Время в игре " + user.PlayTime + "\n" +
		"Уровень подписки " + user.SupportLvl + "\n"

	if user.IsOnline == "true" {
		responseText += "Сейчас онлайн \n"
	} else {
		responseText += "Сейчас не в сети \n"
	}

	if user.IsActive == "true" {
		responseText += "Аккаунт активен \n"
	} else {
		responseText += "Аккаунт не активен \n"
	}

	if user.IsDeleted == "true" {
		responseText += "Аккаунт удалён \n"
	}

	if user.IsBot == "true" {
		responseText += "Это аккаунт бота \n"
	}

	if user.ProfileColor != "null" {
		responseText += "Цвет профиля" + user.ProfileColor + "\n"
	}

	// Отправка данных пользователю
	SendPict(botUrl, update, SendPhoto{
		PhotoUrl: user.AvatarUrl,
		ChatId:   update.Message.Chat.ChatId,
		Caption:  responseText,
	})

}
