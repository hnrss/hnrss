package main

// http://cyber.harvard.edu/rss/rss.html
type RSS struct {
	XMLName       string    `xml:"rss"`
	Version       string    `xml:"version,attr"`
	NSDublinCore  string    `xml:"xmlns:dc,attr"`
	NSAtom        string    `xml:"xmlns:atom,attr"`
	Title         string    `xml:"channel>title"`
	Link          string    `xml:"channel>link"`
	Description   string    `xml:"channel>description"`
	Docs          string    `xml:"channel>docs"`
	Generator     string    `xml:"channel>generator"`
	LastBuildDate string    `xml:"channel>lastBuildDate"`
	AtomLink      AtomLink  `xml:"channel>atom:link"`
	Items         []RSSItem `xml:"channel>item"`
}

type RSSPermalink struct {
	Value       string `xml:",chardata"`
	IsPermaLink string `xml:"isPermaLink,attr"`
}

type RSSItem struct {
	Title       CDATA        `xml:"title"`
	Description CDATA        `xml:"description"`
	Published   string       `xml:"pubDate"`
	Link        string       `xml:"link"`
	Author      string       `xml:"dc:creator"`
	Comments    string       `xml:"comments"`
	Permalink   RSSPermalink `xml:"guid"`
}

func NewRSS(results *AlgoliaSearchResponse, op *OutputParams) *RSS {
	rss := RSS{
		Version:       "2.0",
		NSAtom:        NSAtom,
		NSDublinCore:  NSDublinCore,
		Title:         op.Title,
		Link:          op.Link,
		Description:   "Hacker News RSS",
		Docs:          "https://hnrss.org/",
		Generator:     "go-hnrss " + buildString,
		LastBuildDate: Timestamp("rss", UTCNow()),
		AtomLink:      AtomLink{op.SelfLink, "self", "application/rss+xml"},
	}

	for _, hit := range results.Hits {
		item := RSSItem{
			Title:       CDATA{hit.GetTitle()},
			Link:        hit.GetURL(op.LinkTo),
			Description: CDATA{hit.GetDescription()},
			Author:      hit.Author,
			Comments:    hit.GetPermalink(),
			Published:   Timestamp("rss", hit.GetCreatedAt()),
			Permalink:   RSSPermalink{hit.GetPermalink(), "false"},
		}
		rss.Items = append(rss.Items, item)
	}

	return &rss
}
