package main

import "fmt"

type Sitemap map[string]StringSet

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
