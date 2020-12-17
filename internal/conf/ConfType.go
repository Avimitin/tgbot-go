package conf

import (
	"fmt"
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
