// Package feed2json converts Atom and RSS feeds to JSON feeds.
package feed2json

import (
	"github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

// Convert converts a gofeed.Feed (used by the parser) into a
// gorilla/feeds.JSONFeed (used by the JSONFeed emitter).
func Convert(feed *gofeed.Feed) *feeds.JSONFeed {
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
			Name: output.Author.Name,
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
			ContentHTML:   item.Content,
			Summary:       item.Description,
			PublishedDate: item.PublishedParsed,
			ModifiedDate:  item.UpdatedParsed,
			Tags:          item.Categories,
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
