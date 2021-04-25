package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Avimitin/go-bot/modules/eh"
	tb "gopkg.in/tucnak/telebot.v2"
)

func send(to tb.Recipient, what interface{}, opt ...interface{}) *tb.Message {
	m, err := b.Send(to, what, opt...)
	switch err {
	case nil:
		return m
	case tb.ErrMessageTooLong:
		b.Send(to, "message too long")
	default:
		log.Println("[ERROR]", err)
	}
	return m
}

func edit(msg tb.Editable, what interface{}, opt ...interface{}) *tb.Message {
	m, err := b.Edit(msg, what, opt...)
	switch err {
	case nil:
		return m
	case tb.ErrMessageTooLong:
		b.Edit(msg, "message too long")
	default:
		log.Println("[ERROR]", err)
	}
	return m
}

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

func replaceTag(inputTags []string, c chan string) {
	var tags string
	for _, tag := range inputTags {
		tag = strings.ReplaceAll(tag, " ", "_")
		tag = strings.ReplaceAll(tag, "-", "_")
		tags += "#" + tag + " "
	}
	c <- tags
}

func wrapEHData(m *tb.Message, comment string) (interface{}, interface{}) {
	data, err := eh.GetComic(m.Payload)
	if err != nil {
		return fmt.Sprintf("Request failed: %v", err), nil
	}

	if len(data.Medas) < 1 {
		return "Request failed: comic not found", nil
	}

	metadata := data.Medas[0]

	if metadata.Error != "" {
		return fmt.Sprintf("Request failed: %s", metadata.Error), nil
	}

	tagC := make(chan string)
	go replaceTag(metadata.Tags, tagC)

	var caption string
	caption += fmt.Sprintf("📖标题: <code>%s</code>\n", metadata.TitleJpn)
	caption += fmt.Sprintf("🗂️类别: %s\n", metadata.Category)
	caption += fmt.Sprintf("🏷️标签: %v\n", <-tagC)
	if comment != "" {
		caption += fmt.Sprintf("💬评论: %v", comment)
	}

	menu := &tb.ReplyMarkup{}

	var (
		btnLike    = menu.Text("👍 " + metadata.Rating)
		btnCollect = menu.URL("⭐ 点击收藏",
			fmt.Sprintf(
				"https://e-hentai.org/gallerypopups.php?gid=%d&t=%s&act=addfav",
				metadata.Gid, metadata.Token,
			),
		)
		btnInSite = menu.URL("🐼 里站Link",
			fmt.Sprintf("https://exhentai.org/g/%d/%s/", metadata.Gid, metadata.Token))
		btnOutSite = menu.URL("🔗 表站Link",
			fmt.Sprintf("https://e-hentai.org/g/%d/%s/", metadata.Gid, metadata.Token))
	)

	menu.Inline(
		menu.Row(btnLike, btnCollect),
		menu.Row(btnInSite, btnOutSite),
	)

	return &tb.Photo{
			File:      tb.FromURL(metadata.Thumb),
			Caption:   caption,
			ParseMode: "HTML",
		},
		menu
}
