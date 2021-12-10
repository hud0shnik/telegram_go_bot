package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"tgBot/mods"

	"github.com/spf13/viper"
)

func main() {
	err := mods.InitConfig()
	if err != nil {
		fmt.Println("Config error: ", err)
		return
	}

	botUrl := "https://api.telegram.org/bot" + viper.GetString("token")
	offSet := 0

	for {
		updates, err := getUpdates(botUrl, offSet)
		if err != nil {
			fmt.Println("Something went wrong: ", err)
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
	msg := strings.ToLower(update.Message.Text)

	if msg != "" {
		DanyaFlag := update.Message.Chat.ChatId == viper.GetInt("DanyaChatId")

		switch msg {
		case "/check":
			mods.Check(botUrl, update, DanyaFlag)
			return nil
		case "/git":
			mods.CheckGit(botUrl, update)
			return nil
		case "/weather7":
			mods.SendDailyWeather(botUrl, update, 7)
			return nil
		case "/weather":
			mods.SendCurrentWeather(botUrl, update)
			mods.SendDailyWeather(botUrl, update, 3)
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
		case "/crypto":
			mods.SendCryptoData(botUrl, update)
			return nil
		case "/sun":
			mods.Sun(botUrl, update)
			return nil
		case "/coin":
			mods.Coin(botUrl, update)
			return nil
		case "/start", "/help":
			mods.Help(botUrl, update)
			return nil
		case "/time", "какой сегодня день?":
			mods.SendTime(botUrl, update, DanyaFlag)
			return nil
		case "/ip":
			mods.SendMsg(botUrl, update, "Чтобы узнать страну по ip, отправьте: \n\n/ip 67.77.77.7")
			return nil
		case "owo":
			mods.SendMsg(botUrl, update, "UwU")
			return nil
		case "молодец", "неплохо":
			mods.SendMsg(botUrl, update, "Стараюсь UwU")
			return nil
		}

		lenMsg := len(msg)

		if lenMsg > 2 && update.Message.Text[:2] == "/d" {
			mods.SendMsg(botUrl, update, mods.Dice(update.Message.Text))
			return nil
		}

		if lenMsg > 5 && update.Message.Text[:3] == "/ip" {
			mods.CheckIPAdress(botUrl, update, update.Message.Text[3:])
			return nil
		}

		if msg[lenMsg-1] == '?' {
			mods.Ball8(botUrl, update)
			return nil
		}

		mods.SendMsg(botUrl, update, "OwO")
		return nil

	} else {
		if update.Message.Sticker.File_id != "" {
			mods.SendRandomSticker(botUrl, update)
			return nil
		}
		mods.SendMsg(botUrl, update, "Пока я воспринимаю только текст и стикеры  🤷🏻‍♂️")
		return nil
	}
}
