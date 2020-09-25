package main

import (
	"fmt"
	"github.com/Avimitin/go-bot/utils/modules/hardwareInfo"
	"github.com/Avimitin/go-bot/utils/modules/timer"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

type SendMethod func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error)

var COMMAND = map[string]SendMethod{
	"start": start,
	"help":  help,
	"ping":  ping,
	"sysinfo": sysInfo,
}

func start(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Here is start.")
	m, err = bot.Send(msg)
	return m, err
}

func help(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error) {
	text := "<b>Author</b>:\n\n" +
		"@SaiToAsuKa_kksk\n\n" +
		"<b>Sponsor</b>:\n\n" +
		"暂时还没有赞助，假如你对我的 bot 感兴趣非常欢迎私聊我\n\n" +
		"<b>Guide</b>:\n\n" +
		"大部分功能为管理员专属，目前普通用户可用 /post 功能投稿自己感兴趣的内容\n\n" +
		"<b>Affiliate</b>:\n\n" +
		"我真的真的续不起服务器了QWQ，快点击 AFF 帮助我吧\n\n" +
		"<a href=\"https://www.vultr.com/?ref=8527101-6G\">【VULTR1】</a>\n\n" +
		"<a href=\"https://www.vultr.com/?ref=8527098\">【VULTR2】</a>\n\n" +
		"<b>Group</b>:\n\n" +
		"NSFW 中文水群: @ghs_chat\n\n" +
		"NSFW 本子推荐频道: @hcomic\n\n" +
		"BOT 更新频道: @avimitinbot\n\n" +
		"BOT 反馈群组: @avimitin_studio"

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	m, err = bot.Send(msg)
	return m, err
}

func ping(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error) {
	t := timer.NewTimer()
	msg := tgbotapi.NewMessage(message.Chat.ID, "Pong!(Counting time now...)")
	response, err := bot.Send(msg)
	if err == nil {
		duration := t.StopCounting() / 1000000000.000
		newMsg := tgbotapi.NewEditMessageText(response.Chat.ID, response.MessageID, fmt.Sprintf("bot 到您数据中心的双向延迟大约为 %.5f s", duration))
		response, err = bot.Send(newMsg)
	}
	return response, err
}


func sysInfo(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error) {
	args := strings.Fields(message.Text)
	var text string

	if len(args) != 3 {
		text = "请输入正确的参数数量！"
	} else {
		switch args[1] {
		case "cpu":
			switch args[2] {
			case "model":
				text, _ = hardwareInfo.GetCpuModel()
			case "percent":
				text, _ = hardwareInfo.GetCpuPercent()
			case "load":
				text, _ = hardwareInfo.GetCpuLoad()
			default:
				text = "未知参数。"
			}
		case "disk":
			switch args[2] {
			case "stats":
				text , _ = hardwareInfo.GetDiskUsage("\\")
			default:
				text , err = hardwareInfo.GetDiskUsage(args[2])
				if err != nil {
					text += fmt.Sprintf("\nDescriptions: %v", err)
				}
			}
		default:
			text = "未知参数。"
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	m, err = bot.Send(msg)
	return m, err
}