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

type StatsResponse struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Stars    int    `json:"stars"`
}

type CommitsResponse struct {
	Date     string `json:"date"`
	Username string `json:"username"`
	Commits  int    `json:"commits"`
	Color    int    `json:"color"`
}

// Функции бота
// Вывод списка всех команд
func Help(botUrl string, update Update) {
	SendMsg(botUrl, update, "Привет👋🏻, вот список команд:"+
		"\n\n/ip 67.77.77.7 - узнать страну по ip"+
		"\n\n/git username  date - коммиты пользователя за день"+
		"\n\n/info username - информация об аккаунте на GitHub"+
		"\n\n/crypto - узнать текущий курс криптовалюты SHIB"+
		"\n\n/d20 - кинуть д20, вместо 20 можно поставить любое число"+
		"\n\n/coin - подбросить монетку"+
		"\n\n/meme - мем с реддита (смотреть на свой страх и риск, я за этот контент не отвечаю 😅)"+
		"\n\n/cat и /parrot - картинка кота или попугая "+
		"\n\nТакже можешь позадовать вопросы, я на них отвечу 🙃")
}

// Функция, реализующая бросок n-гранного кубика
func Dice(msg string) string {
	num, err := strconv.Atoi(msg[2:])
	if err != nil {
		return "Это вообще кубик?🤨"
	}
	if num < 1 {
		return "как я по твоему кину такой кубик? Через четвёртое пространство?🤨"
	}
	// D10 - единственный кубик, который имеет грань с номером "0"
	if num == 10 {
		return strconv.Itoa(Random(10))
	}
	return strconv.Itoa(1 + Random(num))
}

// Функция генерации случайных ответов
func Ball8(botUrl string, update Update) {
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

	SendMsg(botUrl, update, answers[Random(len(answers))])
}

// Функция, отвечающая за случайные числа
func Random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

// Функция броска монетки
func Coin(botUrl string, update Update) {
	if Random(2) == 0 {
		SendMsg(botUrl, update, "Орёл")
	} else {
		SendMsg(botUrl, update, "Решка")
	}
}

// Отправка случайного поста с Реддита (мемы, кошки, попугаи)
func SendFromReddit(botUrl string, update Update, subj string) error {
	// Отправка реквеста и обработка респонса
	url := "https://meme-api.herokuapp.com/gimme/" + subj
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Meme API error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response = new(RedditResponse)
	json.Unmarshal(body, &response)

	// Проверка на запрещёнку
	if response.Nsfw || response.Spoiler {
		response.Url = "https://belikebill.ga/billgen-API.php?default=1"
		response.Title = "Картинка оказалась со спойлером или nsfw-контентом, поэтому я заменил её на это"
	}

	// Формирование и отправка результата
	botImageMessage := SendPhoto{
		ChatId:   update.Message.Chat.ChatId,
		PhotoUrl: response.Url,
		Caption:  response.Title,
	}
	SendPict(botUrl, update, botImageMessage)
	return nil
}

// Вывод курса криптовалюты SHIB
func SendCryptoData(botUrl string, update Update) {
	// Отправка реквеста и обработка респонса
	url := "https://api2.binance.com/api/v3/ticker/24hr?symbol=SHIBBUSD"
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Binance API error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response = new(CryptoResponse)
	json.Unmarshal(body, &response)

	// Формирование и отправка результата
	if response.ChangePercent[0] == '-' {
		SendMsg(botUrl, update, "За сегодняшний день курс "+response.Symbol+" упал на "+response.ChangePercent[1:]+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, update, true)
	} else {
		SendMsg(botUrl, update, "За сегодняшний день курс "+response.Symbol+" вырос на "+response.ChangePercent+"%\n"+
			"до отметки в "+response.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, update, false)
	}
}

