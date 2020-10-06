package database

import (
	"database/sql"
	"log"
)

// AddGroups support add authed groups into database.
// It require database's connection, group's ID, group's username.
func AddGroups(db *sql.DB, groupID int64, username string) error {
	stmt, err := db.Prepare("INSERT INTO authgroups (GroupID, GroupUsername) VALUES (?, ?)")
	if err != nil {
		log.Printf("[DATABASE]Error occur when prepare SQL query. Descriptions: %v", err)
		return err
	}
	// 将连接丢回连接池
	defer stmt.Close()

	result, err := stmt.Exec(groupID, username)
	if err != nil {
		log.Printf("[DATABASE]Error occur execute value. Descriptions: %v", err)
		return err
	}

	ID, _ := result.LastInsertId()

	log.Printf("[DATABASE]Successfully insert auth groups, Insert ID: %v\n", ID)
	return nil
}

// DeleteGroups delete the groups record from database.
// Require database connection, group's ID.
func DeleteGroups(db *sql.DB, groupID int64) error {
	stmt, err := db.Prepare("DELETE FROM authgroups WHERE GroupID=?")
	if err != nil {
		log.Printf("[DATABASE]Error occur when prepare SQL query. Descriptions: %v", err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(groupID)
	if err != nil {
		log.Printf("[DATABASE]Error occur execute value. Descriptions: %v", err)
		return err
	}

	ID, err := result.RowsAffected()
	log.Printf("[DATABASE]Successfully delete group's record. Affect %v row.\n", ID)
	return nil
}

// Group's data structure
type group struct {
	GroupID       int64
	GroupUsername string
}

// SearchGroups return a list with all the record in authGroups table.
func SearchGroups(db *sql.DB) ([]group, error) {
	rows, err := db.Query("SELECT GroupID,GroupUsername FROM authgroups")
	if err != nil {
		log.Println("[DATABASE]Error occur when fetching result")
		return nil, err
	}
	defer rows.Close()

	var Groups []group
	for rows.Next() {
		var g group
		err := rows.Scan(&g.GroupID, &g.GroupUsername)
		if err != nil {
			log.Printf("[DATABASE]Error occur when parsing data. Descriptions: %v\n", err)
			return nil, err
		}
		Groups = append(Groups, g)
	}
	return Groups, nil
}

// ChangeGroupsName require original group's ID and new group name to change group name.
func ChangeGroupsName(db *sql.DB, originGroupID int64, groupNameAfter string) error {
	stmt, err := db.Prepare("UPDATE authgroups SET GroupUsername=? WHERE GroupID=?")
	if err != nil {
		log.Printf("[DATABASE]Error occur when prepare SQL query. Descriptions: %v", err)
		return nil
	}
	defer stmt.Close()

	result, err := stmt.Exec(groupNameAfter, originGroupID)
	if err != nil {
		log.Printf("[DATABASE]Error occur when execute value. Descriptions: %v", err)
		return err
	}
	counts, err := result.RowsAffected()
	if err != nil {
		log.Printf("[DATABASE]Error occur when fetching affected rows. Descriptions: %v", err)
		return err
	}
	log.Printf("[DATABASE]Successfully delete records. Affected %v rows", counts)
	return nil
}
