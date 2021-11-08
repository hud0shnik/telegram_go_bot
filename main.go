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
	"time"

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
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
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

//	https://core.telegram.org/bots/api#using-a-local-bot-api-server
func respond(botUrl string, update mods.Update) error {

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
	var sendPict = func(pic mods.SendPhoto) error {
		buf, err := json.Marshal(pic)
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

	if update.Message.Sticker.File_unique_id != "" {
		sendStck(mods.GenerateRandomSticker())
		return nil
	}

	if update.Message.Text == "" {
		sendMsg("Пока я воспринимаю только текст или стикеры, извини 🤷🏻‍♂️")
		return nil
	} else {
		msg := strings.ToLower(update.Message.Text)

		switch msg {
		case "/weather", "w":
			sendMsg(mods.GetWeather())
			return nil
		case "/crypto":
			sendMsg(mods.GetCryptoData("SHIBBUSD") + mods.GetCryptoData("BTCUSDT") + mods.GetCryptoData("ETHUSDT"))
			sendStck(mods.GenerateRandomShibaSticker())
			return nil
		case "/nasa":
			sendPict(mods.GetAstronomyPictureoftheDay(update.Message.Chat.ChatId))
			return nil
		case "/meme":
			sendPict(mods.GetFromReddit(update.Message.Chat.ChatId, "meme"))
			return nil
		case "/parrot":
			sendPict(mods.GetFromReddit(update.Message.Chat.ChatId, "parrot"))
			return nil
		case "/cat":
			sendPict(mods.GetFromReddit(update.Message.Chat.ChatId, "cat"))
			return nil
		case "молодец", "спасибо", "харош", "хорош", "неплохо":
			sendMsg("Стараюсь UwU")
			return nil
		case "/coin":
			sendMsg(mods.Coin())
			return nil
		case "/start", "/help":
			sendMsg(mods.Help())
			return nil
		case "/time", "какой сегодня день?", "сколько времени?":
			currentTime := time.Now().Add(3 * time.Hour)
			if currentTime.Format("01-02") == "11-08" {
				sendMsg("Сегодня день рождения самого умного человека во всей Москве - Дани!!!")
				if update.Message.Chat.ChatId == viper.GetInt("DanyaChatId") {
					sendMsg("🎂 C др, создатель!!! 🥳 🎉")
				}
				sendStck("CAACAgIAAxkBAAINzWGH6G2PfGPH2eRiI-x9fHQ_McDSAAJZAAOtZbwU9LtHF4nhLQkiBA")
			} else {
				sendMsg(currentTime.Format("15:04 2006-01-02"))
				sendStck("CAACAgIAAxkBAAIN6GGH7YzD5gGxsI3XYzLICzRnQ0vWAAKQAgACVp29CjLSqXG41NC1IgQ")
			}
			return nil
		case "owo":
			sendMsg("UwU")
			return nil
		}

		lenMsg := len(msg)
		runeMsg := []rune(msg)

		if lenMsg > 1 && (msg[0] == 'd' || msg[:2] == "/d") {
			sendMsg(mods.Dice(runeMsg))
			return nil
		}
		if lenMsg > 3 && ((msg[lenMsg-1] == '?') || (msg[lenMsg-2] == '?')) {
			sendMsg(mods.Ball8())
			return nil
		}

		sendMsg("OwO")
		return nil
	}
}
