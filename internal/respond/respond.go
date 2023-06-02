package respond

import (
	"strings"
	"tgBot/internal/commands"
	"tgBot/internal/send"
	"tgBot/internal/telegram"
)

// Функция генерации и отправки ответа
func SendRespond(botUrl string, update telegram.Update) {

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
