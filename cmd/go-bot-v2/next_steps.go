package main

import (
	"sync"

	tb "gopkg.in/tucnak/telebot.v2"
)

// contextData store key-value data
type contextData map[string]string

type nextFunc func(*tb.Message, contextData) error

// contextPayload store next function to execute and the needed context data
type contextPayload struct {
	fn   nextFunc
	data map[string]string
}

// userRegistration register user and the next step
type userRegistration map[int]*contextPayload

// chatRegistration register a chat and users belong to the chat
type chatRegistration map[int64]userRegistration

var (
	mu     sync.Mutex
	cmdCtx = make(chatRegistration)
)

// refisNextStep register all the step in a tree structure
func regisNextStep(cid int64, uid int, d contextData, nf nextFunc) {
	mu.Lock()
	defer mu.Unlock()

	payload := &contextPayload{fn: nf, data: d}

	// if the group has register other user
	if uc, ok := cmdCtx[cid]; ok {
		uc[uid] = payload
		return
	}

	// else create a new group registration
	cmdCtx[cid] = userRegistration{uid: payload}
}

func delContext(cid int64, uid int) {
	uc, ok := cmdCtx[cid]
	if !ok {
		return
	}

	_, ok = uc[uid]
	if !ok {
		return
	}

	delete(uc, uid)
}

func getRegis(cid int64, uid int) *contextPayload {
	uc, ok := cmdCtx[cid]
	if !ok {
		return nil
	}

	p, ok := uc[uid]
	if !ok {
		return nil
	}

	return p
}
