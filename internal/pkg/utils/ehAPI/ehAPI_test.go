package ehAPI

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const (
	url = "https://exhentai.org/g/1803679/a0e2c7c327/"
)

func TestGetGid(t *testing.T) {
	gid := getGid(url)
	if len(gid) < 2 {
		t.Fatal("gid length is less than 2")
	}
	if gid[0] != "1803679" {
		t.Fatalf("Got unexpected gid result: %s", gid[0])
	}
	if gid[1] != "a0e2c7c327" {
		t.Fatalf("Got unexpected gid result: %s", gid[1])
	}
}

func TestGRequest(t *testing.T) {
	resp, err := gRequest([][]string{getGid(url)}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(resp), "gmetadata") {
		t.Fatal("Got unwanted result")
	}
}

func TestParseRequest(t *testing.T) {
	resp, err := gRequest([][]string{getGid(url)}, 0)
	if err != nil {
		t.Fatal(err)
	}
	gmd, err := parseResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if len(gmd.GMD) == 0 {
		t.Fatal("Got nil g meta data")
	}
	if gmd.GMD[0].Gid != 1803679 {
		t.Fatal("Got unwanted gid")
	}
	fmt.Printf("%+v\n", gmd.GMD[0])
	resp, err = gRequest([][]string{{"213131", "abcdefg"}}, 0)
	if err != nil {
		t.Fatal(err)
	}
	gmd, err = parseResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if gmd.GMD[0].Error == "" {
		t.Fatal("Parse error failed")
	}
	fmt.Println("Error:", gmd.GMD[0].Error)
}

func isSameType(a interface{}, b string) bool {
	return reflect.TypeOf(a).String() == b
}

func TestGetComic(t *testing.T) {
	gmetadata, err := GetComic([]string{url}, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, data := range gmetadata.GMD {
		if !isSameType(data.Gid, "int64") {
			t.Fatalf("Got unwanted type: %v", reflect.TypeOf(data.Gid))
		}
	}
	fmt.Printf("%+v\n", gmetadata.GMD[0])
}
