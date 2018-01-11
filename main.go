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

type StringSet struct {
	set map[string]bool
}

func NewStringSet() StringSet {
	return StringSet{map[string]bool{}}
}

func (s StringSet) Add(value string) bool {
	_, found := s.set[value]
	s.set[value] = true
	return !found
}

func (s StringSet) List() []string {
	list := make([]string, 0, len(s.set))
	for k := range s.set {
		list = append(list, k)
	}
	return list
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
func crawlPage(host, path string, pages map[string]StringSet) error {
	// Return if this path already exists
	if _, found := pages[path]; found {
		fmt.Println("Page has already been crawled, skipping.")
		return nil
	}
	// Add page to the global state
	pages[path] = NewStringSet()
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
			pages[path].Add(url.Path)
			crawlPage(host, url.Path, pages)
		}
	}
	return nil
}

func printSitemap(pages map[string]StringSet) {
	fmt.Println("Sitemap:")
	for page, linkSet := range pages {
		fmt.Println("/" + page)
		links := linkSet.List()
		for i, link := range links {
			symbol := "├── "
			if i == len(links)-1 {
				symbol = "└── "
			}
			fmt.Println(symbol, "/"+link)
		}
	}
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
		pages := map[string]StringSet{}
		err := crawlPage(host, "", pages)
		if err != nil {
			fmt.Println(err)
		}
		printSitemap(pages)
		return nil
	}

	app.Run(os.Args)
}
