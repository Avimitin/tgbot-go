package bot

import (
	"fmt"
	"strings"
	"time"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var botCMD = command{
	"start": start,
	"ping":  ping,
	"dump":  dump,
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

func dump(m *bapi.Message) error {
	var text = "<b>Message Information</b>\n" +
		"=== <b>CHAT</b> ===\n" +
		"<b>ID:</b> <code>%v</code>\n" +
		"<b>TYPE:</b> <code>%v</code>\n" +
		"<b>USERNAME:</b> <code>%v</code>\n" +
		"=== <b>USER</b> ===\n" +
		"<b>ID:</b> <code>%v</code>\n" +
		"<b>USERNAME:</b> <code>%v</code>\n" +
		"<b>NICKNAME:</b> <code>%v %v</code>\n" +
		"<b>LANGUAGE:</b> <code>%v</code>\n" +
		"=== <b>MSG</b> ===\n" +
		"<b>ID:</b> <code>%v</code>\n" +
		"<b>TEXT:</b> %v"

	if reply := m.ReplyToMessage; reply != nil {
		text = fmt.Sprintf(text,
			reply.Chat.ID, reply.Chat.Type, reply.Chat.UserName,
			reply.From.ID, reply.From.UserName, reply.From.FirstName, reply.From.LastName, reply.From.LanguageCode,
			reply.MessageID, reply.Text)
	} else {
		text = fmt.Sprintf(text,
			m.Chat.ID, m.Chat.Type, m.Chat.UserName,
			m.From.ID, m.From.UserName, m.From.FirstName, m.From.LastName, m.From.LanguageCode,
			m.MessageID, m.Text)
	}

	_, err := sendP(text, m.Chat.ID, "HTML")
	if err != nil {
		return errF("dump", err, "fail to send dump message")
	}
	return nil
}

func kick(m *bapi.Message) error {
	is, err := isAdmin(m.From.ID, m.Chat)
	if err != nil {
		errMsg := "fail to get user permission"
		sendT(errMsg, m.Chat.ID)
		return errF("kick", err, errMsg)
	}
	// if command caller are not admin
	if !is {
		_, err = sendT("YOU ARE NOT ADMIN! DONT TOUCH THIS COMMAND!", m.Chat.ID)
		if err != nil {
			return errF("kick", err, "fail to send kick alert")
		}
		return nil
	}

	if m.ReplyToMessage == nil {
		_, err = sendT("You should reply to a user to kick him.", m.Chat.ID)
		if err != nil {
			return errF("kick", err, "fail to send usage")
		}
		return nil
	}

	userToKick := m.ReplyToMessage.From.ID
	err = kickUser(userToKick, m.Chat.ID, time.Now().Unix()+1)
	if err != nil {
		return errF("kick", err, "fail to kick user")
	}

	_, err = sendT("user has been kick forever", m.Chat.ID)
	if err != nil {
		return errF("kick", err, "fail to send usage")
	}
	return nil
}
