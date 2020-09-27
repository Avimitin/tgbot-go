package auth

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	creator int
	groups []int
}

func Init(cfg *config,filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil { return err }
	return yaml.Unmarshal(file, cfg)
}

func IsCreator(cfg config, uid int) bool {
	return uid == cfg.creator
}

func IsAuthGroups(cfg config, gid int) bool {
	for authGid := range cfg.groups {
		if gid == authGid { return true }
	}
	return false
}
