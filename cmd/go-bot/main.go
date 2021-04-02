package main

import (
	"log"

	"github.com/Avimitin/go-bot/internal/bot"
)

func main() {
	cfg, err := bot.NewJsonConfig(bot.WhereCFG(""))
	if err != nil {
		log.Fatal(err)
	}
	if err := bot.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
