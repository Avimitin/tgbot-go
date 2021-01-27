package conf

import (
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
}
