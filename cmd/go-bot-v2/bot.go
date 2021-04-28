package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	b *tb.Bot
)

func initBot() {
	var err error
	poller := tb.NewMiddlewarePoller(
		&tb.LongPoller{
			Timeout: 15 * time.Second,
		},
		func(u *tb.Update) bool {
			log.Printf("From: %d | Chat: %d | Content: %s\n",
				u.Message.Sender.ID, u.Message.Chat.ID, u.Message.Text)

			return true
		})

	b, err = tb.NewBot(tb.Settings{
		Token:  "",
		Poller: poller,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Establish connection to bot successfully")
}

func main() {
	for cmd, fn := range bc {
		b.Handle(cmd, fn)
	}

	b.Start()
}
