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

type RestResponse struct {
	Result []Update `json:"result"`
}

type SendMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type SendSticker struct {
	ChatId  int    `json:"chat_id"`
	Sticker string `json:"sticker"`
}
type SendPhoto struct {
	ChatId  int    `json:"chat_id"`
	Photo   string `json:"photo"`
	Caption string `json:"caption"`
}

type NasaResponse struct {
	Explanation string `json:"explanation"`
	Url         string `json:"url"`
}

type RedditResponse struct {
	Title   string `json:"title"`
	Url     string `json:"url"`
	Nsfw    bool   `json:"nsfw"`
	Spoiler bool   `json:"spoiler"`
}

type CryptoData struct {
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
	var rs = new(CryptoData)
	json.Unmarshal(body, &rs)

	fmt.Println(rs)
	result := "За последние 24 часа курс " + rs.Symbol + " изменился на " + rs.ChangePercent + "%\n" +
		"до отметки в " + rs.LastPrice + "$\n\n"

	return result
}

func Help() string {
	return "Привет👋🏻, вот список команд:\n\n/weather - показать погоду на Ольховой\n\n/nasa - картинка дня от Nasa" +
		"\n\n/d20 - кинуть д20, вместо 20 можно поставить любое число\n\n/coin - подбросить монетку" +
		"\n\n/meme - мем с реддита (смотреть на свой страх и риск, я за этот контент не отвечаю 😅)\n\n/cat и /parrot - картинка кота или попугая " +
		"\n\n/crypto - узнать текущий курс двух криптовалют (SHIB и BTC)" +
		"\n\nТакже можешь позадовать вопросы, я на них отвечу 🙃"
}

func Dice(runeMsg []rune) string {
	var num int
	if runeMsg[0] == '/' {
		num = MyAtoi(string(runeMsg[2:]))
	} else {
		num = MyAtoi(string(runeMsg[1:]))
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
	rand.Seed(time.Now().Unix())
	answers := []string{
		"Да, конечно",
		"100%",
		"Да",
		"100000000%",
		"Несомненно",
		//
		"Мб",
		"50/50",
		"Скорее да, чем нет",
		"Скорее нет, чем да",
		//
		"Нет, пфф",
		"Да нееееееееееет",
		"Точно нет",
		"0%",
		"Нет",
	}

	return answers[rand.Intn(len(answers))]
}

func MyAtoi(s string) int {
	max := int64(2 << 30)
	signFlag := true
	spaceFlag := true
	sign := 1
	digits := []int{}

	for _, char := range s {
		if char == ' ' && spaceFlag {
			continue
		}
		if signFlag {
			if char == '+' {
				signFlag = false
				spaceFlag = false
				continue
			} else if char == '-' {
				sign = -1
				signFlag = false
				spaceFlag = false
				continue
			}
		}
		if char < '0' || char > '9' {
			break
		}
		spaceFlag, signFlag = false, false
		digits = append(digits, int(char-48))
	}

	var result, place int64
	place, result = 1, 0
	last := -1

	for i, j := range digits {
		if j == 0 {
			last = i
		} else {
			break
		}
	}
	if last > -1 {
		digits = digits[last+1:]
	}

	var rtnMax int64
	if sign > 0 {
		rtnMax = max - 1
	} else {
		rtnMax = max
	}

	digitsLen := len(digits)
	for i := digitsLen - 1; i >= 0; i-- {
		result += int64(digits[i]) * place
		place *= 10
		if digitsLen-i > 10 || result > rtnMax {
			return int(int64(sign) * rtnMax)
		}
	}
	return int(result * int64(sign))
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

func GetFromReddit(chatId int, subj string) SendPhoto {
	var url string
	switch subj {
	case "meme":
		url = "https://meme-api.herokuapp.com/gimme"
	case "parrot":
		url = "https://meme-api.herokuapp.com/gimme/parrots"
	case "cat":
		url = "https://meme-api.herokuapp.com/gimme/cats"
	default:
		url = "https://meme-api.herokuapp.com/gimme/space"
	}
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Meme API error: ", err)
		return SendPhoto{
			ChatId:  chatId,
			Photo:   "https://belikebill.ga/billgen-API.php?default=1",
			Caption: "Meme API error",
		}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(RedditResponse)
	json.Unmarshal(body, &rs)

	if rs.Nsfw || rs.Spoiler {
		rs.Url = "https://belikebill.ga/billgen-API.php?default=1"
		rs.Title = "Мем оказался со спойлером или nsfw-контентом, поэтому вместо него вот тебе картинка с Биллом"
	}

	botImageMessage := SendPhoto{
		ChatId:  chatId,
		Photo:   rs.Url,
		Caption: rs.Title,
	}

	return botImageMessage
}

func GetAstronomyPictureoftheDay(chatId int) SendPhoto {
	InitConfig()
	url := "https://api.nasa.gov/planetary/apod?api_key=" + viper.GetString("nasaToken")
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("Nasa API error: ", err)
		return SendPhoto{
			ChatId:  chatId,
			Photo:   "",
			Caption: "Nasa API error",
		}
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(NasaResponse)
	json.Unmarshal(body, &rs)

	botImageMessage := SendPhoto{
		ChatId:  chatId,
		Photo:   rs.Url,
		Caption: rs.Explanation,
	}
	return botImageMessage
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
