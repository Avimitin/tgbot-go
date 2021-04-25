package main

import (
	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func unwrapMsg(m *tb.Message) string {
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
		"<b>ID:</b> <code>%v</code>\n"

	if !m.IsReply() {
		text = fmt.Sprintf(text,
			m.Chat.ID, m.Chat.Type, m.Chat.Username,
			m.Sender.ID, m.Sender.Username, m.Sender.FirstName,
			m.Sender.LastName, m.Sender.LanguageCode,
			m.ID)
	} else if m.ReplyTo.IsForwarded() {
		text = fmt.Sprintf(text,
			m.Chat.ID, m.Chat.Type, m.Chat.Username,
			m.ReplyTo.OriginalChat.ID, m.ReplyTo.OriginalChat.Type,
			m.ReplyTo.OriginalSender.Username, m.ReplyTo.OriginalSender.FirstName,
			m.ReplyTo.OriginalSender.LastName, m.ReplyTo.OriginalSender.LanguageCode,
			m.ReplyTo.ID)
	} else {
		text = fmt.Sprintf(text,
			m.ReplyTo.Chat.ID, m.ReplyTo.Chat.Type, m.ReplyTo.Chat.Username,
			m.ReplyTo.Sender.ID, m.ReplyTo.Sender.Username, m.ReplyTo.Sender.FirstName,
			m.ReplyTo.Sender.LastName, m.ReplyTo.Sender.LanguageCode,
			m.ReplyTo.ID)
	}
	return text
}
