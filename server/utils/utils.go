package utils

import (
	URL "net/url"
	"strings"
)

// Clean out white spaces
func CleanFields(s string) string {
	return strings.TrimSpace(s)
}

// Validate urls
func IsValidURL(url string) bool {
	_, err := URL.Parse(url)
	return err == nil
}
