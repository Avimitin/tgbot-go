package main

import (
	"time"

	"github.com/Avimitin/go-bot/modules/config"
	"github.com/Avimitin/go-bot/modules/database"
	"github.com/Avimitin/go-bot/modules/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	b      *tb.Bot
	botLog *zerolog.Logger
	DB     database.DataController
)

func middleware(u *tb.Update) bool {
	if u.Message == nil {
		return false
	}

	user, err := DB.GetUser(u.Message.Sender.ID)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	// insert user only when last query has no error
	if user == nil && err == nil {
		user, err = DB.NewUser(u.Message.Sender.ID, PermNormal)
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
		Dict("FROM", zerolog.Dict().
			Str("NAME", u.Message.Sender.FirstName).
			Int("ID", u.Message.Sender.ID),
		).
		Dict("CHAT", zerolog.Dict().
			Str("NAME", u.Message.Chat.FirstName).
			Int64("ID", u.Message.Chat.ID),
		).
		Int("MSGID", u.Message.ID).
		Str("CONTENT", content).
		Send()

	botLog.Trace().Interface("ORIG_MSG", u.Message).Send()

	var perm = PermNormal

	if user != nil {
		perm = user.PermID
	}

	if perm == PermBan {
		return false
	}

	if payload := getRegis(u.Message.Chat.ID, u.Message.Sender.ID); payload != nil {
		log.Trace().Msgf("%d steping next funcion", u.Message.Sender.ID)

		err := payload.fn(u.Message, payload.data)
		if err != nil {
			log.Error().Err(err).Msg("error occur when handle next func")
			send(u.Message.Chat, err.Error())
		}

		log.Trace().Msgf("remaining next step: %+v", cmdCtx)
		// abandon this message as it is been handled gracefully
		return false
	}

	if cmd, ok := msgCommand(u.Message); ok {
		log.Trace().Str("command", cmd).Int32("permission", perm)

		if authPerm(perm, cmd) {
			return true
		}

		send(u.Message.Chat, "permission denied")
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

func initDB(dsn string, logLevel string) {
	var err error
	DB, err = database.NewBotDB(dsn, logLevel)
	if err != nil {
		botLog.Fatal().
			Err(err).
			Msg("can not connect database")
	}

	botLog.Info().Msg("Establish connection to database successfully")
}

func initOwner(id int) {
	u, err := DB.GetUser(id)

	if err != nil {
		botLog.Fatal().Err(err).Msg("initialize owner")
	}

	if u != nil {
		botLog.Trace().Interface("user_detail", u).Msg("owner exist")
		return
	}

	// if owner not exist
	_, err = DB.NewUser(id, PermOwner)
	if err != nil {
		botLog.Fatal().Err(err).Msgf("failed to grant user %d to owner", id)
	}
}

func main() {
	botLog = logger.NewZeroLogger(config.GetBotLogLevel())

	initBot(config.GetBotToken())

	initDB(config.GetDatabaseDSN(), config.GetDatabaseLogLevel())

	initOwner(config.GetOwner())

	for cmd, fn := range bc {
		b.Handle(cmd, fn)
	}

	b.Start()
}
