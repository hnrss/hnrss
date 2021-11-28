package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"astuart.co/goq"
	"github.com/gin-gonic/gin"
)

func BestComments(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	resp, err := algoliaClient.Get("https://news.ycombinator.com/bestcomments")
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

	var oids []string
	for _, id := range parsed.Thing {
		oids = append(oids, fmt.Sprintf("objectID:\"%s\"", id))
	}
	sp.Filters = strings.Join(oids, " OR ")
	sp.Count = strconv.Itoa(len(oids))

	op.Title = "Hacker News: Best Comments"
	op.Link = "https://news.ycombinator.com/bestcomments"

	Generate(c, &sp, &op)
}
