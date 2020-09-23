package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func main() {
	fmt.Println("Bot initializing")

	bot := NewBot()

	bot.Debug = true

	log.Printf("Authorized on accout %s", bot.Self.UserName)

	updateMsg := tgbotapi.NewUpdate(0)
	updateMsg.Timeout = 20

	updates, err := bot.GetUpdatesChan(updateMsg)

	if err != nil {
		log.Printf("Some error occur when getting update.\nDescriptions: %v", err)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			COMMAND[update.Message.Command()](bot, update.Message)
		}
	}
}
