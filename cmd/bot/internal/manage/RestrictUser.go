package manage

import (
	"fmt"
	"github.com/Avimitin/go-bot/utils/modules/timer"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

func ShutTheMouseUp(bot *tgbot.BotAPI, cid int64, uid int, until int64, canSendMessages bool) (tgbot.Message, error) {
	user := tgbot.RestrictChatMemberConfig{
		ChatMemberConfig: tgbot.ChatMemberConfig{
			ChatID: cid,
			UserID: uid,
		},
		CanSendMessages: &canSendMessages,
		CanSendMediaMessages: &canSendMessages,
		CanAddWebPagePreviews: &canSendMessages,
		CanSendOtherMessages: &canSendMessages,
		UntilDate: until,
	}
	_, err := bot.RestrictChatMember(user)

	var msg tgbot.MessageConfig
	if err != nil {
		response := map[string]string{
			"Bad Request: user is an administrator of the chat": "对面是管理员！我没法让他闭嘴qwq",
			"Bad Request: can't remove chat owner": "啊不会吧，不会吧，不会真的有人觉得我权限比群主大吧",
			"Bad Request: not enough rights to restrict/unrestrict chat member": "拜托诶你不给我权限我怎么帮你禁言啦!",
		}
		if responseMsg, ok := response[err.Error()]; ok {
			msg = tgbot.NewMessage(cid, fmt.Sprintf("发生错误啦: %v", responseMsg))
		} else {
			msg = tgbot.NewMessage(cid, fmt.Sprintf("发生错误啦: %v", err))
		}
		return bot.Send(msg)
	}
	msg = tgbot.NewMessage(cid,
		fmt.Sprintf("Restrict User: %v for sending any message until %v", cid, timer.UnixToString(until)))
	return bot.Send(msg)
}