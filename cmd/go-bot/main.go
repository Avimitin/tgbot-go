package main

import (
	"github.com/Avimitin/go-bot/internal/bot"
	"github.com/Avimitin/go-bot/internal/conf"
)

func main() {
	path := conf.WhereCFG("F:/code/golang/go-bot/cfg")
	bot.Run(path, false)
}
