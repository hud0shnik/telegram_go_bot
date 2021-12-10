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

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Sticker struct {
	File_id string `json:"file_id"`
}

type Message struct {
	Chat    Chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker Sticker `json:"sticker"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type TelegramResponse struct {
	Result []Update `json:"result"`
}

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

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func Help(botUrl string, update Update) {
	SendMsg(botUrl, update, "Привет👋🏻, вот список команд:"+
		"\n\n/weather - показать погоду на Ольховой"+
		"\n\n/weather7 - показать погоду на 7 дней"+
		"\n\n/sun - узнать о времени восхода и заката"+
		"\n\n/ip 67.77.77.7 - узнать страну по ip"+
		"\n\n/git - количество коммитов за сегодня"+
		"\n\n/crypto - узнать текущий курс криптовалюты SHIB"+
		"\n\n/time - узнать какое сейчас время"+
		"\n\n/d20 - кинуть д20, вместо 20 можно поставить любое число"+
		"\n\n/coin - подбросить монетку"+
		"\n\n/meme - мем с реддита (смотреть на свой страх и риск, я за этот контент не отвечаю 😅)"+
		"\n\n/cat и /parrot - картинка кота или попугая "+
		"\n\nТакже можешь позадовать вопросы, я на них отвечу 🙃")
}

func Dice(msg string) string {
	num, err := strconv.Atoi(msg[2:])
	if err != nil {
		return "Это вообще кубик?🤨"
	}
	if num < 1 {
		return "как я по твоему кину такой кубик? Через четвёртое пространство?🤨"
	}
	if num == 10 {
		return strconv.Itoa(Random(10))
	}
	return strconv.Itoa(1 + Random(num))
}

func Ball8(botUrl string, update Update) {
	answers := []string{
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

func Random(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}

func Coin(botUrl string, update Update) {
	if Random(2) == 0 {
		SendMsg(botUrl, update, "Орёл")
	} else {
		SendMsg(botUrl, update, "Решка")
	}
}

func SendFromReddit(botUrl string, update Update, subj string) error {
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
	var rs = new(RedditResponse)
	json.Unmarshal(body, &rs)

	if rs.Nsfw || rs.Spoiler {
		rs.Url = "https://belikebill.ga/billgen-API.php?default=1"
		rs.Title = "Мем оказался со спойлером или nsfw-контентом, поэтому вместо него вот тебе картинка с Биллом :^)"
	}

	botImageMessage := SendPhoto{
		ChatId:  update.Message.Chat.ChatId,
		Photo:   rs.Url,
		Caption: rs.Title,
	}

	SendPict(botUrl, update, botImageMessage)
	return nil
}

func SendCryptoData(botUrl string, update Update) {
	url := "https://api2.binance.com/api/v3/ticker/24hr?symbol=SHIBBUSD"
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Binance API error: ", err)
		SendErrorMessage(botUrl, update, 1)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(CryptoResponse)
	json.Unmarshal(body, &rs)

	if rs.ChangePercent[0] == '-' {
		SendMsg(botUrl, update, "За сегодняшний день курс "+rs.Symbol+" упал на "+rs.ChangePercent[1:]+"%\n"+
			"до отметки в "+rs.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, update, true)
	} else {
		SendMsg(botUrl, update, "За сегодняшний день курс "+rs.Symbol+" вырос на "+rs.ChangePercent+"%\n"+
			"до отметки в "+rs.LastPrice+"$\n\n")
		SendRandomShibaSticker(botUrl, update, false)
	}
}

func SendTime(botUrl string, update Update, DanyaFlag bool) {
	currentTime := time.Now().Add(3 * time.Hour)

	switch currentTime.Format("01-02") {
	case "01-01":
		SendMsg(botUrl, update, "С Новым годом!!!")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIWrmGvduu6ERm7-5MIXiO-gyQ060gAA20AA8A2TxO5jCglZ0hJGyME")
		break
	case "01-07":
		SendMsg(botUrl, update, "С Рождеством!!")
		SendStck(botUrl, update, "CAACAgIAAxkBAAIWs2GveCDTQW0YQxSIGKcVVUBTQBhlAAIYAAOtZbwUjebqlxyfJ9IjBA")
		break
	case "11-08":
		SendMsg(botUrl, update, "Сегодня день рождения самого умного человека во всей Москве - Дани!!!")
		if DanyaFlag {
			SendMsg(botUrl, update, "🎂 C др, создатель!!! 🥳 🎉")
		}
		SendStck(botUrl, update, "CAACAgIAAxkBAAINzWGH6G2PfGPH2eRiI-x9fHQ_McDSAAJZAAOtZbwU9LtHF4nhLQkiBA")
		break
	default:
		SendMsg(botUrl, update, currentTime.Format("15:04 2006-01-02"))
		SendStck(botUrl, update, "CAACAgIAAxkBAAIN6GGH7YzD5gGxsI3XYzLICzRnQ0vWAAKQAgACVp29CjLSqXG41NC1IgQ")
	}

}

func Check(botUrl string, update Update, DanyaFlag bool) {
	if DanyaFlag {
		start := time.Now()

		fmt.Println("Start Check() ...")
		SendCurrentWeather(botUrl, update)
		SendDailyWeather(botUrl, update, 3)
		SendCryptoData(botUrl, update)
		SendFromReddit(botUrl, update, "")
		Coin(botUrl, update)
		Help(botUrl, update)
		CheckGit(botUrl, update)
		Sun(botUrl, update)
		SendTime(botUrl, update, DanyaFlag)
		SendMsg(botUrl, update, Dice("/d20"))
		Ball8(botUrl, update)
		SendRandomSticker(botUrl, update)
		SendFromReddit(botUrl, update, "parrots")

		for i := 1; i < 7; i++ {
			SendErrorMessage(botUrl, update, i)
		}

		fmt.Println("That's all!\tTime:", time.Since(start))
		SendMsg(botUrl, update, "Проверка заняла "+time.Since(start).String())
		return
	}
	SendMsg(botUrl, update, "Error 403! Beep Boop... Forbidden! Access denied 🤖")
}

func SendErrorMessage(botUrl string, update Update, errorCode int) {
	result := "err"
	switch errorCode {
	case 1:
		result = "Ошибка работы API"
		break
	case 2:
		result = "Ошибка работы json.Marshal()"
		break
	case 3:
		result = "Ошибка работы метода SendSticker"
		break
	case 4:
		result = "Ошибка работы метода SendPhoto"
		break
	case 5:
		result = "Ошибка работы метода SendMessage"
		break
	case 6:
		result = "Ошибка работы stickers.json"
		break
	}
	result += ", свяжитесь с моим создателем для устранения проблемы \n\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik\nhttps://vk.com/hud0shnik"
	SendMsg(botUrl, update, result)
}

func CheckGit(botUrl string, update Update) {
	resp, err := http.Get("https://github.com/hud0shnik")

	if err != nil {
		fmt.Println("Github error: ", err)
		SendErrorMessage(botUrl, update, 1)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//Добавляю 3 часа из-за того, что сервер находится в другом часовом поясе
	currentDate := string(time.Now().Add(3 * time.Hour).Format("2006-01-02"))

	//Вот так выглядит html одной ячейки:
	//<rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1"></rect>

	if strings.Contains(string(body), "data-date=\""+currentDate+"\" data-level=\"") {
		pageStr, commits := string(body), ""
		i := 0

		for ; i < len(pageStr)-40; i++ {
			if pageStr[i:i+35] == "data-date=\""+currentDate+"\" data-level=\"" {
				//так как количество коммитов стоит перед датой, переставляем i
				i -= 7
				break
			}
		}
		for ; pageStr[i] != '"'; i++ {
			//доводит i до символа "
		}
		for i++; pageStr[i] != '"'; i++ {
			//считывает значение в скобках
			commits += string(pageStr[i])
		}
		for i += 35; pageStr[i] != '"'; i++ {
		}
		dataLevel, _ := strconv.Atoi(pageStr[i+1 : i+2])
		switch dataLevel {
		case 2:
			SendMsg(botUrl, update, "Коммитов за сегодня: "+commits+", неплохо!")
			SendStck(botUrl, update, "CAACAgIAAxkBAAIXWmGyDE1aVXGUY6lcjKxx9bOn0JA1AAJOAAOtZbwUIWzOXysr2zwjBA")
		case 3:
			SendMsg(botUrl, update, "Коммитов за сегодня: "+commits+", отлично!!")
			SendStck(botUrl, update, "CAACAgIAAxkBAAIXXmGyDGGOqHzR7Bu4sxu7BOFSJ5jAAAImAAOtZbwUfEYfHGky-lQjBA")
		case 4:
			SendMsg(botUrl, update, "Коммитов за сегодня: "+commits+", прекрасно!!!")
			SendStck(botUrl, update, "CAACAgIAAxkBAAIXXGGyDFClr69PKZXJo9dlYMbyilXLAAI1AAOtZbwU9aVxXMUw5eAjBA")
		default:
			SendMsg(botUrl, update, "Коммитов за сегодня: "+commits)
			SendStck(botUrl, update, "CAACAgIAAxkBAAIYG2GzRVNm_d_mVDIOaiLXkGukArlTAAJDAAOtZbwU_-iXZG7hfLsjBA")
		}
		return
	}
	SendMsg(botUrl, update, "Коммитов за сегодня пока ещё нет")
	SendStck(botUrl, update, "CAACAgIAAxkBAAIYG2GzRVNm_d_mVDIOaiLXkGukArlTAAJDAAOtZbwU_-iXZG7hfLsjBA")
}

func CheckIPAdress(botUrl string, update Update, IP string) {
	if IP[0] == ' ' {
		IP = IP[1:]
	}
	ipArray := strings.Split(IP, ".")
	if len(ipArray) != 4 {
		SendMsg(botUrl, update, "Не могу считать этот IP")
		return
	}
	for _, ipPart := range ipArray {
		num, err := strconv.Atoi(ipPart)
		if err != nil || num < 0 || num > 255 {
			SendMsg(botUrl, update, "Неправильно набран IP")
			return
		}
		if ipPart != fmt.Sprint(num) {
			SendMsg(botUrl, update, "Неправильно набран IP")
			return
		}
	}
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
	var rs = new(IP2CountryResponse)
	json.Unmarshal(body, &rs)

	if rs.CountryName == "" {
		SendMsg(botUrl, update, "Не могу найти этот IP")
		return
	}
	SendMsg(botUrl, update, "Нашёл! Страна происхождения - "+rs.CountryName+" "+rs.CountryEmoji+
		"\n\nМы не храним IP, которые просят проверить пользователи, весь код бота можно найти на гитхабе.")
	SendStck(botUrl, update, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")
}
