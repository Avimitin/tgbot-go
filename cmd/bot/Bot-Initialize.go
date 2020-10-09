package bot

import (
	"database/sql"
	"github.com/Avimitin/go-bot/cmd/bot/internal/conf"
	"github.com/Avimitin/go-bot/cmd/bot/internal/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func NewBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return bot
}

func NewCFG() *conf.Config {
	config, err := conf.LoadCFG()
	if err != nil {
		log.Panic(err)
	}
	if config.LOADED != true {
		log.Panic("Fail to load config")
	}
	return config
}

func NewDB() *sql.DB {
	DB, err := database.NewDB(cfg)
	if err != nil {
		log.Panic(err)
	}
	return DB
}

func NewAuthGroups() []int64 {
	groups, err := database.SearchGroups(DB)
	if err != nil {
		log.Panic(err)
	}
	groupsID := make([]int64, len(groups))
	for i, group := range groups {
		groupsID[i] = group.GroupID
	}
	return groupsID
}

func NewKeywordReplies() {
	err := LoadKeywordReply(DB, cfg)
	if err != nil {
		log.Panic(err)
	}
}
