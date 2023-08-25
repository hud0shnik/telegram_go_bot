package main

import (
	"os"
	"time"

	"github.com/hud0shnik/telegram_go_bot/internal/config"
	"github.com/hud0shnik/telegram_go_bot/internal/handler"
	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	// Настройка логгера
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	// Инициализация конфига (токенов)
	err := config.InitConfig()
	if err != nil {
		logrus.Fatalf("initConfig error: %s", err)
		return
	}

	// Загрузка переменных окружения
	godotenv.Load()

	// Url бота для отправки и приёма сообщений
	botUrl := "https://api.telegram.org/bot" + os.Getenv("TOKEN")
	offSet := 0

	// Уведомление о старте
	logrus.Info("Bot is running")

	// Цикл работы бота
	for {

		// Получение апдейтов
		updates, err := telegram.GetUpdates(botUrl, offSet)
		if err != nil {
			logrus.Fatalf("getUpdates error: %s", err)
		}

		// Обработка апдейтов
		for _, update := range updates {
			handler.Respond(botUrl, update)
			offSet = update.UpdateId + 1
		}

		// Вывод в консоль для тестов
		// fmt.Println(updates)
	}
}
