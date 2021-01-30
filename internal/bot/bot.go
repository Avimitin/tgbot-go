package bot

import (
	"fmt"
	"log"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	bot *bapi.BotAPI
)

func Run(config *Configuration) error {
	if config.BotToken == "" {
		return fmt.Errorf("bot token is null")
	}
	var err error
	bot, err = bapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Fatal("fail to initialize bot:\n", err)
	}
	bot.Debug = true
	log.Printf("Successfully establish connection to bot: %s", bot.Self.UserName)

	updateChanConfiguration := bapi.NewUpdate(0)
	updateChanConfiguration.Timeout = 15

	updates, err := bot.GetUpdatesChan(updateChanConfiguration)

	for update := range updates {
		if update.Message != nil {
			err = messageHandler(update.Message)
			continue
		}
	}
	return nil
}

func messageHandler(msg *bapi.Message) error {
	if msg.IsCommand() {
		err := commandsHandler(msg)
		if err != nil {
			log.Println("[msgHandler]error occur when handling command:\n", err)
			return err
		}
		return nil
	}
	err := msgTextHandler(msg)
	if err != nil {
		log.Println("[msgHandler]error occur when handling msg text:\n", err)
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
			err = fmt.Errorf("[msgTextHandler]error occur when handling osu link, %v", err)
			log.Println(err)
			return err
		}
		return nil
	}
	return nil
}

func hasOsuDomain(url string) bool {
	const OSUDOMAIN = "https://osu.ppy.sh"
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
