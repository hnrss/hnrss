package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"text/template"
	"time"
)

const (
	hackerNewsItemID = "https://news.ycombinator.com/item?id="
	algoliaSearchURL = "https://hn.algolia.com/api/v1/search_by_date?"
)

var algoliaClient = http.Client{
	Timeout: 10 * time.Second,
}

type AlgoliaSearchResponse struct {
	Hits []AlgoliaSearchHit
}

type AlgoliaSearchHit struct {
	Tags        []string `json:"_tags"`
	ObjectID    string
	Title       string
	URL         string
	Author      string
	CreatedAt   string `json:"created_at"`
	StoryTitle  string `json:"story_title"`
	CommentText string `json:"comment_text"`
	StoryText   string `json:"story_text"`
	NumComments int    `json:"num_comments"`
	Points      int    `json:"points"`
	StoryID     int    `json:"story_id"`
	ParentID    int    `json:"parent_id"`
}

func (hit AlgoliaSearchHit) isComment() bool {
	for _, tag := range hit.Tags {
		if tag == "comment" {
			return true
		}
	}
	return false
}

func (hit AlgoliaSearchHit) isSelfPost() bool {
	return hit.StoryText != ""
}

func (hit AlgoliaSearchHit) GetTitle() string {
	if hit.isComment() {
		return fmt.Sprintf("New comment by %s in \"%s\"", hit.Author, html.UnescapeString(hit.StoryTitle))
	} else {
		return html.UnescapeString(hit.Title)
	}
}

func (hit AlgoliaSearchHit) GetPermalink() string {
	return hackerNewsItemID + hit.ObjectID
}

func (hit AlgoliaSearchHit) GetURL(linkTo string) string {
	if linkTo == "" {
		linkTo = "url"
	}

	if linkTo == "url" && hit.URL != "" {
		return hit.URL
	} else {
		return hit.GetPermalink()
	}
}

func buildTemplateEngine(name string) *template.Template {
	fm := make(template.FuncMap)
	fm["unescapeHTML"] = html.UnescapeString
	t := template.New(name)
	return t.Funcs(fm)
}

func (hit AlgoliaSearchHit) GetDescription() string {
	var (
		b bytes.Buffer
		t = buildTemplateEngine("default")
	)

	if hit.isComment() {
		t = template.Must(t.Parse(`
<p>{{ .CommentText | unescapeHTML }}</p>
`))
	} else if hit.isSelfPost() {
		t = template.Must(t.Parse(`
<p>{{ .StoryText | unescapeHTML }}</p>
<hr>
<p>Comments URL: <a href="{{ .GetPermalink }}">{{ .GetPermalink }}</a></p>
<p>Points: {{ .Points }}</p>
<p># Comments: {{ .NumComments }}</p>
`))
	} else {
		t = template.Must(t.Parse(`
{{ if .URL }}<p>Article URL: <a href="{{ .URL }}">{{ .URL }}</a></p>{{ end }}
<p>Comments URL: <a href="{{ .GetPermalink }}">{{ .GetPermalink }}</a></p>
<p>Points: {{ .Points }}</p>
<p># Comments: {{ .NumComments }}</p>
`))
	}

	t.Execute(&b, hit)
	return b.String()
}

func (hit AlgoliaSearchHit) GetCreatedAt() time.Time {
	if rv, err := time.Parse("2006-01-02T15:04:05.000Z", hit.CreatedAt); err == nil {
		return rv
	} else {
		return UTCNow()
	}
}

func GetResults(params url.Values) (*AlgoliaSearchResponse, error) {
	resp, err := algoliaClient.Get(algoliaSearchURL + params.Encode())
	if err != nil {
		return nil, errors.New("Error getting search results from Algolia")
	}
	defer resp.Body.Close()

	var parsed AlgoliaSearchResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&parsed)
	if err != nil {
		return nil, errors.New("Invalid JSON received from Algolia")
	}

	return &parsed, nil
}
