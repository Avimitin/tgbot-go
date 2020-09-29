package bot

import (
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/auth"
	"github.com/Avimitin/go-bot/utils/modules/hardwareInfo"
	"github.com/Avimitin/go-bot/utils/modules/timer"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

type SendMethod func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error)

var COMMAND = map[string]SendMethod{
	"start": start,
	"help":  help,
	"ping":  ping,
	"sysinfo": sysInfo,
	"authgroups": authGroups,
	"ver": ver,
	"dump": dump,
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
	if !auth.IsCreator(CREATOR, message.From.ID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "您无权使用该命令。")
		m, err = bot.Send(msg)
		return m, err
	}

	args := strings.Fields(message.Text)
	var text string

	if length := len(args); length != 3 {
		text = fmt.Sprintf("请输入正确的参数数量！只需要2个参数但是捕获到%d", length)
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
				text = "未知参数。你是不是想说： /sysinfo cpu percent ?"
			}
		case "disk":
			switch args[2] {
			case "stats":
				text , _ = hardwareInfo.GetDiskUsage("\\")
			default:
				text , err = hardwareInfo.GetDiskUsage(args[2])
				if err != nil {
					text += fmt.Sprintf("\n错误: %v", err)
				}
			}

		case "mem":
			switch args[2] {
			case "stats":
				text, _ = hardwareInfo.GetMemUsage()
			default:
				text = "未知参数，你是不是想说： /sysinfo mem stats ?"
			}

		default:
			text = "未知参数。"
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	m, err = bot.Send(msg)
	return m, err
}

func authGroups (bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error) {
	if !auth.IsCreator(CREATOR, message.From.ID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "您无权使用该命令。")
		m, err = bot.Send(msg)
		return m, err
	}
	var text string

	args := strings.Fields(message.Text)
	if length := len(args); length != 3 {
		text = fmt.Sprintf("请输入正确的参数数量！只需要2个参数但是捕获到%d", length)
	}

	switch args[1] {
	case "add":
		chatID, err := strconv.ParseInt(args[2], 20, 64)
		if err != nil { text = fmt.Sprintf("参数出错了！\n错误：%s", err)}
		newAuthGroups := append(cfg.Groups, chatID)
		cfg.Groups = newAuthGroups
		err = cfg.SaveConfig("F:\\go-bot\\cfg\\auth.yml")
		if err != nil { text = fmt.Sprintf("保存出错了！\n错误：%s", err)}
	default:
		text = "未知参数，您可以输入： /authgroups add 123 增加认证或 /authgroups del 123 删除群组"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	return bot.Send(msg)
}

func ver (bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("当前版本：%s", VERSION))
	return bot.Send(msg)
}

func dump (bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	var text string
	if reply := message.ReplyToMessage; reply != nil {
		text = fmt.Sprintf(
			"<b>Reply To Message Info</b>\n"+
			"|==<b>DATE</b>\n" +
			"|====%v\n" +
			"|==<b>CHAT</b>\n" +
			"|====<b>ID:</b> <code>%v</code>\n" +
			"|====<b>TYPE:</b> <code>%v</code>\n" +
			"|====<b>USERNAME:</b> <code>%v</code>\n" +
			"|==<b>USER</b>\n" +
			"|====<b>ID:</b> <code>%v</code>\n" +
			"|====<b>USERNAME:</b> <code>%v</code>\n" +
			"|====<b>NICKNAME:</b> <code>%v %v</code>\n" +
			"|====<b>LANGUAGE:</b> <code>%v</code>\n" +
			"|==<b>MSG</b>\n" +
			"|====<b>ID:</b> <code>%v</code>\n" +
			"|====<b>TEXT:</b> %v",
			timer.UnixToString(int64(reply.Date)),
			reply.Chat.ID, reply.Chat.Type, reply.Chat.UserName,
			reply.From.ID, reply.From.UserName, reply.From.FirstName, reply.From.LastName, reply.From.LanguageCode,
			reply.MessageID, reply.Text)
	} else {
		text = fmt.Sprintf(
			"<b>Info</b>\n"+
				"|==<b>DATE</b>\n"+
				"|====%v\n"+
				"|==<b>CHAT</b>\n"+
				"|====<b>ID:</b> <code>%v</code>\n"+
				"|====<b>TYPE:</b> <code>%v</code>\n"+
				"|====<b>USERNAME:</b> <code>%v</code>\n"+
				"|==<b>USER</b>\n"+
				"|====<b>ID:</b> <code>%v</code>\n"+
				"|====<b>USERNAME:</b> <code>%v</code>\n"+
				"|====<b>NICKNAME:</b> <code>%v %v</code>\n"+
				"|====<b>LANGUAGE:</b> <code>%v</code>\n"+
				"|==<b>MSG</b>\n"+
				"|====<b>ID:</b> <code>%v</code>\n"+
				"|====<b>TEXT:</b> %v",
			timer.UnixToString(int64(message.Date)),
			message.Chat.ID, message.Chat.Type, message.Chat.UserName,
			message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName, message.From.LanguageCode,
			message.MessageID, message.Text)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	return bot.Send(msg)
}