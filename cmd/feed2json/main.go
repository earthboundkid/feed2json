package main

import (
	"os"

	"github.com/carlmjohnson/exitcode"
	"github.com/carlmjohnson/feed2json/cli"
)

func main() {
	exitcode.Exit(cli.Tool(os.Args[1:]))
}
