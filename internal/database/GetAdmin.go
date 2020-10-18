package database

import (
	"database/sql"
	"fmt"
	"log"
)

func GetAdmin(db *sql.DB, groupName string) ([]int, error) {
	sqlQuery := fmt.Sprintf(`SELECT uid FROM %v WHERE permission="admin"`, groupName)

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Printf("[DATABASE]Error happen when preparing sql query. Descriptions: %v", err)
		return nil, err
	}
	defer stmt.Close()

	response, err := stmt.Query()
	if err != nil {
		log.Printf("[DATABASE]Error happen when get response. Descriptions: %v", err)
		return nil, err
	}

	var adminList []int
	for response.Next() {
		var result int
		err := response.Scan(&result)
		if err != nil {
			log.Printf("[DATABASE]Error happen when parsing result. Descriptions: %v", err)
			return nil, err
		}
		adminList = append(adminList, result)
	}

	return adminList, nil
}
