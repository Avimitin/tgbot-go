package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/internal/conf"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var db *sql.DB

func register(user string, password string, host string, database string, port string) (err error) {
	db, err = sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", user, password, host, port, database))

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
func NewDB(dbs *conf.DBSecret) (*sql.DB, error) {
	if db == nil {
		err := register(dbs.User, dbs.Pwd, dbs.Host, dbs.Database, dbs.Port)
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
