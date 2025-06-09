package service

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"url-analyzer/internal/model"
	"url-analyzer/internal/utils"

	"github.com/PuerkitoBio/goquery"
)

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

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, errors.New("failed to parse HTML document: " + err.Error())
	}

	// Determine HTML version
	htmlVersion := "HTML5" // default assumption
	doc.Find("doctype").Each(func(i int, s *goquery.Selection) {
		if goquery.NodeName(s) == "!doctype" {
			text := strings.ToLower(s.Text())
			if !strings.Contains(text, "html") {
				htmlVersion = "Unknown"
			}
		}
	})

	// Get title
	title := doc.Find("title").Text()

	// Count headings
	headings := make(map[string]int)
	for i := 1; i <= 6; i++ {
		tag := "h" + strconv.Itoa(i)
		headings[tag] = doc.Find(tag).Length()
	}

	// Link analysis
	base, _ := url.Parse(pageURL)
	internal, external, inaccessible := 0, 0, 0

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		link, err := url.Parse(href)
		if err != nil || link.Scheme == "mailto" {
			return
		}
		resolved := base.ResolveReference(link)

		if resolved.Hostname() == base.Hostname() {
			internal++
		} else {
			external++
		}

		if !utils.IsLinkAccessible(resolved.String()) {
			inaccessible++
		}
	})

	// Check for login form
	hasLoginForm := false
	doc.Find("input").Each(func(i int, s *goquery.Selection) {
		t, _ := s.Attr("type")
		if strings.ToLower(t) == "password" {
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
