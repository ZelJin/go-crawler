package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = "Generate a sitemap of a given domain."
	app.Action = func(c *cli.Context) error {
		fmt.Println("Crawling...")
		return nil
	}

	app.Run(os.Args)
}
