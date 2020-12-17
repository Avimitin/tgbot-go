package auth

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type MyError struct {
	info string
}

func (e *MyError) Error() string {
	return e.info
}

func IsMe(me int, uid int) bool {
	return uid == me
}

func getAdmin(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, c chan []int) {
	members, err := bot.GetChatAdministrators(chat.ChatConfig())
	if err != nil {
		log.Print("[ERR]", err)
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
