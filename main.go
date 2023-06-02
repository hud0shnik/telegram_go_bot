package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tgBot/internal/commands"
	"tgBot/internal/config"
	"tgBot/internal/send"

	"github.com/spf13/viper"
)

// Структуры для работы с Telegram API

type telegramResponse struct {
	Result []update `json:"result"`
}

type update struct {
	UpdateId int     `json:"update_id"`
	Message  message `json:"message"`
}

type message struct {
	Chat    chat    `json:"chat"`
	Text    string  `json:"text"`
	Sticker sticker `json:"sticker"`
}

type chat struct {
	ChatId int `json:"id"`
}

type sticker struct {
	File_id string `json:"file_id"`
}

func main() {

	// Инициализация конфига (токенов)
	err := config.InitConfig()
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
func getUpdates(botUrl string, offset int) ([]update, error) {

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
	var restResponse telegramResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

// Функция генерации и отправки ответа
func respond(botUrl string, update update) {

	// Обработчик команд
	if update.Message.Text != "" {

		request := append(strings.Split(update.Message.Text, " "), "", "")

		// Вывод реквеста для тестов
		// fmt.Println("request: \t", request)

		switch request[0] {
		case "/osu":
			commands.SendOsuInfo(botUrl, update.Message.Chat.ChatId, request[1])
		case "/commits":
			commands.SendCommits(botUrl, update.Message.Chat.ChatId, request[1], request[2])
		case "/github":
			commands.SendGithubInfo(botUrl, update.Message.Chat.ChatId, request[1])
		case "/crypto":
			commands.SendCryptoInfo(botUrl, update.Message.Chat.ChatId)
		case "/ip":
			commands.SendIPInfo(botUrl, update.Message.Chat.ChatId, request[1])
		case "/coin":
			commands.FlipCoin(botUrl, update.Message.Chat.ChatId)
		case "/start", "/help":
			commands.Help(botUrl, update.Message.Chat.ChatId)
		case "/d":
			commands.RollDice(botUrl, update.Message.Chat.ChatId, request[1])
		case "OwO":
			send.SendMsg(botUrl, update.Message.Chat.ChatId, "UwU")
		case "Молодец", "молодец":
			send.SendMsg(botUrl, update.Message.Chat.ChatId, "Стараюсь UwU")
		case "Живой?", "живой?":
			send.SendMsg(botUrl, update.Message.Chat.ChatId, "Живой")
			send.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIdGWKu5rpWxb4gn4dmYi_rRJ9OHM9xAAJ-FgACsS8ISQjT6d1ChY7VJAQ")
		case "/check":
			commands.Check(botUrl, update.Message.Chat.ChatId)
		default:
			// Обработчик вопросов
			if update.Message.Text[len(update.Message.Text)-1] == '?' {
				commands.Ball8(botUrl, update.Message.Chat.ChatId)
			} else {
				// Дефолтный ответ
				send.SendMsg(botUrl, update.Message.Chat.ChatId, "OwO")
			}
		}

	} else {

		// Проверка на стикер
		if update.Message.Sticker.File_id != "" {
			send.SendRandomSticker(botUrl, update.Message.Chat.ChatId)
		} else {
			// Если пользователь отправил не сообщение и не стикер:
			send.SendMsg(botUrl, update.Message.Chat.ChatId, "Пока я воспринимаю только текст и стикеры")
			send.SendStck(botUrl, update.Message.Chat.ChatId, "CAACAgIAAxkBAAIaImHkPqF8-PQVOwh_Kv1qQxIFpPyfAAJXAAOtZbwUZ0fPMqXZ_GcjBA")
		}

	}
}
