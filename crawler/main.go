package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type ConcurrentMap struct {
    Map map[string]bool
    mu sync.Mutex
}

func (m *ConcurrentMap) Get(key string) bool {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.Map[key]
}

func (m *ConcurrentMap) Set(key string, isVisited bool) {
    m.mu.Lock()
    m.Map[key] = isVisited
    m.mu.Unlock()
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, concurrentMap *ConcurrentMap, ch chan string) {
	defer close(ch)
	if depth <= 0 {
		return
	}
    
    // is already visited, skip
    if ok := concurrentMap.Get(url); ok {
        return
    }

	body, urls, err := fetcher.Fetch(url)

    // set already visited
    concurrentMap.Set(url, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	ch <- fmt.Sprintf("found: %s %q", url, body)

	result := make([]chan string, len(urls))
	for i, u := range urls {
		result[i] = make(chan string)
		go Crawl(u, depth-1, fetcher, concurrentMap, result[i])
	}

	for i := range result {
		for c := range result[i] {
			ch <- c
		}
	}
}

func main() {
	ch := make(chan string)
    concurrentMap := ConcurrentMap{Map: make(map[string]bool), mu: sync.Mutex{}}
	go Crawl("https://golang.org/", 4, fetcher, &concurrentMap, ch)

	for s := range ch {
		fmt.Println(s)
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
