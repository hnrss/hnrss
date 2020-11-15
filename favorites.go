package main

import (
	"astuart.co/goq"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

const (
	hackerNewsFavoritesBase = "https://news.ycombinator.com/favorites?id="
)

func Favorites(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	resp, err := algoliaClient.Get(hackerNewsFavoritesBase + sp.ID)
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
	sp.Tags = fmt.Sprintf("(story,poll),(%s)", strings.Join(sids, ","))
	sp.Count = strconv.Itoa(len(sids))

	op.Title = fmt.Sprintf("Hacker News - %s's favorites", sp.ID)
	op.Link = hackerNewsFavoritesBase + sp.ID

	Generate(c, &sp, &op)
}
