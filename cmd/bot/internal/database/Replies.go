package database

import (
	"database/sql"
	"fmt"
)

func SetReply(DB *sql.DB, reply string, kid int) error {
	stmt, err := DB.Prepare("INSERT INTO replies (reply, kid) VALUES (?, ?)")
	if err != nil {
		Pln("Error occur when preparing replies. Info:", err.Error())
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(reply, kid)
	if err != nil {
		Pln("Error occur when execute value. Info:", err.Error())
		return err
	}
	row, _ := result.RowsAffected()
	if row > 0 {
		Pln(fmt.Sprintf("Successfully add new reply, Affected %v row.", row))
		return nil
	}
	Pln("Insertion did not affect any rows.")
	return nil
}

func GetReplyWithKey(DB *sql.DB, kid int) ([]string, error) {
	rows, err := DB.Query("SELECT reply FROM replies WHERE kid = ?", kid)
	if err != nil {
		Pln("Error occur when searching reply. Info", err.Error())
		return nil, err
	}

	var reply string
	var replies []string

	for rows.Next() {
		err := rows.Scan(&reply)
		if err != nil {
			Pln("Error occur when scanning result. Info", err.Error())
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}

func PeekReply(DB *sql.DB, reply string) (int, error) {
	rows, err := DB.Query("SELECT rid FROM replies WHERE reply = ?", reply)
	if err != nil {
		Pln("Error occur when fetching reply. Info", err.Error())
		return -1, err
	}
	var rid int
	for rows.Next() {
		err = rows.Scan(&rid)
		if err != nil {
			Pln("Error occur when scanning result. Info", err.Error())
			return -1, err
		}
	}
	return rid, nil
}
