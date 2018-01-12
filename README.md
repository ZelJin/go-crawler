# go-crawler
Simple web crawler written in Go.

## Installation

Make sure you have a working Go environment.
You can find the installation instructions at [golang.org](https://golang.org/doc/install)

To install the crawler, run the following command:

```
$ go get github.com/ZelJin/go-crawler
```

Make sure your `PATH` includes the `$GOPATH/bin` directory:

```
export PATH=$PATH:$GOPATH/bin
```

## Usage

Crawl the Internet!

```
$ go-crawler https://github.com
```

You will receive comprehensive usage instructions after you run `go-crawler` command:

```
$ go-crawler
NAME:
   go-crawler - Generate a sitemap for a given URL.

USAGE:
   go-crawler [-d | -depth] [URL]

VERSION:
   1.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --depth value, -d value    Crawling depth (default: 10)
   --threads value, -t value  Maximum parallel crawling threads (default: 100)
   --help, -h                 show help
   --version, -v              print the version
```
