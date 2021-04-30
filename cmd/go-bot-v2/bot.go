package main

import (
	"fmt"
	"time"

	"github.com/Avimitin/go-bot/modules/database"
	"github.com/Avimitin/go-bot/modules/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	b      *tb.Bot
	botLog *zerolog.Logger
)

func middleware(u *tb.Update) bool {
	if u.Message == nil {
		return false
	}

	user, err := database.DB.GetUser(u.Message.Sender.ID)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	if user == nil {
		_, err := database.DB.NewUser(u.Message.Sender.ID, database.PermNormal)
		if err != nil {
			botLog.Error().
				Err(err).
				Msgf("insert user [%q](%d) failed",
					u.Message.Sender.FirstName, u.Message.Sender.ID)
		}
	}

	var content = u.Message.Text
	if len(content) > 20 {
		content = content[:20] + "..."
	}

	botLog.Info().
		Msgf("From: [%s](%d) | Chat: [%s](%d) | MSGID: %d | Content: %s",
			u.Message.Sender.FirstName, u.Message.Sender.ID,
			u.Message.Chat.FirstName, u.Message.Chat.ID,
			u.Message.ID, content)

	if user != nil && user.PermID == database.PermBan {
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
		botLog.Fatal().
			Err(err).
			Msg("can not connect bot")
	}

	botLog.Info().Msg("Establish connection to bot successfully")
}

func initDB(dsn string) {
	var err error
	database.DB, err = database.NewBotDB(dsn)
	if err != nil {
		botLog.Fatal().
			Err(fmt.Errorf("connect to database %q: %v", dsn, err)).
			Msg("can not connect database")
	}

	botLog.Info().Msg("Establish connection to database successfully")
}

func main() {
	cfg := ReadConfig()

	botLog = logger.NewZeroLogger(cfg.Bot.LogLevel)
	if botLog == nil {
		log.Fatal().Msg("log level not valid")
	}

	initBot(cfg.Bot.Token)

	initDB(cfg.Database.EncodeMySQLDSN())

	for cmd, fn := range bc {
		b.Handle(cmd, fn)
	}

	b.Start()
}
