package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/net/html"
)

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
func crawlPage(host, path string, sitemap Sitemap, depth, maxDepth int) error {
	// Return if this path already exists
	if _, found := sitemap[path]; found {
		fmt.Println("Page has already been crawled, skipping.")
		return nil
	}
	if depth == maxDepth {
		fmt.Println("Maximum depth has been reached, skipping.")
	}
	// Add page to the global state
	sitemap[path] = NewStringSet()
	// Fetch the page
	res, err := fetchPage(host, path)
	if err != nil {
		return err
	}
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
			fmt.Printf("Found link: %v -> %v\n", path, url.Path)
			sitemap[path].Add(url.Path)
			crawlPage(host, url.Path, sitemap, depth+1, maxDepth)
		}
	}
	return nil
}

func main() {
	var depth int

	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = "Generate a sitemap for a given hostname."
	app.UsageText = "go-crawler [hostname]"
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
		fmt.Printf("Crawling %v...\n", host)
		// Create a global storage for pages
		// State is stored in a map, where keys are relative paths,
		// and value is a slice of links to other pages.
		sitemap := Sitemap{}
		err := crawlPage(host, "", sitemap, 0, depth)
		if err != nil {
			fmt.Println(err)
		}
		sitemap.Print()
		return nil
	}

	app.Run(os.Args)
}
