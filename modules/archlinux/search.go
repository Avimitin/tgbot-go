package archlinux

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/Avimitin/go-bot/modules/net"
)

const (
	ARCH_LINUX_SEARCH_BASE_URL = "https://www.archlinux.org/packages/search/json"
	AUR_SEARCH_BASE_URL        = "https://aur.archlinux.org/rpc/"
	SEARCH_LIMIT               = "5"
)

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

func SearchAllAUR(q string) (AURSearchResponse, error) {
	res := AURSearchResponse{}
	v := url.Values{}
	v.Add("v", "5")
	v.Add("type", "search")
	v.Add("by", "name-desc")
	v.Add("arg", q)

	url := fmt.Sprintf("%s?%s", AUR_SEARCH_BASE_URL, v.Encode())
	resp, err := net.Get(url)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return res, fmt.Errorf("Fail to unmarshal data from AUR packages: %w", err)
	}

	return res, err
}

func SearchInfoAUR(q string) (AURInfoResponse, error) {
	res := AURInfoResponse{}
	v := url.Values{}
	v.Add("v", "5")
	v.Add("type", "info")
	v.Add("arg[]", q)

	url := fmt.Sprintf("%s?%s", AUR_SEARCH_BASE_URL, v.Encode())
	resp, err := net.Get(url)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(resp, &res)
	if err != nil {
		return res, fmt.Errorf("Fail to unmarshal data from AUR packages: %w", err)
	}

	return res, err
}
