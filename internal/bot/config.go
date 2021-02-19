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

// Setting contain bot needed information at runtime
type Setting interface {
	Prepare() error    // Prepare initialize setting
	GetUsers() Users   // GetUsers return a set of user
	GetGroups() Groups // GetGroups return a list of certed groups
	Secret() Secret    // Secret return a set of secret store
}

type Users map[int]string
type Groups map[int64]string
type Secret map[string]string

func (s Secret) Set(key string, val string) {
	s[key] = val
}

func (s Secret) Get(key string) (val string) {
	return s[key]
}

type Configuration struct {
	BotToken string           `json:"bot_token"`
	Groups   map[int64]string `json:"groups"`
	Users    map[int]string   `json:"users"`
	mu       sync.Mutex       `json:"-"`
}

func (cfg *Configuration) GetUsers() map[int]string {
	return cfg.Users
}

func (cfg *Configuration) GetGroups() map[int64]string {
	return cfg.Groups
}

func (cfg *Configuration) Secret() map[string]string {
	return map[string]string{"bot_token": cfg.BotToken}
}

func (cfg *Configuration) DumpConfig() error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal %+v failed: %v", cfg, err)
	}
	path := WhereCFG("") + "/config.json"
	err = ioutil.WriteFile(path, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("write %s failed:%v", path, err)
	}
	return nil
}

func (cfg *Configuration) save() {
	err := cfg.DumpConfig()
	if err != nil {
		path := os.Getenv("HOME") + "/config.json.tmp"
		log.Printf("dump config:%v", err)
		log.Printf("saving tmp file to: %s", path)
		byt, err := json.Marshal(cfg)
		if err != nil {
			log.Printf("marshal %+v failed:%v", cfg, err)
			return
		}
		err = ioutil.WriteFile(path, byt, os.ModePerm)
		if err != nil {
			log.Println("saving tmp file failed, program exit:", err)
		}
	}
	log.Println("save successfully")
}

func (cfg *Configuration) Prepare() error {
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

func NewConfig() *Configuration {
	return new(Configuration)
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
