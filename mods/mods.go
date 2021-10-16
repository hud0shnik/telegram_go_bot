package mods

import (
	"math/rand"
	"time"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

func Ball8() string {
	rand.Seed(time.Now().Unix())
	answers := []string{
		"Да, конечно",
		"100%",
		"Да",
		"100000000%",
		"Несомненно",
		//
		"Мб",
		"50/50",
		"Скорее да, чем нет",
		"Скорее нет, чем да",
		//
		"Нет, пфф",
		"Да нееееееееееет",
		"Точно нет",
		"0%",
		"Нет",
	}

	return answers[rand.Intn(len(answers))]
}

func MyAtoi(s string) int {
	max := int64(2 << 30)
	signFlag := true
	spaceFlag := true
	sign := 1
	digits := []int{}

	for _, char := range s {
		if char == ' ' && spaceFlag {
			continue
		}
		if signFlag {
			if char == '+' {
				signFlag = false
				spaceFlag = false
				continue
			} else if char == '-' {
				sign = -1
				signFlag = false
				spaceFlag = false
				continue
			}
		}
		if char < '0' || char > '9' {
			break
		}
		spaceFlag, signFlag = false, false
		digits = append(digits, int(char-48))
	}

	var result, place int64
	place, result = 1, 0
	last := -1

	for i, j := range digits {
		if j == 0 {
			last = i
		} else {
			break
		}
	}
	if last > -1 {
		digits = digits[last+1:]
	}

	var rtnMax int64
	if sign > 0 {
		rtnMax = max - 1
	} else {
		rtnMax = max
	}

	digitsLen := len(digits)
	for i := digitsLen - 1; i >= 0; i-- {
		result += int64(digits[i]) * place
		place *= 10
		if digitsLen-i > 10 || result > rtnMax {
			return int(int64(sign) * rtnMax)
		}
	}
	return int(result * int64(sign))
}

func Coin(n int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(n)
}
