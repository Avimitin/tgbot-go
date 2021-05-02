package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Avimitin/go-bot/modules/database"
	"github.com/Avimitin/go-bot/modules/eh"
	"github.com/Avimitin/go-bot/modules/konachan"
	"github.com/Avimitin/go-bot/modules/net"
	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	PermOwner int32 = iota
	PermAdmin
	PermChannelManager
	PermNormal
	PermBan
)

var (
	// cmdIR store identity requirement for command
	cmdIR = map[string]int32{
		"/setperm": PermAdmin,
		"/ehp":     PermChannelManager,
	}
)

func authPerm(perm int32, cmd string) bool {
	orderPerm, ok := cmdIR[cmd]

	// return true if given command have no identity requirement assigned
	if !ok {
		return true
	}

	return perm <= orderPerm
}

func send(to tb.Recipient, what interface{}, opt ...interface{}) *tb.Message {
	botLog.Trace().
		Str("SEND TO", to.Recipient()).
		Interface("DETAILED", what).
		Send()
	m, err := b.Send(to, what, opt...)
	switch err {
	case nil:
		return m
	case tb.ErrMessageTooLong:
		b.Send(to, "message too long")
	case tb.ErrChatNotFound:
		botLog.Error().Err(err).Msg("bot not in chat")
	default:
		b.Send(to, "operation failed: "+err.Error())
		botLog.Error().Err(err).Send()
	}
	return m
}

func edit(msg tb.Editable, what interface{}, opt ...interface{}) *tb.Message {
	ms, cs := msg.MessageSig()
	botLog.Trace().Msgf("Editing msg %s, %d", ms, cs)
	m, err := b.Edit(msg, what, opt...)
	switch err {
	case nil:
		return m
	case tb.ErrMessageTooLong:
		b.Edit(msg, "message too long")
	case tb.ErrChatNotFound:
		botLog.Error().Err(err).Msg("bot not in chat")
	default:
		b.Edit(msg, "operation failed: "+err.Error())
		botLog.Error().Err(err).Send()
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

func wrapEHData(ehURL string, comment string) (*tb.Photo, *tb.ReplyMarkup, error) {
	data, err := eh.GetComic(ehURL)
	if err != nil {
		return nil, nil, fmt.Errorf("Request failed: %v", err)
	}

	if len(data.Medas) < 1 {
		return nil, nil, fmt.Errorf("Request failed: comic not found")
	}

	metadata := data.Medas[0]

	if metadata.Error != "" {
		return nil, nil, fmt.Errorf("Request failed: %s", metadata.Error)
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
		btnLike    = menu.Data("👍 "+metadata.Rating, "like-button")
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
			File:    tb.FromURL(metadata.Thumb),
			Caption: caption,
		},
		menu,
		nil
}

func getWeather(city string) (string, error) {
	url := "https://wttr.in/" + city + "?format=%l的天气:+%c+温度:%t+湿度:%h+降雨量:%p"
	resp, err := net.Get(url)
	if err != nil {
		return "", fmt.Errorf("get %s weather: %v", city, err)
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, fmt.Sprintf("https://wttr.in/%s.png", city), resp), nil
}

func getMJX() (string, error) {
	var mjxURL string
	rand.Seed(time.Now().UnixNano())

	if rand.Float32() < 0.5 {
		mjxURL = "http://api.vvhan.com/api/tao?type=json"
	} else {
		mjxURL = "http://api.uomg.com/api/rand.img3?format=json"
	}

	botLog.Trace().Str("func", "getMJX").Msgf("requesting url %s", mjxURL)

	data, err := net.Get(mjxURL)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	var mjx = struct {
		Pic    string `json:"pic"`
		Imgurl string `json:"imgurl"`
	}{}
	err = json.Unmarshal(data, &mjx)
	if err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	if mjx.Imgurl != "" {
		return mjx.Imgurl, nil
	}

	if mjx.Pic != "" {
		return (mjx.Pic), nil
	} else {
		return "", fmt.Errorf("fail to fetch pic")
	}
}

func getImage() (string, string, error) {
	baseURL := "https://konachan.com/post.json?limit=50"
	resp, err := net.Get(baseURL)
	if err != nil {
		return "", "", fmt.Errorf("Error occur, please try again later")
	}

	var k []konachan.KonachanResponse

	err = json.Unmarshal(resp, &k)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode msg")
	}

	rand.Seed(time.Now().Unix())
	var i = rand.Intn(50)

	if len(k) >= i {
		return k[i].JpegURL, k[i].FileURL, nil
	}

	if len(k) < i && len(k) > 0 {
		return k[0].JpegURL, k[0].FileURL, nil
	}

	return "", "", fmt.Errorf("api no response")
}

func setPerm(argument string) string {
	args := strings.Fields(argument)
	if len(args) < 2 {
		return "argument not valid, more descriptions refer to /setperm help"
	}

	var (
		idStr = args[0]
		perm  = args[1]
		id    int
	)

	id, convErr := strconv.Atoi(idStr)
	if convErr != nil {
		return fmt.Sprintf("parsed argument %q: %v", idStr, convErr)
	}

	var err error
	var user *database.User
	fn := func(id int, permid int32) {
		user, err = DB.SetUser(id, permid)
	}

	switch perm {
	case "owner", "o":
		fn(id, 0)
	case "admin", "a":
		fn(id, 1)
	case "manager", "m":
		fn(id, 2)
	case "normal", "n":
		fn(id, 3)
	case "ban", "b":
		fn(id, 4)
	default:
		return "argument not valid, more descriptions refer to /setperm help"
	}

	if user == nil {
		return "user not found"
	}

	if err != nil {
		return fmt.Sprintf("failed to set user %d permission: %v", id, err)
	}

	return fmt.Sprintf("user %d permission has set to %q successfully", user.UserID, user.PermDesc)
}

// isCommand returns true if message starts with a "bot_command" entity.
func isCommand(m *tb.Message) bool {
	if m.Entities == nil || len(m.Entities) == 0 {
		return false
	}

	entity := m.Entities[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}

func msgCommand(m *tb.Message) (string, bool) {
	if !isCommand(m) {
		return "", false
	}

	entity := m.Entities[0]
	// match commands like /start@exampleBot
	commandWithAt := m.Text[:entity.Length]

	// get only command
	if i := strings.Index(commandWithAt, "@"); i != -1 {
		commandWithAt = commandWithAt[:i]
	}

	return commandWithAt, true
}

func assertPayload(m *tb.Message, sub ...string) bool {
	if len(m.Payload) == 0 {
		return false
	}

	if len(sub) < 1 {
		return true
	}

	return m.Payload == sub[0]
}
