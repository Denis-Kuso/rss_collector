package fetch

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const MAX_TIMEOUT = 3

type Client struct {
	httpClient http.Client
	apiKey     string
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{
			Timeout: MAX_TIMEOUT * time.Second,
		},
	}
}

var c Client = NewClient()

func fetchEndpoint(c *Client, url string) ([]byte, error) {

	const ERROR_THRESHOLD = 399 // arbitrary for now
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("Invalid response, check your connection.\n")
	}
	defer resp.Body.Close()
	if resp.StatusCode > ERROR_THRESHOLD {
		return nil, fmt.Errorf("Failed on %v; response received:%v\n", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	return data, err
}

func CreateUser(username string) (apiKey string, err error) {
	// request to server at createUser
	URL := "address-to-serve"
	ENDPOINT := "someendpoint"
	data, err := fetchEndpoint(&c, URL+ENDPOINT)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
		return "", err
	}
	// extract apiKey from response
	// save api key/display apiKey
	fmt.Printf("Got data: %v\n", string(data))
	return apiKey, nil
}
