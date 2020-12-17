package bot

import (
	MyTimer "github.com/Avimitin/go-bot/internal/utils/timer"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"sync"
	"testing"
	"time"
)

func TestGoPool(t *testing.T) {
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "foo bar fuck bob dude",
	}
	k := KnRType{"fuck": []string{"fuck you", "fuck him", "fuck me"}}
	ctx := NewContext(&k, nil, nil, nil, 30*time.Second)

	var Wg sync.WaitGroup
	handler := newHandler(time.Second, ctx)

	testFunc := func(t int) {
		timer := MyTimer.NewTimer()
		pureGoTest(&Wg, t, msg.Text, &k)
		Wg.Wait()
		log.Printf("Ding! pure goroutine use %.5f second", timer.StopCounting()/1000000000)
		// ---------
		timer = nil
		timer = MyTimer.NewTimer()
		done := make(chan int32)
		poolGoTest(&Wg, t, handler, msg)
		go func() {
			for {
				select {
				case <-ctx.send:
					Wg.Done()
				case <-done:
					return
				}
			}
		}()
		Wg.Wait()
		done <- 0
		log.Printf("Ding! goroutine with pool use %.5f second", timer.StopCounting()/1000000000)
	}

	log.Println("Test 1, 100 message")
	testFunc(100)
	log.Println("Test 2, 10000 message")
	testFunc(10000)
	log.Println("Test 3, 100000 message")
	testFunc(100000)
	log.Println("Test 4, 1000000 message")
	testFunc(1000000)
	log.Println("Test 5, 10000000 message")
	testFunc(10000000)
}

func poolGoTest(wg *sync.WaitGroup, max int, h *masterHandler, msg *tgbotapi.Message) {
	for t := 0; t < max; t++ {
		wg.Add(1)
		h.submit(msg)
	}
}

func pureGoTest(wg *sync.WaitGroup, max int, s string, k *KnRType) {
	for t := 0; t < max; t++ {
		go func() {
			wg.Add(1)
			_, bingo := RegexKAR(s, k)
			if bingo {
				wg.Done()
			}
		}()
	}
}
