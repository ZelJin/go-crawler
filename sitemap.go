package main

import "fmt"

// Sitemap is a map with page URLs as keys and a set of links
// originating from a particular page as values
type Sitemap map[string]StringSet

// Print function prints a sitemap
func (s Sitemap) Print() {
	fmt.Println("Sitemap:")
	for page, linkSet := range s {
		fmt.Println(page)
		for i, link := range linkSet.List() {
			symbol := "├── "
			if i == linkSet.Length()-1 {
				symbol = "└── "
			}
			fmt.Println(symbol, link)
		}
	}
}
