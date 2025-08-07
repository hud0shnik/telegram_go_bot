package service

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Константы для кнопок
const (
	OSU_BUTTON           = "🎵 Osu! 🎮"
	GITHUB_BUTTON        = "💻 Github 🐙"
	CRYPTO_BUTTON        = "📈 Курс SHIB 🐕"
	COMMITS_BUTTON       = "🔄 Коммиты 💾"
	IP_BUTTON            = "🌐 IP 📡"
	COIN_BUTTON          = "🪙 Монетка ✨"
	DICE_BUTTON          = "🎲 Кубик 🎲"
	CURRENCY_BUTTON      = "💰 Курс JPY 📊"
	NEXT_PAGE_BUTTON     = ">>>"
	PREVIOUS_PAGE_BUTTON = "<<<"
)

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

		case state.WaitingForCurrency:
			s.ConvertJpyToRub(chatID, update.Message.Text)
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
		s.previousButtonPressed(chatID)
	case CURRENCY_BUTTON:
		s.userStates[chatID] = UserState{WaitingForCurrency: true}
		s.currencyButtonPressed(chatID)
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
