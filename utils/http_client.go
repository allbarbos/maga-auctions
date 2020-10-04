package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// HTTPClient is the web client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func MakeRequest(method, uri string, body interface{}) (*http.Request, error) {
	b, err := json.Marshal(body)
	payload := bytes.NewReader(b)

	req, err := http.NewRequest(method, uri, payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	return req, nil
}
