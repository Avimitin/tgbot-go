package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	"gopkg.in/ini.v1"

	"log"
)

func loadCfg() (string, error) {
	cfg, err := ini.Load("cfg/bot.ini")

	if err != nil {
		return "Some error occur.", err
	}

	token := cfg.Section("privacy").Key("token").String()

	return token, nil
}

func NewBot() *tgbotapi.BotAPI {
	token, err := loadCfg()

	if err != nil {
		log.Panicf("Some error occur when loading config.\n Descriptions: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panicf("Some error occur when initialize bot.\n Descriptions: %v", err)
	}

	return bot
}
