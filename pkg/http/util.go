package http

import (
	"io"
	"net/http"
)

func Post(url string, reader io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return nil, err
	}
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buff, nil
}
