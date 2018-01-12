package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

// Page is a thread-safe implementation of a node in a sitemap graph.
// Page has a URL and a list of links to other pages.
type Page struct {
	sync.Mutex
	URL   *url.URL
	Links []*Page
}

// NewPage initializes a Page object with a given URL.
func NewPage(url *url.URL) *Page {
	return &Page{URL: url, Links: []*Page{}}
}

// PrintSitemap function prints a sitemap for a given page.
func (p *Page) PrintSitemap() {
	fmt.Println("Printing sitemap...")
	// Generate a map to track visited nodes
	visited := map[string]bool{}
	traverse(p, "", visited, 0)
}

// Traverse a node in a sitemap graph. Used in PrintSitemap.
func traverse(p *Page, prefix string, visited map[string]bool, depth int) {
	// Check if the node has been visited
	if _, found := visited[p.URL.String()]; found {
		return
	}
	// Mark node as visited
	visited[p.URL.String()] = true
	// Print the URL of the node with its prefix
	prefix = strings.Repeat("  ", depth) + prefix
	fmt.Println(prefix, p.URL.String())
	// Run traverse for every child
	for _, child := range p.Links {
		traverse(child, " ├──", visited, depth+1)
	}
}
