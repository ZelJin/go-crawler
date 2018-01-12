package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

type Page struct {
	sync.Mutex
	URL   *url.URL
	Links []*Page
}

func NewPage(url *url.URL) *Page {
	return &Page{URL: url, Links: []*Page{}}
}

// PrintSitemap function prints a sitemap for a given page
func (p *Page) PrintSitemap() {
	fmt.Println("Printing sitemap...")
	// Generate another map to track visited nodes
	visited := map[string]bool{}
	traverse(p, "", visited, 0)
}

func traverse(p *Page, prefix string, visited map[string]bool, depth int) {
	// Check if the node has been visited
	if _, found := visited[p.URL.String()]; found {
		return
	}
	// Set node to visited
	visited[p.URL.String()] = true
	// Print the node, w.r.t prefix
	prefix = strings.Repeat("  ", depth) + prefix
	fmt.Println(prefix, p.URL.String())
	// Run traverse for every child
	links := p.Links
	for _, child := range links {
		traverse(child, " ├──", visited, depth+1)
	}
}
