package bot

import (
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/auth"
	"github.com/Avimitin/go-bot/cmd/bot/internal/manage"
	"github.com/Avimitin/go-bot/cmd/bot/internal/tools"
	"github.com/Avimitin/go-bot/utils/modules/hardwareInfo"
	"github.com/Avimitin/go-bot/utils/modules/timer"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

type SendMethod func(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error)

var COMMAND = map[string]SendMethod{
	"start":      start,
	"help":       help,
	"ping":       ping,
	"sysinfo":    sysInfo,
	"authgroups": authGroups,
	"ver":        ver,
	"dump":       dump,
	"kick":       kick,
	"shutup":     shutUp,
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
		text = fmt.Sprintf("请输入正确的参数数量！只需要2个参数但是捕获到%d", length-1)
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
				text, _ = hardwareInfo.GetDiskUsage("\\")
			default:
				text, err = hardwareInfo.GetDiskUsage(args[2])
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

func authGroups(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (m tgbotapi.Message, err error) {
	if !auth.IsCreator(CREATOR, message.From.ID) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "您无权使用该命令。")
		m, err = bot.Send(msg)
		return m, err
	}
	var text string

	args := strings.Fields(message.Text)
	if length := len(args); length != 3 {
		text = fmt.Sprintf("请输入正确的参数数量！只需要2个参数但是捕获到%d", length-1)
	}

	switch args[1] {
	case "add":
		chatID, err := strconv.ParseInt(args[2], 20, 64)
		if err != nil {
			text = fmt.Sprintf("参数出错了！\n错误：%s", err)
		}
		newAuthGroups := append(cfg.Groups, chatID)
		cfg.Groups = newAuthGroups
		err = cfg.SaveConfig("F:\\go-bot\\cfg\\auth.yml")
		if err != nil {
			text = fmt.Sprintf("保存出错了！\n错误：%s", err)
		}
	default:
		text = "未知参数，您可以输入： /authgroups add 123 增加认证或 /authgroups del 123 删除群组"
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	return bot.Send(msg)
}

func ver(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("当前版本：%s", VERSION))
	return bot.Send(msg)
}

func dump(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	var text string
	if reply := message.ReplyToMessage; reply != nil {
		text = fmt.Sprintf(
			"<b>Reply To Message Info</b>\n"+
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

func kick(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	var msg tgbotapi.MessageConfig
	isAdmin, err := auth.IsAdmin(message.From.ID, bot, message.Chat)
	// acquire admins list
	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("在获取管理员列表时发生了一些错误：%v", err))
		return bot.Send(msg)
	}
	// check if the command caller is admin or not
	if !isAdmin {
		msg = tgbotapi.NewMessage(message.Chat.ID, "你没有权限，不许乱碰！")
		return bot.Send(msg)
	}
	// see if it's use for reply or just user id
	if reply := message.ReplyToMessage; reply != nil {
		args := strings.Fields(message.Text)
		switch len(args) {
		case 1:
			return manage.KickUser(bot, reply.From.ID, reply.Chat.ID, 0)
		case 2:
			if time, err := strconv.ParseInt(args[1], 10, 64); err == nil {
				return manage.KickUser(bot, reply.From.ID, reply.Chat.ID, time)
			}
			msg = tgbotapi.NewMessage(reply.Chat.ID, "参数错误！")
			return bot.Send(msg)
		default:
			msg = tgbotapi.NewMessage(reply.Chat.ID, fmt.Sprintf("请输入正确的参数数量: 需要0或1个参数但得到了 %v 个参数", len(args)-1))
			return bot.Send(msg)
		}
	}

	msg = tgbotapi.NewMessage(message.Chat.ID, "请回复一名用户的信息来踢出他")
	return bot.Send(msg)
}

func shutUp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	isAdmin, err := auth.IsAdmin(message.From.ID, bot, message.Chat)
	// acquire admins list
	if err != nil {
		return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("在获取管理员列表时发生了一些错误：%v", err))
	}

	// check if the command caller is reply to someone or not
	if reply := message.ReplyToMessage; reply != nil {
		// check permission
		if !isAdmin {
			return tools.SendTextMsg(bot, message.Chat.ID, "你没有权限，不许乱碰！")
		}
		// check arguments
		args := strings.Fields(message.Text)
		// with no args set 180s limits
		if len(args) == 1 {
			until, _ := timer.CalcTime(180, "s")
			return manage.ShutTheMouseUp(bot, message.Chat.ID, reply.From.ID, until, false)
		} else if len(args) == 2 {
			unit := args[1][len(args[1])-1:]
			addStr := args[1][:len(args[1])-1]

			add, err := strconv.ParseInt(addStr, 10, 64)
			if err != nil {
				return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数错误：%v", err))
			}

			until, err := timer.CalcTime(add, unit)
			if err != nil {
				return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数错误：%v", err))
			}
			return manage.ShutTheMouseUp(bot, message.Chat.ID, reply.From.ID, until, false)
		}

		return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数过多：需要0或1个参数但是得到了 %d 个参数", len(args)))
	}

	// If the message is use directly this will ban the sender randomly
	if !isAdmin {
		until := timer.AddRandTimeFromNow()
		return manage.ShutTheMouseUp(bot, message.Chat.ID, message.From.ID, until, false)
	}

	return tools.SendParseTextMsg(bot, message.Chat.ID,
		"Usage: Reply to a member and add a time for until date. "+
			"Support Seconds, minutes, hours, days... as time unit. Or you can just use `rand` as param to get random limit time."+
			"And if limit time is lower than 30s or longer than 366d it means this user is restricted forever.\n"+
			"Exp:\n"+
			"`/shutup 14h`\n"+
			"`/shutup rand`", "markdown")
}
