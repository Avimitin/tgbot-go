package database

import (
	"database/sql"
	"github.com/Avimitin/go-bot/internal/conf"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func register(dsn string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("[ERR]Fail to open database")
		return nil, err
	}
	// Set limit
	if db != nil {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}
	return db, nil
}

// NewDB return a database connection
func NewDB(dbs *conf.DBSecret) (*sql.DB, error) {
	return register(dbs.MySqlURL())
}
