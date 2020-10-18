package bot

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/internal/bot/internal"
	"github.com/Avimitin/go-bot/internal/conf"
	"github.com/Avimitin/go-bot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"strings"
)

func Reply(msg *tgbotapi.Message) (tgbotapi.Message, error) {
	kws, isExist := Regexp(msg.Text, cfg)
	if isExist {
		if len(kws.Replies) == 1 {
			return internal.SendTextMsg(bot, msg.Chat.ID, kws.Replies[0])
		}
		randNum := rand.Intn(len(kws.Replies))
		return internal.SendTextMsg(bot, msg.Chat.ID, kws.Replies[randNum])
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

// Function for set add new keyword and related reply.
func SetKeywordIntoDB(keyword string, reply string) (int, error) {
	return database.AddKeywords(DB, keyword, reply)
}

func SetKeywordIntoCFG(kid int, keyword string, reply string) {
	currentKeywords := &cfg.Keywords
	var currentKAndR *conf.KeywordsReplyType
	for _, kar := range cfg.Keywords {
		if keyword == kar.Keyword.Word {
			currentKAndR = kar
		}
	}
	if currentKAndR != nil {
		currentKAndR.Replies = append(currentKAndR.Replies, reply)
		return
	}
	kw := conf.KeywordType{
		Kid:  kid,
		Word: keyword,
	}

	newKAR := conf.KeywordsReplyType{Keyword: &kw, Replies: []string{reply}}
	*currentKeywords = append(*currentKeywords, &newKAR)
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

func DelKeywordFromCFG(kid int) error {
	for i, kw := range cfg.Keywords {
		if kw.Keyword.Kid == kid {
			leftPart := cfg.Keywords[:i]
			cfg.Keywords = append(leftPart, cfg.Keywords[i+1:]...)
			return nil
		}
	}
	return &KeyWordNotFoundError{}
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

func ListKeywordAndReply() string {
	text := "KEYS AND REPLIES:\n"
	for place, kar := range cfg.Keywords {
		text += fmt.Sprintf("%d. K: %s\n", place, kar.Keyword.Word)
		for _, r := range kar.Replies {
			text += fmt.Sprintf("|--: %s\n", r)
		}
		text += "\n"
	}
	return text
}
