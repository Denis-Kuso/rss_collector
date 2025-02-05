package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrNoAuthHeaderIncluded = errors.New("no authorization header included")
	ErrMalformedAuthHeader  = errors.New("malformed authorization header")
)

// Parses request header and returns the apiKey provided on succes
// and error otherwise
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", fmt.Errorf("%w: %v", ErrMalformedAuthHeader, authHeader)
	}

	return splitAuth[1], nil
}
