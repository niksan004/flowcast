package clients

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {
	// TODO add logging
	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
