package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type SendMethod func(bot *tgbotapi.BotAPI, message *tgbotapi.Message)

var COMMAND = map[string]SendMethod{
	"start": start,
}

func start(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Here is start.")
	bot.Send(msg)
}
