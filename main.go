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
	err := mods.InitConfig()
	if err != nil {
		log.Println("Config error: ", err)
		return
	}
	botUrl := "https://api.telegram.org/bot" + viper.GetString("token")
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

func respond(botUrl string, update mods.Update) error {
	mods.InitConfig()
	//	https://core.telegram.org/bots/api#using-a-local-bot-api-server

	var sendMsg = func(msg string) error {
		botMessage := mods.SendMessage{
			ChatId: update.Message.Chat.ChatId,
			Text:   msg,
		}
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
	var sendStck = func(url string) error {
		botStickerMessage := mods.SendSticker{
			ChatId:  update.Message.Chat.ChatId,
			Sticker: url,
		}
		buf, err := json.Marshal(botStickerMessage)
		if err != nil {
			return err
		}
		_, err = http.Post(botUrl+"/sendSticker", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
		return nil
	}
	if update.Message.Sticker.File_unique_id != "" {
		sendStck(mods.GenerateRandomSticker())
		return nil
	}

	if update.Message.Text == "/nasa" {
		url := "https://api.nasa.gov/planetary/apod?api_key=" + viper.GetString("nasaToken")
		req, _ := http.NewRequest("GET", url, nil)
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println("Nasa API error: ", err)
			return err
		}
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		var rs = new(mods.NasaResponse)
		json.Unmarshal(body, &rs)

		botImageMessage := mods.SendPhoto{
			ChatId:  update.Message.Chat.ChatId,
			Photo:   rs.Hdurl,
			Caption: rs.Explanation,
		}

		buf, err := json.Marshal(botImageMessage)
		if err != nil {
			fmt.Println("Marshal json error: ", err)
			return err
		}
		_, err = http.Post(botUrl+"/sendPhoto", "application/json", bytes.NewBuffer(buf))
		if err != nil {
			fmt.Println("SendPhoto method error: ", err)
			return err
		}
		return nil
	}

	if update.Message.Text == "" {
		sendMsg("Пока я воспринимаю только текст и стикеры, извини 🤷🏻‍♂️")
		return nil
	}
	sendMsg(logic(update.Message.Text))
	return nil
}

func logic(msg string) string {
	msg = strings.ToLower(msg)
	runeMsg := []rune(msg)
	lenMsg := len(msg)

	if lenMsg > 0 && ((msg == "w") || msg == "/weather") {
		return mods.GetWeather()
	}
	if msg == "help" || msg == "/help" || msg == "/start" || msg == "/start start" {
		return "Привет👋🏻, вот список команд:\n\n/weather - показать погоду на Ольховой\n\n/nasa - картинка дня от Nasa\n\n/d20 - кинуть д20, вместо 20 можно поставить любое число\n\n/coin - подброшу монетку\n\nМожешь позадовать вопросы, я на них отвечу"
	}
	if lenMsg > 1 && (msg[0] == 'd' || msg[:2] == "/d") {
		var num int
		if runeMsg[0] == '/' {
			num = mods.MyAtoi(string(runeMsg[2:]))
		} else {
			num = mods.MyAtoi(string(runeMsg[1:]))
		}
		if num <= 0 {
			return "как я по твоему кину такой кубик? Через четвёртое пространство?🤨"
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
	if msg == "coin" || msg == "/coin" {
		if mods.Coin(2) == 0 {
			return "Орёл"
		}
		return "Решка"
	}
	if lenMsg >= 7 && (msg == "молодец" || msg == "спасибо") {
		return "Стараюсь UwU"
	}
	if lenMsg >= 5 && (msg == "харош" || msg == "хорош") {
		return "Стараюсь UwU"
	}
	return "OwO"
}
