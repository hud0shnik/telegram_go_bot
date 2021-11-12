package main

import (
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
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse mods.TelegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

//	https://core.telegram.org/bots/api#using-a-local-bot-api-server
func respond(botUrl string, update mods.Update) error {
	DanyaFlag := update.Message.Chat.ChatId == viper.GetInt("DanyaChatId")

	if update.Message.Sticker.File_unique_id != "" {
		mods.SendStck(botUrl, update, mods.GenerateRandomSticker())
		return nil
	}

	if update.Message.Text == "" {
		mods.SendMsg(botUrl, update, "Пока я воспринимаю только текст или стикеры, извини 🤷🏻‍♂️")
		return nil
	} else {
		msg := strings.ToLower(update.Message.Text)

		switch msg {
		case "/weather":
			mods.SendCurrentWeather(botUrl, update)
			mods.SendDailyWeather(botUrl, update, 3)
			return nil
		case "/check":
			mods.Check(botUrl, update, DanyaFlag)
			return nil
		case "/weather7":
			mods.SendDailyWeather(botUrl, update, 7)
			return nil
		case "/crypto":
			mods.SendCryptoData(botUrl, update)
			return nil
		case "/meme":
			mods.SendFromReddit(botUrl, update, "")
			return nil
		case "/parrot":
			mods.SendFromReddit(botUrl, update, "parrots")
			return nil
		case "/cat":
			mods.SendFromReddit(botUrl, update, "cats")
			return nil
		case "молодец", "спасибо", "харош", "хорош", "неплохо":
			mods.SendMsg(botUrl, update, "Стараюсь UwU")
			return nil
		case "/coin":
			mods.SendMsg(botUrl, update, mods.Coin())
			return nil
		case "/start", "/help":
			mods.Help(botUrl, update)
			return nil
		case "/time", "какой сегодня день?", "сколько времени?":
			mods.GetTime(botUrl, update, DanyaFlag)
			return nil
		case "owo":
			mods.SendMsg(botUrl, update, "UwU")
			return nil
		}

		lenMsg := len(msg)

		if msg[:2] == "/d" {
			mods.SendMsg(botUrl, update, mods.Dice(msg))
			return nil
		}

		if lenMsg > 3 && ((msg[lenMsg-1] == '?') || (msg[lenMsg-2] == '?')) {
			mods.SendMsg(botUrl, update, mods.Ball8())
			return nil
		}

		mods.SendMsg(botUrl, update, "OwO")
		return nil
	}
}
