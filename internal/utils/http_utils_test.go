package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLinkAccessible(t *testing.T) {
	// Accessible mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	assert.True(t, IsLinkAccessible(server.URL), "Expected link to be accessible")

	// Inaccessible (non-resolvable host)
	assert.False(t, IsLinkAccessible("http://nonexistent.localhost"), "Expected link to be inaccessible")
}
