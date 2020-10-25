package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

func HiringCommon(c *gin.Context, query string) {
	params := make(url.Values)
	if query != "" {
		params.Set("query", fmt.Sprintf("\"%s\"", query))
		params.Set("hitsPerPage", "1")
	}
	params.Set("tags", "story,author_whoishiring")

	results, err := GetResults(params)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	if len(results.Hits) < 1 {
		e := errors.New("No whoishiring stories found")
		c.Error(e)
		c.String(http.StatusBadGateway, e.Error())
		return
	}

	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	if query != "" {
		sp.Filters = "parent_id=" + results.Hits[0].ObjectID
	} else {
		var filters []string
		for _, hit := range results.Hits {
			filters = append(filters, "parent_id="+hit.ObjectID)
		}
		sp.Filters = strings.Join(filters, " OR ")
	}
	sp.SearchAttributes = "default"
	op.Title = results.Hits[0].Title
	op.Link = "https://news.ycombinator.com/item?id=" + results.Hits[0].ObjectID

	Generate(c, &sp, &op)
}

func SeekingEmployees(c *gin.Context) {
	HiringCommon(c, "Ask HN: Who is hiring?")
}

func SeekingEmployers(c *gin.Context) {
	HiringCommon(c, "Ask HN: Who wants to be hired?")
}

func SeekingFreelance(c *gin.Context) {
	HiringCommon(c, "Ask HN: Freelancer? Seeking freelancer?")
}

func SeekingAll(c *gin.Context) {
	HiringCommon(c, "")
}
