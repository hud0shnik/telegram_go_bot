package main

import (
	"log"
	"os"

	"github.com/hud0shnik/telegram_go_bot/internal/config"
	"github.com/hud0shnik/telegram_go_bot/internal/handler"
	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
	"github.com/joho/godotenv"
)

func main() {

	// Инициализация конфига (токенов)
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("initConfig error: %s", err)
		return
	}

	// Загрузка переменных окружения
	godotenv.Load()

	// Url бота для отправки и приёма сообщений
	botUrl := "https://api.telegram.org/bot" + os.Getenv("TOKEN")
	offSet := 0

	// Цикл работы бота
	for {

		// Получение апдейтов
		updates, err := telegram.GetUpdates(botUrl, offSet)
		if err != nil {
			log.Fatalf("getUpdates error: %s", err)
		}

		// Обработка апдейтов
		for _, update := range updates {
			handler.SendRespond(botUrl, update)
			offSet = update.UpdateId + 1
		}

		// Вывод в консоль для тестов
		// fmt.Println(updates)
	}
}
