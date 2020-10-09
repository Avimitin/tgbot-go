package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/cmd/bot/internal/conf"
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

func AddKeywords(DB *sql.DB, keyword string, reply string) error {
	// peek if keyword has existed or not.
	kid, err := PeekKeywords(DB, keyword)
	if err != nil {
		return err
	}

	// if keyword has existed, add new reply.
	if kid != -1 {
		return SetReply(DB, reply, kid)
	}

	// If keyword isn't exist, insert keyword first.
	stmt, err := DB.Prepare("INSERT INTO keywords (keyword) VALUES (?)")
	if err != nil {
		Pln("Error occur when preparing insert. Info:", err.Error())
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(keyword)
	if err != nil {
		Pln("Error occur when executive value. Info:", err.Error())
		return err
	}
	ID, _ := result.LastInsertId()
	Pln(fmt.Sprintf("Successfully insert keyword. Insert ID: %v", ID))
	// Then insert reply.
	return SetReply(DB, reply, kid)
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

func DelKeyword(DB *sql.DB, k int) error {
	// Delete reply relate with keyword first.
	stmt, err := DB.Prepare("DELETE FROM replies WHERE kid = ?")
	if err != nil {
		Pln("Error occur when preparing delete keyword. Info:", err.Error())
		return err
	}
	result, err := stmt.Exec(k)
	if err != nil {
		Pln("Error occur when executive value. Info:", err.Error())
		return err
	}

	row, _ := result.RowsAffected()
	Pln(fmt.Sprintf("Delete replies successfully. Affected %v rows", row))

	// Delete keyword
	stmt, err = DB.Prepare("DELETE FROM keywords WHERE kid = ?")
	if err != nil {
		Pln("Error occur when preparing delete keyword. Info:", err.Error())
		return err
	}
	defer stmt.Close()

	result, err = stmt.Exec(k)
	if err != nil {
		Pln("Error occur when executive value. Info:", err.Error())
		return err
	}

	row, _ = result.RowsAffected()
	Pln(fmt.Sprintf("Delete keyword successfully. Affected %v rows", row))
	return nil
}

func FetchKeyword(DB *sql.DB) ([]conf.KeywordType, error) {
	rows, err := DB.Query("SELECT kid,keyword FROM keywords")
	if err != nil {
		Pln("Error occur when preparing query. Info:", err.Error())
		return nil, err
	}
	var k conf.KeywordType
	var keywords []conf.KeywordType

	for rows.Next() {
		err = rows.Scan(&k.Kid, &k.Word)
		if err != nil {
			Pln("Error occur when scanning value. Info:", err.Error())
			return nil, err
		}
		keywords = append(keywords, k)
	}

	return keywords, nil
}
