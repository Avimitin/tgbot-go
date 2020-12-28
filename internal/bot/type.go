package bot

import (
	"database/sql"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// sendPKG is used to communicate with Send goroutine
type sendPKG struct {
	msg     ChTa
	noReply bool
	resp    chan *M
}

// NewSendPKG return a chatTable PKG. If noRep is false it will return
// a PKG with reply channel for getting reply.
func NewSendPKG(msg tgbotapi.Chattable, noRep bool) *sendPKG {
	var resp chan *M
	if !noRep {
		resp = make(chan *M, 1)
	}
	return &sendPKG{
		msg:     msg,
		noReply: noRep,
		resp:    resp,
	}
}

// Context that contain all needed type is used to communicate
// between main goroutine and child goroutine at runtime.
type Context struct {
	kr       *KnRType
	db       *sql.DB
	groups   *map[int64]interface{}
	bot      *tgbotapi.BotAPI
	send     chan *sendPKG
	osuKey   string
	timeOut  time.Duration
	callBack chan *tgbotapi.CallbackQuery
	stop     chan int32
}

func NewContext(
	kr *KnRType,
	db *sql.DB,
	groups *map[int64]interface{},
	bot *B,
	osuKey string,
	timeout time.Duration,
) *Context {
	return &Context{
		kr:       kr,
		db:       db,
		groups:   groups,
		bot:      bot,
		send:     make(chan *sendPKG, 1),
		osuKey:   osuKey,
		timeOut:  timeout,
		callBack: make(chan *tgbotapi.CallbackQuery, 1),
		stop:     make(chan int32),
	}
}

func (ctx *Context) StopAll() {
	close(ctx.stop)
}

func (ctx *Context) Send(msg *sendPKG) {
	ctx.send <- msg
}

func (ctx *Context) Bot() *tgbotapi.BotAPI {
	return ctx.bot
}

func (ctx *Context) KeywordReplies() *KnRType {
	return ctx.kr
}

func (ctx *Context) DB() *sql.DB {
	return ctx.db
}

func (ctx *Context) SetGroup(g int64) {
	(*ctx.groups)[g] = struct{}{}
}

func (ctx *Context) DelGroup(g int64) {
	delete(*ctx.groups, g)
}

func (ctx *Context) Groups() map[int64]interface{} {
	return *ctx.groups
}

func (ctx *Context) IsCertGroup(id int64) bool {
	_, ok := (*ctx.groups)[id]
	return ok
}

func (ctx *Context) CBQuery(q *tgbotapi.CallbackQuery) {
	ctx.callBack <- q
}

// KnRType is a type of map use string value to store list of string
type KnRType map[string][]string

// B is tgbotapi.BotAPI
type B = tgbotapi.BotAPI

// M is tgbotapi.Message
type M = tgbotapi.Message

// C is bot.Context
type C = Context

// CMDMethod is a function need (*M, *C) as argument
type CMDMethod func(message *M, ctx *C)

// ChTa is tgbotapi.Chattable
type ChTa = tgbotapi.Chattable

type cmdFunc map[string]CMDMethod

func (cf cmdFunc) hasCommand(cmd string) bool {
	_, ok := cf[cmd]
	return ok
}

type mjx struct {
	Pic    string `json:"pic"`
	Imgurl string `json:"imgurl"`
}
