package bot

import (
	"fmt"
	"github.com/Avimitin/go-bot/internal/bot/auth"
	"github.com/Avimitin/go-bot/internal/bot/manage"
	"github.com/Avimitin/go-bot/internal/database"
	"github.com/Avimitin/go-bot/internal/utils/hardwareInfo"
	"github.com/Avimitin/go-bot/internal/utils/timer"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	noReply = true
	neReply = false // neReply means need reply
)

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
	"keydel":     KeyDel,
}

// Use for no reply case
func sendText(c *C, cid int64, text string) {
	c.Send(NewSendPKG(SendT(cid, text), noReply))
}

// Use for no reply case
func sendParse(c *C, cid int64, text string, parseM string) {
	c.Send(NewSendPKG(SendP(cid, text, parseM), noReply))
}

func start(message *M, ctx *C) {
	sendText(ctx, message.Chat.ID, "Welcome.")
}

func help(message *M, ctx *C) {
	text := `<b>Author</b>:
@SaiToAsuKa_kksk
<b>Affiliate</b>:
我真的真的续不起服务器了QWQ，快点击 AFF 帮助我吧
<a href=\"https://www.vultr.com/?ref=8527101-6G\">【VULTR1】</a>
<a href=\"https://www.vultr.com/?ref=8527098\">【VULTR2】</a>
<b>Group</b>:
NSFW 本子推荐频道: @hcomic
BOT 更新频道: @avimitinbot
`
	sendParse(ctx, message.Chat.ID, text, "HTML")
}

func ping(message *M, ctx *C) {
	t := timer.NewTimer()
	pkg := NewSendPKG(SendT(message.Chat.ID, "Pong!(Counting time now...)"), neReply)
	ctx.Send(pkg)
	go func() {
		select {
		case <-time.After(ctx.timeOut):
			log.Println("[ERR]Timeout exceed")
			return
		case response := <-pkg.resp:
			duration := t.StopCounting() / 1000000000.000
			pkg.msg = tgbotapi.NewEditMessageText(response.Chat.ID, response.MessageID, fmt.Sprintf("bot 到您数据中心的双向延迟大约为 %.5f s", duration))
			pkg.noReply = noReply
			ctx.Send(pkg)
		}
	}()
}

func sysInfo(message *M, ctx *C) {
	if !auth.IsMe(CREATOR, message.From.ID) {
		sendText(ctx, message.Chat.ID, "You don't have permission.")
		return
	}

	args := strings.Fields(message.Text)
	if len(args) == 1 {
		sendText(ctx, message.Chat.ID, "Miss argument.\nSee '/sysinfo help'")
		return
	}

	var text string
	var err error
	switch args[1] {
	case "cpu":
		if len(args) > 3 {
			switch args[2] {
			case "model":
				text, err = hardwareInfo.GetCpuModel()
			case "percent":
				text, err = hardwareInfo.GetCpuPercent()
			case "load":
				text, err = hardwareInfo.GetCpuLoad()
			case "help":
				text = "model - get cpu model\npercent - get cpu use percent\nload - get cpu load"
			default:
				text = args[2] + " is not a cpu command. See '/sysinfo cpu help'"
			}
		} else {
			text = "Nothing specific.\nTry out /sysinfo cpu help"
		}
	case "disk":
		if len(args) > 3 {
			switch args[2] {
			case "stats":
				text, err = hardwareInfo.GetDiskUsage("/")
			case "help":
				text = "<path> - get disk usage of giving path\nstats - get usage of whole disk"
			default:
				text, err = hardwareInfo.GetDiskUsage(args[2])
			}
		} else {
			text = "Nothing specific.\nTry out /sysinfo disk help"
		}

	case "mem":
		if len(args) == 3 {
			switch args[2] {
			case "stats":
				text, _ = hardwareInfo.GetMemUsage()
			default:
				text = args[2] + " is not a memory command. See '/sysinfo mem help'"
			}
		} else {
			text = "Nothing specific.\nTry out /sysinfo mem help"
		}
	case "help":
		text = "Usage: /sysinfo <cpu/disk/mem> <args>"
	default:
		text = "Unknown argument. See '/sysinfo help' for help message"
	}
	// -----------
	if err != nil {
		sendText(ctx, message.Chat.ID, "Error happen when handle your request: "+err.Error())
	} else {
		sendText(ctx, message.Chat.ID, text)
	}
}

