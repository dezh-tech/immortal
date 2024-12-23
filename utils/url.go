package utils

import (
	"net/url"
	"strings"
)

// helper function for ValidateAuthEvent.
func ParseURL(input string) (*url.URL, error) {
	return url.Parse(
		strings.ToLower(
			strings.TrimSuffix(input, "/"),
		),
	)
}
