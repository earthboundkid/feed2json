# feed2json [![GoDoc](https://godoc.org/github.com/carlmjohnson/feed2json?status.svg)](https://godoc.org/github.com/carlmjohnson/feed2json) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/feed2json)](https://goreportcard.com/report/github.com/carlmjohnson/feed2json)

Given an Atom or RSS feed, creates a comparable JSON feed.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/carlmjohnson/feed2json/...
```
