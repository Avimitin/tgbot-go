package conf

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

func LoadINI() *ini.File {
	cfg, err := ini.Load("F:/code/golang/go-bot/cfg/secret.ini")
	if err != nil {
		log.Printf("[ini]Unable to read ini config file.")
		os.Exit(1)
	}
	return cfg
}

func LoadBotToken() string {
	cfg := LoadINI()
	return cfg.Section("bot").Key("token").String()
}

func LoadDBSecret() *DBSecret {
	cfg := LoadINI()
	dbSec := cfg.Section("DB")
	return &DBSecret{
		User:     dbSec.Key("user").String(),
		Pwd:      dbSec.Key("password").String(),
		Host:     dbSec.Key("host").String(),
		Database: dbSec.Key("database").String(),
		Port:     dbSec.Key("port").String(),
	}
}