func authGroups(message *M, ctx *C) {
	if !auth.IsMe(CREATOR, message.From.ID) {
		sendText(ctx, message.Chat.ID, "不许乱碰！")
	}

	args := strings.Fields(message.Text)
	if length := len(args); length != 3 {
		sendText(ctx, message.Chat.ID, fmt.Sprintf("请输入正确的参数数量！只需要2个参数但是捕获到%d", length-1))
	}

	var text string
	switch args[1] {
	// Add authorized groups ID.
	case "add":
		// Get supergroup's username.
		chatUserName := args[2]

		// Get specific chat username.
		targetChat, err := ctx.Bot().GetChat(tgbotapi.ChatConfig{SuperGroupUsername: chatUserName})
		if err != nil {
			sendText(ctx, message.Chat.ID, fmt.Sprintf("获取群组信息时出现错误\n错误信息：%v", err))
		}

		// Store groups information.
		err = database.AddGroups(ctx.DB(), targetChat.ID, targetChat.UserName)
		if err != nil {
			sendText(ctx, message.Chat.ID, fmt.Sprintf("保存出错了！\n错误：%s", err))
		}
		ctx.SetGroup(targetChat.ID)
		sendText(ctx, message.Chat.ID, "保存认证群组成功")

	// Delete authorized group's record.
	case "del":
		// Convert string arguments to int64.
		chatID, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			sendText(ctx, message.Chat.ID, fmt.Sprintf("参数出错了！\n错误：%s", err))
		}
		ctx.DelGroup(chatID)
		// Delete chat record in database
		err = database.DeleteGroups(ctx.DB(), chatID)
		if err != nil {
			sendText(ctx, message.Chat.ID, fmt.Sprintf("删除出错了！\n错误：%s", err))
		}
		sendText(ctx, message.Chat.ID, "成功删除！")
	// list all groups
	case "list":
		if args[2] == "db" {
			groups, err := database.SearchGroups(ctx.DB())
			if err != nil {
				sendText(ctx, message.Chat.ID, fmt.Sprintf("获取群组信息时发生了一些错误。"))
			}
			for i, group := range *groups {
				text += fmt.Sprintf("%d. GID: %v GNAME: %v\n", i, group.GroupID, group.GroupUsername)
			}
		} else if args[2] == "mem" {
			i := 1
			for key := range ctx.Groups() {
				text += fmt.Sprintf("%d. GID: %v \n", i, key)
				i++
			}
		} else {
			text = "未知参数，你可以输入 /authgroups list mem 或者 db 查询内存或数据库内群组信息。"
		}
	default:
		text = "未知参数，您可以输入： /authgroups add 123 增加认证或 /authgroups del 123 删除群组"
	}

	sendText(ctx, message.Chat.ID, text)
}

func ver(message *M, ctx *C) {
	sendText(ctx, message.Chat.ID, fmt.Sprintf("当前版本：%s", VERSION))
}

