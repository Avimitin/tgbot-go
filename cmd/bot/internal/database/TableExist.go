package database

func TableExist(tableName string) (bool, error) {
	db, err := NewDB()
	if err != nil {
		return false, err
	}
	result, err := db.Query("SHOW TABLES")
	if err != nil {
		return false, err
	}
	for result.Next() {
		var table string
		err = result.Scan(&table)
		if err != nil {
			return false, err
		}

		if tableName == table {
			return true, nil
		}
	}
	return false, nil
}
