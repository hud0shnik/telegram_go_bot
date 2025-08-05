package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/hud0shnik/telegram_go_bot/internal/service"
	"github.com/joho/godotenv"
)

// Путь до файла со стикерами по умолчанию
const STICKERS_DEFAULT = "./assets/stickers.json"

func main() {

	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Загрузка переменных окружения
	godotenv.Load()

	// Получение переменных окружения

	token := os.Getenv("TOKEN")
	if token == "" {
		slog.Error("TOKEN not found")
		return
	}

	debug := false
	if os.Getenv("DEBUG") == "1" {
		debug = true
	}

	stickers := os.Getenv("STICKERS")
	if stickers == "" {
		slog.Info("STICKERS not found, using default", "default", STICKERS_DEFAULT)
		stickers = STICKERS_DEFAULT
	}

	adminChatId, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	if err != nil {
		slog.Error("ADMIN_CHAT_ID not found", "error", err)
		return
	}

	// Инициализация бота
	bot := service.NewBotService(token, debug, stickers, adminChatId)

	// Запуск бота
	bot.Run()
}
