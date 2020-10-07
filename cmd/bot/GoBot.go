package bot

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/CFGLoader"
	"github.com/Avimitin/go-bot/cmd/bot/internal/auth"
	"github.com/Avimitin/go-bot/cmd/bot/internal/tools"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
)

const (
	VERSION = "0.5.8"
	CREATOR = 649191333
)

var (
	DB  *sql.DB
	cfg *CFGLoader.Config
	bot *tgbotapi.BotAPI
)

func Run(CleanMode bool) {
	log.Printf("Fetching config...\n")
	cfg = NewCFG()
	log.Printf("Config loading successfully.\n")

	log.Printf("Bot initializing... Version: %v\n", VERSION)
	bot = NewBot(cfg.BotToken)

	log.Printf("Fetching database connection...\n")
	DB = NewDB()
	defer DB.Close()

	bot.Debug = true

	log.Printf("Authorized on accout %s", bot.Self.UserName)

	updateMsg := tgbotapi.NewUpdate(0)
	updateMsg.Timeout = 20

	updates, err := bot.GetUpdatesChan(updateMsg)

	if err != nil {
		log.Printf("Some error occur when getting update.\nDescriptions: %v", err)
	}

	// 清理模式
	for CleanMode {
		log.Printf("Cleaning MSG...")
		update := <-updates
		if update.Message == nil {
			log.Printf("Cleaning done.")
			os.Exit(0)
		}
		updates.Clear()
	}

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Chat.Type == "supergroup" && !auth.IsAuthGroups(DB, update.Message.Chat.ID) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你们这啥群啊，别乱拉人，爬爬爬！")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("[ERROR] %s", err)
			}

			_, err = bot.LeaveChat(update.Message.Chat.ChatConfig())
			if err != nil {
				log.Printf("[ERROR] %s", err)
			}
		}

		if update.Message.IsCommand() {
			go commandHandler(update.Message)
		}
	}
}

func commandHandler(message *tgbotapi.Message) {
	cmd, hasElem := COMMAND[message.Command()]
	if hasElem {
		_, err := cmd(bot, message)

		if err != nil {
			_, _ = tools.SendParseTextMsg(bot, message.Chat.ID,
				fmt.Sprintf("<b>Some error happen when sending message.</b> \n\nDescriptions: \n\n<code>%v</code>", err),
				"html")
		}
	}
}
