package ehAPI

import (
	"github.com/Avimitin/go-bot/internal/net/browser"
	jsoniter "github.com/json-iterator/go"
	"log"
	"regexp"
	"strings"
)

const (
	apiAddr = "https://api.e-hentai.org/api.php"
)

func getGid(url string) []string {
	urlB := []byte(url)
	pattern := regexp.MustCompile(`https://e.hentai\.org/g/(?P<gid1>\d*)/(?P<gid2>\w*).$`)
	temp := []byte(`$gid1 $gid2`)
	var result []byte
	for _, matches := range pattern.FindAllSubmatchIndex(urlB, -1) {
		result = pattern.Expand(result, temp, urlB, matches)
	}
	return strings.Fields(string(result))
}

// gRequest make a request to api server without namespace field.
func gRequest(gid GIDListType, ns int32) ([]byte, error) {
	if len(gid) == 0 {
		return nil, NilGidErr
	}
	for _, g := range gid {
		if len(g) != 2 {
			return nil, WrongUrlErr
		}
	}
	r := NewGRequest(gid, ns)
	gDataRequestByte, err := jsoniter.Marshal(r)
	if err != nil {
		log.Println("[nnsRequest]Error occur when marshalling data.")
		return nil, err
	}
	resp, err := browser.JsonRequest(apiAddr, gDataRequestByte)
	if err != nil {
		log.Println("[nnsRequest]Error occur when posting data.")
		return nil, err
	}
	return resp, nil
}

func parseResponse(resp []byte) (*GMetaData, error) {
	if resp == nil {
		log.Println("[parseResponse]Receive nil byte data")
		return nil, CantParseNilDataErr
	}
	var gmd GMetaData
	err := jsoniter.Unmarshal(resp, &gmd)
	if err != nil {
		log.Println("[parseResponse]Error occur when unmarshalling data")
		return nil, err
	}
	return &gmd, nil
}

func GetComic(exGalleryUrls []string, ns int32) (*GMetaData, error) {
	if exGalleryUrls == nil {
		return nil, NoneUrlErr
	}
	var gid GIDListType
	for _, url := range exGalleryUrls {
		gid = append(gid, getGid(url))
	}
	resp, err := gRequest(gid, ns)
	if err != nil {
		log.Println("[GetComic]Error occur when making request.", err)
		return nil, err
	}
	gmd, err := parseResponse(resp)
	if err != nil {
		log.Println("[GetComic]Error occur when parsing response.", err)
		return nil, err
	}
	return gmd, nil
}
