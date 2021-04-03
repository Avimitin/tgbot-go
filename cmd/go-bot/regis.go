package main

import (
	"sync"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type msgHandlerFunc func(*bapi.Message) error
type infoMapType map[string]string // infoMapType contain message that store by key to value

func (imt infoMapType) get(key string) string {
	return imt[key]
}

type registration struct {
	registerFunc map[int]msgHandlerFunc
	registerInfo map[int]infoMapType
	mu           sync.RWMutex
}

func NewRegistration() *registration {
	return &registration{
		registerFunc: make(map[int]msgHandlerFunc),
		registerInfo: make(map[int]infoMapType),
	}
}

func (r *registration) getFn(u int) msgHandlerFunc {
	r.mu.RLock()
	fn := r.registerFunc[u]
	r.mu.RUnlock()
	return fn
}

func (r *registration) getInfo(u int) infoMapType {
	r.mu.RLock()
	info := r.registerInfo[u]
	r.mu.RUnlock()
	return info
}

func (r *registration) registerNextFunc(m *bapi.Message, fn msgHandlerFunc, info infoMapType) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.registerFunc[m.From.ID] = fn
	r.registerInfo[m.From.ID] = info
}

func (r *registration) clear(m *bapi.Message) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.registerFunc, m.From.ID)
	delete(r.registerInfo, m.From.ID)
}

// nuclear clear all the data, use with attention!
func (r *registration) nuclear(m *bapi.Message) {
	r = nil
	r = new(registration)
}
