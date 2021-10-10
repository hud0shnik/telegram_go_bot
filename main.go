package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tgBot/mods"

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

func logic(msg string) string {
	msg = strings.ToLower(msg)
	lenMsg := len(msg)

	/*if msg == "w" { // requ.Header.Set("X-Yandex-API-Key", "hrh67--dhjdjeyj-dhje")
		respondWeather, err := http.Get("https://api.weather.yandex.ru/v2/forecast/?lat=55.5692101&lon=37.4588852")
		if err != nil {
			fmt.Println(err)
			return "error1"
		}
		fmt.Println(respondWeather)
		defer respondWeather.Body.Close()
		bodyW, err := ioutil.ReadAll(respondWeather.Body)
		if err != nil {
			fmt.Println(err)
			return "error2"
		}
		fmt.Println(string(bodyW))
		var wrestResponse mods.WResponse
		err = json.Unmarshal(bodyW, &wrestResponse)
		if err != nil {
			fmt.Println(err)
			return "error3"
		}
		fmt.Println(wrestResponse)
		return "workin"
	}*/
	if msg == "help" {
		return "d20 - кинуть д20 (рандомное число от 1 до 20), вместо 20 можно поставить любое число\nПойти ли сегодня в универ? - я отвечу на твой вопрос\ncoin - подброшу монетку (0-орел, 1-решка)"
	}
	if lenMsg > 4 && (msg[:4] == "math") {
		if (lenMsg < 17 && lenMsg > 10) && msg[5:10] == "roman" {
			return mods.IntToRoman(mods.MyAtoi(msg[10:]))
		} // math roman9 -> IX
		return "input: " + strconv.Itoa(mods.MyAtoi(msg[4:]))
	}
	if lenMsg > 1 && (msg[0] == 'd') {
		num := mods.MyAtoi(msg[1:])
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
	if msg == "coin" {
		return strconv.Itoa(mods.Coin(2))
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
