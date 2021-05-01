package main

import (
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

	// insert user only when last query has no error
	if user == nil && err == nil {
		user, err = database.DB.NewUser(u.Message.Sender.ID, PermNormal)
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

		delContext(u.Message.Chat.ID, u.Message.Sender.ID)

		log.Trace().Msgf("remaining next step: %+v", cmdCtx)
		// abandon this message as it is been handled gracefully
		return false
	}

	if cmd, ok := msgCommand(u.Message); ok {
		log.Trace().Str("command", cmd).Int32("permission", perm).Msg("user command request details")

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
	database.DB, err = database.NewBotDB(dsn, logLevel)
	if err != nil {
		botLog.Fatal().
			Err(err).
			Msg("can not connect database")
	}

	botLog.Info().Msg("Establish connection to database successfully")
}

func initOwner(id int) {
	if u, err := database.DB.GetUser(id); err != nil && u != nil {
		log.Trace().Interface("user_info", u).Err(err).Msg("initialize owner")
		return
	}

	_, err := database.DB.SetUser(id, PermOwner)
	if err != nil {
		botLog.Fatal().Err(err).Msgf("failed to grant user %d to owner", id)
	}
}

func main() {
	cfg := ReadConfig()

	botLog = logger.NewZeroLogger(cfg.Bot.LogLevel)

	initBot(cfg.Bot.Token)

	initDB(cfg.Database.EncodeMySQLDSN(), cfg.Database.LogLevel)

	initOwner(cfg.Bot.Owner)

	for cmd, fn := range bc {
		b.Handle(cmd, fn)
	}

	b.Start()
}
