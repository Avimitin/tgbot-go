package bot

import (
	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type command map[string]func(m *bapi.Message) error

// hasCommand return if cmd is executive.
func (c command) hasCommand(cmd string) (func(m *bapi.Message) error, bool) {
	fn, ok := c[cmd]
	return fn, ok
}

func sendT(text string, chatID int64) (bapi.Message, error) {
	msg := bapi.NewMessage(chatID, text)
	return bot.Send(msg)
}

func sendP(text string, chatID int64, parseMode string) (bapi.Message, error) {
	msg := bapi.NewMessage(chatID, text)
	msg.ParseMode = parseMode
	return bot.Send(msg)
}
