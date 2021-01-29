package bot

import (
	"fmt"
	"strings"
	"time"

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
		return errF("start", err, "fail to send message")
	}
	return nil
}

func ping(m *bapi.Message) error {
	now := time.Now()
	msg, err := sendT("pong!", m.Chat.ID)
	if err != nil {
		return errF("ping", err, "fail to send msg")
	}
	delay := now.Sub(now).Milliseconds()
	text := fmt.Sprintf("bot 与 Telegram 服务器的延迟大约为 %d 毫秒", delay)
	_, err = editT(text, m.Chat.ID, msg.MessageID)
	if err != nil {
		return errF("ping", err, "fail to edit msg")
	}
	return nil
}
