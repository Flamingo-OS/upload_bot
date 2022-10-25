package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type BotConfig struct {
	BotToken                string `json:"bot_token"`
	MongoDbConnectionString string `json:"connection_string"`
}

func NewBotConfig() *BotConfig {
	return &BotConfig{}
}

func (ac *BotConfig) ReadConfig(fileName string) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print(err)
	}
	err = json.Unmarshal(b, &ac)
	if err != nil {
		log.Fatal(err)
	}
}
