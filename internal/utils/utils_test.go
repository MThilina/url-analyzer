package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeAndValidateURL(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{"example.com", false},
		{"http://example.com", false},
		{"https://example.com", false},
		{"", true},
		{"ht!tp://invalid-url", true},
	}

	for _, tt := range tests {
		normalized, err := NormalizeAndValidateURL(tt.input)
		if tt.expectError {
			assert.Error(t, err, "Expected error for input: %s", tt.input)
		} else {
			assert.NoError(t, err, "Unexpected error for input: %s", tt.input)
			assert.Contains(t, normalized, "http")
		}
	}
}
