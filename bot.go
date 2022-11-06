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

	// Инициализация конфига (токенов)
	err := mods.InitConfig()
	if err != nil {
		fmt.Println("Config error: ", err)
		return
	}

	// Url бота для отправки и приёма сообщений
	botUrl := "https://api.telegram.org/bot" + viper.GetString("token")
	offSet := 0

	// Цикл работы бота
	for {

		// Получение апдейтов
		updates, err := getUpdates(botUrl, offSet)
		if err != nil {
			fmt.Println("Something went wrong: ", err)
		}

		// Обработка апдейтов
		for _, update := range updates {
			respond(botUrl, update)
			offSet = update.UpdateId + 1
		}

		// Вывод в консоль для тестов
		// fmt.Println(updates)
	}
}

// Функция получения апдейтов
func getUpdates(botUrl string, offset int) ([]mods.Update, error) {

	// Rest запрос для получения апдейтов
	resp, err := http.Get(botUrl + "/getUpdates?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Запись и обработка полученных данных
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

// Функция генерации и отправки ответа
func respond(botUrl string, update mods.Update) error {

	// Обработчик команд
	if update.Message.Text != "" {

		request := append(strings.Split(update.Message.Text, " "), "", "")
		fmt.Println("request: \t", request)

		switch request[0] {
		case "/osu":
			mods.SendOsuInfo(botUrl, update, request[1])
			return nil
		case "/commits":
			mods.SendCommits(botUrl, update, request[1], request[2])
			return nil
		case "/github":
			mods.SendInfo(botUrl, update, request[1])
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
		case "/ip":
			mods.CheckIPAdress(botUrl, update, request[1])
			return nil
		case "/coin":
			mods.Coin(botUrl, update)
			return nil
		case "/start", "/help":
			mods.Help(botUrl, update)
			return nil
		case "OwO":
			mods.SendMsg(botUrl, update, "UwU")
			return nil
		case "Молодец", "молодец":
			mods.SendMsg(botUrl, update, "Стараюсь UwU")
			return nil
		case "Живой?", "живой?":
			mods.SendMsg(botUrl, update, "Живой")
			mods.SendStck(botUrl, update, "CAACAgIAAxkBAAIdGWKu5rpWxb4gn4dmYi_rRJ9OHM9xAAJ-FgACsS8ISQjT6d1ChY7VJAQ")
			return nil
		case "/check":
			mods.Check(botUrl, update)
			return nil
		case "/d":
			mods.SendMsg(botUrl, update, mods.Dice(request[1]))
			return nil
		}

		if update.Message.Text[len(update.Message.Text)-1] == '?' {
			mods.Ball8(botUrl, update)
			return nil
		}

		mods.SendMsg(botUrl, update, "OwO")
		return nil

	} else {

		// Проверка на стикер
		if update.Message.Sticker.File_id != "" {
			mods.SendRandomSticker(botUrl, update)
			return nil
		}

		// Если пользователь отправил не сообщение и не стикер:
		mods.SendMsg(botUrl, update, "Пока я воспринимаю только текст и стикеры")
		mods.SendStck(botUrl, update, "CAACAgIAAxkBAAIaImHkPqF8-PQVOwh_Kv1qQxIFpPyfAAJXAAOtZbwUZ0fPMqXZ_GcjBA")
		return nil

	}
}
