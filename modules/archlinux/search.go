package archlinux

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Avimitin/go-bot/modules/net"
)

const (
	ARCH_LINUX_SEARCH_BASE_URL = "https://www.archlinux.org/packages/search/json"
	SEARCH_LIMIT               = "5"
)

type SearchResponse struct {
	Version int  `json:"version"`
	Limit   int  `json:"limit"`
	Valid   bool `json:"valid"`
	Results []struct {
		Pkgname        string        `json:"pkgname"`
		Pkgbase        string        `json:"pkgbase"`
		Repo           string        `json:"repo"`
		Arch           string        `json:"arch"`
		Pkgver         string        `json:"pkgver"`
		Pkgrel         string        `json:"pkgrel"`
		Epoch          int           `json:"epoch"`
		Pkgdesc        string        `json:"pkgdesc"`
		URL            string        `json:"url"`
		Filename       string        `json:"filename"`
		CompressedSize int           `json:"compressed_size"`
		InstalledSize  int           `json:"installed_size"`
		BuildDate      time.Time     `json:"build_date"`
		LastUpdate     time.Time     `json:"last_update"`
		FlagDate       interface{}   `json:"flag_date"`
		Maintainers    []string      `json:"maintainers"`
		Packager       string        `json:"packager"`
		Groups         []interface{} `json:"groups"`
		Licenses       []string      `json:"licenses"`
		Conflicts      []interface{} `json:"conflicts"`
		Provides       []interface{} `json:"provides"`
		Replaces       []interface{} `json:"replaces"`
		Depends        []string      `json:"depends"`
		Optdepends     []interface{} `json:"optdepends"`
		Makedepends    []string      `json:"makedepends"`
		Checkdepends   []interface{} `json:"checkdepends"`
	} `json:"results"`
	NumPages int `json:"num_pages"`
	Page     int `json:"page"`
}

type query struct {
	all  string
	name string
}

func fmtURL(query query) string {
	val := url.Values{}
	if query.all != "" {
		val.Add("q", query.all)
	}
	if query.name != "" {
		val.Add("name", query.name)
	}

	val.Add("limit", SEARCH_LIMIT)

	return fmt.Sprintf("%s?%s", ARCH_LINUX_SEARCH_BASE_URL, val.Encode())
}

func SearchAll(q string) (SearchResponse, error) {
	query := query{all: q}
	resp := SearchResponse{}
	err := requestAndParse(&query, &resp)
	return resp, err
}

func SearchName(n string) (SearchResponse, error) {
	query := query{name: n}
	resp := SearchResponse{}
	err := requestAndParse(&query, &resp)
	return resp, err
}

func requestAndParse(q *query, sr *SearchResponse) error {
	resp, err := net.Get(fmtURL(*q))

	if err != nil {
		return fmt.Errorf("Send request to Arch Linux package: %w", err)
	}

	err = json.Unmarshal(resp, &sr)

	if err != nil {
		return fmt.Errorf("Fail to unmarshal data from Arch Linux package search")
	}

	return nil
}
