package bot

import (
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"testing"
)

func TestWhereCFG(t *testing.T) {
	const PATH string = "PATH/TO/CFG"
	if path := WhereCFG(PATH); path != PATH {
		t.Fatalf("Want %s got %s", PATH, path)
	}

	err := os.Setenv("BOTCFGPATH", PATH)
	if err != nil {
		t.Fatalf("Error happen when setting config path env")
	}
	if path := WhereCFG(""); path != PATH {
		t.Fatalf("Env test fail. Want %s got %s", PATH, path)
	}
	err = os.Unsetenv("BOTCFGPATH")
	if err != nil {
		t.Log(err)
	}

	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	dir := u.HomeDir + "/.config/go-bot"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	if path := WhereCFG(""); path != dir {
		t.Fatalf("home path test failed, want %s got %s", dir, path)
	}
}

func TestNewConfig(t *testing.T) {
	const PATH string = "../../cfg"
	config := newConfigFromGivenPath(PATH)
	if config == nil {
		t.Errorf("got nil config")
	}
	if config.BotToken == "" {
		t.Errorf("can't not read config")
	}
}

func TestDumpCFG(t *testing.T) {
	cfg := Configuration{
		BotToken:     "abc:123",
		Groups: []int64{123, 456},
	}

	os.Setenv("BOTCFGPATH", "../../cfg")
	if os.Getenv("BOTCFGPATH") != "../../cfg" {
		t.Fatalf("fail to set environment variable")
	}

	err := cfg.DumpConfig()
	if err != nil {
		t.Fatal(err)
	}

	tCfg := NewConfig()
	if tCfg.BotToken != cfg.BotToken {
		t.Errorf("got %+v", tCfg)
	}

	os.Unsetenv("BOTCFGPATH")
}

func TestConfigSave(t *testing.T) {
	cfg := Configuration{
		BotToken: "abc:123",
		certedGroups: map[int64]interface{}{
			1234: struct{}{},
		},
	}

	os.Setenv("BOTCFGPATH", "../../cfg")
	cfg.save()
	_, err := ioutil.ReadFile("../../cfg/config.json")
	if err != nil {
		t.Fatal(err)
	}

	nCFG := NewConfig()

	if !reflect.DeepEqual(cfg.certedGroups, nCFG.certedGroups) {
		t.Fatalf("got %+v", nCFG)
	}

	os.Unsetenv("BOTCFGPATH")
}
