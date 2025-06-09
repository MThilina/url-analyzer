package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeURL_Valid(t *testing.T) {
	resp, err := AnalyzeURL("https://example.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Title, "Example")
	assert.Contains(t, resp.HTMLVersion, "HTML")
}

func TestAnalyzeURL_Invalid(t *testing.T) {
	_, err := AnalyzeURL("http://nonexistent.localhost")
	assert.Error(t, err)
}
