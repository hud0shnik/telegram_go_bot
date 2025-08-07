package service

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

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
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи количество граней")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) nextButtonPressed(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Да, тут ещё есть кнопки)")
	msg.ReplyMarkup = s.createPersistentKeyboardSecondPage()
	s.api.Send(msg)
}

func (s *BotService) previousButtonPressed(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Возвращаю")
	msg.ReplyMarkup = s.createPersistentKeyboardFirstPage()
	s.api.Send(msg)
}

func (s *BotService) currencyButtonPressed(chatID int64) {
	s.userStates[chatID] = UserState{WaitingForCurrency: true}
	msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введи сумму в JPY")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	s.api.Send(msg)
}

func (s *BotService) createPersistentKeyboardFirstPage() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(OSU_BUTTON),
			tgbotapi.NewKeyboardButton(CURRENCY_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(IP_BUTTON),
			tgbotapi.NewKeyboardButton(CRYPTO_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(COIN_BUTTON),
			tgbotapi.NewKeyboardButton(NEXT_PAGE_BUTTON),
		),
	)
}
func (s *BotService) createPersistentKeyboardSecondPage() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DICE_BUTTON),
			tgbotapi.NewKeyboardButton(COMMITS_BUTTON),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(GITHUB_BUTTON),
			tgbotapi.NewKeyboardButton(PREVIOUS_PAGE_BUTTON),
		),
	)
}
