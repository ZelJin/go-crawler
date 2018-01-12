package main

import (
	"fmt"
	"sync"
)

// Sitemap is thread-safe a map with page URLs as keys and a set of links
// originating from a particular page as values
type Sitemap struct {
	sync.Mutex
	sitemap map[string]Page
}

func NewSitemap() *Sitemap {
	return &Sitemap{
		sitemap: map[string]Page{},
	}
}

func (s *Sitemap) Set(key string, value Page) {
	s.Lock()
	defer s.Unlock()
	s.sitemap[key] = value
}

func (s *Sitemap) Get(key string) (Page, bool) {
	s.Lock()
	defer s.Unlock()
	value, found := s.sitemap[key]
	return value, found
}

func (s *Sitemap) GetSitemap() map[string]Page {
	s.Lock()
	defer s.Unlock()
	return s.sitemap
}

// Print function prints a sitemap
func (s *Sitemap) Print() {
	fmt.Println("Sitemap:")
	sitemap := s.GetSitemap()
	for _, page := range sitemap {
		fmt.Println(page.Path)
		for i, link := range page.LinkSet.List() {
			symbol := "├── "
			if i == page.LinkSet.Length()-1 {
				symbol = "└── "
			}
			fmt.Println(symbol, link)
		}
	}
}
