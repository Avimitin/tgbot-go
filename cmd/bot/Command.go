package bot

import (
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/auth"
	"github.com/Avimitin/go-bot/cmd/bot/internal/database"
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
	"unshutup":   unShutUp,
	"keyadd":     keyAdd,
	"keylist":    KeyList,
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
		return tools.SendTextMsg(bot, message.Chat.ID, "不许乱碰！")
	}

	args := strings.Fields(message.Text)
	if length := len(args); length != 3 {
		return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("请输入正确的参数数量！只需要2个参数但是捕获到%d", length-1))
	}

	var text string
	switch args[1] {
	// Add authorized groups ID.
	case "add":
		// Get supergroup's username.
		chatUserName := args[2]

		// Get specific chat username.
		targetChat, err := bot.GetChat(tgbotapi.ChatConfig{SuperGroupUsername: chatUserName})
		if err != nil {
			return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("获取群组信息时出现错误\n错误信息：%v", err))
		}

		// Store groups information.
		err = database.AddGroups(DB, targetChat.ID, targetChat.UserName)
		if err != nil {
			return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("保存出错了！\n错误：%s", err))
		}

		// If all things behind has done, store groups id into memory.
		// This cycle is for appending value with order.
		for i, group := range cfg.Groups {
			if targetChat.ID < group {
				current := append([]int64{targetChat.ID}, cfg.Groups[i:]...)
				cfg.Groups = append(cfg.Groups[:i], current...)
				return tools.SendTextMsg(bot, message.Chat.ID, "保存认证群组成功")
			}
		}

		// If targetChat's ID is the biggest just insert it.
		cfg.Groups = append(cfg.Groups, targetChat.ID)
		return tools.SendTextMsg(bot, message.Chat.ID, "保存认证群组成功")

	// Delete authorized group's record.
	case "del":
		// Convert string arguments to int64.
		i, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数出错了！\n错误：%s", err))
		}

		// Search group is exist or not, if exist, delete it from memory and database.
		if int(i) > len(cfg.Groups) {
			return tools.SendTextMsg(bot, message.Chat.ID, "找不到指定序号的群组。")
		}
		chatID := cfg.Groups[i]
		cfg.Groups = append(cfg.Groups[:i], cfg.Groups[i+1:]...)

		// Delete chat record in database
		err = database.DeleteGroups(DB, chatID)
		if err != nil {
			return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("删除出错了！\n错误：%s", err))
		}
		return tools.SendTextMsg(bot, message.Chat.ID, "成功删除！")
	// list all groups
	case "list":
		if args[2] == "db" {
			groups, err := database.SearchGroups(DB)
			if err != nil {
				return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("获取群组信息时发生了一些错误。"))
			}
			for i, group := range groups {
				text += fmt.Sprintf("%d. GID: %v GNAME: %v\n", i, group.GroupID, group.GroupUsername)
			}
		} else if args[2] == "mem" {
			for i, group := range cfg.Groups {
				text += fmt.Sprintf("%d. GID: %v \n", i, group)
			}
		} else {
			text = "未知参数，你可以输入 /authgroups list mem 或者 db 查询内存或数据库内群组信息。"
		}
	default:
		text = "未知参数，您可以输入： /authgroups add 123 增加认证或 /authgroups del 123 删除群组"
	}

	return tools.SendTextMsg(bot, message.Chat.ID, text)
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
	isAdmin, err := auth.IsAdmin(bot, message.From.ID, message.Chat)
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
	isAdmin, err := auth.IsAdmin(bot, message.From.ID, message.Chat)
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
			// init until time
			var until int64
			if args[1] == "rand" {
				until = timer.AddRandTimeFromNow()
			} else {
				unit := args[1][len(args[1])-1:]
				addStr := args[1][:len(args[1])-1]
				// convert string to int64
				add, err := strconv.ParseInt(addStr, 10, 64)
				if err != nil {
					return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数错误：%v", err))
				}
				// add time from now.unix()
				until, err = timer.CalcTime(add, unit)
				if err != nil {
					return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数错误：%v", err))
				}
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
		"Usage: *Reply* to a member and add a time for until date. "+
			"Support Seconds, minutes, hours, days... as time unit. Or you can just use `rand` as param to get random limit time."+
			"And if limit time is lower than 30s or longer than 366d it means this user is restricted forever.\n"+
			"Exp:\n"+
			"`/shutup 14h`\n"+
			"`/shutup rand`", "markdown")
}

func unShutUp(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	if reply := message.ReplyToMessage; reply != nil {
		return manage.OpenMouse(bot, message.Chat.ID, reply.From.ID, true)
	}
	args := strings.Split(message.Text, "/unshutup ")
	if len(args) == 1 {
		return tools.SendTextMsg(bot, message.Chat.ID, "请回复一个用户的信息或者输入他的UID来解封")
	}
	if len(args) == 2 {
		uid, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("参数错误！需要用户的 USERID"))
		}
		return manage.OpenMouse(bot, message.Chat.ID, int(uid), true)
	}
	return tools.SendTextMsg(bot, message.Chat.ID, "参数过多！")
}

func keyAdd(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	isAdmin, err := auth.IsAdmin(bot, message.From.ID, message.Chat)
	if err != nil {
		return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("获取管理员列表时发生错误：%v", err))
	}
	if !isAdmin {
		return tools.SendTextMsg(bot, message.Chat.ID, "不许乱碰！")
	}
	order := strings.Split(message.Text, "/keyadd ")
	if len(order) != 2 {
		return tools.SendParseTextMsg(bot, message.Chat.ID, "使用 `/keyadd key=reply` 增加关键词。", "markdown")
	}
	args := strings.SplitN(order[1], "=", 2)
	keyword := args[0]
	reply := args[1]
	c := make(chan int, 1)

	// Make new goroutine for add record.
	// First add into database
	go func(c chan int) {
		kid, err := SetKeywordIntoDB(keyword, reply)
		if err != nil {
			c <- -1
		}
		c <- kid
	}(c)

	// Then add record into memory.
	kid := <-c
	go SetKeywordIntoCFG(kid, keyword, reply)

	return tools.SendParseTextMsg(bot, message.Chat.ID,
		fmt.Sprintf("我已经学会啦！当你说 *%s* 的时候， 我会回复 *%s*。", keyword, reply), "markdown")
}

func KeyList(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	return tools.SendTextMsg(bot, message.Chat.ID, ListKeywordAndReply())
}

func KeyDel(bot *tgbotapi.BotAPI, message *tgbotapi.Message) (tgbotapi.Message, error) {
	order := strings.SplitN(message.Text, "/keydel ", 2)
	args := strings.Fields(order[1])
	successCount := 0
	c := make(chan string, 2)
	d := make(chan bool, 1)
	// Delete keyword in database
	go func(c chan string, d chan bool) {
		kw := <-c
		err := DelKeyword(kw)
		if err == nil {
			d <- true
		}
	}(c, d)

	// If success add into successCount
	go func(d chan bool) {
		success := <-d
		if success {
			successCount += 1
		}
	}(d)

	// Passing value to methods.
	for _, arg := range args {
		c <- arg
	}
	return tools.SendTextMsg(bot, message.Chat.ID, fmt.Sprintf("成功删除 %d 个关键词。", successCount))
}
