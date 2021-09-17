package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"tgBot/mods"
)

func main() {
	botToken := "tokennnnnnnnnnnnnnnnnnnnn"
	//https://api.telegram.org/bot<token>/METHOD_NAME
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offSet := 0
	for {
		updates, err := getUpdates(botUrl, offSet)
		if err != nil {
			log.Println("Something went wrong: ", err)
		}
		for _, update := range updates {
			respond(botUrl, update)
			offSet = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]mods.Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse mods.RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}
func logic(msg string) string {
	if len(msg) > 4 && (msg[:4] == "math" || msg[:4] == "Math") {
		return "input: " + strconv.Itoa(mods.MyAtoi(msg[4:]))
	}
	if len(msg) > 4 && (msg[:4] == "coin" || msg[:4] == "Coin") {
		return mods.Coin(mods.MyAtoi(msg[4:]))
	}
	if len(msg) > 3 && (msg[0] == 'q' || msg[0] == 'Q') {
		return mods.Ball8()
	}
	if len(msg) >= 3 && msg[:3] == "OwO" {
		return "UwU"
	}
	if msg == "coin" || msg == "Coin" {
		return mods.Coin(2)
	}
	return "OwO"
}
func respond(botUrl string, update mods.Update) error {
	var botMessage mods.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = logic(update.Message.Text)

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}
