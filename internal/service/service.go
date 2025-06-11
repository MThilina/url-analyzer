package service

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"url-analyzer/internal/model"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type linkCheckResult struct {
	index int
	ok    bool
}

func checkLinkAccessible(link string, index int, ch chan<- linkCheckResult, client *http.Client) {
	resp, err := client.Head(link)
	if err != nil || resp.StatusCode >= 400 {
		ch <- linkCheckResult{index: index, ok: false}
		return
	}
	ch <- linkCheckResult{index: index, ok: true}
}

// detectHTMLVersionFromBytes scans the full HTML buffer for a DOCTYPE and returns the version.
func detectHTMLVersionFromBytes(buf []byte) string {
	version := "Unknown"
	z := html.NewTokenizer(bytes.NewReader(buf))
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.DoctypeToken {
			dt := strings.ToUpper(string(z.Text()))
			switch {
			case strings.Contains(dt, "HTML 4.01"):
				version = "HTML 4.01"
			case strings.Contains(dt, "XHTML"):
				version = "XHTML"
			case strings.Contains(dt, "HTML"):
				// e.g. <!DOCTYPE html> → HTML5
				version = "HTML5"
			default:
				version = strings.TrimSpace(string(z.Text()))
			}
			break
		}
	}
	return version
}

// AnalyzeURL fetches pageURL, detects its HTML version, parses content,
// and summarizes title, headings, link counts, and login‐form presence.
func AnalyzeURL(pageURL string) (*model.AnalyzeResponse, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, errors.New("unable to fetch the URL: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("received non-OK HTTP status: " + resp.Status)
	}

	// Read entire body into memory
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response body: " + err.Error())
	}

	// Detect HTML version from full buffer
	htmlVersion := detectHTMLVersionFromBytes(bodyBytes)

	// Parse document with GoQuery
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, errors.New("failed to parse HTML document: " + err.Error())
	}

	// Extract title
	title := doc.Find("title").Text()

	// Count headings h1–h6
	headings := make(map[string]int)
	for i := 1; i <= 6; i++ {
		tag := "h" + strconv.Itoa(i)
		headings[tag] = doc.Find(tag).Length()
	}

	// Resolve and collect all <a href=""> links
	base, _ := url.Parse(pageURL)
	var links []string
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		u, err := url.Parse(href)

		if err != nil {
			return
		}
		resolved := base.ResolveReference(u)
		links = append(links, resolved.String())
	})

	// Check links in parallel
	results := make(chan linkCheckResult, len(links))
	for i, link := range links {
		go checkLinkAccessible(link, i, results, client)
	}

	internal, external, inaccessible := 0, 0, 0
	for i := 0; i < len(links); i++ {
		res := <-results
		u, _ := url.Parse(links[res.index])
		if u.Hostname() == base.Hostname() {
			internal++
		} else {
			external++
		}
		if !res.ok {
			inaccessible++
		}
	}

	// Detect presence of a login form (password inputs)
	hasLoginForm := false
	doc.Find("input[type]").Each(func(i int, s *goquery.Selection) {
		if strings.EqualFold(s.AttrOr("type", ""), "password") {
			hasLoginForm = true
		}
	})

	return &model.AnalyzeResponse{
		HTMLVersion:  htmlVersion,
		Title:        title,
		Headings:     headings,
		HasLoginForm: hasLoginForm,
		Links: model.LinkSummary{
			Internal:     internal,
			External:     external,
			Inaccessible: inaccessible,
		},
	}, nil
}
