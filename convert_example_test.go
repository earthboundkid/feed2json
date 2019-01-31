package feed2json_test

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/carlmjohnson/feed2json"
)

func ExampleConvert() {
	var from, to bytes.Buffer
	from.WriteString(`
<?xml version="1.0"?>
<rss version="2.0">
   <channel>
      <title>Liftoff News</title>
      <link>http://liftoff.msfc.nasa.gov/</link>
      <description>Liftoff to Space Exploration.</description>
      <language>en-us</language>
      <pubDate>Tue, 10 Jun 2003 04:00:00 GMT</pubDate>
      <lastBuildDate>Tue, 10 Jun 2003 09:41:01 GMT</lastBuildDate>
      <docs>http://blogs.law.harvard.edu/tech/rss</docs>
      <generator>Weblog Editor 2.0</generator>
      <managingEditor>editor@example.com</managingEditor>
      <webMaster>webmaster@example.com</webMaster>
      <item>
         <title>Star City</title>
         <link>http://liftoff.msfc.nasa.gov/news/2003/news-starcity.asp</link>
         <description>How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, language and protocol at Russia's &lt;a href="http://howe.iki.rssi.ru/GCTC/gctc_e.htm"&gt;Star City&lt;/a&gt;.</description>
         <pubDate>Tue, 03 Jun 2003 09:39:21 GMT</pubDate>
         <guid>http://liftoff.msfc.nasa.gov/2003/06/03.html#item573</guid>
      </item>
   </channel>
</rss>
    `)
	if err := feed2json.Convert(&from, &to); err == nil {
		from.Reset()
		json.Indent(&from, to.Bytes(), "", "  ")
		fmt.Println(from.String())
	}
	// Output:
	// {
	//   "version": "https://jsonfeed.org/version/1",
	//   "title": "Liftoff News",
	//   "home_page_url": "http://liftoff.msfc.nasa.gov/",
	//   "description": "Liftoff to Space Exploration.",
	//   "author": {},
	//   "items": [
	//     {
	//       "id": "http://liftoff.msfc.nasa.gov/2003/06/03.html#item573",
	//       "url": "http://liftoff.msfc.nasa.gov/news/2003/news-starcity.asp",
	//       "title": "Star City",
	//       "content_html": "How do Americans get ready to work with Russians aboard the International Space Station? They take a crash course in culture, language and protocol at Russia's <a href=\"http://howe.iki.rssi.ru/GCTC/gctc_e.htm\">Star City</a>.",
	//       "date_published": "2003-06-03T09:39:21Z"
	//     }
	//   ]
	// }
}
