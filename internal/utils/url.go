// internal/utils/url.go
package utils

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func NormalizeAndValidateURL(raw string) (string, error) {
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return "", errors.New("invalid URL")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errors.New("invalid scheme")
	}

	if parsed.Host == "" || strings.Contains(parsed.Host, "!") {
		return "", errors.New("invalid host")
	}

	return parsed.String(), nil
}

func IsLinkAccessible(link string) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(link)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode < 400
}
