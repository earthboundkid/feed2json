// Package feed2json converts Atom and RSS feeds to JSON feeds.
package feed2json

import (
	"bytes"
	"encoding/json"
	"regexp"

	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

// Detects if <tag> or &code; is in text to guess if it's HTML
var htmlishRe = regexp.MustCompile(`<\w.*?>|&\w+;`)

// convertObject converts a gofeed.Feed (used by the parser) into a
// gorilla/feeds.JSONFeed (used by the JSONFeed emitter).
func convertObject(feed *gofeed.Feed) *feeds.JSONFeed {
	if feed == nil {
		return nil
	}
	output := feeds.JSONFeed{
		Version:     "https://jsonfeed.org/version/1",
		Title:       feed.Title,
		HomePageUrl: feed.Link,
		FeedUrl:     feed.FeedLink,
		Description: feed.Description,
	}
	if feed.Author != nil {
		output.Author = &feeds.JSONAuthor{
			Name: feed.Author.Name,
		}
	}
	if feed.Image != nil {
		output.Icon = feed.Image.URL
	}
	output.Items = make([]*feeds.JSONItem, len(feed.Items))
	for i := range feed.Items {
		if feed.Items[i] == nil {
			continue
		}
		item := feed.Items[i]
		jsonItem := &feeds.JSONItem{
			Id:            item.GUID,
			Url:           item.Link,
			Title:         item.Title,
			PublishedDate: item.PublishedParsed,
			ModifiedDate:  item.UpdatedParsed,
			Tags:          item.Categories,
		}
		if item.Content == "" {
			if htmlishRe.MatchString(item.Description) {
				jsonItem.ContentHTML = item.Description
			} else {
				jsonItem.ContentText = item.Description
			}
		} else {
			jsonItem.ContentHTML = item.Content
			jsonItem.Summary = item.Description
		}

		if item.Image != nil {
			jsonItem.Image = item.Image.URL
		}
		if item.Author != nil {
			jsonItem.Author = &feeds.JSONAuthor{
				Name: item.Author.Name,
			}
		}
		for _, enc := range item.Enclosures {
			if enc == nil {
				continue
			}
			jsonItem.Attachments = append(jsonItem.Attachments, feeds.JSONAttachment{
				Url:      enc.URL,
				MIMEType: enc.Type,
				// TODO convert length
			})
		}
		output.Items[i] = jsonItem
	}
	return &output
}

// Convert takes an XML feed from one buffer and turns it into a JSON feed
// in the other buffer.
func Convert(from, to *bytes.Buffer) (err error) {
	p := gofeed.NewParser()
	xmlfeed, err := p.Parse(from)
	if err != nil {
		return err
	}
	jsonfeed := convertObject(xmlfeed)
	enc := json.NewEncoder(to)
	enc.SetEscapeHTML(false)
	err = enc.Encode(jsonfeed)
	return
}
