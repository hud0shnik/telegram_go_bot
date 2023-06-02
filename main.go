package main

import (
	"log"

	"github.com/hud0shnik/telegram_go_bot/internal/config"
	"github.com/hud0shnik/telegram_go_bot/internal/handler"
	"github.com/hud0shnik/telegram_go_bot/internal/telegram"

	"github.com/spf13/viper"
)

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
