package browser

import (
	"io/ioutil"
	"net/http"
)

// Browse use given url to fetch data
func Browse(url string) (string, error) {
	client := newClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", UA)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), nil
}
