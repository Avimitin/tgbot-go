package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/CFGLoader"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func register() (*sql.DB, error) {
	cfg, err := CFGLoader.LoadCFG()
	if err != nil {
		return nil, err
	}

	database, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", cfg.DBCfg.User, cfg.DBCfg.Password, cfg.DBCfg.Host, cfg.DBCfg.Database))

	if err != nil {
		return nil, err
	}

	// Set limit
	if database != nil {
		database.SetConnMaxLifetime(time.Minute * 3)
		database.SetMaxOpenConns(10)
		database.SetMaxIdleConns(10)
	}

	return database, nil
}

func NewDB() (*sql.DB, error) {
	DB, err := register()
	if err != nil {
		return nil, err
	}
	return DB, nil
}
