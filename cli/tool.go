package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/carlmjohnson/errutil"
	"github.com/carlmjohnson/feed2json"
	"github.com/carlmjohnson/flagext"
)

// Tool is the command line tool for cmd/feed2json
func Tool(args []string) (err error) {
	fl := flag.NewFlagSet("feed2json", flag.ContinueOnError)
	src := flagext.FileOrURL(flagext.StdIO, nil)
	fl.Var(src, "src", "source `file or URL`")
	dst := flagext.FileWriter(flagext.StdIO)
	fl.Var(dst, "dst", "destination `file`")
	fl.DurationVar(&http.DefaultClient.Timeout, "timeout", 5*time.Second, "timeout for URL sources")
	fl.Usage = func() {
		fmt.Fprintf(fl.Output(),
			`feed2json converts an Atom or RSS feed into a JSON feed.

	feed2json [opts]

Options:
`)
		fl.PrintDefaults()
	}

	if err := fl.Parse(args); err != nil {
		return flag.ErrHelp
	}

	var from, to bytes.Buffer
	if _, err = from.ReadFrom(src); err != nil {
		return err
	}
	if err = feed2json.Convert(&from, &to); err != nil {
		return err
	}
	_, err = io.Copy(dst, &to)
	defer errutil.Defer(&err, dst.Close)
	return err
}
