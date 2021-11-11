package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type WeatherAPIResponse struct {
	Current Current `json:"current"`
	Daily   []Day   `json:"daily"`
}

type Day struct {
	Dt         int64         `json:"dt"`
	Sunrise    int           `json:"sunrise"`
	Sunset     int           `json:"sunset"`
	Temp       Temp          `json:"temp"`
	Feels_like Temp          `json:"feels_like"`
	Wind_speed float32       `json:"wind_speed"`
	Weather    []WeatherInfo `json:"weather"`
	Humidity   int           `json:"humidity"`
}

type Temp struct {
	/*
		"night 0,1,2,3,4,5",
		"morning 6,7,8,9,10,11",
		"day 12,13,14,15,16,17",
		"evening 18,19,20,21,22,23"
	*/
	Day     float32 `json:"day"`
	Night   float32 `json:"night"`
	Evening float32 `json:"eve"`
	Morning float32 `json:"morn"`
	Min     float32 `json:"min"`
	Max     float32 `json:"max"`
}

type Current struct {
	Sunrise    int           `json:"sunrise"`
	Sunset     int           `json:"sunset"`
	Temp       float32       `json:"temp"`
	Feels_like float32       `json:"feels_like"`
	Humidity   int           `json:"humidity"`
	Wind_speed float32       `json:"wind_speed"`
	Weather    []WeatherInfo `json:"weather"`
}

type WeatherInfo struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func SendWeather(botUrl string, update Update) error {
	InitConfig()
	fmt.Println("update weather ...")
	file, err := os.Create("weather/weather.json")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	url := "https://api.openweathermap.org/data/2.5/onecall?lat=55.5692101&lon=37.4588852&lang=ru&exclude=minutely,alerts&units=metric&appid=" + viper.GetString("weatherToken")
	req, _ := http.NewRequest("GET", url, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("weather API error")
		return err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(WeatherAPIResponse)
	json.Unmarshal(body, &rs)
	file.WriteString(string(body))
	fmt.Println("weather.json Updated!")

	n := 1

	result := "Погода на " + time.Unix(rs.Daily[n].Dt, 0).Format("02/01/2006") + ":\n \n" +
		"Погода - " + rs.Daily[n].Weather[0].Description +
		"\n🌡Температура: " + strconv.Itoa(int(rs.Daily[n].Temp.Morning)) +
		"\n🤔Ощущается как: " + strconv.Itoa(int(rs.Daily[n].Feels_like.Morning)) + "°" +
		"\n💨Ветер: " + fmt.Sprintf("%v", rs.Daily[n].Wind_speed) + " м/с" +
		"\n💧Влажность воздуха: " + strconv.Itoa(rs.Daily[n].Humidity) + "%"

	SendMsg(botUrl, update, result)
	return nil
}
