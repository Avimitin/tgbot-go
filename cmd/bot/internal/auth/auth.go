package auth

import (
	"database/sql"
	"github.com/Avimitin/go-bot/cmd/bot/internal/conf"
	"github.com/Avimitin/go-bot/cmd/bot/internal/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type MyError struct {
	info string
}

func (e *MyError) Error() string {
	return e.info
}

func IsCreator(creator int, uid int) bool {
	return uid == creator
}

func DBIsAuthGroups(DB *sql.DB, gid int64) bool {
	groups, err := database.SearchGroups(DB)
	if err != nil {
		return false
	}

	id := binarySearch(gid, groups)
	if id != -1 {
		return true
	}
	return false
}

func CFGIsAuthGroups(cfg *conf.Config, gid int64) bool {
	if SearchGroupsAlg(cfg, gid) != -1 {
		return true
	}
	return false
}

func SearchGroupsAlg(cfg *conf.Config, gid int64) int {
	lo, hi := 0, len(cfg.Groups)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if gid < cfg.Groups[mid] {
			hi = mid - 1
		} else if gid > cfg.Groups[mid] {
			lo = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func getAdmin(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, c chan []int) {
	members, err := bot.GetChatAdministrators(chat.ChatConfig())
	if err != nil {
		log.Print("[GET_ADMIN_ERROR]", err)
		c <- nil
		close(c)
	}
	admins := make([]int, len(members))
	for i, member := range members {
		admins[i] = member.User.ID
	}
	c <- admins
}

func IsAdmin(bot *tgbotapi.BotAPI, uid int, chat *tgbotapi.Chat) (bool, error) {
	c := make(chan []int)
	go getAdmin(bot, chat, c)
	admins := <-c
	if admins == nil {
		return false, &MyError{info: "Error fetching admin"}
	}

	for _, admin := range admins {
		if uid == admin {
			return true, nil
		}
	}
	return false, nil
}

func binarySearch(target int64, groups []database.Group) int {
	var lo, hi = 0, len(groups) - 1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if target < groups[mid].GroupID {
			hi = mid - 1
		} else if target > groups[mid].GroupID {
			lo = mid + 1
		} else {
			return mid
		}
	}
	return -1
}
