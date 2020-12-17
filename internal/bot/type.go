package bot

import (
	"database/sql"
	"math"
	"sync"
	"sync/atomic"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// masterHandler is a goroutine pool to have multiplexing handler
type masterHandler struct {
	cap         int32         //cap is limit of applied
	applied     int32         //applied is how many handler is applied
	expiredTime time.Duration //expiredTime is a duration of killing handler
	idling      []*msgHandler //idling is a slice of idling worker
	lock        sync.Mutex
	ctx         *Context
}

func newHandler(exp time.Duration, ctx *Context) *masterHandler {
	return newCapHandler(math.MaxInt32, exp, ctx)
}

func newCapHandler(cap int32, exp time.Duration, ctx *Context) *masterHandler {
	m := &masterHandler{
		cap:         cap,
		expiredTime: exp,
		ctx:         ctx,
		idling:      make([]*msgHandler, 0),
	}
	go m.listenAndKill()
	return m
}

func (m *masterHandler) submit(msg *tgbotapi.Message) {
	h := m.getHandler()
	h.task <- msg
}

func (m *masterHandler) apply() *msgHandler {
	atomic.AddInt32(&m.applied, 1)
	return &msgHandler{
		main: m,
		task: make(chan *tgbotapi.Message, 1),
	}
}

// idle append the worker at the last of idling
func (m *masterHandler) idle(h *msgHandler) {
	m.lock.Lock()
	h.lastIdleTime = time.Now()
	m.idling = append(m.idling, h)
	m.lock.Unlock()
}

func (m *masterHandler) isRunOut() bool {
	return m.applied >= m.cap
}

// getHandler use LIFO queue to get free handler
func (m *masterHandler) getHandler() *msgHandler {
	test := func() bool {
		m.lock.Lock()
		allWorking := len(m.idling) == 0
		m.lock.Unlock()
		return allWorking
	}
	for test() {
		if m.applied < m.cap {
			h := m.apply()
			h.Run()
			return h
		}
	}
	m.lock.Lock()
	last := len(m.idling) - 1
	h := m.idling[last]
	m.idling[last] = nil
	m.idling = m.idling[:last]
	h.Run()
	return h
}

func (m *masterHandler) listenAndKill() {
	heartBeat := time.NewTicker(m.expiredTime)
	defer heartBeat.Stop()
	for range heartBeat.C {
		now := time.Now()
		m.lock.Lock()
		if len(m.idling) == 0 {
			m.lock.Unlock()
			continue
		}
		var needExp int
		for i, mh := range m.idling {
			// Queue is arrange by LIFO, so those handler is arranged from old to new.
			if now.Sub(mh.lastIdleTime) <= m.expiredTime {
				break
			}
			needExp = i
			mh.task <- nil
			m.idling[i] = nil
		}
		if needExp == len(m.idling) {
			m.idling = m.idling[:0]
		} else {
			m.idling = m.idling[needExp+1:]
		}
		m.lock.Unlock()
	}
}

// msgHandler is a child handler to handle all the tgbotapi.Message
type msgHandler struct {
	main         *masterHandler
	task         chan *tgbotapi.Message
	lastIdleTime time.Time
}

func (mh *msgHandler) Run() {
	go func() {
		for msg := range mh.task {
			// Stop signal
			if msg == nil {
				atomic.AddInt32(&mh.main.applied, -1)
				return
			}
			// Do jobs
			if msg.IsCommand() {
				doCMD(mh.main.ctx, msg)
			} else {
				doRegex(mh.main.ctx, msg)
			}
			// After job is done, idle handler itself.
			mh.main.idle(mh)
		}
	}()
}

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
	kr      *KnRType
	db      *sql.DB
	groups  *map[int64]interface{}
	bot     *tgbotapi.BotAPI
	send    chan *sendPKG
	timeOut time.Duration
}

func NewContext(
	kr *KnRType,
	db *sql.DB,
	groups *map[int64]interface{},
	bot *B,
	timeout time.Duration,
) *Context {
	return &Context{
		kr:      kr,
		db:      db,
		groups:  groups,
		bot:     bot,
		send:    make(chan *sendPKG),
		timeOut: timeout,
	}
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

// KnRType is a type of map use string value to store list of string
type KnRType map[string][]string

// B is tgbotapi.BotAPI
type B = tgbotapi.BotAPI

// M is tgbotapi.Message
type M = tgbotapi.Message

// C is bot.Context
type C = Context

// SendMethod is a function need (*M, *C) as argument
type SendMethod func(message *M, ctx *C)

// ChTa is tgbotapi.Chattable
type ChTa = tgbotapi.Chattable
