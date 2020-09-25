package auth

var creator int = 649191333

func IsCreator(uid int) bool {
	return uid == creator
}
