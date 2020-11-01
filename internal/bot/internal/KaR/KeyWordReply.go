package KaR

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/internal/conf"
	"github.com/Avimitin/go-bot/internal/database"
	"math/rand"
	"time"
)

func RegexKAR(msg string, d *conf.BotData) (string, bool) {
	kw, e := d.KAR[msg]
	if e {
		if length := len(kw); length > 1 {
			rand.Seed(time.Now().Unix())
			return kw[rand.Intn(length)], true
		}
		return kw[0], true
	}
	return "", false
}

// Function for set add new keyword and related reply.
func SetKeywordIntoDB(keyword string, reply string, DB *sql.DB) (int, error) {
	return database.AddKeywords(DB, keyword, reply)
}

func SetKeywordIntoCFG(keyword string, reply string, d *conf.BotData) {
	kw, e := d.KAR[keyword]
	if e {
		kw = append(kw, reply)
	}
	r := []string{reply}
	d.KAR[keyword] = r
}

func SetKeyword(keyword string, reply string, d *conf.BotData, db *sql.DB) string {
	SetKeywordIntoCFG(keyword, reply, d)
	id, err := SetKeywordIntoDB(keyword, reply, db)
	if id == -1 && err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Insert successful, id: %d", id)
}

type KeyWordNotFoundError struct{}

func (k *KeyWordNotFoundError) Error() string {
	return "Word not found."
}

func DelKeywordAtDB(keyword string, db *sql.DB) error {
	ID, err := database.PeekKeywords(db, keyword)
	if err != nil {
		return err
	}
	if ID == -1 {
		return &KeyWordNotFoundError{}
	}
	return database.DelKeyword(db, ID)
}

func DelKeywordAtCFG(keyword string, d *conf.BotData) error {
	_, e := d.KAR[keyword]
	if e {
		delete(d.KAR, keyword)
	}
	return &KeyWordNotFoundError{}
}

func DelKeyword(keyword string, d *conf.BotData, db *sql.DB) error {
	err := DelKeywordAtCFG(keyword, d)
	if err != nil {
		return err
	}
	err = DelKeywordAtDB(keyword, db)
	if err != nil {
		return err
	}
	return nil
}

// This method is used for loading all the keyword and reply.
func LoadKeywordReply(db *sql.DB, d *conf.BotData) error {
	kts, err := database.FetchKeyword(db)
	if err != nil {
		return err
	}
	for _, kt := range kts {
		replies, err := database.GetReplyWithKey(db, kt.I)
		if err != nil {
			return err
		}
		d.KAR[kt.K] = replies
	}
	return nil
}

func ListKeywordAndReply(d *conf.BotData) string {
	text := "KEYS AND REPLIES:\n"
	for key, word := range d.KAR {
		text += fmt.Sprintf("%s. K: %s\n", key, word)
		text += "\n"
	}
	return text
}

func DelReplies(keyword string, reply string, d *conf.BotData, db *sql.DB) error {
	err := DelRepliesAtData(keyword, reply, d)
	if err != nil {
		return err
	}
	err = DelRepliesAtDatabase(reply, db)
	if err != nil {
		return err
	}
	return nil
}

func DelRepliesAtData(keyword string, reply string, d *conf.BotData) error {
	replies, ok := d.KAR[keyword]
	if !ok {
		return &KeyWordNotFoundError{}
	}
	for i, r := range replies {
		if r == reply {
			d.KAR[keyword] = append(replies[:i], replies[i+1:]...)
			return nil
		}
	}
	return &KeyWordNotFoundError{}
}

func DelRepliesAtDatabase(reply string, db *sql.DB) error {
	i, err := database.PeekReply(db, reply)
	if err != nil {
		return err
	}
	err = database.DelReply(db, i)
	if err != nil {
		return err
	}
	return nil
}
