package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html" // go get golang.org/x/net/html
	"net/http"
)

// TextBlock represents a block of text extracted from a webpage.
type TextBlock struct {
	URL  string
	Text string
}

// fetchURL fetches the HTML content of a URL.
func fetchURL(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

// extractTextBlocks extracts text from HTML nodes.
func extractTextBlocks(n *html.Node) string {
	var text strings.Builder
	if n.Type == html.TextNode {
		text.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text.WriteString(extractTextBlocks(c))
	}
	return text.String()
}

// worker scrapes a URL and sends the result to the results channel.
func worker(url string, results chan<- TextBlock, wg *sync.WaitGroup) {
	defer wg.Done()

	doc, err := fetchURL(url)
	if err != nil {
		log.Printf("Error fetching %s: %v\n", url, err)
		return
	}

	text := extractTextBlocks(doc)
	results <- TextBlock{URL: url, Text: text}
}

// readURLsFromFile reads URLs from a text file.
func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			urls = append(urls, url)
		}
	}
	return urls, scanner.Err()
}

func main() {
	// Read URLs from file
	urls, err := readURLsFromFile("urls.txt")
	if err != nil {
		log.Fatalf("Error reading URLs: %v\n", err)
	}

	// Channel for results and WaitGroup for synchronization
	results := make(chan TextBlock)
	var wg sync.WaitGroup

	// Start workers
	for _, url := range urls {
		wg.Add(1)
		go worker(url, results, &wg)
	}

	// Goroutine to close results channel after all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Timeout
	timeout := 30 * time.Second
	done := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		close(done)
	}()

	// Collect results
	var report []TextBlock
	for {
			select {
			case result, ok := <-results:
					if !ok {
							// Channel closed, all workers done
							fmt.Println("\n=== Scraping Report ===")
							for _, block := range report {
									fmt.Printf("\nURL: %s\nText:\n%s\n\n", block.URL, block.Text)
							}
							fmt.Println("All workers completed.")
							return
					}
					report = append(report, result)
					fmt.Printf("Extracted text from: %s\n", result.URL)
			case <-done:
					fmt.Println("\n=== Scraping Report (Timeout Reached) ===")
					for _, block := range report {
							fmt.Printf("\nURL: %s\nText:\n%s\n\n", block.URL, block.Text)
					}
					fmt.Println("Timeout reached.")
					return
			}
	}
}
