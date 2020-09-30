package manage

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func KickUser(bot *tgbotapi.BotAPI, uid int, cid int64, restrictTime int64) (tgbotapi.Message, error) {
	kickUser := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: tgbotapi.ChatMemberConfig{
			ChatID:             cid,
			SuperGroupUsername: "",
			ChannelUsername:    "",
			UserID:             uid,
		},
		UntilDate: restrictTime,
	}

	_, err := bot.KickChatMember(kickUser)
	var msg tgbotapi.MessageConfig

	if err != nil {
		msg = tgbotapi.NewMessage(cid, fmt.Sprintf("发生错误：%v\n\nHint: 您是不是忘记给Bot设置管理员权限了？", err))
	} else {
		msg = tgbotapi.NewMessage(cid, fmt.Sprintf("成功踢出 %v", uid))
	}

	return bot.Send(msg)
}
