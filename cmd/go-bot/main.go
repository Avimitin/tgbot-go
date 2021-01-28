package main

import (
	"github.com/Avimitin/go-bot/internal/bot"
)

func main() {
	cfg := new(bot.Configuration)
	cfg.BotToken = "some token"
	bot.Run(cfg)
}
