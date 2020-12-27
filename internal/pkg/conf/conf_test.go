package conf

import (
	"os"
	"os/user"
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
	dir := u.HomeDir + "/.config/avimi-bot"
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	if path := WhereCFG(""); path != dir {
		t.Fatalf("home path test failed, want %s got %s", dir, path)
	}

}

func TestLoadINI(t *testing.T) {
	file := LoadINI(WhereCFG("F:\\code\\golang\\go-bot\\cfg"))
	if file == nil {
		t.Fatal("get null value")
	}
	if !file.Section("bot").HasKey("token") {
		t.Fatal("Can't fetch token")
	}
}

func TestLoadBotToken(t *testing.T) {
	token := LoadBotToken(WhereCFG(""))
	if token != "114514:qwerty" {
		t.Fatalf("unwanted token")
	}
}

func TestLoadDBSecret(t *testing.T) {
	db := LoadDBSecret(WhereCFG(""))
	if db.MySqlURL() != "tgbot:tgbot@tcp(127.0.0.1:3306)/tgbotDB" {
		t.Fatalf("Unwanted db secret, got %s", db.MySqlURL())
	}
}
