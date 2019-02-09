package main

import (
	"os"

	"github.com/carlmjohnson/errors"
	"github.com/carlmjohnson/feed2json/cli"
)

func main() {
	os.Exit(errors.Execute(cli.Tool, nil))
}
