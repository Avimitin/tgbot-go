package net

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type testS struct {
	A string `json:"a"`
	B int    `json:"b"`
}

func Test_JsonPost(t *testing.T) {
	http.HandleFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		reqByte, err := ioutil.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("read body:%v", err)
		}
		var ts testS
		err = json.Unmarshal(reqByte, &ts)
		if err != nil {
			t.Fatalf("Unmarshal:%v", err)
		}
		if ts.A != "test" || ts.B != 114514 {
			t.Errorf("Want A=test and B=114514 Got %+v", ts)
		}
		rw.Write([]byte(`{"ok": true}`))
	})
	go http.ListenAndServe(":11451", nil)

	resp, err := PostJSON("http://127.0.0.1:11451/test", strings.NewReader(`{"a":"test", "b": 114514}`))
	if err != nil {
		t.Fatalf("jsonpost:%v", err)
	}
	if string(resp) != "ok" {
		t.Fatal("get othe message")
	}
}
