package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/urfave/cli"
	"golang.org/x/net/html"
)

// CrawlParams is a struct that contains page crawling parameters.
// It consists of a page that is being crawled and the remaining
// crawling depth.
type CrawlParams struct {
	Page  *Page
	Depth int
}

// Extract links from "href" attributes of an <a> tag.
// We don't need css or js for a sitemap.
func extractLinks(body io.ReadCloser) (links []string) {
	tokenizer := html.NewTokenizer(body)
	for {
		switch tokenType := tokenizer.Next(); tokenType {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

// Crawl a webpage.
func crawlPage(params *CrawlParams, visited *StringSet, queue chan *CrawlParams, workers chan bool) {
	defer func() {
		<-workers
		fmt.Printf("Finished processing url %s\n", params.Page.URL.String())
		fmt.Println(len(workers), len(queue))
		if len(workers) == 0 && len(queue) == 0 {
			// TODO: Possible race condition?
			close(queue)
		}
	}()

	// Return if reached max depth
	if params.Depth <= 0 {
		fmt.Println("Maximum depth has been reached, skipping.")
		return
	}

	// Lock the page to prevent simultaneous reads and writes
	params.Page.Lock()
	defer params.Page.Unlock()

	// Check if the page has been visited
	// If not, add the page to the "visited" set
	added := visited.Add(params.Page.URL.String())
	if !added {
		fmt.Println("Page has already been crawled, skipping.")
		return
	}

	// Fetch the page
	res, err := http.Get(params.Page.URL.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract links
	defer res.Body.Close()
	links := extractLinks(res.Body)
	for _, link := range links {
		linkURL, err := url.Parse(link)
		if err != nil {
			fmt.Printf("Failed to parse a link: %s\n", err)
			return
		}
		// If the hostname is the same or the link is relative
		if linkURL.Hostname() == params.Page.URL.Hostname() || linkURL.Hostname() == "" {
			// Resolve reference if the link is relative
			if !linkURL.IsAbs() {
				linkURL = params.Page.URL.ResolveReference(linkURL)
			}
			// Break if the page links to itself
			if linkURL.String() == params.Page.URL.String() {
				break
			}
			fmt.Printf("Found link: %v -> %v\n", params.Page.URL, linkURL)
			childPage := NewPage(linkURL)
			// Check if the maximum queue capacity is reached.
			if len(queue) == cap(queue) {
				panic("Deadlock! Queue limit reached.")
			}
			params.Page.Links = append(params.Page.Links, childPage)
			queue <- &CrawlParams{childPage, params.Depth - 1}
		}
	}
	return
}

func main() {
	var depth, workers, queue int

	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = "Generate a sitemap for a given URL."
	app.UsageText = "go-crawler [-d | -depth] [URL]"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "depth, d",
			Value:       10,
			Usage:       "Crawling depth",
			Destination: &depth,
		},
		cli.IntFlag{
			Name:        "workers, w",
			Value:       100,
			Usage:       "Maximum parallel crawling workers",
			Destination: &workers,
		},
		cli.IntFlag{
			Name:        "queue, q",
			Value:       100000,
			Usage:       "Maximum crawling queue",
			Destination: &queue,
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			cli.ShowAppHelp(c)
			return nil
		}
		url, err := url.Parse(c.Args().Get(0))
		if err != nil {
			fmt.Printf("Error parsing url: %s\n", err)
			return nil
		}
		if !url.IsAbs() {
			fmt.Println("Please specify an absolute URL (like https://monzo.com)")
			return nil
		}

		// Start recording time
		start := time.Now()
		// Generate a string set to track visited websites
		visited := NewStringSet()

		// Init channels
		workers := make(chan bool, workers)
		queue := make(chan *CrawlParams, queue)

		rootPage := NewPage(url)
		queue <- &CrawlParams{rootPage, depth}

		for params := range queue {
			workers <- true
			fmt.Printf("Started processing url %s\n", params.Page.URL.String())
			go crawlPage(params, visited, queue, workers)
		}

		rootPage.PrintSitemap()
		elapsed := time.Since(start)
		fmt.Printf("Crawling took %s\n", elapsed)
		return nil
	}

	app.Run(os.Args)
}
