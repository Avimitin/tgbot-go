package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BotDB struct {
	cnct *gorm.DB
}

// NewBotDB return a abstract packaging database
func NewBotDB(dsn string) (*BotDB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", dsn, err)
	}
	return &BotDB{cnct: db}, nil
}
