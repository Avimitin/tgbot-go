package main

import (
	"errors"
	"fmt"
	"log"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot      *bapi.BotAPI
	setting  SettingsGetter
	registry = NewRegistration()
)

func Run(s SettingsGetter) error {
	if s == nil {
		return errors.New("no setting")
	}
	setting = s

	botToken := setting.Secret().Get("bot_token")
	if botToken == "" {
		return errors.New("no bot token")
	}
	var err error
	bot, err = bapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("connect to api server: %v", err)
	}
	bot.Debug = false
	log.Printf("Successfully establish connection to bot: %s", bot.Self.UserName)

	updateChanConfiguration := bapi.NewUpdate(0)
	updateChanConfiguration.Timeout = 15

	updates, err := bot.GetUpdatesChan(updateChanConfiguration)
	if err != nil {
		return fmt.Errorf("get update chan:%v", err)
	}
	safeExit()

	for update := range updates {
		if update.Message != nil {
			go func() {
				e := messageHandler(update.Message)
				if e != nil {
					log.Println(e)
				}
			}()
		}
	}
	return nil
}

func messageHandler(msg *bapi.Message) error {
	log.Printf(
		"CHAT:[%q](%d) FROM:[%q](%d) MSG:[%q] ISCMD:[%v]",
		msg.Chat.FirstName, msg.Chat.ID,
		msg.From.FirstName, msg.From.ID,
		msg.Text,
		msg.IsCommand(),
	)

	// identify
	switch msg.Chat.Type {
	case "supergroup", "group":
		if perm := setting.GetGroups().Get(msg.Chat.ID); perm == "" {
			sendT("unauthorized groups", msg.Chat.ID)
			leaveGroup(msg.Chat)
			return nil
		} else if perm == permBanned {
			log.Printf("banned chat [%s](%d) keep using the bot", msg.Chat.FirstName, msg.Chat.ID)
			leaveGroup(msg.Chat)
			return nil
		}
	case "private":
		if perm := setting.GetUsers().Get(msg.From.ID); perm == permBanned {
			return nil
		}
	}

	if fn := registry.getFn(msg.From.ID); fn != nil {
		return fn(msg)
	}

	if msg.IsCommand() {
		err := handleCommand(msg)
		if err != nil {
			err = fmt.Errorf("handle command:[%d]%s:%v", msg.From.ID, msg.Command(), err)
			log.Println(err)
			return err
		}
		return nil
	}
	err := handleMsgText(msg)
	if err != nil {
		err = fmt.Errorf("handle msg:[%s]%s :%s", msg.From.FirstName, msg.Text, err)
		log.Println(err)
		return err
	}
	return nil
}

func handleCommand(msg *bapi.Message) error {
	cmd := msg.Command()
	if fn, ok := botCMD.hasCommand(cmd); ok {
		err := fn(msg)
		if err != nil {
			log.Printf("do cmd %s: %v", cmd, err)
			return err
		}
	}
	return nil
}

func handleMsgText(msg *bapi.Message) error {
	if hasOsuDomain(msg.Text) {
		err := handleOsuLink(msg.Text)
		if err != nil {
			err = fmt.Errorf("handle %s: %v", msg.Text, err)
			log.Println(err)
			return err
		}
		return nil
	}
	return nil
}

func hasOsuDomain(url string) bool {
	const OSUDOMAIN = "https://osu.ppy.sh"
	if len(url) < len(OSUDOMAIN) {
		return false
	}
	for i := 0; i < len(OSUDOMAIN); i++ {
		if url[i] != OSUDOMAIN[i] {
			return false
		}
	}
	return true
}

func handleOsuLink(url string) error {
	return nil
}
