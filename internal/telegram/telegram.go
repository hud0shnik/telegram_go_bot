package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Интерфейс отправки сообщений через Telegram-api
type Sender interface {
	SendMessage(chatId int64, text string) error
	SendSticker(chatId int64, url string) error
	SendPhoto(chatId int64, url string, text string) error
	SendRandomSticker(chatId int64, stickersFilePath string) error
}

type Telegram struct {
	api *tgbotapi.BotAPI
}

func NewTelegram(api *tgbotapi.BotAPI) Sender {
	return &Telegram{api: api}
}

// SendMessage отправляет сообщение
// chatId - id чата
// text - текст сообщения (HTML разрешён)
func (s *Telegram) SendMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"
	_, err := s.api.Send(msg)
	if err != nil {
		slog.Error("SendMessage error", "error", err)
		return fmt.Errorf("SendMessage error: %w", err)
	}
	return nil
}

// SendSticker отправляет стикер
// chatId - id чата
// url - url стикера (или его уникальный идентификатор)
func (s *Telegram) SendSticker(chatId int64, url string) error {
	sticker := tgbotapi.NewSticker(chatId, tgbotapi.FileURL(url))
	_, err := s.api.Send(sticker)
	if err != nil {
		slog.Error("SendSticker error", "error", err)
		return fmt.Errorf("SendSticker error: %w", err)
	}
	return nil
}

// SendPhoto отправляет фото
// chatId - id чата
// url - url фото
// text - текст под фото (HTML разрешён)
func (s *Telegram) SendPhoto(chatId int64, url string, text string) error {
	photo := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(url))
	photo.Caption = text
	photo.ParseMode = "HTML"
	_, err := s.api.Send(photo)
	if err != nil {
		slog.Error("SendPhoto error", "error", err)
		return fmt.Errorf("SendPhoto error: %w", err)
	}
	return nil
}

// SendRandomSticker отправляет случайный стикер
// chatId - id чата
// stickersFilePath - путь до файла со стикерами
func (s *Telegram) SendRandomSticker(chatId int64, stickersFilePath string) error {

	// Открытие json файла со стикерами
	file, err := os.Open(stickersFilePath)
	if err != nil {
		slog.Error("os.Open error", "error", err)
		return fmt.Errorf("os.Open error: %w", err)
	}
	defer file.Close()

	// Запись стикеров в массив
	stickers := [359]string{}
	body, err := io.ReadAll(file)
	if err != nil {
		slog.Error("io.ReadAll error", "error", err)
		return fmt.Errorf("io.ReadAll error: %w", err)
	}
	err = json.Unmarshal(body, &stickers)
	if err != nil {
		slog.Error("json.Unmarshal error", "error", err)
		return fmt.Errorf("json.Unmarshal error: %w", err)
	}

	// Отправка случайного стикера
	s.SendSticker(chatId, stickers[rand.Intn(len(stickers))])

	return nil
}
