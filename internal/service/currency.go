package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

// erApiResponse структура для хранения ответа от Open Exchange Rates
type erApiResponse struct {
	Rates struct {
		JPY float32 `json:"JPY"`
	} `json:"rates"`
}

// cbrApiResponse структура для хранения ответа от ЦБ
type cbrApiResponse struct {
	Valute struct {
		JPY struct {
			Nominal int     `json:"Nominal"`
			Value   float32 `json:"Value"`
		} `json:"JPY"`
	} `json:"Valute"`
}

// ConvertJpyToRub конвертирует JPY в RUB
// Отправляет сообщение с результатом
func (s *BotService) ConvertJpyToRub(chatId int64, jpyString string) {

	jpy, err := strconv.ParseFloat(jpyString, 32)
	if err != nil {
		s.SendMessage(chatId, "Внутренняя ошибка")
		slog.Error("strconv.ParseFloat error", "error", err)
		return
	}

	cbrCurrency, err := getCbrCurrency()
	if err == nil {
		s.SendMessage(chatId, fmt.Sprintf("%.2f ₽", float32(jpy)/cbrCurrency))
		return
	}

	slog.Error("Ошибка со стороны API ЦБ, использую Open Exchange Rates", "error", err)
	erCurrency, err := getOerCurrency()
	if err != nil {
		s.SendMessage(chatId, "Внутренняя ошибка")
		slog.Error("getOerCurrency error", "error", err)
		return
	}

	s.SendMessage(chatId, fmt.Sprintf("%.2f ₽", erCurrency/float32(jpy)))
}

// getCbrCurrency получает курс валюты JPY у ЦБ через прокси-апи
// Возвращает курс валюты JPY к RUB
func getCbrCurrency() (float32, error) {

	// Отправка запроса
	resp, err := http.Get("https://www.cbr-xml-daily.ru/daily_json.js")
	if err != nil {
		return 0, fmt.Errorf("http.Get error: %w", err)
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var response = new(cbrApiResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("json.Unmarshal error: %w", err)
	}

	// Проверка на дефолтное значение
	if response.Valute.JPY.Value == 0 {
		return 0, fmt.Errorf("response.Valute.Value == 0")
	}

	return float32(response.Valute.JPY.Nominal) / response.Valute.JPY.Value, nil
}

// getOerCurrency получает курс валюты JPY у Open Exchange Rates
// Возвращает курс валюты JPY к RUB
func getOerCurrency() (float32, error) {

	// Отправка запроса
	resp, err := http.Get("https://open.er-api.com/v6/latest/RUB")
	if err != nil {
		return 0, fmt.Errorf("http.Get error: %w", err)
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var response = new(erApiResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return 0, fmt.Errorf("json.Unmarshal error: %w", err)
	}

	// Проверка на дефолтное значение
	if response.Rates.JPY == 0 {
		return 0, fmt.Errorf("response.Rates.JPY == 0")
	}

	return response.Rates.JPY, nil
}
