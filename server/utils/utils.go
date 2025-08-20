package utils

import (
	URL "net/url"
	"regexp"
	"strings"
)

// Clean out white spaces
func CleanFields(s string) string {
	return strings.TrimSpace(s)
}

// Validate urls
func IsValidURL(url string) bool {
	parsedUrl, err := URL.ParseRequestURI(url)
	return err == nil && parsedUrl.Host != ""
}

func GenerateInternalUrl(url string, href string) string {
	generatedUrl := ""
	// Check href contains section in the path
	regex := regexp.MustCompile(`^#.+`)
	if regex.MatchString(href) {
		generatedUrl = url + href
	} else {
		parsedUrl, err := URL.Parse(url)
		if err != nil {
			return href
		}
		parsedUrl.Path = ""
		parsedUrl.RawQuery = ""
		parsedUrl.Fragment = ""
		generatedUrl = parsedUrl.String() + href
	}
	return generatedUrl
}
