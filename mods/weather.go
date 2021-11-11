package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type WeatherAPIResponse struct {
	Current Current `json:"current"`
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

func GetWeather() string {
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
		return "weather error"
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var rs = new(WeatherAPIResponse)
	json.Unmarshal(body, &rs)
	file.WriteString(string(body))
	fmt.Println("weather.json Updated!")

	fmt.Println(rs)

	return strconv.Itoa(int(rs.Current.Feels_like))
}
