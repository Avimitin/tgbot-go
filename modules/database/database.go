package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BotDB struct {
	cnct *gorm.DB
}

// NewBotDB return a abstract packaging database
func NewBotDB(dsn string, logLevel string) (*BotDB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(transLogLevel(logLevel)),
	})
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", dsn, err)
	}

	db.AutoMigrate(&User{})

	return &BotDB{cnct: db}, nil
}

func transLogLevel(logLevel string) logger.LogLevel {
	switch logLevel {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	}

	return logger.Silent
}
