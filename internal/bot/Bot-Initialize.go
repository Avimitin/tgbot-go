package bot

import (
	"database/sql"
	"github.com/Avimitin/go-bot/internal/bot/internal/KaR"
	"log"

	"github.com/Avimitin/go-bot/internal/conf"
	"github.com/Avimitin/go-bot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

//
func NewBot() *tgbotapi.BotAPI {
	token := conf.LoadBotToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	return bot
}

func NewData() *conf.BotData {
	var botData *conf.BotData = &conf.BotData{
		KAR:    make(conf.KeywordsReplyType),
		Groups: make([]int64, 5),
	}
	err := KaR.LoadKeywordReply(DB, botData)
	if err != nil {
		log.Panic(err)
	}
	return botData
}

// NewSB
func NewDB() *sql.DB {
	s := conf.LoadDBSecret()
	DB, err := database.NewDB(s)
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
	err := KaR.LoadKeywordReply(DB, data)
	if err != nil {
		log.Panic(err)
	}
}
