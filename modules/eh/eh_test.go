package eh

import "testing"

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

func TestGetComic(t *testing.T) {
	gmetadata, err := GetComic(url)
	if err != nil {
		t.Fatal(err)
	}
	got := gmetadata.Medas[0].Gid
	if got != 1803679 {
		t.Errorf("got %v != want, err: %v", got, gmetadata.Medas[0].Error)
	}
}
