// main.go
package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: ./crawler URL maxConcurrency maxPages")
		os.Exit(1)
	}
	rawBaseURL := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid maxConcurrency:", os.Args[2])
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid maxPages:", os.Args[3])
		os.Exit(1)
	}

	parsedBase, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Invalid URL: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Parsed base URL:", parsedBase.String())
	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            parsedBase,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}

	cfg.wg.Add(1)
	cfg.concurrencyControl <- struct{}{} // use a slot
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	args := os.Args[1:]

	fmt.Println("\nCrawled Pages Report:")
	for url, count := range cfg.pages {
		fmt.Printf("%s — visited %d time(s)\n", url, count)
	}

	fmt.Println("\nPages visited:")
	for url, count := range cfg.pages {
		fmt.Printf("%s — visited %d time(s)\n", url, count)
	}

	baseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", baseURL)

	pages := make(map[string]int)
	crawlPage(baseURL, baseURL, pages)

	fmt.Println("\nCrawled pages:")
	for url, count := range pages {
		fmt.Printf("%s: %d\n", url, count)
	}
}
