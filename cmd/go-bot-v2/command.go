package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	}
)

func handleErr(e error) {
	log.Println(e)
}

func send(to tb.Recipient, what interface{}, opt ...interface{}) *tb.Message {
	m, err := b.Send(to, what, opt...)
	if err != nil {
		handleErr(fmt.Errorf("sending %#v: %v", what, err))
	}
	return m
}

func edit(msg tb.Editable, what interface{}, opt ...interface{}) *tb.Message {
	m, err := b.Edit(msg, what, opt...)
	if err != nil {
		handleErr(fmt.Errorf("editing msg to %#v: %v", what, err))
	}
	return m
}

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
	send(m.Chat, text, &tb.SendOptions{ParseMode: "HTML"})
}

func getWeather(city string) (string, error) {
	url := "https://wttr.in/" + city + "?format=%l的天气:+%c+温度:%t+湿度:%h+降雨量:%p"
	resp, err := net.Get(url)
	if err != nil {
		return "", fmt.Errorf("get %s weather: %v", city, err)
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, fmt.Sprintf("https://wttr.in/%s.png", city), resp), nil
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

	var data []byte
	var err error
	rand.Seed(time.Now().UnixNano())

	if rand.Float32() < 0.5 {
		data, err = net.Get("http://api.vvhan.com/api/tao?type=json")
	} else {
		data, err = net.Get("http://api.uomg.com/api/rand.img3?format=json")
	}
	if err != nil {
		edit(msg, "request failed:"+err.Error())
		return
	}

	var mjx = struct {
		Pic    string `json:"pic"`
		Imgurl string `json:"imgurl"`
	}{}
	err = json.Unmarshal(data, &mjx)
	if err != nil {
		edit(msg, "unmarshal failed:"+err.Error())
		return
	}

	editURL := func(url string) {
		edit(msg, fmt.Sprintf(
			`<a href="tg://user?id=%d">%s</a>, the <a href="%s">pic</a> you request have arrived.`,
			m.Sender.ID, m.Sender.FirstName, url), &tb.SendOptions{ParseMode: "html"})
	}
	if mjx.Imgurl != "" {
		editURL(mjx.Imgurl)
	} else if mjx.Pic != "" {
		editURL(mjx.Pic)
	} else {
		edit(msg, "fail to fetch pic")
	}
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
