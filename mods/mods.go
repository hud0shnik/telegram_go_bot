package mods

import (
	"math/rand"
	"time"
)

//																							WEATHER API

/*
   clear — ясно.
   partly-cloudy — малооблачно.
   cloudy — облачно с прояснениями.
   overcast — пасмурно.
   drizzle — морось.
   light-rain — небольшой дождь.
   rain — дождь.
   moderate-rain — умеренно сильный дождь.
   heavy-rain — сильный дождь.
   continuous-heavy-rain — длительный сильный дождь.
   showers — ливень.
   wet-snow — дождь со снегом.
   light-snow — небольшой снег.
   snow — снег.
   snow-showers — снегопад.
   hail — град.
   thunderstorm — гроза.
   thunderstorm-with-rain — дождь с грозой.
   thunderstorm-with-hail — гроза с градом.
*/
type Fact struct {
	Temp       int    `json:"temp"`
	Feels_like int    `json:"feels_like"`
	Condition  string `json:"condition"`
	Wind_speed int    `json:"wind_speed"`
	Humidity   int    `json:"humidity"`
	Season     string `json:"season"`
}
type Forecast struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}
type WeatherResponse struct {
	Now       int      `json:"now"`
	Now_dt    int      `json:"now_dt"`
	Facts     Fact     `json:"fact"`
	Forecasts Forecast `json:"forecast"`
}

//
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
		"Да, кончно",
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

func IntToRoman(num int) string {
	roman := []string{"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"}
	romanValues := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
	result := "" // int: 1994 - 1000 -> 994 - 900 -> 94 -90 -> 4 - 4
	// result: "" -> ""+"M" -> "M"+"CM" -> "MCM"+"XC" -> "MCMXC"+ "IV"
	isNegative := num < 0
	if isNegative {
		num *= -1
	}
	if num == 0 {
		return "0"
	}
	for i := 0; num > 0; {
		for romanValues[i] > num {
			i++
		}
		result += roman[i]
		num -= romanValues[i]
	}
	if isNegative {
		return "-" + result
	}
	return result
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
