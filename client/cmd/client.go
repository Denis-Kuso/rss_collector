package cmd

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("Connection error")
	ErrNotFound        = errors.New("Not found")
	ErrInvalidResponse = errors.New("Invalid server response")
	ErrInvalidRequest  = errors.New("Malformed request")
)

func newClient() *http.Client {
	const TIMEOUT = 5
	c := &http.Client{
		Timeout: TIMEOUT * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   time.Second,
			ResponseHeaderTimeout: time.Second,
		},
	}
	return c
}

func sendReq(url, method, apiKey, contentType string, expStatus int, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if apiKey != "" {
		req.Header.Add("Authorization", "ApiKey "+apiKey)
	}

	r, err := newClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read body: %w", err)
	}
	if r.StatusCode != expStatus {
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	return msg, nil
}
