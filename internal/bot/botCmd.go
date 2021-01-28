package bot

import (
	"fmt"
	"log"
	"strings"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var botCMD = command{
	"start": start,
}

func cmdArgv(msg *bapi.Message) []string {
	args := strings.Fields(msg.Text)
	if len(args) > 1 {
		args = args[1:]
		return args
	}
	return nil
}

func start(m *bapi.Message) error {
	username := m.From.UserName
	if username == "" {
		username = m.From.FirstName
	}
	userID := m.From.ID
	userLink := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, userID, username)
	_, err := sendT("Hi, "+userLink, m.Chat.ID)
	if err != nil {
		log.Println("[start]", err)
		return fmt.Errorf("%v", err)
	}
	return nil
}
