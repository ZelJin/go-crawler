package main

import (
	"fmt"
	"strings"
	"sync"
)

// Sitemap is thread-safe a map with page URLs as keys and a set of links
// originating from a particular page as values
type Sitemap struct {
	sync.Mutex
	sitemap map[string]*Page
}

func NewSitemap() *Sitemap {
	return &Sitemap{
		sitemap: map[string]*Page{},
	}
}

func (s *Sitemap) Set(key string, value *Page) {
	s.Lock()
	defer s.Unlock()
	s.sitemap[key] = value
}

func (s *Sitemap) Get(key string) (*Page, bool) {
	s.Lock()
	defer s.Unlock()
	value, found := s.sitemap[key]
	return value, found
}

func (s *Sitemap) GetSitemap() map[string]*Page {
	s.Lock()
	defer s.Unlock()
	return s.sitemap
}

// Print function prints a sitemap
func (s *Sitemap) Print() {
	fmt.Println("Printing sitemap...")
	sitemap := s.GetSitemap()
	// Generate another map to track visited nodes
	visited := map[string]bool{}
	rootPage := sitemap[""]
	traverse(rootPage, "", sitemap, visited, 0)
}

func traverse(p *Page, prefix string, sitemap map[string]*Page, visited map[string]bool, depth int) {
	// Check if the node has been visited
	if _, found := visited[p.Path]; found {
		return
	}
	// Set node to visited
	visited[p.Path] = true
	// Print the node, w.r.t prefix
	prefix = strings.Repeat(" ", 2*depth) + prefix
	fmt.Println(prefix, p.Path)
	// Run traverse for every child
	for i, path := range p.LinkSet.List() {
		childPrefix := "├── "
		if i == p.LinkSet.Length()-1 {
			childPrefix = "└── "
		}
		traverse(sitemap[path], childPrefix, sitemap, visited, depth+1)
	}
}
