package bot

import (
	"fmt"
	"github.com/Avimitin/go-bot/internal/pkg/utils/ehAPI"
	"github.com/Avimitin/go-bot/internal/pkg/utils/osuAPI"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Avimitin/go-bot/internal/pkg/conf"
	"github.com/Avimitin/go-bot/internal/pkg/database"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	VERSION = "0.5.8"
	CREATOR = 649191333
)

func Run(cfgPath string, CleanMode bool) {
	ctx := newCTX(cfgPath)
	bot := ctx.Bot()

	updateMsg := tgbotapi.NewUpdate(0)
	updateMsg.Timeout = 20

	updates, err := bot.GetUpdatesChan(updateMsg)

	if err != nil {
		log.Printf("Some error occur when getting update.\nDescriptions: %v", err)
	}

	// 清理模式
	if CleanMode {
		log.Printf("[INFO]Cleaning MSG...")
		updates.Clear()
		log.Printf("[INFO]All message has clear")
	}
	go sendHandler(bot, ctx)
	go callBackQueryHandler(ctx)

	for update := range updates {

		if update.Message == nil {
			if update.CallbackQuery != nil {
				ctx.CBQuery(update.CallbackQuery)
			}
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Chat.Type == "supergroup" && !ctx.IsCertGroup(update.Message.Chat.ID) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你们这啥群啊，别乱拉人，爬爬爬！")
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("[ERR] %s", err)
			}

			_, err = bot.LeaveChat(update.Message.Chat.ChatConfig())
			if err != nil {
				log.Printf("[ERR] Error happen when leave chat %s", err)
			}
			continue
		}

		go msgHandler(ctx, update.Message)
	}
}

func msgHandler(ctx *Context, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		if hasCmdLimit(msg) {
			return
		}
		doCMD(ctx, msg)
	} else if ok, url := isLink(msg); ok {
		urlHandler(ctx, msg.Chat.ID, url)
	} else {
		doRegex(ctx, msg)
	}
}

func isLink(m *M) (bool, string) {
	if m.Entities == nil {
		return false, ""
	}
	for _, entitle := range *m.Entities {
		if entitle.Type == "url" {
			url := m.Text[entitle.Offset:entitle.Length]
			return true, url
		}
	}
	return false, ""
}

func hasCmdLimit(msg *M) bool {
	// have group ↓
	if canDoCmd, ok := cmdDoAble[msg.Chat.ID]; ok {
		// have cmd limit ↓
		if can, ok := canDoCmd[msg.Command()]; ok {
			// cmd disable ↓
			if !can {
				return true
			}
		}
	}
	return false
}

// doCMD do a command
func doCMD(ctx *Context, msg *tgbotapi.Message) {
	if fn, ok := BaseCommand[msg.Command()]; ok {
		fn(msg, ctx)
	} else if fn, ok = PubCommand[msg.Command()]; ok {
		fn(msg, ctx)
	}
}

// doRegex do a regexp job
func doRegex(ctx *Context, msg *tgbotapi.Message) {
	if rpy, ok := RegexKAR(&msg.Text, ctx.KeywordReplies()); ok {
		sendText(ctx, msg.Chat.ID, rpy)
	}
}

func newCTX(path string) *Context {
	// ------bot setting------
	token := conf.LoadBotToken(path)
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("[ERR]%v", err)
	}
	log.Printf("[INFO]Successfully established connection to bot: %s", bot.Self.UserName)
	//------db setting------
	dbs := conf.LoadDBSecret(path)
	db, err := database.NewDB(dbs)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO]Successfully established connection to database")
	//------Cert group loading------
	groups, err := database.SearchGroups(db)
	if err != nil {
		log.Fatal(err)
	}
	groupsSet := make(map[int64]interface{})
	for _, group := range groups {
		groupsSet[group.GroupID] = struct{}{}
	}
	log.Printf("[INFO]Successfully load all certed groups")
	//------Keyword and replies loading------
	k, err := Load(db)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO]Successfully load all keyword and replies")
	key := conf.LoadOSUAPI(path)
	ctx := NewContext(k, db, &groupsSet, bot, key, 30*time.Second)
	return ctx
}

func urlHandler(ctx *C, cid int64, url string) {
	if strings.Contains(url, "osu.ppy.sh") {
		osuURLHandler(ctx, cid, url)
	} else if strings.Contains(url, "hentai.org") {
		exURLHandler(ctx, cid, url)
	}
}

