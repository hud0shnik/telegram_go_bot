package handler

import (
	"strings"

	"github.com/hud0shnik/telegram_go_bot/internal/commands"
	"github.com/hud0shnik/telegram_go_bot/internal/send"
	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
)

// Функция отправки ответа
func Respond(botUrl string, update telegram.Update) {

	// Запись айди чата
	chatId := update.Message.Chat.ChatId

	// Обработчик команд
	if update.Message.Text != "" {

		request := append(strings.Split(update.Message.Text, " "), "", "")

		// Вывод реквеста для тестов
		// fmt.Println("request: \t", request)

		switch request[0] {
		case "/osu":
			commands.SendOsuInfo(botUrl, chatId, request[1])
		case "/commits":
			commands.SendCommits(botUrl, chatId, request[1], request[2])
		case "/github":
			commands.SendGithubInfo(botUrl, chatId, request[1])
		case "/crypto":
			commands.SendCryptoInfo(botUrl, chatId)
		case "/ip":
			commands.SendIPInfo(botUrl, chatId, request[1])
		case "/coin":
			commands.FlipCoin(botUrl, chatId)
		case "/start", "/help":
			commands.Help(botUrl, chatId)
		case "/d":
			commands.RollDice(botUrl, chatId, request[1])
		case "OwO":
			send.SendMsg(botUrl, chatId, "UwU")
		case "Молодец", "молодец":
			send.SendMsg(botUrl, chatId, "Стараюсь UwU")
		case "Живой?", "живой?":
			send.SendMsg(botUrl, chatId, "Живой")
			send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIdGWKu5rpWxb4gn4dmYi_rRJ9OHM9xAAJ-FgACsS8ISQjT6d1ChY7VJAQ")
		case "/check":
			commands.Check(botUrl, chatId)
		default:
			// Обработчик вопросов
			if update.Message.Text[len(update.Message.Text)-1] == '?' {
				commands.Ball8(botUrl, chatId)
			} else {
				// Дефолтный ответ
				send.SendMsg(botUrl, chatId, "OwO")
			}
		}

	} else {

		// Проверка на стикер
		if update.Message.Sticker.File_id != "" {
			send.SendRandomSticker(botUrl, chatId)
		} else {
			// Если пользователь отправил не сообщение и не стикер:
			send.SendMsg(botUrl, chatId, "Пока я воспринимаю только текст и стикеры")
			send.SendStck(botUrl, chatId, "CAACAgIAAxkBAAIaImHkPqF8-PQVOwh_Kv1qQxIFpPyfAAJXAAOtZbwUZ0fPMqXZ_GcjBA")
		}

	}
}
