package tools

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func SendTextMsg(bot *tgbotapi.BotAPI, chatID int64, text string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	return bot.Send(msg)
}

func SendParseTextMsg (bot *tgbotapi.BotAPI, chatID int64, text string, parse string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parse
	return bot.Send(msg)
}
