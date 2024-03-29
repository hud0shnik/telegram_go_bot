package commands

import (
	"strconv"

	"github.com/hud0shnik/telegram_go_bot/internal/telegram"
	"github.com/hud0shnik/telegram_go_bot/internal/utils"

	"github.com/spf13/viper"
)

// Структуры для работы API

// Функция броска монетки
func FlipCoin(botUrl string, chatId int) {
	if utils.Random(2) == 0 {
		telegram.SendMsg(botUrl, chatId, "Орёл")
		return
	}
	telegram.SendMsg(botUrl, chatId, "Решка")
}

// Функция вывода списка всех команд
func Help(botUrl string, chatId int) {
	telegram.SendMsg(botUrl, chatId, "Привет👋🏻, вот список команд:\n\n"+
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
func RollDice(botUrl string, chatId int, parameter string) {

	// Проверка параметра
	if parameter == "" {
		telegram.SendMsg(botUrl, chatId, "Пожалуйста укажи количество граней\nНапример /d <b>20</b>")
		return
	}

	// Считывание числа граней
	num, err := strconv.Atoi(parameter)
	if err != nil || num < 1 {
		telegram.SendMsg(botUrl, chatId, "Это вообще кубик?🤨")
		return
	}

	// Проверка на d10 (единственный кубик, который имеет грань со значением "0")
	if num == 10 {
		telegram.SendMsg(botUrl, chatId, strconv.Itoa(utils.Random(10)))
		return
	}

	// Бросок
	telegram.SendMsg(botUrl, chatId, strconv.Itoa(1+utils.Random(num)))

}

// Функция генерации случайных ответов
func Ball8(botUrl string, chatId int) {

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
	telegram.SendMsg(botUrl, chatId, answers[utils.Random(10)])

}

// Функция проверки всех команд
func Check(botUrl string, chatId int) {

	// Проверка на мой id
	if chatId == viper.GetInt("AdminChatId") {

		// Вызов функций для тестирования
		SendOsuInfo(botUrl, chatId, "hud0shnik")
		SendCommits(botUrl, chatId, "hud0shnik", "")
		SendGithubInfo(botUrl, chatId, "hud0shnik")
		SendCryptoInfo(botUrl, chatId)
		SendIPInfo(botUrl, chatId, "67.77.77.7")
		telegram.SendRandomSticker(botUrl, chatId)

	} else {

		// Вывод для других пользователей
		telegram.SendMsg(botUrl, chatId, "Error 403! Beep Boop... Forbidden! Access denied 🤖")

	}

}
