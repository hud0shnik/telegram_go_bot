package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Sticker struct {
	File_id        string `json:"file_id"`
	File_unique_id string `json:"file_unique_id"`
	Emoji          string `json:"emoji"`
	Is_animated    bool   `json:"is_animated"`
	Set_name       string `json:"set_name"`
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

func GetCryptoData(ticker string) string {
	url := "https://api2.binance.com/api/v3/ticker/24hr?symbol=" + ticker
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Binance API error: ", err)
		return "error"
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(CryptoResponse)
	json.Unmarshal(body, &rs)

	for i := len(rs.LastPrice) - 1; ; {
		if rs.LastPrice[i] == '0' {
			rs.LastPrice = rs.LastPrice[:i]
		} else {
			break
		}
		i--
	}

	result := "За последние 24 часа курс " + rs.Symbol + " изменился на " + rs.ChangePercent + "%\n" +
		"до отметки в " + rs.LastPrice + "$\n\n"

	return result
}

func Help() string {
	return "Привет👋🏻, вот список команд:" +
		"\n\n/weather - показать погоду на Ольховой" +
		"\n\n/weather7 - показать погоду на 7 дней" +
		"\n\n/crypto - узнать текущий курс трёх криптовалют (SHIB, BTC и ETH)" +
		"\n\n/time - узнать какое сейчас время" +
		"\n\n/d20 - кинуть д20, вместо 20 можно поставить любое число" +
		"\n\n/coin - подбросить монетку" +
		"\n\n/meme - мем с реддита (смотреть на свой страх и риск, я за этот контент не отвечаю 😅)" +
		"\n\n/cat и /parrot - картинка кота или попугая " +
		"\n\nТакже можешь позадовать вопросы, я на них отвечу 🙃"
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

func Ball8() string {
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

	return answers[Random(len(answers))]
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
	SendMsg(botUrl, update, GetCryptoData("SHIBBUSD")+GetCryptoData("BTCUSDT")+GetCryptoData("ETHUSDT"))
	SendStck(botUrl, update, GenerateRandomShibaSticker())
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

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
