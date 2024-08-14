package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	// Define flags for the CLI arguments
	language := flag.String("l", "go", "Programming Language to use")
	searchTerm := flag.String("s", "slice", "Search term to look for")

	// Parse the flags
	flag.Parse()

	// Validate the flags
	if *language == "" || *searchTerm == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Context with timeout for the HTTP request
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Construct the URL
	url := fmt.Sprintf("https://cht.sh/%s/%s", *searchTerm, *language)

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create request: %v\n", err)
		os.Exit(1)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to make request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Failed to get response: %v\n", resp.Status)
		os.Exit(1)
	}

	// Determine if the response is HTML
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		// Parse the HTML and extract text
		doc, err := html.Parse(resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing HTML: %v\n", err)
			os.Exit(1)
		}
		extractText(doc)
	} else {
		// If it's not HTML, assume it's plain text
		if _, err := io.Copy(os.Stdout, resp.Body); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading response: %v\n", err)
			os.Exit(1)
		}
	}

}

// extractText recursively extracts and prints text from an HTML node
func extractText(n *html.Node) {
	if n.Type == html.TextNode {
		fmt.Print(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c)
	}
}
