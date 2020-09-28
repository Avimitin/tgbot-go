package bot

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

const (
	VERSION = 0.01
	CREATOR = 649191333
)

var (
	cfg = NewCFG()
	bot = NewBot(cfg.BotToken)
)

func Run() {
	fmt.Printf("Bot initializing... Version: %v\n", VERSION)

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
			go commandHandler(update.Message)
		}
	}
}

func commandHandler(message *tgbotapi.Message) {
	cmd, hasElem := COMMAND[message.Command()]
	if hasElem {
		_, err := cmd(bot, message)

		if err != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(
				"<b>Some error happen when sending message.</b> \n\nDescriptions: \n\n<code>%v</code>", err))
			msg.ParseMode = "HTML"
			_, _ = bot.Send(msg)
		}
	}
}