func osuURLHandler(ctx *C, cid int64, url string) {
	bm := osuAPI.GetBeatMapByURL(ctx.osuKey, url)
	if bm == nil {
		log.Println("[osuURLHandler]Got nil beatmap.")
		return
	}
	photo := tgbotapi.NewPhotoShare(cid, fmt.Sprintf("https://assets.ppy.sh/beatmaps/%s/covers/cover.jpg", bm.BeatmapsetID))
	photo.Caption = osuBeatMapCaptionTemplate(bm)
	photo.ParseMode = "HTML"
	//--MakeButton--
	downloadURL := fmt.Sprintf("https://osu.ppy.sh/beatmapsets/%s/download", bm.BeatmapsetID)
	downloadBTN := tgbotapi.InlineKeyboardButton{
		Text: "📎 Download link",
		URL:  &downloadURL,
	}
	photo.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{downloadBTN},
		},
	}

	ctx.Send(NewSendPKG(photo, noReply))
}

func exURLHandler(ctx *C, cid int64, url string) {
	gmd, err := ehAPI.GetComic([]string{url}, 0)
	if err != nil {
		log.Println("[exURLHandler]Error occur when sending request to eh api:", err)
		return
	}
	metaData := gmd.GMD[0]
	photo := tgbotapi.NewPhotoShare(cid, metaData.Thumb)
	var tags string
	for i, tag := range metaData.Tags {
		if i == 6 {
			break
		}
		tags += "#" + tag + " "
	}
	photo.Caption = fmt.Sprintf(`
#NSFW
Title: %s
Category: %s
Like: %s ⭐
Tags: %s ...
`, metaData.TitleJpn, metaData.Category, metaData.Rating, tags)
	ctx.Send(NewSendPKG(photo, noReply))
}

func callBackQueryHandler(ctx *C) {
	queryDataFunc := map[string]func(*Context, *tgbotapi.CallbackQuery){
		"osu":  osuDataHandler,
		"osuu": osuUserModeHandler,
	}

	for {
		select {
		case <-ctx.stop:
			return
		case query := <-ctx.callBack:
			prefixIndex := strings.Index(query.Data, ":")
			prefix := query.Data[:prefixIndex]
			if fn, ok := queryDataFunc[prefix]; ok {
				query.Data = query.Data[prefixIndex+1:]
				go fn(ctx, query)
			} else {
				log.Println("[INFO]Got an unrecognizable query data:", query.Data)
			}
		}
	}
}

func osuDataHandler(ctx *C, query *tgbotapi.CallbackQuery) {
	cb := tgbotapi.NewCallback(query.ID, "Requesting...")
	_, err := ctx.Bot().AnswerCallbackQuery(cb)
	if err != nil {
		log.Println("[osuDataHandler]Error occur when answering callback query:", err)
	}
	newCaption := tgbotapi.NewEditMessageCaption(query.Message.Chat.ID, query.Message.MessageID, "Processing...")
	ctx.Send(NewSendPKG(newCaption, noReply))
	// query data
	bm := osuAPI.GetBeatMap(ctx.osuKey, query.Data)
	text := osuBeatMapCaptionTemplate(bm)
	newCaption = tgbotapi.NewEditMessageCaption(query.Message.Chat.ID, query.Message.MessageID, text)
	newCaption.ParseMode = "HTML"
	// get button information
	bms := osuAPI.GetBeatMapByBeatMapSet(ctx.osuKey, bm.BeatmapsetID)
	newCaption.ReplyMarkup = makeOSUButton(bms)
	ctx.Send(NewSendPKG(newCaption, noReply))
}

func osuUserModeHandler(ctx *C, query *tgbotapi.CallbackQuery) {
	cb := tgbotapi.NewCallback(query.ID, "Requesting...")
	_, err := ctx.Bot().AnswerCallbackQuery(cb)
	if err != nil {
		log.Println("[osuDataHandler]Error occur when answering callback query:", err)
	}
	mode, err := strconv.Atoi(string(query.Data[0]))
	if err != nil {
		sendText(ctx, query.Message.Chat.ID, "Fail to convert user's mode")
		return
	}
	u, err := osuAPI.GetUser(ctx.osuKey, query.Data[1:], "string", mode)
	if err != nil {
		sendText(ctx, query.Message.Chat.ID, "Fail to fetch user's information: "+err.Error())
		return
	}
	newCaption := tgbotapi.NewEditMessageCaption(query.Message.Chat.ID, query.Message.MessageID, osuUserCaptionTemplate(u))
	newCaption.ParseMode = "HTML"
	newCaption.ReplyMarkup = osuModeButton(u)
	ctx.Send(NewSendPKG(newCaption, noReply))
}
