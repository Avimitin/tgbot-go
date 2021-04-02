package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	_, err := bot.KickChatMember(userToKick)
	if err != nil {
		return fmt.Errorf("kick %d:%v", userID, err)
	}
	return nil
}

func editUserPermissions(user int, chat int64, untilDate int64, notBan bool) error {
	_, err := bot.RestrictChatMember(bapi.RestrictChatMemberConfig{
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
		return fmt.Errorf("restrict %d: %v", user, err)
	}
	return nil
}

func leaveGroup(chat *bapi.Chat) error {
	_, err := bot.LeaveChat(chat.ChatConfig())
	if err != nil {
		log.Printf("leave group [%s](%d) failed: %v", chat.FirstName, chat.ID, err)
		return fmt.Errorf("leave %d: %v", chat.ID, err)
	}
	return nil
}

func safeExit() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("exit and save config", setting.Update())
		os.Exit(1)
	}()
}

type msgHandlerFunc func(*bapi.Message) error
type infoMapType map[string]string // infoMapType contain message that store by key to value

func (imt infoMapType) get(key string) string {
	return imt[key]
}

type registration struct {
	registerFunc map[int]msgHandlerFunc
	registerInfo map[int]infoMapType
	mu           sync.RWMutex
}

func NewRegistration() *registration {
	return &registration{
		registerFunc: make(map[int]msgHandlerFunc),
		registerInfo: make(map[int]infoMapType),
	}
}

func (r *registration) getFn(u int) msgHandlerFunc {
	r.mu.RLock()
	fn := r.registerFunc[u]
	r.mu.RUnlock()
	return fn
}

func (r *registration) getInfo(u int) infoMapType {
	r.mu.RLock()
	info := r.registerInfo[u]
	r.mu.RUnlock()
	return info
}

func (r *registration) registerNextFunc(m *bapi.Message, fn msgHandlerFunc, info infoMapType) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registerFunc[m.From.ID] = fn
	r.registerInfo[m.From.ID] = info
}

func (r *registration) clear(m *bapi.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.registerFunc, m.From.ID)
	delete(r.registerInfo, m.From.ID)
}

// nuclear clear all the data, use with attention!
func (r *registration) nuclear(m *bapi.Message) {
	r = nil
	r = new(registration)
}
