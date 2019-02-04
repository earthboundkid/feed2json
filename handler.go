package feed2json

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

// URLExtractor is a user provided callback that determines a URL for an XML feed
// based on a request
type URLExtractor = func(*http.Request) *url.URL

// ExtractURLFromParam is a URLExtractor that extracts a URL from the query
// param specified by name.
func ExtractURLFromParam(name string) URLExtractor {
	return func(r *http.Request) *url.URL {
		u, _ := url.Parse(r.URL.Query().Get(name))
		return u
	}
}

// URLValidator is a user provided callback that determines whether the URL
// for an XML feed is valid for Handler.
type URLValidator = func(*url.URL) bool

// ValidateHost is a URLValidator that approves of URLs where the hostname
// is in the names list.
func ValidateHost(names ...string) URLValidator {
	if len(names) == 0 {
		return func(u *url.URL) bool {
			return true
		}
	}
	m := map[string]struct{}{}
	for _, name := range names {
		m[name] = struct{}{}
	}
	return func(u *url.URL) bool {
		if u == nil {
			return false
		}
		_, ok := m[u.Host]
		return ok
	}
}

// Middleware wraps an http.Handler in a http.Handler.
type Middleware = func(http.Handler) http.Handler

// Logger is a user provided callback that matches the fmt/log.Printf calling
// conventions.
type Logger = func(format string, v ...interface{})

// Handler is an http.Handler that extracts and validates a URL for a request,
// then requests is with the provided http.Client. Responses from Handler are
// wrapped by the user provided middleware, if any.
//
// c if nil defaults to http.DefaultClient. l if nil defaults to log.Printf.
func Handler(x URLExtractor, v URLValidator, c *http.Client, l Logger, ms ...Middleware) http.Handler {
	if c == nil {
		c = http.DefaultClient
	}
	if l == nil {
		l = log.Printf
	}
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			handleErr(w, l, http.StatusMethodNotAllowed,
				"method %q not allowed", r.Method)
			return
		}
		u := x(r)
		if u == nil || !v(u) {
			handleErr(w, l, http.StatusBadRequest,
				"bad url requested: %q", u)
			return
		}
		rsp, err := c.Get(u.String())
		if err != nil {
			handleErr(w, l, http.StatusBadGateway,
				"error connecting to %q: %v", u, err)
			return
		}
		defer rsp.Body.Close()
		var from, to bytes.Buffer
		if _, err = from.ReadFrom(rsp.Body); err != nil {
			handleErr(w, l, http.StatusBadGateway,
				"error reading %q: %v", u, err)
			return
		}
		if err = Convert(&from, &to); err != nil {
			handleErr(w, l, http.StatusInternalServerError,
				"error converting %q: %v", u, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err = io.Copy(w, &to); err != nil {
			l("error completing write of %q: %v", u, err)
			return
		}
		l("[%d] converted %q", http.StatusOK, u)
	})
	for _, m := range ms {
		h = m(h)
	}
	return h
}

func handleErr(w http.ResponseWriter, l Logger, code int, format string, v ...interface{}) {
	statusText := http.StatusText(code)
	http.Error(w, statusText, code)
	v = append([]interface{}{code}, v...)
	l("[%d] "+format, v...)
}
