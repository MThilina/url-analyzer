// internal/utils/url.go
package utils

import (
	"errors"
	"net/url"
	"strings"
)

func NormalizeAndValidateURL(raw string) (string, error) {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || parsed.Host == "" {
		return "", errors.New("invalid URL")
	}
	return parsed.String(), nil
}
