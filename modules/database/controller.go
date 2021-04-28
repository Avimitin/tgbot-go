package database

// DataController can controll local data store
type DataController interface {
	GetUser(id int) (*User, error)
	NewUser(id int, permID int32) (*User, error)
	SetUser(id int, permID int32) (*User, error)
}
