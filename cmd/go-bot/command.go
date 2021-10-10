package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Avimitin/go-bot/modules/config"
	"github.com/Avimitin/go-bot/modules/currency"
	"github.com/Avimitin/go-bot/modules/mark"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

type botCommands map[string]func(*tb.Message)

var (
	bc = botCommands{
		"/start":    cmdHello,
		"/ping":     cmdPing,
		"/dump":     cmdDump,
		"/weather":  cmdWeather,
		"/mjx":      cmdMJX,
		"/ghs":      cmdGhs,
		"/eh":       cmdEH,
		"/ehp":      cmdEHPost,
		"/setperm":  cmdSetPerm,
		"/repeat":   cmdRepeat,
		"/remake":   cmdRemake,
		"/exchange": cmdExchange,
		"/mark":     cmdAddMark,
		"/lsmark":   cmdGetMark,
		"/delmark":  cmdDelMark,
		"/collect":  cmdCollectMessage,
		"/me":       cmdMe,
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

	pht, opt, err := wrapEHData(m.Payload, nil)
	if err != nil {
		log.Error().Err(err).Str("FUNC", "cmdEH").Msg("wrap eh data failed")
		send(m.Chat, err.Error())
	}

	send(m.Chat, pht, &tb.SendOptions{ParseMode: "HTML", ReplyMarkup: opt})
}

func cmdEHPost(m *tb.Message) {
	if !assertPayload(m) {
		send(m.Chat, "Usage: /ehp <URL>")
		return
	}

	send(m.Chat, "any comment?")
	regisNextStep(m.Chat.ID, m.Sender.ID, contextData{"url": m.Payload}, handleComicMetadata)
}

func handleComicMetadata(m *tb.Message, p contextData) error {
	defer delContext(m.Chat.ID, m.Sender.ID)

	ehURL, ok := p["url"]
	if !ok {
		send(m.Chat, "internal error: url missing")
		return fmt.Errorf("internal error: url missing")
	}

	msg := send(m.Chat, "requesting...")

	pht, opt, err := wrapEHData(ehURL, m)

	if err != nil {
		return err
	}

	success := send(
		&tb.Chat{
			ID: config.GetEhPostChannelID(),
		},

		pht,

		&tb.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: opt,
		},
	)

	if success == nil {
		edit(msg, "post failed")
	} else {
		edit(msg, "post success")
	}

	return nil
}

func cmdSetPerm(m *tb.Message) {
	if m.Payload == "" || m.Payload == "help" {
		send(m.Chat, "Usage: /setperm <ID> <PERM>\n\n"+
			"permission options: o(wner)|a(dmin)|m(anager)|n(ormal)|b(an)")
		return
	}

	send(m.Chat, setPerm(m.Payload))
}

func cmdRepeat(m *tb.Message) {
	if !m.IsReply() {
		send(m.Chat, "use reply")
		return
	}

	send(m.Chat, encodeEntity(m.ReplyTo), &tb.SendOptions{ParseMode: "HTML"})
}

func cmdRemake(m *tb.Message) {
	user := fmt.Sprintf(`<a href="tg://user?id=%d">%s</a>`, m.Sender.ID, m.Sender.FirstName)

	send(m.Chat,
		user+" 你复活啦？你知道你寄吧谁吗？",
		&tb.SendOptions{ParseMode: "HTML"},
	)
}

func cmdExchange(m *tb.Message) {
	if m.Payload == "" || m.Payload == "help" {
		send(m.Chat, "Usage: /exchange 100 cny usd")
		return
	}

	args := strings.Fields(m.Payload)

	if len(args) < 3 {
		send(m.Chat, "Usage: /exchange 100 cny usd")
		return
	}

	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		send(m.Chat, fmt.Sprintf("Unsupport amount: %q", args[0]))
		return
	}

	msg := send(m.Chat, "requesting...")
	result, err := currency.CalculateExchange(amount, args[1], args[2])
	if err != nil {
		send(m.Chat, fmt.Sprintf("calculate failed: %v", err))
		return
	}

	edit(msg, fmt.Sprintf("%s %s\n=\n%s %f", args[1], args[0], args[2], result))
}

