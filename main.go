package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"tgBot/mods"
	"time"

	"github.com/spf13/viper"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Println("Config error: ", err)
		return
	}
	botToken := viper.GetString("token")
	//https://api.telegram.org/bot<token>/METHOD_NAME
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offSet := 0
	for {
		updates, err := getUpdates(botUrl, offSet)
		if err != nil {
			log.Println("Something went wrong: ", err)
		}
		for _, update := range updates {
			respond(botUrl, update)
			offSet = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]mods.Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse mods.RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}
func updateWeatherJson() {
	fmt.Println("update weather")
	file, err := os.Create("weather/weather.json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	url := "https://api.weather.yandex.ru/v2/forecast/?lat=55.5692101&lon=37.4588852&lang=ru_RU"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Yandex-API-Key", "keyyyyyyyyy")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("weather API error")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(mods.WeatherResponse)
	json.Unmarshal(body, &rs)
	file.WriteString(string(body))
	fmt.Println("WeatherJson Updated!")
}

func getWeather() string {
	curTime := time.Now().Unix()
	fileU, err := os.Open("weather/lastWeatherUpdate.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fileU.Close()

	var last = new(mods.WeatherUpdate)
	bodyU, _ := ioutil.ReadAll(fileU)
	json.Unmarshal(bodyU, &last)

	var timeSinceLastUpdate int64 = (curTime - int64(last.LastWeatherUpdate))
	fmt.Println("seconds since last update: " + strconv.Itoa(int(curTime-int64(last.LastWeatherUpdate))))
	if timeSinceLastUpdate > 3600 { //3600 seconds = 1 hour
		fileU, err := os.Create("weather/lastWeatherUpdate.json")
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer fileU.Close()
		fileU.Write([]byte("{\"lastWeatherUpdate\": " + strconv.Itoa(int(curTime)) + "}"))
		updateWeatherJson()
	}
	res, err := os.Open("weather/weather.json")
	if err != nil {
		return "error1"
	}
	defer res.Close()
	weatherContent, _ := ioutil.ReadAll(res)
	var rs = new(mods.WeatherResponse)
	json.Unmarshal(weatherContent, &rs)
	return "Погода на Ольховой сейчас\n \n" + "Температура: " + strconv.Itoa(rs.Facts.Temp) + "°" + "\nОщущается как: " + strconv.Itoa(rs.Facts.Feels_like) + "°"
}

func logic(msg string) string {
	msg = strings.ToLower(msg)
	runeMsg := []rune(msg)
	lenMsg := len(msg)
	if lenMsg > 0 && ((runeMsg[0] == 'п') || msg[0] == 'w') {
		return getWeather()
	}
	if msg == "help" {
		return "погода или weather - показать погоду на Ольховой\nd20 - кинуть д20 (рандомное число от 1 до 20), вместо 20 можно поставить любое число\nМожешь позадовать вопросы, я на них отвечу\ncoin - подброшу монетку (0-орел, 1-решка)"
	}
	if lenMsg > 4 && (msg[:4] == "math") {
		if (lenMsg < 17 && lenMsg > 10) && msg[5:10] == "roman" {
			return mods.IntToRoman(mods.MyAtoi(msg[10:]))
		} // math roman9 -> IX
		return "input: " + strconv.Itoa(mods.MyAtoi(msg[4:]))
	}
	if lenMsg > 1 && ((runeMsg[0] == 'д') || msg[0] == 'd') {
		num := mods.MyAtoi(string(runeMsg[1:]))
		if num <= 0 {
			return "как я по твоему кину такой кубик? Через четвёртое пространство?"
		}
		if num == 10 {
			return strconv.Itoa(mods.Coin(10))
		}
		return strconv.Itoa(1 + mods.Coin(num))
	}
	if lenMsg > 3 && ((msg[lenMsg-1] == '?') || (msg[lenMsg-2] == '?')) {
		return mods.Ball8()
	}
	if lenMsg >= 3 && msg[:3] == "owo" {
		return "UwU"
	}
	if msg == "coin" || msg == "монетка" || msg == "монета" {
		return strconv.Itoa(mods.Coin(2))
	}
	if lenMsg >= 7 && (msg == "молодец" || msg == "спасибо") {
		return "Стараюсь UwU"
	}
	if lenMsg >= 5 && (msg == "харош" || msg == "хорош") {
		return "Стараюсь UwU"
	}
	return "OwO"
}

func respond(botUrl string, update mods.Update) error {
	var botMessage mods.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = logic(update.Message.Text)

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
