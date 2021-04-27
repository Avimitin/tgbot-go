package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/Avimitin/go-bot/modules/net"
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

	baseURL := "https://konachan.com/post.json?limit=50"
	resp, err := net.Get(baseURL)
	if err != nil {
		edit(msg, "Error occur, please try again later")
		return
	}

	var k []KonachanResponse

	err = json.Unmarshal(resp, &k)
	if err != nil {
		edit(msg, "failed to decode msg")
		return
	}

	rand.Seed(time.Now().Unix())
	var i = rand.Intn(50)
	var picURL string
	var fileURL string

	if len(k) < i && len(k) > 0 {
		picURL = k[0].JpegURL
		fileURL = k[0].FileURL
	} else if len(k) >= i {
		picURL = k[i].JpegURL
		fileURL = k[i].FileURL
	} else {
		edit(msg, "api no response")
		return
	}

	b.Delete(msg)
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
