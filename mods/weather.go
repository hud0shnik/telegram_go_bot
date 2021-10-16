package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Fact struct {
	Temp       int     `json:"temp"`
	Feels_like int     `json:"feels_like"`
	Condition  string  `json:"condition"`
	Wind_speed float32 `json:"wind_speed"`
	Humidity   int     `json:"humidity"`
	Season     string  `json:"season"`
}
type ForecastsPart struct {
	Part_name  string  `json:"part_name"`
	Temp       int     `json:"temp_avg"`
	Feels_like int     `json:"feels_like"`
	Condition  string  `json:"condition"`
	Humidity   int     `json:"humidity"`
	Wind_speed float32 `json:"wind_speed"`
	Wind_dir   string  `json:"wind_dir"`
}
type Forecasts struct {
	Week    int             `json:"week"`
	Sunrise string          `json:"sunrise"`
	Sunset  string          `json:"sunset"`
	Parts   []ForecastsPart `json:"parts"`
}
type WeatherResponse struct {
	Now       int       `json:"now"`
	Now_dt    string    `json:"now_dt"`
	Facts     Fact      `json:"fact"`
	Forecasts Forecasts `json:"forecast"`
}

type WeatherUpdate struct {
	LastWeatherUpdate int `json:"lastWeatherUpdate"`
}

func weatherSeasonToEmoji(season string) string {
	switch season {
	case "summer":
		return "☀️"
	case "autumn":
		return "🍁"
	case "winter":
		return "☃️"
	case "spring":
		return "🌱"
	default:
		return "🤔"
	}
}
func weatherConditionTranslate(eng string) string {
	switch eng { //	"смайлики: ☀️ 🌤 ⛅️ 🌥 ☁️ 🌦 🌧 ⛈ 🌩 🌨"
	case "clear":
		return "Ясно☀️"
	case "partly-cloudy":
		return "Малооблачно🌤"
	case "cloudy":
		return "Облачно с прояснениями🌥"
	case "overcast":
		return "Пасмурно☁️"
	case "drizzle":
		return "Моросит☔️"
	case "light-rain":
		return "Небольшой дождь🌦"
	case "rain":
		return "Дождь🌦"
	case "moderate-rain":
		return "Умеренно сильный дождь🌧"
	case "heavy-rain":
		return "Сильный дождь🌧🌧"
	case "continuous-heavy-rain":
		return "Длительный сильный дождь🌧🌧"
	case "showers":
		return "Ливень🌧🌧🌧"
	case "wet-snow":
		return "Дождь со снегом🌧🌨"
	case "light-snow":
		return "Небольшой снег🌨"
	case "snow":
		return "Снег🌨"
	case "snow-showers":
		return "Снегопад🌨"
	case "hail":
		return "Град🌧❄️"
	case "thunderstorm":
		return "Гроза🌩"
	case "thunderstorm-with-rain":
		return "Дождь с грозой⛈"
	case "thunderstorm-with-hail":
		return "Гроза с градом⛈⛈"
	default:
		return "Темпоральный дождь"
	}
}

func GetWeather() string {
	/*
		"night 0,1,2,3,4,5",
		"morning 6,7,8,9,10,11",
		"day 12,13,14,15,16,17",
		"evening 18,19,20,21,22,23"
	*/
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

	return createWeatherString(*rs)
}

func createWeatherString(rs WeatherResponse) string {

	return weatherSeasonToEmoji(rs.Facts.Season) + "Погода на Ольховой:\n \n" +
		"Сейчас - " + weatherConditionTranslate(rs.Facts.Condition) + ", \nчерез 3 часа - " + weatherConditionTranslate(rs.Forecasts.Parts[0].Condition) + ", \nа через 9 часов - " + weatherConditionTranslate(rs.Forecasts.Parts[1].Condition) + ".\n" +
		"\n🌡Температура: " + strconv.Itoa(rs.Facts.Temp) + "°" + " -> " + strconv.Itoa(rs.Forecasts.Parts[0].Temp) + "°" + " -> " + strconv.Itoa(rs.Forecasts.Parts[1].Temp) + "°" +
		"\n🤔Ощущается как: " + strconv.Itoa(rs.Facts.Feels_like) + "°" + " -> " + strconv.Itoa(rs.Forecasts.Parts[0].Feels_like) + "°" + " -> " + strconv.Itoa(rs.Forecasts.Parts[1].Feels_like) + "°" +
		"\n💨Ветер: " + fmt.Sprintf("%v", rs.Facts.Wind_speed) + " м/с" + " -> " + fmt.Sprintf("%v", rs.Forecasts.Parts[0].Wind_speed) + " м/с" + " -> " + fmt.Sprintf("%v", rs.Forecasts.Parts[1].Wind_speed) + " м/с" +
		"\n💧Влажность воздуха: " + strconv.Itoa(rs.Facts.Humidity) + "%" + " -> " + strconv.Itoa(rs.Forecasts.Parts[0].Humidity) + "%" + " -> " + strconv.Itoa(rs.Forecasts.Parts[1].Humidity) + "%"

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

	req.Header.Add("X-Yandex-API-Key", "может хранить его где-нибудь в другом месте?")
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
