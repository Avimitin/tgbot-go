package manage

import (
	"fmt"
	"github.com/Avimitin/go-bot/internal/utils/timer"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func ShutTheMouseUp(bot *tgbot.BotAPI, cid int64, uid int, until int64, canSendMessages bool) *tgbot.MessageConfig {
	user := tgbot.RestrictChatMemberConfig{
		ChatMemberConfig: tgbot.ChatMemberConfig{
			ChatID: cid,
			UserID: uid,
		},
		CanSendMessages:       &canSendMessages,
		CanSendMediaMessages:  &canSendMessages,
		CanAddWebPagePreviews: &canSendMessages,
		CanSendOtherMessages:  &canSendMessages,
		UntilDate:             until,
	}
	_, err := bot.RestrictChatMember(user)
	// handle error
	var text string
	if err != nil {
		// remake response
		response := map[string]string{
			"Bad Request: user is an administrator of the chat":                 "对面是管理员！我没法让他闭嘴qwq",
			"Bad Request: can't remove chat owner":                              "啊不会吧，不会吧，不会真的有人觉得我权限比群主大吧",
			"Bad Request: not enough rights to restrict/unrestrict chat member": "拜托诶你不给我权限我怎么帮你禁言啦!",
		}

		if responseMsg, ok := response[err.Error()]; ok {
			text = responseMsg
		} else {
			text = fmt.Sprintf("发生错误啦: %v", err)
		}
		msg := tgbot.NewMessage(cid, text)
		return &msg
	}

	if until-time.Now().Unix() < 31 || until-time.Now().Unix() > 885427200 {
		text = fmt.Sprintf("用户: %v 被永久禁言", cid)
	} else {
		text = fmt.Sprintf("用户: %v 直到 %v 都不准说话", uid, timer.UnixToString(until))
	}

	return &tgbot.MessageConfig{
		BaseChat:              tgbot.BaseChat{ChatID: cid, ReplyToMessageID: 0},
		Text:                  text,
		DisableWebPagePreview: false,
	}
}

func OpenMouse(bot *tgbot.BotAPI, cid int64, uid int, canSendMessages bool) *tgbot.MessageConfig {
	user := tgbot.RestrictChatMemberConfig{
		ChatMemberConfig: tgbot.ChatMemberConfig{
			ChatID: cid,
			UserID: uid,
		},
		CanSendMessages:       &canSendMessages,
		CanSendMediaMessages:  &canSendMessages,
		CanAddWebPagePreviews: &canSendMessages,
		CanSendOtherMessages:  &canSendMessages,
	}
	_, err := bot.RestrictChatMember(user)
	if err != nil {
		// remake response
		response := map[string]string{
			"Bad Request: not enough rights to restrict/unrestrict chat member": "拜托诶你不给我权限我怎么帮你解封人啦!",
			"Bad Request: user not found":                                       "没在这个群里找到这个人",
		}

		var text string
		if responseMsg, ok := response[err.Error()]; ok {
			text = responseMsg
		} else {
			text = fmt.Sprintf("发生错误啦: %v", err)
		}
		return &tgbot.MessageConfig{
			BaseChat: tgbot.BaseChat{ChatID: cid, ReplyToMessageID: 0},
			Text:     text,
		}

	}
	return &tgbot.MessageConfig{
		BaseChat: tgbot.BaseChat{ChatID: cid, ReplyToMessageID: 0},
		Text:     fmt.Sprintf("%v 已经被解封啦", uid),
	}
}
