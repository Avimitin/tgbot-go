package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"gopkg.in/ini.v1"
)

func main() {
	fmt.Println("Bot initializing")
	cfg, err := ini.Load("cfg/bot.ini")
	if err != nil {
		log.Printf("Config loading error >>> %v\n", err)
	}
	token := cfg.Section("privacy").Key("token").String()

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on accout %s", bot.Self.UserName)

	updateMsg := tgbotapi.NewUpdate(0)
	updateMsg.Timeout = 20

	updates, err := bot.GetUpdatesChan(updateMsg)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "echo start."
			case "about":
				msg.Text = "I am a bot written by Golang."
			case "help":
				msg.Text = "Now I only have\n\n" +
							"/start  -------  For starting me\n\n" +
							"/about  -------  Information about me\n\n" +
							"/help   -------  Get some help information\n\n"
			default:
				msg.Text = "I Don't know about this command."
			}

			bot.Send(msg)
		}

	}
}