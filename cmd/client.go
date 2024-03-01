package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnection      = errors.New("Connection error")
	ErrNotFound        = errors.New("Not found")
	ErrInvalidResponse = errors.New("Invalid server response")
)

const (
	TIMEOUT  = 3
	ROOT_URL = "http://www.localhost:8080" //TODO: change
)

func newClient() *http.Client {
	c := &http.Client{
		Timeout: TIMEOUT * time.Second,
	}
	return c
}

var c *http.Client = newClient()

func fetchEndpoint(c *http.Client, endpoint string) ([]byte, error) {

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s\n", ErrConnection, err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, ErrConnection
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	return data, err
}

func sendReq(url, method, contentType string, expStatus int, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	// set header
	req.Header.Add("Authorization: apikey", "ApiKeyFromUser") //TODO HOW IS APIKEY provided?
	r, err := newClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("Cannot read body: %w\n", err)
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
