// crawler.go
package main

import (
	"fmt"
	"net/url"
	"sync"
)

// Config holds shared state and tools for concurrency
type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

// addPageVisit tracks page visits and returns true if this is the first time seeing the URL
func (cfg *config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, exists := cfg.pages[normalizedURL]
	if exists {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}

// crawlPage visits a single page, extracts links, and recursively spawns goroutines for them
func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()
	defer func() { <-cfg.concurrencyControl }()
	// --- EARLY RETURN IF MAX PAGES REACHED ---
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	parsedCurrent, err := url.Parse(rawCurrentURL)
	if err != nil || parsedCurrent.Host != cfg.baseURL.Host {
		return // skip pages on different domains or bad URLs
	}

	normalized, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("normalize error: %v\n", err)
		return
	}

	if !cfg.addPageVisit(normalized) {
		return // already visited
	}

	fmt.Println("Crawling:", normalized)

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("getHTML error: %v\n", err)
		return
	}

	links, err := getURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		fmt.Printf("getURLs error: %v\n", err)
		return
	}

	for _, link := range links {
		cfg.wg.Add(1)
		go func(l string) {
			cfg.concurrencyControl <- struct{}{}
			cfg.crawlPage(l)
		}(link)
	}
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("invalid base URL: %v\n", err)
		return
	}

	currURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("invalid current URL: %v\n", err)
		return
	}

	// Make sure the current URL is on the same domain
	if currURL.Host != baseURL.Host {
		return
	}

	normURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
	fmt.Printf("error normalizing URL: %v\n", err)
	return
	}

	// Already seen this page, increment count and return
	if _, ok := pages[normURL]; ok {
		pages[normURL]++
		return
	}

	fmt.Println("Crawling:", normURL)
	pages[normURL] = 1

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error fetching page: %v\n", err)
		return
	}

	// Extract URLs from HTML
	urls, err := getURLsFromHTML(html, rawCurrentURL)
	if err != nil {
		fmt.Printf("error extracting URLs: %v\n", err)
		return
	}

	// Recursively crawl discovered URLs
	for _, link := range urls {
		crawlPage(rawBaseURL, link, pages)
	}
}

