package bot

import (
	"errors"
	"log"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendT(text string, chatID int64) (bapi.Message, error) {
	msg := bapi.NewMessage(chatID, text)
	return bot.Send(msg)
}

func sendP(text string, chatID int64, parseMode string) (bapi.Message, error) {
	msg := bapi.NewMessage(chatID, text)
	msg.ParseMode = parseMode
	return bot.Send(msg)
}

func editT(newText string, chatID int64, msgID int) (bapi.Message, error) {
	msg := bapi.NewEditMessageText(chatID, msgID, newText)
	return bot.Send(msg)
}

func errF(where string, err error, more string) error {
	log.Printf("[%s]%v", where, err)
	return errors.New(more)
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
