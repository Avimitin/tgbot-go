package auth

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type dbCfg struct {
	User     string `yaml:"user"`
	Password string `yaml:password`
	Host     string `yaml:"host"`
	Database string `yaml:database`
}

type Config struct {
	LOADED   bool    `yaml:"LOADED"`
	Groups   []int64 `yaml:"groups"`
	BotToken string  `yaml:"bot_token"`
	DBCfg    dbCfg   `yaml:"db_cfg"`
}

func NewCFG() (Config, error) {
	cfg := Config{}
	err := decode(&cfg, "F:\\go-bot\\cfg\\auth.yml")
	return cfg, err
}

func decode(cfg *Config, filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, cfg)
}
