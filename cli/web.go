package cli

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/carlmjohnson/feed2json"
	"github.com/carlmjohnson/flagext"
	"github.com/go-chi/cors"
	"github.com/peterbourgon/ff"
	"github.com/pseidemann/finish"
)

// Web is the command line tool for cmd/feed2jsonweb
func Web(args []string) error {
	var (
		addr         string
		path         string
		param        string
		readTimeout  time.Duration
		writeTimeout time.Duration
		maxAge       time.Duration
		hosts        flagext.Strings
		corsOrigins  = flagext.Strings{"*"}
	)
	{
		fl := flag.NewFlagSet("feed2jsonweb", flag.ContinueOnError)
		fl.DurationVar(&readTimeout, "read-timeout", 1*time.Second, "timeout for reading request headers")
		fl.DurationVar(&writeTimeout, "write-timeout", 2*time.Second, "timeout for writing response")
		fl.DurationVar(&http.DefaultClient.Timeout, "request-timeout", 1*time.Second, "timeout for fetching XML")
		fl.DurationVar(&maxAge, "max-age", 5*time.Minute, "set Cache-Control: public, max-age header")
		port := fl.String("port", "8080", "port `number` to listen on")
		host := fl.String("host", "127.0.0.1", "host `name` to listen for")
		fl.StringVar(&path, "url-path", "/", "serve requests on this path")
		fl.StringVar(&param, "param", "url", "expect URL in this query param")
		fl.Var(&hosts, "allow-host", "require requested URLs to be on `host`")
		fl.Var(&corsOrigins, "cors-origin", "allow these CORS origins")

		fl.Usage = func() {
			fmt.Fprintf(fl.Output(),
				`feed2jsonweb is an HTTP server that converts Atom and RSS feeds to JSON feeds

Usage:

    feed2jsonweb [opts]


Options:

`)
			fl.PrintDefaults()
			fmt.Fprintf(fl.Output(),
				`
Note: -allow-host and -cors-origin can be passed multiple times to set more hosts and origins. Options can also be passed as environmental variables (CAPITALIZED_WITH_UNDERSCORES) prefixed by FEED2JSON_.
`,
			)

		}

		if err := ff.Parse(fl, args, ff.WithEnvVarPrefix("FEED2JSON")); err != nil {
			return flag.ErrHelp
		}
		addr = net.JoinHostPort(*host, *port)
	}

	http.Handle(path, feed2json.Handler(
		feed2json.ExtractURLFromParam(param),
		feed2json.ValidateHost(hosts...),
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
			if maxAge == 0 {
				return next
			}

			seconds := float64(maxAge) / float64(time.Second)
			header := fmt.Sprintf("public, max-age=%.0f", seconds)

			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Cache-Control", header)
				next.ServeHTTP(w, r)
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != path {
					log.Printf("[%d] Not found %q", http.StatusNotFound, r.URL)
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				}
				next.ServeHTTP(w, r)
			})
		},
		cors.New(cors.Options{
			AllowedOrigins: corsOrigins,
		}).Handler,
	))

	srv := http.Server{
		Addr:              addr,
		ReadHeaderTimeout: readTimeout,
		WriteTimeout:      writeTimeout,
	}

	fin := finish.New()
	fin.Add(&srv)

	go func() {
		log.Printf("Serving on %s", addr)
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	fin.Wait()
	return nil
}
