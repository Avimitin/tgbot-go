package bot

import (
	"sync"

	bapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type command map[string]func(m *bapi.Message) error

// hasCommand return if cmd is executive.
func (c command) hasCommand(cmd string) (func(m *bapi.Message) error, bool) {
	fn, ok := c[cmd]
	return fn, ok
}

// Users store user permission, it is safe for simultaneous
// use by mulitiple goroutine
type Users struct {
	userPermMap map[int]string
	m           sync.RWMutex
}

// Get return the given user's permission
func (u *Users) Get(user int) (perm string) {
	u.m.RLock()
	defer u.m.RUnlock()
	if p, ok := u.userPermMap[user]; ok {
		return p
	}
	return ""
}

func (u *Users) Traverse() map[int]string {
	u.m.RLock()
	defer u.m.RUnlock()
	return u.userPermMap
}

// Set set the given user and permission into map
func (u *Users) Set(user int, perm string) {
	u.m.Lock()
	defer u.m.Unlock()
	u.userPermMap[user] = perm
}

// Groups is a map that store group information
type Groups struct {
	groupPermMap map[int64]string
	m            sync.RWMutex
}

// Get return the given group's permission
func (g *Groups) Get(group int64) (perm string) {
	g.m.RLock()
	defer g.m.RUnlock()
	if p, ok := g.groupPermMap[group]; ok {
		return p
	}
	return ""
}

// Set set the given group and permission into map
func (g *Groups) Set(group int64, perm string) {
	g.m.Lock()
	defer g.m.Unlock()
	g.groupPermMap[group] = perm
}

func (g *Groups) Traverse() map[int64]string {
	g.m.RLock()
	defer g.m.RUnlock()
	return g.groupPermMap
}

// Secret store any secret like token or password
// bot need to use at runtime
type Secret struct {
	ma map[string]string
	mu sync.Mutex
}

// Set set the given key and val into map
func (s *Secret) Set(key string, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ma[key] = val
}

// Get return the relative val with given key
func (s *Secret) Get(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.ma[key]; ok {
		return v
	}
	return ""
}