func dump(message *M, ctx *C) {
	var text = "<b>Reply To Message Info</b>\n" +
		"<b>DATE</b>\n" +
		"%v\n" +
		"=== === ===" +
		"<b>CHAT</b>\n" +
		"<b>ID:</b> <code>%v</code>\n" +
		"<b>TYPE:</b> <code>%v</code>\n" +
		"<b>USERNAME:</b> <code>%v</code>\n" +
		"<b>USER</b>\n" +
		"<b>ID:</b> <code>%v</code>\n" +
		"<b>USERNAME:</b> <code>%v</code>\n" +
		"<b>NICKNAME:</b> <code>%v %v</code>\n" +
		"<b>LANGUAGE:</b> <code>%v</code>\n" +
		"=== === ===" +
		"<b>MSG</b>\n" +
		"<b>ID:</b> <code>%v</code>\n" +
		"<b>TEXT:</b> %v"
	if reply := message.ReplyToMessage; reply != nil {
		text = fmt.Sprintf(text,
			timer.UnixToString(int64(reply.Date)),
			reply.Chat.ID, reply.Chat.Type, reply.Chat.UserName,
			reply.From.ID, reply.From.UserName, reply.From.FirstName, reply.From.LastName, reply.From.LanguageCode,
			reply.MessageID, reply.Text)
	} else {
		text = fmt.Sprintf(text,
			timer.UnixToString(int64(message.Date)),
			message.Chat.ID, message.Chat.Type, message.Chat.UserName,
			message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName, message.From.LanguageCode,
			message.MessageID, message.Text)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	pkg := NewSendPKG(msg, noReply)
	ctx.Send(pkg)
}

func kick(message *M, ctx *C) {
	isAdmin, err := auth.IsAdmin(ctx.Bot(), message.From.ID, message.Chat)
	// acquire admins list
	if err != nil {
		sendText(ctx, message.Chat.ID, fmt.Sprintf("在获取管理员列表时发生了一些错误：%v", err))
		return
	}
	// check if the command caller is admin or not
	if !isAdmin {
		sendText(ctx, message.Chat.ID, "你没有权限，不许乱碰！")
		return
	}
	// Get kick until date
	if reply := message.ReplyToMessage; reply != nil {
		args := strings.Fields(message.Text)
		switch len(args) {
		case 1:
			pkg := NewSendPKG(manage.KickUser(ctx.Bot(), reply.From.ID, reply.Chat.ID, 0), noReply)
			ctx.Send(pkg)
			return
		case 2:
			if time, err := strconv.ParseInt(args[1], 10, 64); err == nil {
				pkg := NewSendPKG(manage.KickUser(ctx.Bot(), reply.From.ID, reply.Chat.ID, time), noReply)
				ctx.Send(pkg)
				return
			}
			sendText(ctx, reply.Chat.ID, "参数错误！")
		default:
			sendText(ctx, message.Chat.ID, fmt.Sprintf("请输入正确的参数数量: 需要0或1个参数但得到了 %v 个参数", len(args)-1))
		}
	}

	sendText(ctx, message.Chat.ID, "请回复一名用户的信息来踢出他")
}

func shutUp(message *M, ctx *C) {
	isAdmin, err := auth.IsAdmin(ctx.Bot(), message.From.ID, message.Chat)
	// acquire admins list
	if err != nil {
		sendText(ctx, message.Chat.ID, fmt.Sprintf("在获取管理员列表时发生了一些错误：%v", err))
		return
	}

	// If the command is not used by admin, bot will ban the sender in random time.
	if !isAdmin {
		until := timer.AddRandTimeFromNow()
		msg := manage.ShutTheMouseUp(ctx.Bot(), message.Chat.ID, message.From.ID, until, false)
		pkg := NewSendPKG(msg, noReply)
		ctx.Send(pkg)
		return
	}

	// check if the command caller is reply to someone or not
	if reply := message.ReplyToMessage; reply != nil {
		// check arguments
		args := strings.Fields(message.Text)
		// with no args set 180s limits
		if len(args) == 1 {
			until, _ := timer.CalcTime(180, "s")
			msg := manage.ShutTheMouseUp(ctx.Bot(), message.Chat.ID, reply.From.ID, until, false)
			pkg := NewSendPKG(msg, noReply)
			ctx.Send(pkg)
			return
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
					sendText(ctx, message.Chat.ID, fmt.Sprintf("参数错误：%v", err))
					return
				}
				// add time from now.unix()
				until, err = timer.CalcTime(add, unit)
				if err != nil {
					sendText(ctx, message.Chat.ID, fmt.Sprintf("参数错误：%v", err))
					return
				}
			}
			msg := manage.ShutTheMouseUp(ctx.Bot(), message.Chat.ID, reply.From.ID, until, false)
			pkg := NewSendPKG(msg, noReply)
			ctx.Send(pkg)
			return
		}

		sendText(ctx, message.Chat.ID, fmt.Sprintf("参数过多：需要0或1个参数但是得到了 %d 个参数", len(args)))
	}

	sendParse(ctx, message.Chat.ID,
		"Usage: *Reply* to a member and add a time argument as until date. "+
			"Support Seconds, minutes, hours, days... as time unit. Or you can just use `rand` as param to get random limit time."+
			"And if limit time is lower than 30s or longer than 366d it means this restriction will until forever.\n"+
			"Exp:\n"+
			"`/shutup 14h`\n"+
			"`/shutup rand`", "markdown")
}

