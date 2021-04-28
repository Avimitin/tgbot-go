package main

import (
	"encoding/json"
	"fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type botCommands map[string]func(*tb.Message)

var (
	bc = botCommands{
		"/start":   cmdHello,
		"/ping":    cmdPing,
		"/dump":    cmdDump,
		"/weather": cmdWeather,
		"/mjx":     cmdMJX,
		"/ghs":     cmdGhs,
		"/eh":      cmdEH,
		"/setperm": cmdSetPerm,
	}
)

func cmdHello(m *tb.Message) {
	send(m.Chat,
		"你好!\nHello!\nBonjour!\nHallo!\nこんにちは!\nHola!\nCiao!")
}

func cmdPing(m *tb.Message) {
	now := time.Now()
	msg := send(m.Chat, "pong!")
	delay := time.Now().Sub(now).Milliseconds()
	edit(msg, fmt.Sprintf("it takes bot %d ms to send pong to your DC", delay))
}

func cmdDump(m *tb.Message) {
	if m.Payload != "" && m.Payload == "json" {
		data, err := json.MarshalIndent(m, "", "  ")
		if err == nil {
			send(m.Chat, string(data))
			return
		}

		send(m.Chat, err.Error())
		return
	}

	send(m.Chat, unwrapMsg(m), &tb.SendOptions{ParseMode: "HTML"})
}

func cmdWeather(m *tb.Message) {
	if m.Payload == "" {
		send(m.Chat, "no specific city")
		return
	}

	msg := send(m.Chat, "requesting api...")

	var city = m.Payload
	weather, err := getWeather(city)
	if err != nil {
		edit(msg, fmt.Sprintf("failed to fetch %s weather: %v", city, err))
		return
	}

	edit(msg, weather, &tb.SendOptions{ParseMode: "html"})
}

func cmdMJX(m *tb.Message) {
	msg := send(m.Chat, "requesting api...")

	url, err := getMJX()

	if err != nil {
		edit(msg, fmt.Sprintf("request failed: %v", err))
		return
	}

	edit(msg, fmt.Sprintf(
		`<a href="tg://user?id=%d">%s</a>, the <a href="%s">pic</a> you request have arrived.`,
		m.Sender.ID, m.Sender.FirstName, url), &tb.SendOptions{ParseMode: "html"})
}

func cmdGhs(m *tb.Message) {
	msg := send(m.Chat, "requesting...")
	defer b.Delete(msg)

	picURL, fileURL, err := getImage()
	if err != nil {
		send(m.Chat, "request failed: "+err.Error())
		return
	}

	send(m.Chat, &tb.Photo{
		File: tb.FromURL(picURL),
		Caption: fmt.Sprintf(
			"Source: [Click](%s) \\| Original File: [Click](%s)", picURL, fileURL,
		)},
		&tb.SendOptions{ParseMode: "markdownv2"},
	)
}

func cmdEH(m *tb.Message) {
	if m.Payload == "" {
		send(m.Chat, "Usage: /eh <URL>")
		return
	}
	msg := send(m.Chat, "handling...")
	defer b.Delete(msg)
	pht, opt := wrapEHData(m, "")
	send(m.Chat, pht, opt, &tb.SendOptions{ParseMode: "HTML"})
}

func cmdSetPerm(m *tb.Message) {
	if m.Payload == "" || m.Payload == "help" {
		send(m.Chat, "Usage: /setperm <ID> <PERM>\n\n"+
			"permission options: o(wner)|a(dmin)|m(anager)|n(ormal)|b(an)")
		return
	}

	send(m.Chat, setPerm(m.Payload))
}
