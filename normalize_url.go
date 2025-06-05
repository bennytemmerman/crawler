// normalize_url.go
package main

import (
	"net/url"
	"strings"
)

// normalizeURL takes a URL string and returns a normalized version of it.
func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(parsedURL.Host)
	if strings.HasPrefix(host, "www.") {
		host = strings.TrimPrefix(host, "www.")
	}

	path := parsedURL.EscapedPath()
	if path == "/" {
		path = ""
	}

	normalized := host + path

	// Include query string if it exists
	if parsedURL.RawQuery != "" {
		normalized += "?" + parsedURL.RawQuery
	}

	return normalized, nil
}

