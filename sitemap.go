package main

import "fmt"

// Sitemap is a map with page URLs as keys and a set of links
// originating from a particular page as values
type Sitemap map[string]Page

// Print function prints a sitemap
func (s Sitemap) Print() {
	fmt.Println("Sitemap:")
	for _, page := range s {
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
