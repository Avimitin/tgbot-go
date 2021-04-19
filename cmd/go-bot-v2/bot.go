package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	b *tb.Bot
)

func main() {
	var err error
	b, err = tb.NewBot(tb.Settings{
		Token: "",
		Poller: &tb.LongPoller{
			Timeout: 10 * time.Second,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	var bc botCommands

	for cmd, fn := range bc {
		b.Handle(cmd, fn)
	}

	b.Start()
}
