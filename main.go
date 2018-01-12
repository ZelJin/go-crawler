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

type CrawlState struct {
	Path  string
	Depth int
}

type Page struct {
	Path    string
	LinkSet StringSet
}

// Fetch page using http.Get
func fetchPage(host, path string) (*http.Response, error) {
	url := url.URL{Scheme: "http", Host: host, Path: path}
	return http.Get(url.String())
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

// Crawl a page specified by domain and relative path
func crawlPage(host string, sitemap *Sitemap, state *CrawlState, chQueue chan *CrawlState, chFinished chan bool) {
	defer func() { chFinished <- true }()

	// Return if this path already exists
	_, found := sitemap.Get(state.Path)
	if found {
		fmt.Println("Page has already been crawled, skipping.")
		return
	}
	// Return if reached max depth
	if state.Depth < 0 {
		fmt.Println("Maximum depth has been reached, skipping.")
		return
	}

	// Add empty page to global state to prevent it from
	// being crawled multiple times
	page := &Page{state.Path, NewStringSet()}
	sitemap.Set(state.Path, page)

	// Fetch the page
	fmt.Printf("Crawling %v\n", state.Path)
	res, err := fetchPage(host, state.Path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Extract links
	defer res.Body.Close()
	links := extractLinks(res.Body)
	for _, link := range links {
		url, err := url.Parse(link)
		if err != nil {
			fmt.Println()
			return
		}
		// If hostname is the same or link is relative
		if url.Hostname() == host || url.Hostname() == "" {
			fmt.Printf("Found link: %v -> %v\n", state.Path, url.Path)
			// Create a new routine and populate LinkSet
			page.LinkSet.Add(url.Path)
			chQueue <- &CrawlState{url.Path, state.Depth - 1}
		}
	}
	// Add completed page to the global state
	sitemap.Set(state.Path, page)
	return
}

func main() {
	var depth int

	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = "Generate a sitemap for a given hostname."
	app.UsageText = "go-crawler [-d | --depth] [hostname]"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "d, depth",
			Value:       10,
			Usage:       "Maximum crawling depth",
			Destination: &depth,
		},
	}

	app.Action = func(c *cli.Context) error {
		if !c.Args().Present() {
			fmt.Println("Please specify a hostname.")
			return nil
		}
		host := c.Args().Get(0)
		fmt.Printf("Started crawling %v\n", host)

		// Start recording time
		start := time.Now()
		sitemap := NewSitemap()

		// Init queues
		chQueue := make(chan *CrawlState, 100)
		chFinished := make(chan bool)

		// Init routines count
		openRoutines := 1
		finishedRoutines := 0
		go crawlPage(host, sitemap, &CrawlState{"", depth}, chQueue, chFinished)

		for finishedRoutines < openRoutines {
			select {
			case state := <-chQueue:
				openRoutines++
				go crawlPage(host, sitemap, state, chQueue, chFinished)
			case <-chFinished:
				finishedRoutines++
				fmt.Printf("Crawling routine finished, %v left\n", openRoutines-finishedRoutines)
			}
		}
		sitemap.Print()
		elapsed := time.Since(start)
		fmt.Printf("Crawling took %s\n", elapsed)
		return nil
	}

	app.Run(os.Args)
}
