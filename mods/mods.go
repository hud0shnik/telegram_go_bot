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

func Coin() string {
	if Random(2) == 0 {
		return "Орёл"
	}
	return "Решка"
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

	for i := len(rs.LastPrice) - 1; ; i-- {
		if rs.LastPrice[i] == '0' {
			rs.LastPrice = rs.LastPrice[:i]
		} else {
			break
		}
	}

	SendMsg(botUrl, update, "За сегодняшний день курс "+rs.Symbol+" изменился на "+rs.ChangePercent+"%\n"+
		"до отметки в "+rs.LastPrice+"$\n\n")

	SendRandomShibaSticker(botUrl, update)
}

func GetTime(botUrl string, update Update, DanyaFlag bool) {
	currentTime := time.Now().Add(3 * time.Hour)
	if currentTime.Format("01-02") == "11-08" {
		SendMsg(botUrl, update, "Сегодня день рождения самого умного человека во всей Москве - Дани!!!")
		if DanyaFlag {
			SendMsg(botUrl, update, "🎂 C др, создатель!!! 🥳 🎉")
		}
		SendStck(botUrl, update, "CAACAgIAAxkBAAINzWGH6G2PfGPH2eRiI-x9fHQ_McDSAAJZAAOtZbwU9LtHF4nhLQkiBA")
	} else {
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
		SendMsg(botUrl, update, Coin())
		Help(botUrl, update)
		CheckGit(botUrl, update)
		Sun(botUrl, update)
		GetTime(botUrl, update, DanyaFlag)
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
	SendMsg(botUrl, update, "Error 403! Beep Bop... Forbidden! Access denied 🤖")
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
	currentDate := string(time.Now().Add(3 * time.Hour).Format("2006-01-02"))

	//Вот так выглядит html одной ячейки:
	//<rect width="11" height="11" x="-36" y="75" class="ContributionCalendar-day" rx="2" ry="2" data-count="1" data-date="2021-12-03" data-level="1"></rect>

	if strings.Contains(string(body), "data-date=\""+currentDate+"\" data-level=\"") {
		pageStr, commits := string(body), ""
		i := 0

		for ; i < len(pageStr)-40; i++ {
			if pageStr[i:i+35] == "data-date=\""+currentDate+"\" data-level=\"" {
				i -= 7
				break
			}
		}
		for ; pageStr[i] != '"'; i++ {
			//доводит i до символа "
		}
		for i++; pageStr[i] != '"'; i++ {
			commits += string(pageStr[i])
		}
		SendMsg(botUrl, update, "Коммитов за сегодня: "+commits)
		return
	}
	SendMsg(botUrl, update, "Коммитов за сегодня: 0")
}
