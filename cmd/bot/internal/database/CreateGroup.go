package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Use this methods to create a new table
func CreateGroup(db *sql.DB, groupName string) error {
	sqlQuery := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %v (
		uid 		INT NOT NULL,
		username 	VARCHAR(255) NOT NULL, 
		permission  VARCHAR(10) NOT NULL DEFAULT 'member',
		PRIMARY KEY (uid)
		) DEFAULT CHARSET=utf8`, groupName)

	stmt, err := db.Prepare(sqlQuery)
	defer stmt.Close()
	if err != nil {
		log.Printf("[Database]%v\n", err)
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Printf("[Database]%v\n", err)
		return err
	}

	return nil
}
