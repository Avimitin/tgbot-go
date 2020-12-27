package database

import (
	"database/sql"
	"log"
)

func TableExist(db *sql.DB, tableName string) (bool, error) {
	result, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Printf("[DATABASE]Error fetching tables. Descriptions: %v", err)
		return false, err
	}

	for result.Next() {
		var table string
		err = result.Scan(&table)
		if err != nil {
			log.Printf("[DATABASE]Error happen when parsing result. Descriptions: %v", err)
			return false, err
		}

		if tableName == table {
			return true, nil
		}
	}
	return false, nil
}
