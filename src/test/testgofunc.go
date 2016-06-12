package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {

	var urls = []string{
		"http://www.golang.org/",
		"http://www.google.com/",
		"http://www.somestupidname.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go printUrl(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}

func printUrl(url string) (url1 string) {
	defer wg.Done()
	fmt.Println(url)

	return url1
}
