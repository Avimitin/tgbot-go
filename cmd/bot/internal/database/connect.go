package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/CFGLoader"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db setUpDatabase

type setUpDatabase struct {
	dataBaseConn *sql.DB // Here is the actual database connection
	hasSetUp     bool    // Here record if the database connection had set up or not
}

func register() error {
	cfg, err := CFGLoader.LoadCFG()
	if err != nil {
		return err
	}

	db.dataBaseConn, err = sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8", cfg.DBCfg.User, cfg.DBCfg.Password, cfg.DBCfg.Host, cfg.DBCfg.Database))

	if err != nil {
		return err
	}

	// Set limit
	if db.dataBaseConn != nil {
		db.dataBaseConn.SetConnMaxLifetime(time.Minute * 3)
		db.dataBaseConn.SetMaxOpenConns(10)
		db.dataBaseConn.SetMaxIdleConns(10)
	}
	db.hasSetUp = true
	return nil
}

// This is the function that return database connection for other modules
func NewDB() (*sql.DB, error) {
	if !db.hasSetUp {
		err := register()
		if err != nil {
			return nil, err
		}
	}
	return db.dataBaseConn, nil
}
