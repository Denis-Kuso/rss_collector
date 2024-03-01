package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
  ErrConnection = errors.New("Connection error")
  ErrNotFound = errors.New("Not found")
  ErrInvalidResponse = errors.New("Invalid server response")
)

const (
  TIMEOUT = 3
  URL = "http://www.localhost:8080"//TODO: change
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
