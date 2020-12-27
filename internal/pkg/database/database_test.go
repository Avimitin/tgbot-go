package database

import (
	"database/sql"
	"fmt"
	"github.com/Avimitin/go-bot/internal/pkg/conf"
	"testing"
)

func getDB(t *testing.T) *sql.DB {
	data := conf.LoadDBSecret(conf.WhereCFG(""))
	testDB, err := NewDB(data)
	if err != nil {
		t.Fatal(err)
	}
	return testDB
}

func TestNewDB(t *testing.T) {
	err := getDB(t).Ping()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestTableExist(t *testing.T) {
	isExist, err := TableExist(getDB(t), "user")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Printf("isExist: %v\n", isExist)
}

func TestGetAdmin(t *testing.T) {
	admins, err := GetAdmin(getDB(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, admin := range *admins {
		if admin == 649191333 {
			return
		}
	}
	t.Logf("Admin 649191333 not found")
}

func TestAddGroups(t *testing.T) {
	err := AddGroups(getDB(t), 123456789, "TestGroup")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteGroups(t *testing.T) {
	// Delete not exist groups
	err := DeleteGroups(getDB(t), 131231)
	if err != nil {
		t.Fatal(err)
	}
	// Delete exist groups
	err = DeleteGroups(getDB(t), 123456789)
	if err != nil {
		t.Fatal(err)
	}
	groups, err := SearchGroups(getDB(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range groups {
		if group.GroupID == 123456789 {
			t.Fatal("Unwanted group still exist")
		}
	}
}

func TestSearchGroups(t *testing.T) {
	groups, err := SearchGroups(getDB(t))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(groups)
}

func TestChangeGroupsName(t *testing.T) {
	err := ChangeGroupsName(getDB(t), 123456789, "testGroups2")
	if err != nil {
		t.Fatal(err)
	}
	groups, err := SearchGroups(getDB(t))
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range groups {
		if group.GroupID == 123456789 && group.GroupUsername != "testGroups2" {
			t.Fatal("GroupName change fail")
		}
	}
}

func TestPln(t *testing.T) {
	Pln("THIS IS A TESTING MESSAGE", "THIS IS ERROR MESSAGE")
}

func TestAddKeywords(t *testing.T) {
	ID, err := AddKeywords(getDB(t), "test", "testReply")
	if ID == -1 || err != nil {
		t.Fail()
	}
}

func TestRenameKeywords(t *testing.T) {
	id, _ := PeekKeywords(getDB(t), "test")
	err := RenameKeywords(getDB(t), id, "test1")
	if err != nil {
		t.Fail()
	}
}

func TestPeekKeywords(t *testing.T) {
	id, err := PeekKeywords(getDB(t), "test1")
	if err != nil {
		t.Fatal(err)
	}
	if id < 0 {
		t.Fatal("Can't fetch wanted keyword")
	}
}

func TestGetReplyWithKey(t *testing.T) {
	id, err := PeekKeywords(getDB(t), "test1")
	replies, err := GetReplyWithKey(getDB(t), id)
	if err != nil {
		t.Fatal(err)
	}
	for _, reply := range replies {
		if reply == "testReply" {
			return
		}
	}
	t.Fatal("Can't got wanted replies")
}

func TestSetReplyAndPeekReply(t *testing.T) {
	id, err := PeekKeywords(getDB(t), "test1")
	err = SetReply(getDB(t), "testReply2", id)
	if err != nil {
		t.Fatal(err)
	}
	if id, err = PeekReply(getDB(t), "testReply2"); id == -1 || err != nil {
		t.Fatal("Fail to get reply")
	}
}

func TestDelKeywordAndReply(t *testing.T) {
	id, err := PeekKeywords(getDB(t), "test1")
	err = DelReplyByKeyword(getDB(t), id)
	err = DelKeyword(getDB(t), id)
	if err != nil {
		t.Fail()
	}
	if id, err = PeekKeywords(getDB(t), "test1"); id != -1 || err != nil {
		t.Fail()
	}
}

func TestFetchKeyword(t *testing.T) {
	k, err := FetchKeyword(getDB(t))
	if err != nil {
		t.Fail()
	}
	fmt.Println(k)
}
