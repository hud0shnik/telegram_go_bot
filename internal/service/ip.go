package service

import (
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// Структура респонса ip-api
type ipApiResponse struct {
	Status      string  `json:"status"`
	CountryName string  `json:"country"`
	Region      string  `json:"regionName"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
}

// Список зарезервированных IP
var reservedIPs = map[string]net.IPNet{
	// Приватные сети
	"Локальная сеть 10.x":  mustParseCIDR("10.0.0.0/8"),
	"Локальная сеть 172.x": mustParseCIDR("172.16.0.0/12"),
	"Локальная сеть 192.x": mustParseCIDR("192.168.0.0/16"),

	// Служебные адреса
	"Loopback":          mustParseCIDR("127.0.0.0/8"),
	"Link-local":        mustParseCIDR("169.254.0.0/16"),
	"Multicast":         mustParseCIDR("224.0.0.0/4"),
	"240.0.0.0/4":       mustParseCIDR("240.0.0.0/4"),
	"Широковещательный": mustParseCIDR("255.255.255.255/32"),

	// Тестовые сети
	"Тестовая сеть 1": mustParseCIDR("192.0.2.0/24"),
	"Тестовая сеть 2": mustParseCIDR("198.51.100.0/24"),
	"Тестовая сеть 3": mustParseCIDR("203.0.113.0/24"),

	// Специальные
	"CGNAT":         mustParseCIDR("100.64.0.0/10"),
	"Тестовая сеть": mustParseCIDR("198.18.0.0/15"),

	// IPv6
	"IPv6 Loopback":     mustParseCIDR("::1/128"),
	"IPv6 Link-local":   mustParseCIDR("fe80::/10"),
	"IPv6 Локальная":    mustParseCIDR("fc00::/7"),
	"IPv6 Multicast":    mustParseCIDR("ff00::/8"),
	"IPv6 Документация": mustParseCIDR("2001:db8::/32"),
}

func mustParseCIDR(cidr string) net.IPNet {
	_, ipNet, _ := net.ParseCIDR(cidr)
	return *ipNet
}

// Функция поиска местоположения по IP
func (s *BotService) SendIPInfo(chatId int64, IP string) {

	// Проверка корректности ввода
	ip := net.ParseIP(IP)
	if ip == nil {
		s.SendMessage(chatId, "Неправильно набран IP")
		s.SendMessage(chatId, "Вот примеры:\n\n104.16.229.229\n151.101.65.140\n2a02:6b8::")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	for name, reservedIP := range reservedIPs {
		if reservedIP.Contains(ip) {
			s.SendMessage(chatId, "Похоже это зарезервированный IP: "+name)
			s.SendSticker(chatId, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")
			return
		}
	}

	// Формирование url
	apiUrl := "http://ip-api.com/json/" + IP

	// Отправка запроса API
	resp, err := http.Get(apiUrl)
	if err != nil {
		s.SendMessage(chatId, "Внутренняя ошибка")
		slog.Error("http.Get error", "error", err, "request", apiUrl)
		return
	}
	defer resp.Body.Close()

	// Запись респонса
	body, _ := io.ReadAll(resp.Body)
	var response = new(ipApiResponse)
	err = json.Unmarshal(body, &response)
	if err != nil {
		slog.Error("json.Unmarshal err", "error", err, "request", apiUrl)
		s.SendMessage(chatId, "Внутренняя ошибка")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	// Вывод сообщения для респонса без страны
	if response.Status != "success" {
		s.SendMessage(chatId, "Не могу найти этот IP")
		s.SendSticker(chatId, "CAACAgIAAxkBAAIY4mG13Vr0CzGwyXA1eL3esZVCWYFhAAJIAAOtZbwUgHOKzxQtAAHcIwQ")
		return
	}

	var result = "Нашёл! Страна происхождения - " + response.CountryName + "\n"

	if response.Region != "" {
		result += "Регион - " + response.Region + "\n"
	}

	if response.Lat != 0 {
		result += "Местоположение: " + strconv.FormatFloat(response.Lat, 'f', -1, 64) + ", " + strconv.FormatFloat(response.Lon, 'f', -1, 64) + "\n"
	}

	if response.Timezone != "" {
		result += "Временная зона - " + response.Timezone + "\n"
	}

	if response.Org != "" || response.AS != "" {
		result += "Организация - " + strings.Join([]string{response.Org, response.AS}, " ") + "\n"
	}

	if response.ISP != "" {
		result += "Провайдер - " + response.ISP + "\n"
	}

	result += "\n<i>Мы не храним IP, которые просят проверить пользователи, весь код бота можно найти на гитхабе.</i>"

	// Вывод результатов поиска
	s.SendMessage(chatId, result)
	s.SendSticker(chatId, "CAACAgIAAxkBAAIXqmGyGtvN_JHUQVDXspAX5jP3BvU9AAI5AAOtZbwUdHz8lasybOojBA")

}
