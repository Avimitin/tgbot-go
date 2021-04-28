package intergration

import (
	"fmt"
	"os"
	"testing"

	"github.com/Avimitin/go-bot/modules/database"
)

var (
	model *database.BotDB
)

func envFB(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

func TestMain(m *testing.M) {
	var err error

	model, err = database.NewBotDB(
		envFB("GO_BOT_MYSQL_DSN",
			"root:goBotDB@tcp(127.0.0.1:3306)/goBotDB?charset=utf8mb4&parseTime=True&loc=Local"),
	)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m.Run()
}

func TestDatabaseConnection(t *testing.T) {
	var (
		uid        = 1234567
		perm int32 = 0
	)
	t.Run("test new user", func(t *testing.T) {
		u, err := model.NewUser(uid, perm)
		if err != nil {
			t.Fatal(err)
		}

		if u == nil {
			t.Errorf("user is nil")
		}
	})

	t.Run("test get user", func(t *testing.T) {
		u, err := model.GetUser(uid)
		if err != nil {
			t.Fatal(err)
		}

		if u.UserID != uid && u.PermID != perm {
			t.Errorf("get %#v is not wanted", u)
		}
	})
}
