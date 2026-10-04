package agent

import (
	"io"
	"net/http"
	"project/internal/protocol"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

type httpAc struct {
	protocol.HttpArgs
}

func (action *httpAc) Run() (protocol.Result, error) {
	// build request
	req, err := http.NewRequest(action.Method, action.Url, strings.NewReader(action.Body))
	if err != nil {
		return &protocol.HttpResult{}, err
	}
	req.Header = action.Header

	// send request
	resp, err := httpClient.Do(req)
	if err != nil {
		return &protocol.HttpResult{}, err
	}
	defer resp.Body.Close()

	// read response body limited to 1 MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return &protocol.HttpResult{}, err
	}

	return &protocol.HttpResult{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       string(body),
	}, nil
}
