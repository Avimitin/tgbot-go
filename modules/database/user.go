package database

//User describe a user
type User struct {
	ID int

	// PermDesc describe what PermID mean
	PermDesc string
	PermID   int32
}
