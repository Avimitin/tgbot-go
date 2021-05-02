package main

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml"
)

// Configuration contain bot setting and database setting
type Configuration struct {
	Bot      BotSetting
	Database DBSetting
}

// BotSetting contain bot configuration at runtime
type BotSetting struct {
	Token    string
	LogLevel string
	Owner    int
}

// DBSetting contain database configuration at runtime
type DBSetting struct {
	User     string
	Password string
	Protocol string
	Addr     string
	Name     string
	Params   string
	LogLevel string
}

// EncodeMySQLDSN use DBSetting to encode a mysql data source link
func (d *DBSetting) EncodeMySQLDSN() string {
	var buf bytes.Buffer

	if len(d.User) > 0 {
		buf.WriteString(d.User)
		if len(d.Password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(d.Password)
		}
		buf.WriteByte('@')
	}

	if len(d.Protocol) > 0 {
		buf.WriteString(d.Protocol)
		if len(d.Addr) > 0 {
			buf.WriteByte('(')
			buf.WriteString(d.Addr)
			buf.WriteByte(')')
		}
	}

	buf.WriteByte('/')
	buf.WriteString(d.Name)
	buf.WriteByte('?')

	if len(d.Params) > 0 {
		buf.WriteString(d.Params)
	} else {
		buf.WriteString("charset=utf8mb4&parseTime=True&loc=Local")
	}

	return buf.String()
}

var cfg *Configuration

// ReadConfig parse toml document from $GOBOT_CONFIG_PATH or $HOME/.config/go-bot/config.toml
func ReadConfig() *Configuration {
	configPath := os.Getenv("GOBOT_CONFIG_PATH")
	if configPath == "" {
		configPath = path.Join(os.Getenv("HOME"), ".config", "go-bot", "config.toml")
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("reading config from %q: %v\n", configPath, err)
		os.Exit(1)
	}

	var cfg Configuration
	err = toml.Unmarshal(configFile, &cfg)
	if err != nil {
		fmt.Printf("decode config %q: %v\n", configFile, err)
		os.Exit(1)
	}

	return &cfg
}

func init() {
	cfg = ReadConfig()
}
