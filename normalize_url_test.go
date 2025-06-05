// normalize_url_test.go
package main

import (
	"testing"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		expected string
	}{
		{
			name:     "remove scheme",
			inputURL: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			name:     "remove www",
			inputURL: "https://www.boot.dev/learn",
			expected: "boot.dev/learn",
		},
		{
			name:     "keep query params",
			inputURL: "https://blog.boot.dev/path?sort=asc",
			expected: "blog.boot.dev/path?sort=asc",
		},
		{
			name:     "strip trailing slash on root",
			inputURL: "https://www.boot.dev/",
			expected: "boot.dev",
		},
		{
			name:     "no path",
			inputURL: "https://www.boot.dev",
			expected: "boot.dev",
		},
		{
			name:     "mixed case host",
			inputURL: "https://WWW.Boot.Dev/Path",
			expected: "boot.dev/Path",
		},
		{
			name:     "http scheme",
			inputURL: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := normalizeURL(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