// Функция только для меня, проверка всех комманд
func Check(botUrl string, update Update) {
	if update.Message.Chat.ChatId == viper.GetInt("DanyaChatId") {
		start := time.Now()

		fmt.Println("Start Check() ...")
		SendCryptoData(botUrl, update)
		SendFromReddit(botUrl, update, "")
		Coin(botUrl, update)
		Help(botUrl, update)
		SendCommits(botUrl, update, "hud0shnik")
		SendMsg(botUrl, update, Dice("/d20"))
		Ball8(botUrl, update)
		SendRandomSticker(botUrl, update)
		SendFromReddit(botUrl, update, "parrots")

		for i := 1; i < 7; i++ {
			SendErrorMessage(botUrl, update, i)
		}

		SendMsg(botUrl, update, "Проверка заняла "+time.Since(start).String())
		return
	}
	SendMsg(botUrl, update, "Error 403! Beep Boop... Forbidden! Access denied 🤖")
}

// Обработчик ошибок
func SendErrorMessage(botUrl string, update Update, errorCode int) {
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

	// При возникновении ошибки, бот меня оповестит (анонимно)
	var updateDanya Update
	updateDanya.Message.Chat.ChatId = viper.GetInt("DanyaChatId")
	SendMsg(botUrl, updateDanya, "Дань, тут у одного из пользователей "+result+", надеюсь он скоро тебе о ней напишет.")

	// Вывод ошибки пользователю
	// И просьба связаться со мной для её устранения
	result += ", свяжитесь с моим создателем для устранения проблемы \n\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik"
	SendMsg(botUrl, update, result)
}

// Вывод информации о пользователе GitHub
func SendStats(botUrl string, update Update, parametrs string) {

	// Отправка запроса моему API и обработка респонса
	resp, err := http.Get("https://hud0shnikgitapi.herokuapp.com/info/" + parametrs)
	if err != nil {
		fmt.Println("GithubGoAPI error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(StatsResponse)
	json.Unmarshal(body, &user)

	// Отправка данных пользователю
	SendPict(botUrl, update, SendPhoto{
		PhotoUrl: user.Avatar,
		ChatId:   update.Message.Chat.ChatId,
		Caption: "Информация о " + user.Username + ":\n" +
			"Имя - " + user.Name + "\n" +
			"Поставленных звезд " + strconv.Itoa(user.Stars) + "⭐\n" +
			"Cсылка на аватар:\n " + user.Avatar,
	})
}

// Вывод количества коммитов пользователя GitHub
func SendCommits(botUrl string, update Update, parametrs string) {
	i, date := 0, ""

	// Поиск конца юзернейма и начала даты
	for ; i < len(parametrs); i++ {
		if parametrs[i] == ' ' {
			break
		}
	}

	// Если дата задана, записывает её
	if i != len(parametrs) {
		date = parametrs[i+1:]
	}

	// Отправка запроса моему API и обработка респонса
	resp, err := http.Get("https://hud0shnikgitapi.herokuapp.com/commits/" + parametrs[:i] + "/" + date)
	if err != nil {
		fmt.Println("GithubGoAPI error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var user = new(CommitsResponse)
	json.Unmarshal(body, &user)

	// Меняет date на "сегодня" для дальнейшего вывода
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

// Получение местоположения по IP адрессу
func CheckIPAdress(botUrl string, update Update, IP string) {
	if IP == "127.0.0.1" {
		SendMsg(botUrl, update, "Айпишник локалхоста, ага")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIYLGGzR7310Hqf8K2sljgcQF8kgOpYAAJTAAOtZbwUo9c59oswVBQjBA")
		return
	}
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

	// Отправка реквеста и обработка респонса
	SendMsg(botUrl, update, "Ищу...")
	url := "https://api.ip2country.info/ip?" + IP
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("IP2Country API error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var response = new(IP2CountryResponse)
	json.Unmarshal(body, &response)

	// Вывод результатов поиска
	if response.CountryName == "" {
		SendMsg(botUrl, update, "Не могу найти этот IP")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}
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
