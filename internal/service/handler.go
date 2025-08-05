package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Константы для кнопок
const (
	OSU_BUTTON           = "Osu"
	GITHUB_BUTTON        = "Github"
	CRYPTO_BUTTON        = "Курс SHIB"
	COMMITS_BUTTON       = "Коммиты"
	IP_BUTTON            = "IP"
	COIN_BUTTON          = "Монетка"
	DICE_BUTTON          = "Кубик"
	NEXT_PAGE_BUTTON     = ">>>"
	PREVIOUS_PAGE_BUTTON = "<<<"
)

// Структура для хранения состояния пользователя
type UserState struct {
	WaitingForOsuNickname    bool
	WaitingForGithubNickname bool
	WaitingForCommitNickname bool
	WaitingForIP             bool
	WaitingForDice           bool
}

func (s *BotService) handleMessage(update tgbotapi.Update) {

	chatID := update.Message.Chat.ID

	// Проверяем состояние пользователя
	if state, exists := s.userStates[chatID]; exists {
		switch {
		case state.WaitingForOsuNickname:
			s.SendOsuInfo(chatID, update.Message.Text)
			delete(s.userStates, chatID)

		case state.WaitingForGithubNickname:
			s.SendGithubInfo(chatID, update.Message.Text)
			delete(s.userStates, chatID)

		case state.WaitingForCommitNickname:
			s.SendCommits(chatID, update.Message.Text)
			delete(s.userStates, chatID)

		case state.WaitingForIP:
			s.SendIPInfo(chatID, update.Message.Text)
			delete(s.userStates, chatID)

		case state.WaitingForDice:
			s.SendDice(chatID, update.Message.Text)
			delete(s.userStates, chatID)
		}
	}

	// Обрабатываем нажатия кнопок
	switch update.Message.Text {
	case OSU_BUTTON:
		s.userStates[chatID] = UserState{WaitingForOsuNickname: true}
		s.osuButtonPressed(chatID)
	case GITHUB_BUTTON:
		s.userStates[chatID] = UserState{WaitingForGithubNickname: true}
		s.githubButtonPressed(chatID)
	case CRYPTO_BUTTON:
		s.SendCryptoInfo(chatID)
	case COMMITS_BUTTON:
		s.userStates[chatID] = UserState{WaitingForCommitNickname: true}
		s.commitsButtonPressed(chatID)
	case IP_BUTTON:
		s.userStates[chatID] = UserState{WaitingForIP: true}
		s.ipButtonPressed(chatID)
	case COIN_BUTTON:
		s.SendCoin(chatID)
	case DICE_BUTTON:
		s.userStates[chatID] = UserState{WaitingForDice: true}
		s.diceButtonPressed(chatID)
	case NEXT_PAGE_BUTTON:
		s.nextButtonPressed(chatID)
	case PREVIOUS_PAGE_BUTTON:
		s.sendKeyboard(chatID)
	case "OwO", "UwU", "owo", "uwu":
		s.SendMessage(chatID, "UwU")
	case "Молодец", "молодец", "Молодец!", "молодец!":
		s.SendMessage(chatID, "Стараюсь UwU")
	case "Живой?", "живой?", "живой", "Живой":
		s.SendMessage(chatID, "Живой")
		s.SendSticker(chatID, "CAACAgIAAxkBAAIdGWKu5rpWxb4gn4dmYi_rRJ9OHM9xAAJ-FgACsS8ISQjT6d1ChY7VJAQ")
	case "/check":
		s.SendCheck(chatID)
	default:
		s.sendKeyboard(chatID)
	}
}

func (s *BotService) osuButtonPressed(chatID int64) {
	s.userStates[chatID] = UserState{WaitingForOsuNickname: true}
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи ник osu!")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) githubButtonPressed(chatID int64) {
	s.userStates[chatID] = UserState{WaitingForGithubNickname: true}
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи ник github")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) commitsButtonPressed(chatID int64) {
	s.userStates[chatID] = UserState{WaitingForCommitNickname: true}
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи ник github")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) ipButtonPressed(chatID int64) {
	s.userStates[chatID] = UserState{WaitingForIP: true}
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи ip")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) sendKeyboard(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Чем ещё могу помочь?")
	msg.ReplyMarkup = s.createPersistentKeyboardFirstPage()
	s.api.Send(msg)
}

func (s *BotService) diceButtonPressed(chatID int64) {
	s.userStates[chatID] = UserState{WaitingForDice: true}
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи граней")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) nextButtonPressed(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Да, тут ещё есть кнопки)")
	msg.ReplyMarkup = s.createPersistentKeyboardSecondPage()
	s.api.Send(msg)

}

func (s *BotService) createPersistentKeyboardFirstPage() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(OSU_BUTTON),
			tgbotapi.NewKeyboardButton(GITHUB_BUTTON),
			tgbotapi.NewKeyboardButton(CRYPTO_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(COMMITS_BUTTON),
			tgbotapi.NewKeyboardButton(IP_BUTTON),
			tgbotapi.NewKeyboardButton(COIN_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DICE_BUTTON),
			tgbotapi.NewKeyboardButton(NEXT_PAGE_BUTTON),
		),
	)
}
func (s *BotService) createPersistentKeyboardSecondPage() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(PREVIOUS_PAGE_BUTTON),
		),
	)
}
