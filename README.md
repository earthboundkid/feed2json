# feed2json [![GoDoc](https://godoc.org/github.com/carlmjohnson/feed2json?status.svg)](https://godoc.org/github.com/carlmjohnson/feed2json) [![Go Report Card](https://goreportcard.com/badge/github.com/carlmjohnson/feed2json)](https://goreportcard.com/report/github.com/carlmjohnson/feed2json)

Given an Atom or RSS feed, creates a comparable JSON feed.

## Installation

First install [Go](http://golang.org).

If you just want to install the binary to your current directory and don't care about the source code, run

```bash
GOBIN="$(pwd)" GOPATH="$(mktemp -d)" go get github.com/carlmjohnson/feed2json/...
```

## Screenshots
```bash
$ feed2json -h
feed2json converts an Atom or RSS feed into a JSON feed.

    feed2json [opts]

Options:
  -dst file
        destination file (default stdout)
  -src file or URL
        source file or URL (default stdin)
  -timeout duration
        timeout for URL sources (default 5s)

$ feed2json -src 'https://jsonfeed.org/xml/rss.xml' | json-tidy
{
    "description": "JSON Feed is a pragmatic syndication format for blogs, microblogs, and other time-based content.",
    "home_page_url": "https://jsonfeed.org/",
    "items": [
        {
            "content_html": "<p>We — Manton Reece and Brent Simmons — have noticed that JSON has become the developers’ choice for APIs, and that developers will often go out of their way to avoid XML. JSON is simpler to read and write, and it’s less prone to bugs.</p>\n\n<p>So we developed JSON Feed, a format similar to <a href=\"http://cyber.harvard.edu/rss/rss.html\">RSS</a> and <a href=\"https://tools.ietf.org/html/rfc4287\">Atom</a> but in JSON. It reflects the lessons learned from our years of work reading and publishing feeds.</p>\n\n<p><a href=\"https://jsonfeed.org/version/1\">See the spec</a>. It’s at version 1, which may be the only version ever needed. If future versions are needed, version 1 feeds will still be valid feeds.</p>\n\n<h4>Notes</h4>\n\n<p>We have a <a href=\"https://github.com/manton/jsonfeed-wp\">WordPress plugin</a> and, coming soon, a JSON Feed Parser for Swift. As more code is written, by us and others, we’ll update the <a href=\"https://jsonfeed.org/code\">code</a> page.</p>\n\n<p>See <a href=\"https://jsonfeed.org/mappingrssandatom\">Mapping RSS and Atom to JSON Feed</a> for more on the similarities between the formats.</p>\n\n<p>This website — the Markdown files and supporting resources — <a href=\"https://github.com/brentsimmons/JSONFeed\">is up on GitHub</a>, and you’re welcome to comment there.</p>\n\n<p>This website is also a blog, and you can subscribe to the <a href=\"https://jsonfeed.org/xml/rss.xml\">RSS feed</a> or the <a href=\"https://jsonfeed.org/feed.json\">JSON feed</a> (if your reader supports it).</p>\n\n<p>We worked with a number of people on this over the course of several months. We list them, and thank them, at the bottom of the <a href=\"https://jsonfeed.org/version/1\">spec</a>. But — most importantly — <a href=\"http://furbo.org/\">Craig Hockenberry</a> spent a little time making it look pretty. :)</p>",
            "date_published": "2017-05-17T15:02:12Z",
            "id": "https://jsonfeed.org/2017/05/17/announcing_json_feed",
            "title": "Announcing JSON Feed",
            "url": "https://jsonfeed.org/2017/05/17/announcing_json_feed"
        }
    ],
    "title": "JSON Feed",
    "version": "https://jsonfeed.org/version/1"
}
```
