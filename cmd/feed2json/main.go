package main

import (
	"os"

	"github.com/carlmjohnson/errors"
	"github.com/carlmjohnson/feed2json"
)

func main() {
	os.Exit(errors.Execute(feed2json.CLI, nil))
}
