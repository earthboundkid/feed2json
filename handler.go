package feed2json

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

type URLExtractor = func(*http.Request) *url.URL

func ExtractURLFromParam(name string) URLExtractor {
	return func(r *http.Request) *url.URL {
		u, _ := url.Parse(r.URL.Query().Get(name))
		return u
	}
}

type URLValidator = func(*url.URL) bool

func ValidateHost(name string) URLValidator {
	return func(u *url.URL) bool {
		if u == nil {
			return false
		}
		return u.Host == name || name == ""
	}
}

type Middleware = func(http.Handler) http.Handler

type Logger = func(format string, v ...interface{})

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
