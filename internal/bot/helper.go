package bot

import (
	"fmt"
	"log"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendT(text string, chatID int64) bapi.Message {
	msg := bapi.NewMessage(chatID, text)
	resp, err := bot.Send(msg)
	if err != nil {
		log.Printf("send %s...: \n%v", msg.Text[0:15], err)
	}
	return resp
}

func sendP(text string, chatID int64, parseMode string) bapi.Message {
	msg := bapi.NewMessage(chatID, text)
	msg.ParseMode = parseMode
	resp, err := bot.Send(msg)
	if err != nil {
		log.Printf("send %s...: \n%v", msg.Text[0:15], err)
	}
	return resp
}

func editT(newText string, chatID int64, msgID int) (bapi.Message, error) {
	msg := bapi.NewEditMessageText(chatID, msgID, newText)
	return bot.Send(msg)
}

func editP(newText string, chatID int64, msgID int, parseMode string) (bapi.Message, error) {
	msg := bapi.NewEditMessageText(chatID, msgID, newText)
	msg.ParseMode = parseMode
	return bot.Send(msg)
}

func errF(where string, err error, more string) error {
	return fmt.Errorf("%s: %s: %w", where, more, err)
}

func isAdmin(userID int, chat *bapi.Chat) (bool, error) {
	admins, err := bot.GetChatAdministrators(chat.ChatConfig())
	if err != nil {
		return false, errF("isAdmin", err, "fail to fetch chat admins")
	}
	for _, admin := range admins {
		if admin.User.ID == userID {
			return true, nil
		}
	}
	return false, nil
}

func kickUser(userID int, chatID int64, untilDate int64) error {
	userToKick := bapi.KickChatMemberConfig{
		UntilDate: untilDate,
		ChatMemberConfig: bapi.ChatMemberConfig{
			ChatID: chatID,
			UserID: userID,
		},
	}
	resp, err := bot.KickChatMember(userToKick)
	if err != nil {
		return errF("kickUser", err, "fail to send kick request: "+resp.Description)
	}
	return nil
}

func editUserPermissions(user int, chat int64, untilDate int64, notBan bool) error {
	resp, err := bot.RestrictChatMember(bapi.RestrictChatMemberConfig{
		ChatMemberConfig: bapi.ChatMemberConfig{
			ChatID: chat,
			UserID: user,
		},
		CanSendMessages:       &notBan,
		CanSendMediaMessages:  &notBan,
		CanSendOtherMessages:  &notBan,
		CanAddWebPagePreviews: &notBan,
		UntilDate:             untilDate,
	})

	if err != nil {
		return errF("limitUser", err, "fail to restrict chat member: "+resp.Description)
	}
	return nil
}
