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
	app.UsageText = "go-crawler [domain]"
	app.Version = "1.0.0"
	app.Action = func(c *cli.Context) error {
		domain := c.Args().Get(0)
		fmt.Printf("Crawling %v...\n", domain)
		return nil
	}

	app.Run(os.Args)
}
