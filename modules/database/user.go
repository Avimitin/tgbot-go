package database

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	permission = map[int32]string{
		0: "owner",
		1: "admin",
		2: "channelManager",
		3: "normal",
		4: "ban",
	}
)

//User describe a user
type User struct {
	gorm.Model

	UserID int

	// PermDesc describe what PermID mean
	PermDesc string
	PermID   int32
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

func (db *BotDB) SetUser(id int, permID int32) (*User, error) {
	user, err := db.GetUser(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	user.PermID = permID
	user.PermDesc = permission[permID]

	result := db.cnct.Save(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("updating user %d: %v", id, err)
	}

	return user, nil
}
