package auth

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func (cfg *Config) SaveConfig (filename string) error {
	yml, err := yaml.Marshal(cfg)
	if err != nil { return err }
	return ioutil.WriteFile(filename, yml, 0600)
}
