package conf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// DBSecret store database DSN information
type DBSecret struct {
	user     string
	pwd      string
	host     string
	database string
	port     string
}

// MySqlURL return formatted mysql tcp link
func (db *DBSecret) MySqlURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.user, db.pwd, db.host, db.port, db.database)
}

type Configuration struct {
	CertGroup map[int64]interface{}
	Users     map[int]string
}

// NewConfiguration return initialized configuration
func NewConfiguration(cfgPath string) (*Configuration, error) {
	files, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("can't read configuration JSON file:\n%v", err)
	}
	return readConfigurationFromBytes(bytes.NewReader(files))
}

func readConfigurationFromBytes(file io.Reader) (*Configuration, error) {
	fileByte, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("fail to read file bytes:\n%v", err)
	}
	configJSONStruct := struct {
		CertGroup []int64        `json:"cert_group"`
		Users     map[int]string `json:"users"`
	}{}
	err = json.Unmarshal(fileByte, &configJSONStruct)
	if err != nil {
		return nil, fmt.Errorf("Fail to unmarshal configuration data:\n%v", err)
	}
	if len(configJSONStruct.CertGroup) == 0 || configJSONStruct.Users == nil {
		return nil, fmt.Errorf("fail to initialize configuration file")
	}
	cfg := new(Configuration)
	cfg.CertGroup = make(map[int64]interface{})
	for _, group := range configJSONStruct.CertGroup {
		cfg.CertGroup[group] = struct{}{}
	}
	cfg.Users = configJSONStruct.Users
	return cfg, nil
}
