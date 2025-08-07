package service

import (
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
)

type BotService struct {
	telegram.Sender
	api                    *tgbotapi.BotAPI
	updates                tgbotapi.UpdatesChannel
	debug                  bool
	userStates             map[int64]UserState
	randomStickersFilePath string
	adminChatId            int64
}

// Структура для хранения состояния пользователя
type UserState struct {
	WaitingForOsuNickname    bool
	WaitingForGithubNickname bool
	WaitingForCommitNickname bool
	WaitingForIP             bool
	WaitingForDice           bool
	WaitingForCurrency       bool
}

func NewBotService(token string, debug bool, randomStickersFilePath string, adminChatId int64) *BotService {

	// Инициализация бота
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error("NewBotAPI error", "error", err)
		return nil
	}
	slog.Info("Authorized", "username", api.Self.UserName)

	api.Debug = debug
	sender := telegram.NewTelegram(api)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	service := &BotService{
		api:                    api,
		updates:                api.GetUpdatesChan(u),
		debug:                  debug,
		userStates:             make(map[int64]UserState),
		randomStickersFilePath: randomStickersFilePath,
		adminChatId:            adminChatId,
		Sender:                 sender,
	}

	return service
}

func (s *BotService) Run() {
	for update := range s.updates {
		if update.Message.Text != "" {
			s.handleMessage(update)
		} else if update.Message.Sticker != nil {
			s.SendRandomSticker(update.Message.Chat.ID, s.randomStickersFilePath)
		} else {
			s.SendMessage(update.Message.Chat.ID, "Пока я воспринимаю только текст и стикеры")
			s.SendSticker(update.Message.Chat.ID, "CAACAgIAAxkBAAIaImHkPqF8-PQVOwh_Kv1qQxIFpPyfAAJXAAOtZbwUZ0fPMqXZ_GcjBA")
		}
	}
}
