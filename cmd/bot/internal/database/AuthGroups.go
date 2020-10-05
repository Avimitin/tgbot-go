package database

import "log"

func AddGroups(groupID int64, username string) error {
	db, err := NewDB()
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO authgroups (GroupID, GroupUsername) VALUES (?, ?)")
	if err != nil {
		return err
	}
	// 将连接丢回连接池
	defer stmt.Close()
	result, err := stmt.Exec(groupID, username)
	if err != nil {
		return err
	}
	ID, _ := result.LastInsertId()
	log.Printf("[DATABASE]Successfully Insert, Insert ID: %v\n", ID)
	return nil
}
