package bot

import (
	"errors"
	"fmt"
	"log"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot      *bapi.BotAPI
	setting  Setting
	registry = NewRegistration()
)

func Run(s Setting) error {
	if s == nil {
		return errors.New("setting not initialized yet")
	}
	setting = s
	if err := setting.Prepare(); err != nil {
		return fmt.Errorf("prepare setting: %v", err)
	}

	botToken := setting.Secret().Get("bot_token")
	if botToken == "" {
		return fmt.Errorf("bot token is null")
	}
	var err error
	bot, err = bapi.NewBotAPI(botToken)
	if err != nil {
		return fmt.Errorf("connect to api server: %v", err)
	}
	bot.Debug = true
	log.Printf("Successfully establish connection to bot: %s", bot.Self.UserName)

	updateChanConfiguration := bapi.NewUpdate(0)
	updateChanConfiguration.Timeout = 15

	updates, err := bot.GetUpdatesChan(updateChanConfiguration)
	if err != nil {
		return fmt.Errorf("get update chan:%v", err)
	}

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
	// identify
	switch msg.Chat.Type {
	case "supergroup", "group":
		if _, ok := setting.GetGroups()[msg.Chat.ID]; !ok {
			sendT("unauthorized groups, contact @avimibot", msg.Chat.ID)
			leaveGroup(msg.Chat)
			return nil
		}
	case "private":
		if perm := setting.GetUsers()[msg.From.ID]; perm == permBanned {
			return nil
		}
	}

	if fn := registry.getFn(msg.From.ID); fn != nil {
		return fn(msg)
	}

	if msg.IsCommand() {
		err := commandsHandler(msg)
		if err != nil {
			err = fmt.Errorf("handle command:[%d]%s:%v", msg.From.ID, msg.Command(), err)
			log.Println(err)
			return err
		}
		return nil
	}
	err := msgTextHandler(msg)
	if err != nil {
		err = fmt.Errorf("handle msg:[%s]%s :%s", msg.From.FirstName, msg.Text, err)
		log.Println(err)
		return err
	}
	return nil
}

func commandsHandler(msg *bapi.Message) error {
	if fn, ok := botCMD.hasCommand(msg.Command()); ok {
		err := fn(msg)
		if err != nil {
			log.Println("[cmdHandler]", err)
			return err
		}
	}
	return nil
}

func msgTextHandler(msg *bapi.Message) error {
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
