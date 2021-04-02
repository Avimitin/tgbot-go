package net

import (
	"encoding/json"
	"strings"
	"testing"
)

type testS struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func Test_JsonPost(t *testing.T) {
	artical := `{"title": "foo", "body": "bar", "userId": 10}`
	resp, err := PostJSON("https://jsonplaceholder.typicode.com/posts", strings.NewReader(artical))
	if err != nil {
		t.Fatalf("jsonpost:%v", err)
	}
	var ts *testS
	err = json.Unmarshal(resp, &ts)
	if err != nil {
		t.Fatal(err)
	}
	if ts.ID != 101 {
		t.Errorf("Want id = 101 got %+v", ts)
	}
}
