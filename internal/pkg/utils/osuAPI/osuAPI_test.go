package osuAPI

import (
	"fmt"
	"github.com/Avimitin/go-bot/internal/pkg/conf"
	"testing"
)

func TestBuildURL(t *testing.T) {
	url := buildURL("get_beatmaps", map[string]string{"k": "api", "u": "Doormat", "type": "string"})
	if url != "https://osu.ppy.sh/api/get_beatmaps?k=api&u=Doormat&type=string" {
		t.Fatal("Get != Want")
	}
}

func TestGetBeatMap(t *testing.T) {
	bm := GetBeatMap(conf.LoadOSUAPI(conf.WhereCFG("")), "1658251")
	if bm == nil {
		t.Fatal("Got nil beatmap")
	}
	if bm.Creator != "Doormat" {
		t.Fatalf("Got unwanted creator")
	}
	fmt.Printf("%+v\n", bm)
}

func TestGetUser(t *testing.T) {
	user, err := GetUser(conf.LoadOSUAPI(conf.WhereCFG("")), "avimitin", "string", std)
	if err != nil {
		t.Fatal(err)
	}
	if user == nil {
		t.Fatal("Got nil user")
	}
	if user.Username != "avimitin" {
		t.Fatal("Got unexpected user")
	}
	fmt.Printf("%+v\n", user)
}

func TestGetBeatMapByURL(t *testing.T) {
	bm := GetBeatMapByURL(conf.LoadOSUAPI(conf.WhereCFG("")), "https://osu.ppy.sh/beatmapsets/34348#osu/111680")
	if bm == nil {
		t.Fatal("Can't recognize url")
	}
	if bm.BeatmapsetID != "34348" {
		t.Fatal("Can't get beatmap")
	}
}
