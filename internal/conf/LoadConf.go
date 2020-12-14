package conf

import (
	"gopkg.in/ini.v1"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

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
	u, err := user.Current()
	if err != nil {
		log.Fatalf("[ERR]conf/LoadConf.go read user error: %v", err)
	}
	files, err := ioutil.ReadDir(u.HomeDir + "/.config")
	if err != nil {
		log.Fatalf("[ERR]conf/LoadConf.go read directory %s error: %v", u.HomeDir+"/.config", err)
	}
	for _, file := range files {
		if path = file.Name(); path == "avimi-bot" && file.IsDir() {
			return u.HomeDir + "/.config/" + path
		}
	}
	return ""
}

// LoadINI load file with given uri
func LoadINI(path string) *ini.File {
	cfg, err := ini.Load(path + "/cfg.ini")
	if err != nil {
		log.Fatalf("[ERROR]Read %s error", path+"/cfg.ini")
	}
	return cfg
}

// LoadBotToken return bots token define in cfg.ini
func LoadBotToken(path string) string {
	cfg := LoadINI(path)
	return cfg.Section("bot").Key("token").String()
}

// LoadDBSecret return database information pre-define in cfg.ini
func LoadDBSecret(path string) *DBSecret {
	cfg := LoadINI(path)
	dbSec := cfg.Section("DB")
	return &DBSecret{
		user:     dbSec.Key("user").String(),
		pwd:      dbSec.Key("password").String(),
		host:     dbSec.Key("host").String(),
		database: dbSec.Key("database").String(),
		port:     dbSec.Key("port").String(),
	}
}
