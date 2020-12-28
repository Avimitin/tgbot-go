package osuAPI

import (
	"errors"
	"github.com/Avimitin/go-bot/internal/pkg/browser"
	jsoniter "github.com/json-iterator/go"
	"log"
	"regexp"
)

const (
	std = iota
	taiko
	ctb
	mania
	apiAddr     = "https://osu.ppy.sh/api/"
	defaultMode = -1
)

var (
	InvalidParamsErr  = errors.New("invalid params input")
	InvalidUsrTypeErr = errors.New("invalid user's type")
	InvalidModeErr    = errors.New("invalid mode")
	InvalidUsrNameErr = errors.New("invalid user name")
)

func buildURL(method string, params map[string]string) string {
	if params == nil {
		return ""
	}
	url := apiAddr + method + "?"
	for param, value := range params {
		url += param + "=" + value + "&"
	}
	return url[:len(url)-1]
}

func apiRequest(method string, params map[string]string) ([]byte, error) {
	if params == nil {
		return nil, InvalidParamsErr
	}
	url := buildURL(method, params)
	return browser.Browse(url)
}

func getBeatMapsBySpecificMapID(key string, mapID string) []Beatmap {
	data, err := apiRequest("get_beatmaps", map[string]string{"k": key, "b": mapID})
	if err != nil {
		log.Println("[getBeatMaps]Error occur when making request to osu api:", err)
		return nil
	}
	var beatmaps []Beatmap
	err = jsoniter.Unmarshal(data, &beatmaps)
	if err != nil {
		log.Println("[getBeatMaps]Error occur when unmarshalling data:", err)
		return nil
	}
	return beatmaps
}

// GetBeatMapByBeatMapSet return set of beatmaps by given setID
func GetBeatMapByBeatMapSet(key string, setID string) []Beatmap {
	data, err := apiRequest("get_beatmaps", map[string]string{"k": key, "s": setID})
	if err != nil {
		log.Println("[getBeatMaps]Error occur when making request to osu api:", err)
		return nil
	}
	var beatmaps []Beatmap
	err = jsoniter.Unmarshal(data, &beatmaps)
	if err != nil {
		log.Println("[getBeatMaps]Error occur when unmarshalling data:", err)
		return nil
	}
	return beatmaps
}

// GetBeatMap return a map by given specific beatmapset_id and beatmap_id
func GetBeatMap(key string, mapID string) *Beatmap {
	beatmaps := getBeatMapsBySpecificMapID(key, mapID)
	if beatmaps == nil {
		log.Println("[GetBeatMap]Get nil beatmaps")
		return nil
	}
	return &(beatmaps[0])
}

// GetBeatMapByURL return a map with given url
func GetBeatMapByURL(key string, url string) *Beatmap {
	pattern := regexp.MustCompile(`osu\.ppy\.sh/beatmapsets/[0-9]+#[a-z]+/([0-9]+).*`)
	if !pattern.MatchString(url) {
		return nil
	}
	matches := pattern.FindStringSubmatch(url)[1]
	return GetBeatMap(key, matches)
}

//-------------
// USER API
//-------------

func getUserWithMode(key string, user string, userType string, mode string) []User {
	var params map[string]string
	if mode == "default" {
		params = map[string]string{"k": key, "u": user, "type": userType}
	} else {
		params = map[string]string{"k": key, "u": user, "type": userType, "m": mode}
	}
	data, err := apiRequest("get_user", params)
	if err != nil {
		log.Println("[GetUser]Error occur when making request:", err)
		return nil
	}
	var users []User
	err = jsoniter.Unmarshal(data, &users)
	if err != nil {
		log.Println("[GetUser]Error occur when unmarshalling data")
		return nil
	}
	return users
}

// GetUser return a user with given user's username or user_id,
// require userType to specify if u is a user_id or a username.
// Use "string" for usernames or "id" for user_ids.
func GetUser(key string, user string, userType string, mode int) (*User, error) {
	if userType != "string" && userType != "id" {
		return nil, InvalidUsrTypeErr
	}
	var m string
	switch mode {
	case defaultMode:
		m = "default"
	case std:
		m = "0"
	case taiko:
		m = "1"
	case ctb:
		m = "2"
	case mania:
		m = "3"
	default:
		return nil, InvalidModeErr
	}
	users := getUserWithMode(key, user, userType, m)
	if len(users) == 0 {
		return nil, InvalidUsrNameErr
	}
	return &(users[0]), nil
}
