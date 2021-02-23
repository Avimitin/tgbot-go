package main

import (
	"log"

	"github.com/Avimitin/go-bot/internal/bot"
)

func main() {
	cfg := bot.NewConfig()
	if err := bot.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
