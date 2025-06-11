package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDetectHTMLVersionFromBytes verifies DOCTYPE detection.
func TestDetectHTMLVersionFromBytes(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"HTML5 doctype", `<!DOCTYPE html><html></html>`, "HTML5"},
		{"HTML 4.01 Strict", `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN"><html></html>`, "HTML 4.01"},
		{"XHTML Transitional", `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"><html></html>`, "XHTML"},
		{"No doctype", `<html></html>`, "Unknown"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := detectHTMLVersionFromBytes([]byte(tc.input))
			if got != tc.expected {
				t.Errorf("detectHTMLVersionFromBytes() = %q; want %q", got, tc.expected)
			}
		})
	}
}

// TestAnalyzeURL tests the end-to-end behavior of AnalyzeURL.
func TestAnalyzeURL(t *testing.T) {
	const htmlBody = `<!DOCTYPE html>
<html>
  <head><title>My Title</title></head>
  <body>
    <h1>Heading1</h1>
    <h2>Heading2</h2>
    <a href="/foo">Foo</a>
    <a href="https://bar.example.com">Bar</a>
    <form><input type="password"/></form>
  </body>
</html>`

	// Server only for the test host; external HEAD will fail.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(htmlBody))
	}))
	defer ts.Close()

	resp, err := AnalyzeURL(ts.URL)
	if err != nil {
		t.Fatalf("AnalyzeURL() error: %v", err)
	}

	// HTMLVersion
	if resp.HTMLVersion != "HTML5" {
		t.Errorf("HTMLVersion = %q; want \"HTML5\"", resp.HTMLVersion)
	}

	// Title
	if resp.Title != "My Title" {
		t.Errorf("Title = %q; want \"My Title\"", resp.Title)
	}

	// Headings
	if resp.Headings["h1"] != 1 {
		t.Errorf("h1 count = %d; want 1", resp.Headings["h1"])
	}
	if resp.Headings["h2"] != 1 {
		t.Errorf("h2 count = %d; want 1", resp.Headings["h2"])
	}

	// Login form
	if !resp.HasLoginForm {
		t.Error("HasLoginForm = false; want true")
	}

	// Link summary: 1 internal ("/foo"), 1 external ("bar.example.com"), and external is inaccessible.
	links := resp.Links
	if links.Internal != 1 {
		t.Errorf("Internal links = %d; want 1", links.Internal)
	}
	if links.External != 1 {
		t.Errorf("External links = %d; want 1", links.External)
	}
	if links.Inaccessible != 1 {
		t.Errorf("Inaccessible links = %d; want 1", links.Inaccessible)
	}
}
