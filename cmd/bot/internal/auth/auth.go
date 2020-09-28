package auth

func IsCreator (creator int,uid int) bool {
	return uid == creator
}

func (cfg Config) IsAuthGroups (gid int64) bool {
	for _,authGid := range cfg.Groups {
		return gid == authGid
	}
	return false
}
