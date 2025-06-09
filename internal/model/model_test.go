package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzeResponseSerialization(t *testing.T) {
	resp := AnalyzeResponse{
		HTMLVersion:  "HTML5",
		Title:        "Test Page",
		Headings:     map[string]int{"h1": 1, "h2": 2},
		HasLoginForm: true,
		Links: LinkSummary{
			Internal:     5,
			External:     3,
			Inaccessible: 1,
		},
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "Test Page")
	assert.Contains(t, string(data), "HTML5")
	assert.Contains(t, string(data), "h1")
	assert.Contains(t, string(data), "inaccessible")
}
