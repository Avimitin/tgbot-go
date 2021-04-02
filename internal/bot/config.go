package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const (
	// User permission
	permBanned = "ban"
	permNormal = "norm"
	permComic  = "comic"
	permAdmin  = "admin"
)

// SettingsGetter contain bot needed information at runtime
type SettingsGetter interface {
	GetUsers() *Users    // GetUsers return a set of user
	GetGroups() *Groups  // GetGroups return a list of certed groups
	Secret() *Secret     // Secret return a set of secret store
	Update(string) error // Update update setting data store
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

// JsonConfig store the local bot configuration at json file type
// BotToken, Groups, Users type are exposed for parsed into json, but it's
// not recommended to use them directly. You should use the packaging method.
type JsonConfig struct {
	BotToken string           `json:"bot_token"`
	Groups   map[int64]string `json:"groups"`
	Users    map[int]string   `json:"users"`

	s  *Secret      `json:"-"`
	g  *Groups      `json:"-"`
	u  *Users       `json:"-"`
	mu sync.RWMutex `json:"-"`
}

// NewJsonConfig read file at given path value `p` and decode it to
// JsonConfig. Must ensure `p` is a valid file path and is a json type file.
func NewJsonConfig(p string) (*JsonConfig, error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("read %s:%v", p, err)
	}
	var cfg *JsonConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("decode %s:%v", data, err)
	}

	cfg.u = &Users{userPermMap: make(map[int]string)}
	for k, v := range cfg.Users {
		cfg.u.Set(k, v)
	}

	cfg.g = &Groups{groupPermMap: make(map[int64]string)}
	for k, v := range cfg.Groups {
		cfg.g.Set(k, v)
	}

	cfg.s = &Secret{ma: make(map[string]string)}
	cfg.s.Set("bot_token", cfg.BotToken)

	cfg.Users = nil
	cfg.Groups = nil

	return cfg, nil
}

// GetUsers return the Users type
func (cfg *JsonConfig) GetUsers() *Users {
	return cfg.u
}

// GetGroups return the Groups type
func (cfg *JsonConfig) GetGroups() *Groups {
	return cfg.g
}

// GetSecret return the Secret type
func (cfg *JsonConfig) Secret() *Secret {
	return cfg.s
}

func (cfg *JsonConfig) dumpConfig(path string) error {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()

	cfg.Users = cfg.u.Traverse()
	cfg.Groups = cfg.g.Traverse()

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal %+v failed: %v", cfg, err)
	}
	err = ioutil.WriteFile(path, data, 0640)
	if err != nil {
		return fmt.Errorf("write %s failed:%v", path, err)
	}

	cfg.Users = nil
	cfg.Groups = nil
	return nil
}

func (cfg *JsonConfig) Update(p string) error {
	err := cfg.dumpConfig(p)
	if err != nil {
		return fmt.Errorf("save config %+v to %s failed", cfg, p)
	}
	return nil
}

// WhereCFG give the config loader specific config path.
// If p is given, it will return given path. Else this function will
// find config from environment variable "BOTCFGPATH" or user's home directory.
// If can't found config from this place, return a null string value.
func WhereCFG(p string) (path string) {
	// if p had given, use p as path
	if p != "" {
		return p
	}

	// if path is specific in environment variable, use env as path
	if path = os.Getenv("BOTCFGPATH"); path != "" {
		return path
	}

	// if config path exist in user's home directory, use it as path
	if userHomePath := os.Getenv("HOME"); userHomePath != "" {
		return userHomePath + "/.config/go-bot"
	}
	return ""
}
