package service

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"url-analyzer/internal/model"

	"github.com/PuerkitoBio/goquery"
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

	defer resp.Body.Close()

	htmlVersion := "HTML5"
	title := doc.Find("title").Text()

	headings := make(map[string]int)
	for i := 1; i <= 6; i++ {
		tag := "h" + strconv.Itoa(i)
		headings[tag] = doc.Find(tag).Length()
	}

	base, _ := url.Parse(pageURL)
	var links []string

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		link, err := url.Parse(href)
		if err != nil || link.Scheme == "mailto" {
			return
		}
		resolved := base.ResolveReference(link)
		links = append(links, resolved.String())
	})

	results := make(chan linkCheckResult, len(links))
	for i, link := range links {
		go checkLinkAccessible(link, i, results, client)
	}

	internal, external, inaccessible := 0, 0, 0
	for i := 0; i < len(links); i++ {
		result := <-results
		linkURL, _ := url.Parse(links[result.index])
		if linkURL.Hostname() == base.Hostname() {
			internal++
		} else {
			external++
		}
		if !result.ok {
			inaccessible++
		}
	}

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
