package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

// SendT return a message config that only contain text
func SendT(chat int64, text string) ChTa {
	return tgbotapi.NewMessage(chat, text)
}

// SendP return a message config that parse message with given parse mode
func SendP(chatID int64, text string, parse string) ChTa {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parse
	return msg
}

// sendHandler will apply a individual goroutine
// to manage all the send and receive stuff.
// For closing this goroutine, send a nil sendPKG.
func sendHandler(bot *B, c *C) {
	for {
		select {
		case <-c.stop:
			return
		case msgP := <-c.send:
			log.Println("[INFO]Receive new send request:", msgP)
			if msgP == nil {
				return
			}
			resp, err := bot.Send(msgP.msg)
			if err != nil {
				log.Println("[ERR]sendHandler got error:", err)
			}
			if !msgP.noReply {
				msgP.resp <- &resp
			}
		}
	}
}