func cmdAddMark(m *tb.Message) {
	if !m.IsReply() {
		send(m.Chat, "reply to a message to mark it")
		return
	}

	description := m.ReplyTo.Text
	if len(description) > 15 {
		description = description[:15] + "..."
	}

	chatID := fmt.Sprintf("%d", m.Chat.ID)
	if m.Chat.Type == tb.ChatGroup || m.Chat.Type == tb.ChatChannelPrivate {
		chatID = chatID[4:]
	}

	link := fmt.Sprintf("https://t.me/c/%s/%d", chatID, m.ReplyTo.ID)
	mark.AddMark(m.Sender.FirstName, link, description)

	send(m.Chat,
		fmt.Sprintf("%q marked", description))
}

func cmdGetMark(m *tb.Message) {
	result, err := mark.GetMark(m.Sender.FirstName)
	if err != nil {
		send(m.Chat, fmt.Sprintf("getting mark: %v", err))
		return
	}

	var text string
	for _, entry := range result {
		text += fmt.Sprintf(`[%d] %q <a href="%s">link</a>`, entry.ID, entry.Description, entry.Link)
		text += "\n"
	}

	send(m.Chat, text, &tb.SendOptions{ParseMode: "HTML"})
}

func cmdDelMark(m *tb.Message) {
	targetID, err := strconv.ParseInt(m.Payload, 10, 32)
	if err != nil {
		send(m.Chat,
			fmt.Sprintf("Usage example: /delmark {{target id (int32)}}"))
		return
	}

	err = mark.DelMark(m.Sender.FirstName, int32(targetID))

	if err != nil {
		send(m.Chat,
			fmt.Sprintf("delete mark: %v", err))
		return
	}

	send(m.Chat, "delete success")
}

// cmdCollectMessage force bot enter collect mode and record all the
// message from user.
func cmdCollectMessage(m *tb.Message) {
	send(m.Chat, "你可以开始发消息给我了。\n输入 /end_collect 结束录入。")
	regisNextStep(m.Chat.ID, m.Sender.ID, contextData{}, collectMessage)
}

type msgInfoForCollect struct {
	text     string
	username string
}

var (
	collectedData   = map[int][]msgInfoForCollect{}
	collectDataLock = sync.Mutex{}
)

func collectMessage(m *tb.Message, p contextData) error {
	collectDataLock.Lock()
	defer collectDataLock.Unlock()

	currentUserCollectData, exist := collectedData[m.Sender.ID]

	if !exist {
		currentUserCollectData = make([]msgInfoForCollect, 0, 2)
	}

	if strings.HasPrefix(m.Text, "/end_collect") {
		defer delContext(m.Chat.ID, m.Sender.ID)
		defer func() {
			delete(collectedData, m.Sender.ID)
		}()

		send(m.Chat, "录入结束, 正在合成中...")

		collectedDataLiteral := "Collected messages:\n\n"
		for _, msg := range currentUserCollectData {
			collectedDataLiteral += fmt.Sprintf("<b>%s:</b> %s", msg.username, msg.text)
			collectedDataLiteral += "\n---\n"
		}

		send(m.Chat, collectedDataLiteral, &tb.SendOptions{ParseMode: "HTML"})
		return nil
	}

	var username string
	if m.OriginalSender != nil && m.OriginalSender.FirstName != "" {
		username = m.OriginalSender.FirstName
	} else if m.OriginalSenderName != "" {
		username = m.OriginalSenderName
	} else {
		username = "Anoynomous"
	}

	currentUserCollectData = append(currentUserCollectData,
		msgInfoForCollect{
			username: username,
			text:     m.Text,
		},
	)

	collectedData[m.Sender.ID] = currentUserCollectData
	regisNextStep(m.Chat.ID, m.Sender.ID, contextData{}, collectMessage)
	return nil
}

func cmdMe(m *tb.Message) {
	userLink := createUserLink(m.Sender.FirstName, m.Sender.ID)

	if m.Payload != "" {
		send(m.Chat, fmt.Sprintf("%s %s 了", userLink, m.Payload), newHTMLParseMode())
	} else {
		send(m.Chat, fmt.Sprintf("你是 %s", userLink), newHTMLParseMode())
	}
}
