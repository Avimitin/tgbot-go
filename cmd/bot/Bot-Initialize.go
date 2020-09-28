package bot

import (
	"github.com/Avimitin/go-bot/cmd/bot/internal/auth"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil { panic(err) }
	return bot
}

func NewCFG() *auth.Config {
	config, err := auth.NewCFG()
	if err != nil { panic(err) }
	if config.LOADED != true { panic("Fail to load config")}
	return &config
}
