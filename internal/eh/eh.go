package eh

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/Avimitin/go-bot/internal/net"
)

const (
	addr = "https://api.e-hentai.org/api.php"
)

func getGid(url string) []string {
	pattern := regexp.MustCompile(`e.hentai\.org/g/(\d+)/(\w+).$`)
	if !pattern.MatchString(url) {
		return nil
	}
	matches := pattern.FindStringSubmatch(url)
	return matches[1:]
}

// gRequest make a request to api server without namespace field.
func gRequest(gid []string) ([]byte, error) {
	if len(gid) < 2 {
		return nil, errors.New("nil gid data")
	}

	requestJsonData := struct {
		Method    string   `json:"method"`
		GIDList   []string `json:"gidlist"`
		Namespace int32    `json:"namespace"`
	}{
		Method:    "gdata",
		GIDList:   gid,
		Namespace: 0,
	}
	request, err := json.Marshal(requestJsonData)
	if err != nil {
		err = fmt.Errorf("marshal %+v:%v", requestJsonData, err)
		return nil, err
	}
	resp, err := net.PostJSON(addr, bytes.NewReader(request))
	if err != nil {
		err = fmt.Errorf("post %v to %s:%v", request, addr, err)
		return nil, err
	}
	return resp, nil
}

func parseResponse(resp []byte) (*GMetaData, error) {
	var err error
	if resp == nil {
		err = fmt.Errorf("get nil response")
		return nil, err
	}
	var gmd *GMetaData
	err = json.Unmarshal(resp, &gmd)
	if err != nil {
		err = fmt.Errorf("unmarshall %s: %v", resp, err)
		return nil, err
	}
	return gmd, nil
}

// GetComic return comic meta data by given url
func GetComic(url string) (*GMetaData, error) {
	if url == "" {
		return nil, errors.New("invalid url input")
	}
	gid := getGid(url)
	resp, err := gRequest(gid)
	if err != nil {
		err = fmt.Errorf("request eh api: %v", err)
		return nil, err
	}
	gmd, err := parseResponse(resp)
	if err != nil {
		err = fmt.Errorf("parse %s:%v", resp, err)
		return nil, err
	}
	return gmd, nil
}
