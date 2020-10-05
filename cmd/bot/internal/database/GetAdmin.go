package database

import (
	"fmt"
)

func GetAdmin(groupName string) ([]int, error) {
	db, err := NewDB()
	if err != nil {
		return nil, err
	}

	sqlQuery := fmt.Sprintf(`SELECT uid FROM %v WHERE permission="admin"`, groupName)

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	response, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var adminList []int
	for response.Next() {
		var result int
		err := response.Scan(&result)
		if err != nil {
			return nil, err
		}
		adminList = append(adminList, result)
	}

	return adminList, nil
}
