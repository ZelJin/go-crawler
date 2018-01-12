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
func crawlPage(host string, sitemap Sitemap, state CrawlState, chQueue chan CrawlState, chWrite chan Page, chFinished chan bool) error {
	// Return if this path already exists
	if _, found := sitemap[state.Path]; found {
		fmt.Println("Page has already been crawled, skipping.")
		return nil
	}
	if state.Depth < 0 {
		fmt.Println("Maximum depth has been reached, skipping.")
		return nil
	}
	// Add page to the global state
	page := Page{state.Path, NewStringSet()}
	chWrite <- page
	// Fetch the page
	fmt.Printf("Crawling %v\n", state.Path)
	res, err := fetchPage(host, state.Path)
	if err != nil {
		return err
	}

	defer func() { chFinished <- true }()

	// Extract links
	defer res.Body.Close()
	links := extractLinks(res.Body)
	for _, link := range links {
		url, err := url.Parse(link)
		if err != nil {
			return err
		}
		//fmt.Println(url.Hostname())
		if url.Hostname() == host || url.Hostname() == "" {
			fmt.Printf("Found link: %v -> %v\n", state.Path, url.Path)
			page.LinkSet.Add(url.Path)
			chQueue <- CrawlState{url.Path, state.Depth - 1}
		}
	}
	// Add completed page to the global state
	chWrite <- page
	return nil
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

		sitemap := Sitemap{}

		chQueue := make(chan CrawlState, 100)
		chWrite := make(chan Page, 1)
		chFinished := make(chan bool)

		openRoutines := 1
		finishedRoutines := 0
		go crawlPage(host, sitemap, CrawlState{"", depth}, chQueue, chWrite, chFinished)

		for finishedRoutines < openRoutines {
			select {
			case state := <-chQueue:
				openRoutines++
				go crawlPage(host, sitemap, state, chQueue, chWrite, chFinished)
			case page := <-chWrite:
				sitemap[page.Path] = page
			case <-chFinished:
				finishedRoutines++
				fmt.Printf("Crawling routine finished, %v left.", openRoutines-finishedRoutines)
			}
		}

		sitemap.Print()

		elapsed := time.Since(start)
		fmt.Printf("Crawling took %s\n", elapsed)

		return nil
	}

	app.Run(os.Args)
}
