package validate

import (
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

const (
	MAX_USERNAME_LENGTH = 35
	MAX_LIMIT           = 100
)

func ValidateUsername(username string) bool {
	runeUsername := []rune(username)
	n := len(runeUsername)
	if n > MAX_USERNAME_LENGTH || n == 0 {
		return false
	}
	return true
}

func IsUrl(providedURL string) bool {
	u, err := url.Parse(providedURL)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

func ValidLimit(limit string) bool {
	i, err := strconv.Atoi(limit)
	if err != nil {
		return false
	}
	if (i < 0) || (i > MAX_LIMIT) {
		return false
	}
	return true
}

func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	if err != nil {
		return false
	}
	return true
}
