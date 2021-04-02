package net

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	UA = "Mozilla/5.0 (X11; Linux x86_64)" +
		" AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
)

// PostJSON post given data with "Application/json" header,
// return []byte if make request successfully.
func PostJSON(url string, data io.Reader) ([]byte, error) {
	client := http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return nil, fmt.Errorf("post json to %s: %v", url, err)
	}
	req.Header.Set("Content-Type", "Application/json")
	req.Header.Set("User-Agent", UA)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do json request: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s:%v", body, err)
	}
	return body, nil
}

// Get make request to the given url, return response body.
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get %s: %v", url, err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %q:%v", resp.Body, err)
	}
	return body, nil
}
