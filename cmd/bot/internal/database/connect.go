package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/CFGLoader"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	DB       sql.DB
	hasSetUp bool
)

func register(database *sql.DB) error {
	cfg, err := CFGLoader.LoadCFG()
	if err != nil {
		return err
	}
	database, err = sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", cfg.DBCfg.User, cfg.DBCfg.Password, cfg.DBCfg.Host, cfg.DBCfg.Database))

	// Set limit
	if database != nil {
		database.SetConnMaxLifetime(time.Minute * 3)
		database.SetMaxOpenConns(10)
		database.SetMaxIdleConns(10)
	}

	return nil
}

func connect() (*sql.DB, error) {
	err := register(&DB)
	if err != nil {
		return nil, err
	}
	hasSetUp = true
	return &DB, nil
}

func NewDB() (*sql.DB, error) {
	if !hasSetUp {
		return connect()
	}
	return &DB, nil
}
