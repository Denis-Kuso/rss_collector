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

func sendReq(url, method, contentType string, expStatus int, body io.Reader) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	// set header
	req.Header.Add("Authorization: apikey", "ApiKeyFromUser") //TODO HOW IS APIKEY provided?
	r, err := newClient().Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != expStatus {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s", err, msg)
	}
	return nil
}
