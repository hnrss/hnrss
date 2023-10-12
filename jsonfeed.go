package main

import "fmt"

const hackerNewsAuthorFormatStr = "https://news.ycombinator.com/user?id=%s"

// https://jsonfeed.org/version/1
type JSONFeed struct {
	Version     string         `json:"version"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Link        string         `json:"home_page_url"`
	Items       []JSONFeedItem `json:"items"`
}

type JSONFeedItem struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	ContentHTML string             `json:"content_html,omitempty"`
	URL         string             `json:"url"`
	ExternalURL string             `json:"external_url"`
	Published   string             `json:"date_published"`
	Author      JSONFeedItemAuthor `json:"author"`
}

type JSONFeedItemAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func NewJSONFeed(results *AlgoliaSearchResponse, op *OutputParams) (*JSONFeed, error) {
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
			URL:         hit.GetURL(op.LinkTo),
			ExternalURL: hit.GetPermalink(),
			Published:   Timestamp("jsonfeed", hit.GetCreatedAt()),
			Author: JSONFeedItemAuthor{
				Name: hit.Author,
				URL:  fmt.Sprintf(hackerNewsAuthorFormatStr, hit.Author),
			},
		}

		if op.Description != descriptionDisabledFlag {
			desc, err := hit.GetDescription()
			if err != nil {
				return nil, err
			}

			item.ContentHTML = desc
		}

		j.Items = append(j.Items, item)
	}
	return &j, nil
}
