package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/conf"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func register(cfg *conf.Config) (err error) {
	db, err = sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", cfg.DBCfg.User, cfg.DBCfg.Password, cfg.DBCfg.Host, cfg.DBCfg.Database))

	if err != nil {
		return err
	}

	// Set limit
	if db != nil {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}

	return nil
}

// This is the function that return database connection for other modules.
// Require database's config to fetch .
func NewDB(dbCFG *conf.Config) (*sql.DB, error) {
	if db == nil {
		err := register(dbCFG)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
