package database

import (
	"database/sql"
)

func PeekKeywords(DB *sql.DB, keyword string) (bool, error) {
	rows, err := DB.Query("SELECT kid FROM keywords WHERE keywords=?", keyword)
	if err != nil {
		Pln("Error Occur when preparing sql query. INFO:", err.Error())
		return false, err
	}
	defer rows.Close()
	var peek int
	for rows.Next() {
		err = rows.Scan(&peek)
	}
	if err != nil {
		Pln("Error Occur when scanning data. INFO:", err.Error())
		return false, err
	}

	if peek != 0 {
		return true, nil
	}
	return false, nil
}
