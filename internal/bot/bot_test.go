package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"testing"
)

func TestIsURL(t *testing.T) {
	msg := tgbotapi.Message{
		Entities: &[]tgbotapi.MessageEntity{
			{
				Length: 10,
				Offset: 0,
				Type:   "url",
			},
		},
		Text: "osu.ppy.sh",
	}
	if ok, url := isLink(&msg); ok {
		fmt.Println(url)
		return
	}
	t.Fatal("Can't recognize url")
}
