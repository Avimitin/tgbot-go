package auth

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	LOADED		bool		`yaml:"LOADED"`
	Groups   	[]int64		`yaml:"groups"`
	BotToken 	string		`yaml:"bot_token"`
}

func NewCFG() (Config, error) {
	cfg := Config{}
	err := decode(&cfg, "F:\\go-bot\\cfg\\auth.yml")
	return cfg, err
}

func decode(cfg *Config,filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil { return err }
	return yaml.Unmarshal(file, cfg)
}

