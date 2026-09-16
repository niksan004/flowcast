package clients

import (
	"net/http"
	"time"
)

type HTTPClient struct {
	client *http.Client
}

func NewHTTPClient() *HTTPClient {

	return &HTTPClient{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}
