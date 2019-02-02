package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/carlmjohnson/errors"
	"github.com/carlmjohnson/feed2json"
	"github.com/pseidemann/finish"
)

func main() {
	os.Exit(errors.Execute(webCLI, nil))
}

func webCLI(args []string) error {
	fl := flag.NewFlagSet("feed2jsonweb", flag.ContinueOnError)
	readTimeout := fl.Duration("read-timeout", 1*time.Second, "timeout for reading request headers")
	writeTimeout := fl.Duration("write-timeout", 1*time.Second, "timeout for writing response")

	addr := fl.String("address", ":8080", "listen on `host:port`")
	path := fl.String("path", "/", "serve requests on this path")
	param := fl.String("param", "url", "expect URL in this query param")
	host := fl.String("host", "", "require URLs to be on this host")

	fl.Usage = func() {
		fmt.Fprintf(fl.Output(),
			`feed2jsonweb serve

	feed2jsonweb [opts]

Options:
`)
		fl.PrintDefaults()
	}

	if err := fl.Parse(args); err != nil {
		return flag.ErrHelp
	}

	http.Handle(*path, feed2json.Handler(
		feed2json.ExtractURLFromParam(*param),
		feed2json.ValidateHost(*host),
		nil,
		nil,
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				next.ServeHTTP(w, r)
				log.Printf("%s for %q in %v", r.URL, r.UserAgent(), time.Since(start))
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != *path {
					log.Printf("[%d] Not found %q", http.StatusNotFound, r.URL)
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
	))

	srv := http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: *readTimeout,
		WriteTimeout:      *writeTimeout,
	}

	fin := finish.New()
	fin.Add(&srv)

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fin.Wait()
	return nil
}
