package browser

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

// Browse use given url to fetch data
func Browse(url string) (string, error) {
	client := NewClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[Browse]Error occur when making request.")
		return "", err
	}
	req.Header.Add("User-Agent", UA)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[Browse]Error occur when posting request.")
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[Browse]Error occur when reading response data.")
		return "", err
	}
	return string(body), nil
}

// JsonRequest POST a request to the given url.
func JsonRequest(url string, data []byte) ([]byte, error) {
	client := NewClient()
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		log.Println("[JsonRequest]Error occur when making post request.")
		return nil, err
	}
	req.Header.Set("Content-Type", "Application/json")
	req.Header.Set("User-Agent", UA)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[JsonRequest]Error occur when posting request.")
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[JsonRequest]Error occur when reading response data.")
		return nil, err
	}
	return body, nil
}
