package browser

import (
	"net/http"
	"time"
)

// UA is user-agent
const (
	UA string = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"
)

func NewClient() *http.Client {
	return NewClientWithTimeout(60 * time.Second)
}

func NewClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}
