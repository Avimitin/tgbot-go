package bot

import (
	"log"
	"os"
	"time"

	"github.com/Avimitin/go-bot/internal/conf"
	"github.com/Avimitin/go-bot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	VERSION = "0.5.8"
	CREATOR = 649191333
)

func Run(cfgPath string, CleanMode bool) {
	ctx := newCTX(cfgPath)
	bot := ctx.Bot()

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
	sendHandler(bot, ctx)

	for update := range updates {

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Chat.Type == "supergroup" && !ctx.IsCertGroup(update.Message.Chat.ID) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你们这啥群啊，别乱拉人，爬爬爬！")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("[ERR] %s", err)
			}

			_, err = bot.LeaveChat(update.Message.Chat.ChatConfig())
			if err != nil {
				log.Printf("[ERR] Error happen when leave chat %s", err)
			}
			continue
		}

		msgHandler(ctx, update.Message)
	}
}

func msgHandler(ctx *Context, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		if hasCmdLimit(msg) {
			return
		}
		doCMD(ctx, msg)
	} else {
		doRegex(ctx, msg)
	}
}

func hasCmdLimit(msg *M) bool {
	// have group ↓
	if canDoCmd, ok := cmdDoAble[msg.Chat.ID]; ok {
		// have cmd limit ↓
		if can, ok := canDoCmd[msg.Command()]; ok {
			// cmd disable ↓
			if !can {
				return true
			}
		}
	}
	return false
}

// doCMD do a command
func doCMD(ctx *Context, msg *tgbotapi.Message) {
	if fn, ok := COMMAND[msg.Command()]; ok {
		fn(msg, ctx)
	}
}

// doRegex do a regexp job
func doRegex(ctx *Context, msg *tgbotapi.Message) {
	if rpy, ok := RegexKAR(&msg.Text, ctx.KeywordReplies()); ok {
		sendText(ctx, msg.Chat.ID, rpy)
	}
}

func newCTX(path string) *Context {
	// ------bot setting------
	token := conf.LoadBotToken(path)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("[ERR]%v", err)
	}
	log.Printf("[INFO]Successfully established connection to bot: %s", bot.Self.UserName)
	//------db setting------
	dbs := conf.LoadDBSecret(path)
	db, err := database.NewDB(dbs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO]Successfully established connection to database")
	//------Cert group loading------
	groups, err := database.SearchGroups(db)
	if err != nil {
		log.Fatal(err)
	}
	groupsSet := make(map[int64]interface{})
	for _, group := range groups {
		groupsSet[group.GroupID] = struct{}{}
	}
	log.Printf("[INFO]Successfully load all certed groups")
	//------Keyword and replies loading------
	k, err := Load(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO]Successfully load all keyword and replies")

	ctx := NewContext(k, db, &groupsSet, bot, 30*time.Second)
	return ctx
}
