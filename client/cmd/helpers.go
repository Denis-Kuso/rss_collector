package cmd

import (
	"net/url"
	"strconv"
)

func isUrl(providedURL string) bool {
	u, err := url.Parse(providedURL)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

// Checks the validity of the username provided.
// Max length of username is 35 characters. White space cannot be used.
func validateUsername(username string) bool {
	runeUsername := []rune(username)
	n := len(runeUsername)
	if n > MAX_USERNAME_LENGTH || n == 0 {
		return false
	}
	return true
}

// check validity of provided id
// asserting other properties of uuid left to server
func isValidID(id string) bool {
	const UUID_LENGTH = 36 // 4 hypens and 32 chars
	if len([]rune(id)) != UUID_LENGTH {
		return false
	}
	return true
}

func validLimit(limit string) bool {
	const MAX_LIMIT = 100
	i, err := strconv.Atoi(limit)
	if err != nil {
		return false
	}
	if (i < 0) || (i > MAX_LIMIT) {
		return false
	}
	return true
}
