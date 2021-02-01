package main

import (
	"github.com/Avimitin/go-bot/internal/bot"
)

func main() {
	cfg := bot.NewConfig()
	bot.Run(cfg)
}
