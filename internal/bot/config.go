package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	GetUsers() *Users   // GetUsers return a set of user
	GetGroups() *Groups // GetGroups return a list of certed groups
	Secret() *Secret    // Secret return a set of secret store
	Update() error      // Update update setting data store
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

type JsonConfig struct {
	BotToken string           `json:"bot_token"`
	Groups   map[int64]string `json:"groups"`
	Users    map[int]string   `json:"users"`
	mu       sync.RWMutex     `json:"-"`
}

func (cfg *JsonConfig) GetUsers() Users {
	return cfg.Users
}

func (cfg *JsonConfig) GetGroups() Groups {
	return cfg.Groups
}

func (cfg *JsonConfig) Secret() Secret {
	return map[string]string{"bot_token": cfg.BotToken}
}

func (cfg *JsonConfig) dumpConfig(path string) error {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal %+v failed: %v", cfg, err)
	}
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write %s failed:%v", path, err)
	}
	return nil
}

func (cfg *JsonConfig) Update() error {
	err := cfg.dumpConfig(WhereCFG("") + "/config.json")
	if err != nil {
		path := os.Getenv("HOME") + "/config.json.tmp"
		err = cfg.dumpConfig(path)
		if err != nil {
			log.Println("save failed")
			return fmt.Errorf("save config %+v to %s failed", cfg, path)
		}
	}
	log.Println("save successfully")
	return nil
}

func (cfg *JsonConfig) Prepare() error {
	cfg.BotToken = ""
	cfg.Groups = make(map[int64]string)
	cfg.Users = make(map[int]string)
	cfgPath := WhereCFG("")
	if cfgPath == "" {
		return errors.New("no config")
	}
	cfgPath = cfgPath + "/config.json"
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return fmt.Errorf("read %s:%v", cfgPath, err)
	}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return fmt.Errorf("decode %s:%v", data, err)
	}
	return nil
}

func NewConfig() *JsonConfig {
	return new(JsonConfig)
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
