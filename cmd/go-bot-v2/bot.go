package main

import (
	"log"
	"time"

	"github.com/Avimitin/go-bot/modules/database"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	b *tb.Bot
)

func middleware(u *tb.Update) bool {
	if u.Message == nil {
		return false
	}

	user, err := database.DB.GetUser(u.Message.Sender.ID)
	if err != nil {
		log.Printf("[Error]get sender %d id: %v", u.Message.Sender.ID, err)
		return true
	}

	if user == nil {
		_, err := database.DB.NewUser(u.Message.Sender.ID, database.PermNormal)
		if err != nil {
			log.Printf(
				"[Error]Insert user [%q](%d) failed: %v\n",
				u.Message.Sender.FirstName, u.Message.Sender.ID, err,
			)
		}
	}

	var content = u.Message.Text
	if len(content) > 10 {
		content = content[:10] + "..."
	}

	log.Printf("From: %d | Chat: %d | Content: %s | Perm: %s\n",
		u.Message.Sender.ID, u.Message.Chat.ID, content, user.PermDesc)

	if user.PermID == database.PermBan {
		return false
	}

	return true
}

func initBot(token string) {
	var err error
	poller := tb.NewMiddlewarePoller(
		&tb.LongPoller{Timeout: 15 * time.Second},
		middleware,
	)

	b, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: poller,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Establish connection to bot successfully")
}

func initDB(dsn string) {
	var err error
	database.DB, err = database.NewBotDB(dsn)
	if err != nil {
		log.Fatalf("connect to database %q: %v", dsn, err)
	}
}

func main() {
	cfg := ReadConfig()

	initBot(cfg.Bot.Token)
	initDB(cfg.Database.EncodeMySQLDSN())

	for cmd, fn := range bc {
		b.Handle(cmd, fn)
	}

	b.Start()
}
