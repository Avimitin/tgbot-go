package conf

import (
	"encoding/json"
	"fmt"
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
	CertGroup []int64        `json:"cert_group"`
	Users     map[int]string `json:"users"`
}

func NewConfiguration(cfgPath string) (*Configuration, error) {
	files, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("can't read configuration JSON file: %v", err)
	}
	var cfg *Configuration
	err = json.Unmarshal(files, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Fail to unmarshal configuration data: %v", err)
	}
	return cfg, nil
}
