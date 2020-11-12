package main

import (
	"astuart.co/goq"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func Classic(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	resp, err := algoliaClient.Get("https://news.ycombinator.com/classic")
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}
	defer resp.Body.Close()

	var parsed ItemList
	err = goq.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	var sids []string
	for _, id := range parsed.Thing {
		sids = append(sids, "story_"+id)
	}
	sp.Tags = fmt.Sprintf("story,(%s)", strings.Join(sids, ","))
	sp.Count = strconv.Itoa(len(sids))

	op.Title = "Hacker News: Classic"
	op.Link = "https://news.ycombinator.com/classic"

	Generate(c, &sp, &op)
}
