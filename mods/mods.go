package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//																							WEATHER API
type Fact struct {
	Temp       int    `json:"temp"`
	Feels_like int    `json:"feels_like"`
	Condition  string `json:"condition"`
	Wind_speed int    `json:"wind_speed"`
	Humidity   int    `json:"humidity"`
	Season     string `json:"season"`
}
type Parts struct {
	Day     Fact `json:"day"`
	Night   Fact `json:"night"`
	Morning Fact `json:"morning"`
	Evening Fact `json:"evening"`
}
type Forecasts struct {
	Week    int    `json:"week"`
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
	Parts   Parts  `json:"parts"`
}
type WeatherResponse struct {
	Now       int       `json:"now"`
	Now_dt    int       `json:"now_dt"`
	Facts     Fact      `json:"fact"`
	Forecasts Forecasts `json:"forecasts"`
}
type WeatherUpdate struct {
	LastWeatherUpdate int `json:"lastWeatherUpdate"`
}

func weatherConditionTranslate(eng string) string {
	switch eng {
	case "clear":
		return "ясно"
	case "partly-cloudy":
		return "Малооблачно"
	case "cloudy":
		return "Облачно с прояснениями"
	case "overcast":
		return "Пасмурно"
	case "drizzle":
		return "Моросит"
	case "light-rain":
		return "Небольшой дождь"
	case "rain":
		return "Дождь"
	case "moderate-rain":
		return "Умеренно сильный дождь"
	case "heavy-rain":
		return "Сильный дождь"
	case "continuous-heavy-rain":
		return "Длительный сильный дождь"
	case "showers":
		return "Ливень"
	case "wet-snow":
		return "Дождь со снегом"
	case "light-snow":
		return "Небольшой снег"
	case "snow":
		return "Снег"
	case "snow-showers":
		return "Снегопад"
	case "hail":
		return "Град"
	case "thunderstorm":
		return "Гроза"
	case "thunderstorm-with-rain":
		return "Дождь с грозой"
	case "thunderstorm-with-hail":
		return "Гроза с градом"
	default:
		return "Темпоральный дождь"
	}
}

func GetWeather() string {
	curTime := time.Now().Unix()
	fileU, err := os.Open("weather/lastWeatherUpdate.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer fileU.Close()

	var last = new(WeatherUpdate)
	bodyU, _ := ioutil.ReadAll(fileU)
	json.Unmarshal(bodyU, &last)
	var timeSinceLastUpdate int64 = (curTime - int64(last.LastWeatherUpdate))
	fmt.Println("seconds since last update: " + strconv.Itoa(int(curTime-int64(last.LastWeatherUpdate))))
	if timeSinceLastUpdate > 3600 { //3600 seconds = 1 hour
		fileU, err := os.Create("weather/lastWeatherUpdate.json")
		if err != nil {
			fmt.Println("Unable to create file:", err)
			os.Exit(1)
		}
		defer fileU.Close()
		fileU.Write([]byte("{\"lastWeatherUpdate\": " + strconv.Itoa(int(curTime)) + "}"))
		UpdateWeatherJson()
	}
	res, err := os.Open("weather/weather.json")
	if err != nil {
		return "error1"
	}
	defer res.Close()
	weatherContent, _ := ioutil.ReadAll(res)
	var rs = new(WeatherResponse)
	json.Unmarshal(weatherContent, &rs)

	condition := weatherConditionTranslate(rs.Facts.Condition)

	result := "Погода на Ольховой сейчас:\n \n" + condition +
		"\nТемпература: " + strconv.Itoa(rs.Facts.Temp) + "°" +
		"\nОщущается как: " + strconv.Itoa(rs.Facts.Feels_like) + "°" +
		"\nВетер: " + strconv.Itoa(rs.Facts.Wind_speed) + " м/с" +
		"\nВлажность воздуха: " + strconv.Itoa(rs.Facts.Humidity) + " %"
	fmt.Println(rs.Forecasts)
	return result
}

func UpdateWeatherJson() {
	fmt.Println("update weather")
	file, err := os.Create("weather/weather.json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	url := "https://api.weather.yandex.ru/v2/informers?lat=55.5692101&lon=37.4588852&lang=ru_RU"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-Yandex-API-Key", "77777777777777777777")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("weather API error")
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(WeatherResponse)
	json.Unmarshal(body, &rs)
	file.WriteString(string(body))
	fmt.Println("weather.json Updated!")
}

//																						------------------------------
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
