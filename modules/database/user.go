package database

import "gorm.io/gorm"

const (
	PermOwner = iota
	PermAdmin
	PermChannelManager
	PermNormal
	PermBan
)

var (
	permission = map[int32]string{
		PermOwner:          "owner",
		PermAdmin:          "admin",
		PermChannelManager: "channelManager",
		PermNormal:         "normal",
		PermBan:            "ban",
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
