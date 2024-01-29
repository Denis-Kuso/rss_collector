package auth

import (
	"errors"
	"net/http"
	"fmt"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

// Parses request header and returns the apiKey provided on succes
// and error otherwise
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	fmt.Printf("Header: %s\n", authHeader)
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	for indx, item := range splitAuth{
		fmt.Printf("Auth element %v: %v\n",indx, item)
	}
	if len(splitAuth) < 2 || splitAuth[0] != "ranac" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}
