package bot

import (
	"database/sql"
	"log"
	"os"

	"github.com/Avimitin/go-bot/internal/auth"
	"github.com/Avimitin/go-bot/internal/bot/internal/KaR"
	"github.com/Avimitin/go-bot/internal/conf"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	VERSION = "0.5.8"
	CREATOR = 649191333
)

var (
	DB   *sql.DB
	bot  *tgbotapi.BotAPI
	data *conf.BotData
)

func Run(CleanMode bool) {
	// Bot init
	log.Printf("Bot initializing... Version: %v", VERSION)
	bot = NewBot()
	bot.Debug = true
	log.Printf("Authorized on accout %s", bot.Self.UserName)

	// DB init
	log.Print("Fetching database connection...")
	DB = NewDB()
	defer DB.Close()

	log.Print("Initializing data")
	data = NewData()

	log.Print("Fetching authorized groups...")
	data.Groups = NewAuthGroups()

	log.Print("Fetching keywords and replies...")
	NewKeywordReplies()

	updateMsg := tgbotapi.NewUpdate(0)
	updateMsg.Timeout = 20

	updates, err := bot.GetUpdatesChan(updateMsg)

	if err != nil {
		log.Printf("Some error occur when getting update.\nDescriptions: %v", err)
	}

	// 清理模式
	for CleanMode {
		log.Printf("Cleaning MSG...")
		updates.Clear()
		os.Exit(0)
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Chat.Type == "supergroup" && !auth.CFGIsAuthGroups(data, update.Message.Chat.ID) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你们这啥群啊，别乱拉人，爬爬爬！")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("[ERROR] %s", err)
			}

			_, err = bot.LeaveChat(update.Message.Chat.ChatConfig())
			if err != nil {
				log.Printf("[ERROR] Error happen when leave chat %s", err)
			}
		}
		msgHandler(update.Message)
	}
}

func commandHandler(message *tgbotapi.Message, ctx *conf.Context) {
	cmd, hasElem := COMMAND[message.Command()]
	if hasElem {
		_, err := cmd(bot, message)
		if err != nil {
			ctx.AppendError(err.Error())
			ctx.Done <- false
			return
		}
		ctx.Done <- true
	}
}

func msgHandler(message *tgbotapi.Message, cfg *conf.Config) {
	if message.IsCommand() {
		go commandHandler(message, cfg.Context())
	} else {
		go keywordHandler(message, cfg.Context())
	}

	select {
	case done := <-cfg.Context().Done:
		cfg.SetOcpyThread(0)
		if !done {
			SendTextMsg(bot, message.Chat.ID, cfg.Context().LatestError())
		}
	}
}

func keywordHandler(message *tgbotapi.Message, ctx *conf.Context) {
	reply, e := KaR.RegexKAR(message.Text, data)
	if e {
		_, err := SendTextMsg(bot, message.Chat.ID, reply)
		if err != nil {
			ctx.AppendError(err.Error())
			ctx.Done <- false
		}
		ctx.Done <- true
	}
}
