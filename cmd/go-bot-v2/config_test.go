package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	cfg := ReadConfig()
	if cfg == nil {
		t.Error("cfg is nil")
	}

	if cfg.Bot.Token != "1234567" {
		t.Errorf("token is not wanted, got %s", cfg.Bot.Token)
	}

	if dsn := cfg.Database.EncodeMySQLDSN(); dsn !=
		"root:password@tcp(127.0.0.1:3306)/bot_db?param=value&param2=value2" {
		t.Errorf("dsn is not wanted, got %s", dsn)
	}
}