func unShutUp(message *M, ctx *C) {
	isAdmin, err := auth.IsAdmin(ctx.Bot(), message.From.ID, message.Chat)
	// acquire admins list
	if err != nil {
		sendText(ctx, message.Chat.ID, fmt.Sprintf("在获取管理员列表时发生了一些错误：%v", err))
		return
	}
	// check if the command caller is admin or not
	if !isAdmin {
		sendText(ctx, message.Chat.ID, "你没有权限，不许乱碰！")
		return
	}

	if reply := message.ReplyToMessage; reply != nil {
		msg := manage.OpenMouse(ctx.Bot(), message.Chat.ID, reply.From.ID, true)
		ctx.Send(NewSendPKG(msg, noReply))
		return
	}

	args := strings.Split(message.Text, "/unshutup ")
	if len(args) == 1 {
		sendText(ctx, message.Chat.ID, "请回复一个用户的信息或者输入他的UID来解封")
	}
	if len(args) == 2 {
		uid, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			sendText(ctx, message.Chat.ID, fmt.Sprintf("参数错误！需要用户的 USERID"))
		}
		msg := manage.OpenMouse(ctx.Bot(), message.Chat.ID, int(uid), true)
		ctx.Send(NewSendPKG(msg, noReply))
		return
	}
	sendText(ctx, message.Chat.ID, "参数过多！")
}

func keyAdd(message *M, ctx *C) {
	if !(message.From.ID == CREATOR) {
		sendText(ctx, message.Chat.ID, "不许乱碰！")
		return
	}

	order := strings.Split(message.Text, "/keyadd ")
	if len(order) != 2 {
		sendParse(ctx, message.Chat.ID, "使用 `/keyadd key=reply` 增加关键词。", "markdown")
		return
	}

	args := strings.SplitN(order[1], "=", 2)
	keyword := args[0]
	reply := args[1]
	err := Set(keyword, reply, ctx)
	if err != nil {
		sendText(ctx, message.Chat.ID, "Error happen:"+err.Error())
		return
	}

	sendParse(ctx, message.Chat.ID,
		fmt.Sprintf("我已经学会啦！当你说 *%s* 的时候， 我会回复 *%s*。", keyword, reply), "markdown")
}

func KeyList(message *M, ctx *C) {
	if message.From.ID != CREATOR {
		sendText(ctx, message.Chat.ID, "憋搁这乱碰")
	}
	sendText(ctx, message.Chat.ID, ListKeywordAndReply(ctx.KeywordReplies()))
}

func KeyDel(message *M, ctx *C) {
	if message.From.ID != CREATOR {
		sendText(ctx, message.Chat.ID, "憋搁这乱碰")
	}
	order := strings.SplitN(message.Text, "/keydel ", 2)
	args := strings.Fields(order[1])
	var failure int
	// Delete keyword in database
	for _, arg := range args {
		err := Del(arg, ctx)
		if err != nil {
			failure++
		}
	}
	if failure == 0 {
		sendText(ctx, message.Chat.ID, "给定关键词全部删除成功")
		return
	}
	sendText(ctx, message.Chat.ID, fmt.Sprintf("有 %d 个关键词删除失败。", failure))
}
