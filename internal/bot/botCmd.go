package bot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Avimitin/go-bot/internal/net"
	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var botCMD = command{
	"start":     start,
	"ping":      ping,
	"dump":      dump,
	"kick":      kick,
	"shutup":    shutUp,
	"disshutup": disShutUp,
	"weather":   weather,
	"mjx":       mjx,
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
	sendP("Hi, "+userLink, m.Chat.ID, "HTML")
	return nil
}

func ping(m *bapi.Message) error {
	now := time.Now()
	msg := sendT("pong!", m.Chat.ID)
	current := time.Now()
	delay := current.Sub(now).Milliseconds()
	text := fmt.Sprintf("bot 与 Telegram 服务器的延迟大约为 %d 毫秒", delay)
	_, err := editT(text, m.Chat.ID, msg.MessageID)
	if err != nil {
		sendT(fmt.Sprintf("edit %s failed: %v", msg.Text, err), m.Chat.ID)
		return errF("ping", err, "edit fail")
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

	sendP(text, m.Chat.ID, "HTML")
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
		sendT("YOU ARE NOT ADMIN! DONT TOUCH THIS COMMAND!", m.Chat.ID)
		return nil
	}

	if m.ReplyToMessage == nil {
		sendT("You should reply to a user to kick him.", m.Chat.ID)
		return nil
	}

	userToKick := m.ReplyToMessage.From.ID
	err = kickUser(userToKick, m.Chat.ID, time.Now().Unix()+1)
	if err != nil {
		return errF("kick", err, "fail to kick user")
	}

	sendT("user has been kick forever", m.Chat.ID)
	return nil
}

func punishNoPermissionUser(m *bapi.Message) error {
	var err error
	respMsg := sendT("generating....", m.Chat.ID)

	var minLimit, maxLimit int64 = 60, 300
	rand.Seed(time.Now().Unix())
	randTime := rand.Int63n(maxLimit-minLimit) + minLimit
	err = editUserPermissions(m.From.ID, m.Chat.ID, time.Now().Unix()+randTime, false)
	if err != nil {
		sendT("fail to limit user:"+err.Error(), m.Chat.ID)
		return errF("shutUp", err, "fail to limit user")
	}

	respMsg, err = editT(
		fmt.Sprintf("Boom, you get a %d mins ban", randTime), m.Chat.ID, respMsg.MessageID)
	if err != nil {
		return errF("shutUp", err, "fail to edit message")
	}
	return nil
}

func shutUp(m *bapi.Message) error {
	is, err := isAdmin(m.From.ID, m.Chat)
	if err != nil {
		sendT("fail to fetch admins, please try again later", m.Chat.ID)
		return errF("shutUp", err, "fail to fetch admins")
	}
	// if user is not admin
	if !is {
		return punishNoPermissionUser(m)
	}

	if m.ReplyToMessage == nil {
		sendT("reply to a user to use this command", m.Chat.ID)
	}

	err = editUserPermissions(m.From.ID, m.Chat.ID, time.Now().Unix()+1, false)
	if err != nil {
		sendT("fail to limit user: "+err.Error(), m.Chat.ID)
		return errF("shutUp", err, "fail to limit user")
	}
	sendT("user has been forever muted", m.Chat.ID)
	return nil
}

func disShutUp(m *bapi.Message) error {
	is, err := isAdmin(m.From.ID, m.Chat)
	if err != nil {
		sendT("fail to fetch admin list", m.Chat.ID)
		return err
	}
	if !is {
		return punishNoPermissionUser(m)
	}

	if m.ReplyToMessage == nil {
		sendT("Reply to a user to recover his permission", m.Chat.ID)
		return nil
	}
	err = editUserPermissions(m.From.ID, m.Chat.ID, 0, true)
	if err != nil {
		sendT("recover user priviledge: "+err.Error(), m.Chat.ID)
		return errF("disShutUp", err, "")
	}
	sendT("User has recovered", m.Chat.ID)
	return nil
}

func weather(m *bapi.Message) error {
	argv := cmdArgv(m)
	if argv == nil {
		sendT("Gib me a city name", m.Chat.ID)
		return nil
	}

	respMsg := sendT("requesting API server...", m.Chat.ID)

	city := argv[0]
	text, err := getWeatherContext(city)
	if err != nil {
		_, err = editT("fetch weather failed: "+err.Error(), m.Chat.ID, respMsg.MessageID)
		if err != nil {
			return errF("weather", err, "edit failed")
		}
	}
	_, err = editP(text, m.Chat.ID, respMsg.MessageID, "HTML")
	if err != nil {
		return errF("weather", err, "edit failed")
	}
	return nil
}

func getWeatherContext(city string) (string, error) {
	url := "https://wttr.in/" + city + "?format=%l的天气:+%c+温度:%t+湿度:%h+降雨量:%p"
	resp, err := net.Get(url)
	if err != nil {
		return "", errF("getWeatherContext", err, "get city "+city)
	}
	return fmt.Sprintf(`<a href="%s">%s</a>`, fmt.Sprintf("https://wttr.in/%s.png", city), resp), nil
}

func mjx(m *bapi.Message) error {
	msg := sendT("requesting API server...", m.Chat.ID)

	rand.Seed(time.Now().UnixNano())
	var data []byte
	var err error
	edit := func(whatToEdit string) error {
		_, editErr := editT(whatToEdit, m.Chat.ID, msg.MessageID)
		if editErr != nil {
			sendT("edit msg failed:"+editErr.Error(), m.Chat.ID)
		}
		return err
	}

	if rand.Float32() < 0.5 {
		data, err = net.Get("http://api.vvhan.com/api/tao?type=json")
	} else {
		data, err = net.Get("http://api.uomg.com/api/rand.img3?format=json")
	}
	if err != nil {
		return edit("request failed:" + err.Error())
	}

	var mjx struct {
		Pic    string `json:"pic"`
		Imgurl string `json:"imgurl"`
	}
	err = json.Unmarshal(data, &mjx)
	if err != nil {
		return edit("unmarshal failed:" + err.Error())
	}

	editURL := func(url string) {
		_, editErr := editP(
			fmt.Sprintf(
				`<a href="tg://user?id=%d">%s</a>, the <a href="%s">pic</a> you request have arrived.`, m.From.ID, m.From.FirstName, url),
			m.Chat.ID, msg.MessageID, "HTML")
		if editErr != nil {
			sendT("edit failed:"+err.Error(), m.Chat.ID)
		}
	}
	if mjx.Imgurl != "" {
		editURL(mjx.Imgurl)
	} else if mjx.Pic != "" {
		editURL(mjx.Pic)
	} else {
		return edit("fail to fetch pic")
	}
	return nil
}
