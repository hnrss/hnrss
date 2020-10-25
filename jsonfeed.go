package main

// https://jsonfeed.org/version/1
type JSONFeed struct {
	Version     string         `json:"version"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Link        string         `json:"home_page_url"`
	Items       []JSONFeedItem `json:"items"`
}

type JSONFeedItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	ContentHTML string `json:"content_html"`
	URL         string `json:"url"`
	ExternalURL string `json:"external_url"`
	Published   string `json:"date_published"`
	Author      string `json:"author"`
}

func NewJSONFeed(results *AlgoliaSearchResponse, op *OutputParams) *JSONFeed {
	j := JSONFeed{
		Version:     "https://jsonfeed.org/version/1",
		Title:       op.Title,
		Link:        op.Link,
		Description: "Hacker News RSS",
	}
	for _, hit := range results.Hits {
		item := JSONFeedItem{
			ID:          hit.GetPermalink(),
			Title:       hit.GetTitle(),
			ContentHTML: hit.GetDescription(),
			URL:         hit.GetURL(op.LinkTo),
			ExternalURL: hit.GetPermalink(),
			Published:   Timestamp("jsonfeed", hit.GetCreatedAt()),
			Author:      hit.Author,
		}
		j.Items = append(j.Items, item)
	}
	return &j
}
