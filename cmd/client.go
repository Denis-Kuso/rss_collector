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

const (
	ROOT_URL = "http://localhost:8080/v1" //TODO: change
)

func newClient() *http.Client {
	const TIMEOUT  = 5
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

var c *http.Client = newClient()

func fetchEndpoint(c *http.Client, endpoint string) ([]byte, error) {

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s\n", ErrConnection, err)
	}
	fmt.Printf("Sending request with url: %v\n", endpoint)
	resp, err := c.Do(req)
	if err != nil {
		return nil, ErrConnection
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	fmt.Printf("Called url: %v using post,got: %v\n", endpoint, resp.StatusCode)
	return data, err
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
		req.Header.Add("Authorization", "ApiKey "+apiKey) //TODO will default header allow this?
	}

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
		fmt.Printf("Got status: %v, expected: %v\n", r.StatusCode, expStatus)
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}

	return msg, nil
}
