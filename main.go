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
	initConfig()
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
	if msg == "help" {
		return "d20 - кинуть д20 (рандомное число от 1 до 20), вместо 20 можно поставить любое число\nq пойти ли сегодня в универ? - я отвечу на твой вопрос\ncoin - подброшу монетку (1-орел, 2-решка)"
	}
	if len(msg) > 4 && (msg[:4] == "math") {
		if (len(msg) < 17 && len(msg) > 10) && msg[5:10] == "roman" {
			return mods.IntToRoman(mods.MyAtoi(msg[10:]))
		} // math roman9 -> IX
		return "input: " + strconv.Itoa(mods.MyAtoi(msg[4:]))
	}
	if len(msg) > 1 && (msg[0] == 'd') {
		if mods.MyAtoi(msg[1:]) <= 0 {
			return "как я по твоему кину такой кубик? Через четвёртое пространство?"
		}
		return mods.Coin(1 + mods.MyAtoi(msg[1:]))
	}
	if len(msg) > 3 && (msg[0] == 'q') {
		return mods.Ball8()
	}
	if len(msg) >= 3 && msg[:3] == "owo" {
		return "UwU"
	}
	if msg == "coin" {
		return mods.Coin(2)
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
