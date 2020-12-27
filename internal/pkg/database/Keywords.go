package database

import (
	"database/sql"
	"fmt"
)

func PeekKeywords(DB *sql.DB, keyword string) (int, error) {
	rows, err := DB.Query("SELECT kid FROM keywords WHERE keyword=?", keyword)
	if err != nil {
		Pln("Error Occur when preparing sql query. INFO:", err.Error())
		return -1, err
	}
	defer rows.Close()
	var peek int
	for rows.Next() {
		err = rows.Scan(&peek)
	}
	if err != nil {
		Pln("Error Occur when scanning data. INFO:", err.Error())
		return -1, err
	}

	if peek != 0 {
		return peek, nil
	}
	return -1, nil
}

func AddKeywords(DB *sql.DB, keyword string, reply string) (int, error) {
	stmt, err := DB.Prepare("INSERT INTO keywords (keyword) VALUES (?)")
	if err != nil {
		Pln("Error occur when preparing insert. Info:", err.Error())
		return -1, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(keyword)
	if err != nil {
		Pln("Error occur when executive value. Info:", err.Error())
		return -1, err
	}
	ID, _ := result.LastInsertId()
	Pln(fmt.Sprintf("Successfully insert keyword. Insert ID: %v", ID))
	// Then use insert reply.
	return int(ID), SetReply(DB, reply, int(ID))
}

func RenameKeywords(DB *sql.DB, kid int, newKeyword string) error {
	stmt, err := DB.Prepare("UPDATE keywords SET keyword = ? WHERE kid = ?")
	if err != nil {
		Pln("Error occur when preparing update keyword. Info:", err.Error())
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newKeyword, kid)
	if err != nil {
		Pln("Error occur when executive value. Info:", err.Error())
		return err
	}
	row, _ := result.RowsAffected()
	Pln(fmt.Sprintf("Rename keyword successfully. Affected %v rows", row))
	return nil
}

// DelKeyword delete keyword
func DelKeyword(DB *sql.DB, k int) error {
	// Delete keyword
	stmt, err := DB.Prepare("DELETE FROM keywords WHERE kid = ?")
	if err != nil {
		Pln("Error occur when preparing delete keyword. Info:", err.Error())
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(k)
	if err != nil {
		Pln("Error occur when executive value. Info:", err.Error())
		return err
	}

	row, _ := result.RowsAffected()
	Pln(fmt.Sprintf("Delete keyword successfully. Affected %v rows", row))
	return nil
}

type KT struct {
	K string
	I int
}

// FetchKeyword return list of keyword and it's id
func FetchKeyword(DB *sql.DB) (*[]KT, error) {
	rows, err := DB.Query("SELECT kid, keyword FROM keywords")
	if err != nil {
		Pln("Error occur when preparing query. Info:", err.Error())
		return nil, err
	}

	var k KT
	var ks []KT
	for rows.Next() {
		err = rows.Scan(&k.I, &k.K)
		if err != nil {
			Pln("Error occur when scanning value. Info:", err.Error())
			return nil, err
		}
		ks = append(ks, k)
	}
	return &ks, nil
}
