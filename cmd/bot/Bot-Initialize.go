package bot

import (
	"database/sql"
	"github.com/Avimitin/go-bot/cmd/bot/internal/CFGLoader"
	"github.com/Avimitin/go-bot/cmd/bot/internal/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	return bot
}

func NewCFG() *CFGLoader.Config {
	config, err := CFGLoader.LoadCFG()
	if err != nil {
		panic(err)
	}
	if config.LOADED != true {
		panic("Fail to load config")
	}
	return &config
}

func NewDB() *sql.DB {
	DB, err := database.NewDB()
	if err != nil {
		panic(err)
	}
	return DB
}
