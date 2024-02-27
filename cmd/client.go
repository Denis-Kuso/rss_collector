package cmd

import (
  "errors"
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
)

func newClient() *http.Client {
  c := &http.Client{
    Timeout: TIMEOUT * time.Second,
  }
  return c
}
