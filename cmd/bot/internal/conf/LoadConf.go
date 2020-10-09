package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadCFG() (*Config, error) {
	cfg := &Config{}
	err := decode(cfg, "F:\\go-bot\\cfg\\auth.yml")
	return cfg, err
}

func decode(cfg *Config, filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, cfg)
}
