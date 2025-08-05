package service

import (
	"math/rand"
	"strconv"
)

// Функция броска монетки
func (s *BotService) SendCoin(chatId int64) {
	if rand.Intn(2) == 0 {
		s.SendMessage(chatId, "Орёл")
		return
	}
	s.SendMessage(chatId, "Решка")
}

// Функция вывода списка всех команд
func (s *BotService) SendHelp(chatId int64) {
	s.SendMessage(chatId, "Привет👋🏻, вот список команд:\n\n"+
		"/commits <u>username</u> <u>date</u> - коммиты пользователя за день\n\n"+
		"/github <u>username</u> - информация о пользователе GitHub\n\n"+
		"/osu <u>username</u> - информация о пользователе Osu\n\n"+
		"/ip <u>ip_address</u> - узнать страну по ip\n\n"+
		"/crypto - узнать текущий курс криптовалюты SHIB\n\n"+
		"/d <b>n</b> - кинуть <b>n</b>-гранную кость\n\n"+
		"/coin - бросить монетку\n\n"+
		"Также можешь позадавать вопросы, я на них отвечу 🙃")

}

// Функция броска n-гранного кубика
func (s *BotService) SendDice(chatId int64, parameter string) {

	// Проверка параметра
	if parameter == "" {
		s.SendMessage(chatId, "Пожалуйста укажи количество граней\nНапример /d <b>20</b>")
		return
	}

	// Считывание числа граней
	num, err := strconv.Atoi(parameter)
	if err != nil || num < 1 {
		s.SendMessage(chatId, "Это вообще кубик?🤨")
		return
	}

	// Проверка на d10 (единственный кубик, который имеет грань со значением "0")
	if num == 10 {
		s.SendMessage(chatId, strconv.Itoa(rand.Intn(10)))
		return
	}

	// Бросок
	s.SendMessage(chatId, strconv.Itoa(1+rand.Intn(num)))

}

// Функция генерации случайных ответов
func (s *BotService) SendBall8(chatId int64) {

	// Массив ответов
	answers := [10]string{
		"Да, конечно!",
		"100%",
		"Да.",
		"100000000%",
		"Точно да!",
		"Нет, пфф",
		"Нееееееееееет",
		"Точно нет!",
		"Нет, нет, нет",
		"Нет.",
	}

	// Выбор случайного ответа
	s.SendMessage(chatId, answers[rand.Intn(10)])

}

// Функция проверки всех команд
func (s *BotService) SendCheck(chatId int64) {

	// Проверка на мой id
	if chatId == s.adminChatId {

		// Вызов функций для тестирования
		s.SendOsuInfo(chatId, "hud0shnik")
		s.SendCommits(chatId, "hud0shnik")
		s.SendGithubInfo(chatId, "hud0shnik")
		s.SendCryptoInfo(chatId)
		s.SendIPInfo(chatId, "67.77.77.7")
		s.SendRandomSticker(chatId, s.randomStickersFilePath)

	} else {

		// Вывод для других пользователей
		s.SendMessage(chatId, "Error 403! Beep Boop... Forbidden! Access denied 🤖")

	}

}
