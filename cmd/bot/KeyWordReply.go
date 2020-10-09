package bot

import (
	"database/sql"
	"github.com/Avimitin/go-bot/cmd/bot/internal/conf"
	"github.com/Avimitin/go-bot/cmd/bot/internal/database"
	"github.com/Avimitin/go-bot/cmd/bot/internal/tools"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"strings"
)

func Reply(msg *tgbotapi.Message) (tgbotapi.Message, error) {
	kws, isExist := Regexp(msg.Text, cfg)
	if isExist {
		if len(kws.Replies) == 1 {
			return tools.SendTextMsg(bot, msg.Chat.ID, kws.Replies[0])
		}
		randNum := rand.Intn(len(kws.Replies))
		return tools.SendTextMsg(bot, msg.Chat.ID, kws.Replies[randNum])
	}
	return tgbotapi.Message{}, nil
}

func Regexp(s string, c *conf.Config) (*conf.KeywordsReplyType, bool) {
	for _, kw := range c.Keywords {
		if strings.Contains(s, kw.Keyword.Word) {
			return kw, true
		}
	}
	return nil, false
}

func SetKeyword(keyword string, reply string) error {
	return database.AddKeywords(DB, keyword, reply)
}

type KeyWordNotFoundError struct{}

func (k *KeyWordNotFoundError) Error() string {
	return "Word not found."
}

func DelKeyword(keyword string) error {
	ID, err := database.PeekKeywords(DB, keyword)
	if err != nil {
		return err
	}
	if ID == -1 {
		return &KeyWordNotFoundError{}
	}
	return database.DelKeyword(DB, ID)
}

// This method is used for loading all the keyword and reply.
func LoadKeywordReply(db *sql.DB, c *conf.Config) error {
	// Get all the keyword.
	keywords, err := database.FetchKeyword(db)
	if err != nil {
		return err
	}

	// For all keyword, get their related replies.
	// And append into config.
	for _, keyword := range keywords {
		replies, err := database.GetReplyWithKey(db, keyword.Kid)
		if err != nil {
			return err
		}
		k := keyword
		keywordAndReply := conf.KeywordsReplyType{Keyword: &k, Replies: replies}
		c.Keywords = append(c.Keywords, &keywordAndReply)
	}
	return nil
}
