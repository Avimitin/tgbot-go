package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/auth"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func connect() (*sql.DB, error) {
	cfg, err := auth.NewCFG()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", cfg.DBCfg.User, cfg.DBCfg.Password, cfg.DBCfg.Host, cfg.DBCfg.Database))

	// Set limit
	if db != nil {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}

	return db, db.Ping()
}
