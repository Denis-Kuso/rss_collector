package main

import (
	"net/url"
	"strconv"
)

func isURL(providedURL string) bool {
	u, err := url.Parse(providedURL)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

// Checks the validity of the username provided.
// Max length of username is 35 characters. White space cannot be used.
func validateUsername(username string) bool {
	maxUsernameLength := 35
	runeUsername := []rune(username)
	n := len(runeUsername)
	if n > maxUsernameLength || n == 0 {
		return false
	}
	return true
}

// check validity of provided id
// asserting other properties of uuid left to server
func isValidID(id string) bool {
	const UUIDchars = 36 // 4 hypens and 32 chars
	if len([]rune(id)) != UUIDchars {
		return false
	}
	return true
}

func validLimit(limit string) bool {
	const maxLimit = 100
	i, err := strconv.Atoi(limit)
	if err != nil {
		return false
	}
	if (i < 0) || (i > maxLimit) {
		return false
	}
	return true
}
