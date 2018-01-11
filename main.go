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
	fmt.Println("Fetching ", url.String())
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
						fmt.Println("Found href: ", attr.Val)
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}

// Crawl a page specified by domain and relative path
func crawlPage(host, path string, pages map[string][]string) error {
	// Return if this path already exists
	if _, found := pages[path]; found {
		fmt.Println("Page has already been crawled, skipping.")
		return nil
	}
	// Add page to the global state
	pages[path] = []string{}
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
			pages[path] = append(pages[path], url.Path)
			crawlPage(host, url.Path, pages)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = "Generate a sitemap for a given hostname."
	app.UsageText = "go-crawler [hostname]"
	app.Version = "1.0.0"

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
		pages := map[string][]string{}
		err := crawlPage(host, "", pages)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("+%v", pages)
		return nil
	}

	app.Run(os.Args)
}
