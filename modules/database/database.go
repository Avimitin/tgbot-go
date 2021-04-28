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

func (db *BotDB) GetUser(id int) (*User, error) {
	var u User

	result := db.cnct.Where("user_id = ?", id).First(&u)

	switch result.Error {
	case nil:
		break
	case gorm.ErrRecordNotFound:
		return nil, nil
	default:
		return nil, fmt.Errorf("get user %d: %w", id, result.Error)
	}

	return &u, nil
}

func (db *BotDB) NewUser(id int, permID int32) (*User, error) {
	u := &User{
		UserID:   id,
		PermID:   permID,
		PermDesc: permission[permID],
	}

	result := db.cnct.Create(&u)

	if result.Error != nil {
		return nil, fmt.Errorf("create user %d: %w", id, result.Error)
	}

	return u, nil
}
