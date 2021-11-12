package mods

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func SendRandomShibaSticker(botUrl string, update Update) {
	stickers := []string{
		"CAACAgIAAxkBAAIM7mF7830wgmsiYJ5xHTEZjHgJ_YphAAKRAQACK15TC92mC_kqIE5PIQQ",
		"CAACAgIAAxkBAAIM8mF785AXsxybm8IbstiOBA8vc7ujAAKHAQACK15TC3gn1k2Gf2lgIQQ",
		"CAACAgIAAxkBAAIM8GF784o9uWLTWhdCbaiH3xebHlDpAAKKAQACK15TCxtDbMsAAT60RCEE",
	}

	SendStck(botUrl, update, stickers[Random(len(stickers))])
}

func SendRandomSticker(botUrl string, update Update) error {
	fileU, err := os.Open("mods/stickers.json")
	if err != nil {
		fmt.Println(err)
		SendErrorMessage(botUrl, update, 6)
		os.Exit(1)
	}
	defer fileU.Close()

	bodyU, _ := ioutil.ReadAll(fileU)
	stickers := []string{}

	json.Unmarshal(bodyU, &stickers)

	SendStck(botUrl, update, stickers[Random(len(stickers))])
	return nil
}
