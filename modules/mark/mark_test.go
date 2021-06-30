package mark

import "testing"

func TestAddMark(t *testing.T) {
	const (
		user        = "example"
		msgLink     = "https://example.com"
		description = "example desc"
	)
	err := AddMark(user, msgLink, description)

	if err != nil {
		t.Fatal(err)
	}

	err = AddMark(user, msgLink, description)
	if err != nil {
		t.Fatal(err)
	}

	if innerDB.amount != 2 {
		t.Errorf("database is not wanted size, got %d", innerDB.amount)
	}

	if length := len(innerDB.db[user]); length != 2 {
		t.Fatalf("no enough data in db, got %d", length)
	}

	for _, entry := range innerDB.db[user] {
		if entry.Link != msgLink {
			t.Errorf("got %s want %s", entry.Link, msgLink)
		}
	}
}

func TestGetMark(t *testing.T) {
	const (
		user = "example"
		link = "test"
		desc = "desc"
	)

	AddMark(user, link, desc)

	links, err := GetMark(user)

	if err != nil {
		t.Fatal(err)
	}

	if linksLength := len(links); linksLength <= 0 {
		t.Errorf("unexpect links size: %d", linksLength)
	}

	for _, result := range links {
		if result.Link == link {
			return
		}
	}

	t.Errorf("can't fetch expected link, got %v", links)

	_, err = GetMark("invalid")
	if err == nil {
		t.Errorf("want error")
	}
}

func TestDelMark(t *testing.T) {
	const (
		user  = "user"
		link  = "link"
		link2 = "link2"
		desc  = "desc"
	)
	AddMark(user, link, desc)
	AddMark(user, link2, desc)

	originalLength := len(innerDB.db[user])
	links, err := GetMark(user)
	if err != nil {
		t.Fatal(err)
	}
	err = DelMark(user, int32(len(links)-1))
	if err != nil {
		t.Errorf("delete user: %v", err)
	}

	if len(innerDB.db[user]) != originalLength-1 {
		t.Errorf("no thing be deleted, length = %d origin = %d", len(innerDB.db[user]), originalLength)
	}
}
