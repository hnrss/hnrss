package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	NSDublinCore = "http://purl.org/dc/elements/1.1/"
	NSAtom       = "http://www.w3.org/2005/Atom"
	SiteURL      = "https://hnrss.org"
)

type CDATA struct {
	Value string `xml:",cdata"`
}

func Timestamp(fmt string, input time.Time) string {
	switch fmt {
	case "rss":
		return input.Format(time.RFC1123Z)
	case "atom", "jsonfeed":
		return input.Format(time.RFC3339)
	case "http":
		return input.Format(http.TimeFormat)
	default:
		return input.Format(time.RFC1123Z)
	}
}

func UTCNow() time.Time {
	return time.Now().UTC()
}

func ParseRequest(c *gin.Context, sp *SearchParams, op *OutputParams) {
	err := c.ShouldBindQuery(sp)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if strings.Contains(sp.Query, " OR ") {
		sp.Query = strings.Replace(sp.Query, " OR ", " ", -1)

		var q []string
		for _, f := range strings.Fields(sp.Query) {
			q = append(q, fmt.Sprintf("\"%s\"", f))
		}
		sp.Query = strings.Join(q, " ")
		sp.OptionalWords = strings.Join(q, " ")
	}

	err = c.ShouldBindQuery(op)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	op.Format = c.GetString("format")
	op.SelfLink = SiteURL + c.Request.URL.String()
}

func Generate(c *gin.Context, sp *SearchParams, op *OutputParams) {
	if op.Format == "" {
		op.Format = "rss"
	}

	results, err := GetResults(sp.Values())
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	c.Header("X-Algolia-URL", algoliaSearchURL+sp.Values().Encode())

	if len(results.Hits) > 0 {
		item := results.Hits[0]

		recent := item.GetCreatedAt()
		c.Header("Last-Modified", Timestamp("http", recent))

		if c.Request.URL.Path == "/item" {
			if sp.Query != "" {
				op.Title = fmt.Sprintf("Hacker News - \"%s\": \"%s\"", item.StoryTitle, sp.Query)
			} else {
				op.Title = fmt.Sprintf("Hacker News: New comments on \"%s\"", item.StoryTitle)
			}
		}
	}

	switch op.Format {
	case "rss":
		rss := NewRSS(results, op)
		c.XML(http.StatusOK, rss)
	case "atom":
		atom := NewAtom(results, op)
		c.XML(http.StatusOK, atom)
	case "jsonfeed":
		jsonfeed := NewJSONFeed(results, op)
		c.JSON(http.StatusOK, jsonfeed)
	}
}
