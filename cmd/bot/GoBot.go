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

	updates, _:= bot.GetUpdatesChan(updateMsg)

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