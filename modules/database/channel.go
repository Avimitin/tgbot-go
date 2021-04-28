package database

import (
	"time"

	"gorm.io/gorm"
)

// Post contain channel post information
type Post struct {
	gorm.Model

	MsgID      int
	Content    string
	SendedDate time.Time
	Link       string
}
