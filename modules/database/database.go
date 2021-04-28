package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BotDB struct {
	cnct *gorm.DB
}

// NewBotDB return a abstract packaging database
func NewBotDB(dsn string) (*BotDB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", dsn, err)
	}

	db.AutoMigrate(&User{})

	return &BotDB{cnct: db}, nil
}
