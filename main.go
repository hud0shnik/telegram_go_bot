package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tgBot/mods"

	"github.com/spf13/viper"
)

func main() {

	// Инициализация конфига (токенов)
	err := initConfig()
	if err != nil {
		log.Fatalf("initConfig error: %s", err)
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
			log.Fatalf("getUpdates error: %s", err)
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
func respond(botUrl string, update mods.Update) {

	// Обработчик команд
	if update.Message.Text != "" {

		request := append(strings.Split(update.Message.Text, " "), "", "")

		// Вывод реквеста для тестов
		// fmt.Println("request: \t", request)

		switch request[0] {
		case "/osu":
			mods.SendOsuInfo(botUrl, update.Message.Chat.ChatId, request[1])
		case "/commits":
			mods.SendCommits(botUrl, update.Message.Chat.ChatId, request[1], request[2])
		case "/github":
			mods.SendGithubInfo(botUrl, update.Message.Chat.ChatId, request[1])
		case "/crypto":
			mods.SendCryptoInfo(botUrl, update.Message.Chat.ChatId)
		case "/ip":
			mods.SendIPInfo(botUrl, update.Message.Chat.ChatId, request[1])
		case "/coin":
			mods.FlipCoin(botUrl, update.Message.Chat.ChatId)
		case "/start", "/help":
			mods.Help(botUrl, update.Message.Chat.ChatId)
		case "/d":
			mods.SendMsg(botUrl, update.Message.Chat.ChatId, mods.Dice(request[1]))
		case "OwO":
			mods.SendMsg(botUrl, update.Message.Chat.ChatId, "UwU")
		case "Молодец", "молодец":
			mods.SendMsg(botUrl, update.Message.Chat.ChatId, "Стараюсь UwU")
		case "Живой?", "живой?":
			mods.SendMsg(botUrl, update.Message.Chat.ChatId, "Живой")
			mods.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIdGWKu5rpWxb4gn4dmYi_rRJ9OHM9xAAJ-FgACsS8ISQjT6d1ChY7VJAQ")
		case "/check":
			mods.Check(botUrl, update.Message.Chat.ChatId)
		default:
			// Обработчик вопросов
			if update.Message.Text[len(update.Message.Text)-1] == '?' {
				mods.Ball8(botUrl, update.Message.Chat.ChatId)
			} else {
				// Дефолтный ответ
				mods.SendMsg(botUrl, update.Message.Chat.ChatId, "OwO")
			}
		}

	} else {

		// Проверка на стикер
		if update.Message.Sticker.File_id != "" {
			mods.SendRandomSticker(botUrl, update.Message.Chat.ChatId)
		} else {
			// Если пользователь отправил не сообщение и не стикер:
			mods.SendMsg(botUrl, update.Message.Chat.ChatId, "Пока я воспринимаю только текст и стикеры")
			mods.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIaImHkPqF8-PQVOwh_Kv1qQxIFpPyfAAJXAAOtZbwUZ0fPMqXZ_GcjBA")
		}

	}
}

// Функция инициализации конфига (всех токенов)
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()

}
