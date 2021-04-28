package database

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
	ID int

	// PermDesc describe what PermID mean
	PermDesc string
	PermID   int32
}
