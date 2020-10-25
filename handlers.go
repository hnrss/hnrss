package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Newest(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "(story,poll)"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Newest: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Newest"
	}
	op.Link = "https://news.ycombinator.com/newest"

	Generate(c, &sp, &op)
}

func Frontpage(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "front_page"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Front Page: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Front Page"
	}
	op.Link = "https://news.ycombinator.com/"

	Generate(c, &sp, &op)
}

func Newcomments(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	if sp.Query != "" {
		sp.SearchAttributes = "default"
		op.Title = fmt.Sprintf("Hacker News - New Comments: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: New Comments"
	}
	op.Link = "https://news.ycombinator.com/newcomments"

	Generate(c, &sp, &op)
}

func AskHN(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "ask_hn"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Ask HN: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Ask HN"
	}
	op.Link = "https://news.ycombinator.com/ask"

	Generate(c, &sp, &op)
}

func ShowHN(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "show_hn"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Show HN: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Show HN"
	}
	op.Link = "https://news.ycombinator.com/shownew"

	Generate(c, &sp, &op)
}

func Polls(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "poll"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Polls: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Polls"
	}
	op.Link = "https://news.ycombinator.com/"

	Generate(c, &sp, &op)
}

func Jobs(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "job"
	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - Jobs: \"%s\"", sp.Query)
	} else {
		op.Title = "Hacker News: Jobs"
	}
	op.Link = "https://news.ycombinator.com/jobs"

	Generate(c, &sp, &op)
}

func UserAll(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	tags := []string{"(story,comment,poll)", "author_" + sp.ID}
	sp.Tags = strings.Join(tags, ",")

	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - %s: \"%s\"", sp.ID, sp.Query)
	} else {
		op.Title = fmt.Sprintf("Hacker News: %s", sp.ID)
	}
	op.Link = "https://news.ycombinator.com/user?id=" + sp.ID

	Generate(c, &sp, &op)
}

func UserThreads(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	tags := []string{"comment", "author_" + sp.ID}
	sp.Tags = strings.Join(tags, ",")

	if sp.Query != "" {
		sp.SearchAttributes = "default"
		op.Title = fmt.Sprintf("Hacker News - %s threads: \"%s\"", sp.ID, sp.Query)
	} else {
		op.Title = fmt.Sprintf("Hacker News: %s threads", sp.ID)
	}
	op.Link = "https://news.ycombinator.com/threads?id=" + sp.ID

	Generate(c, &sp, &op)
}

func UserSubmitted(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	tags := []string{"(story,poll)", "author_" + sp.ID}
	sp.Tags = strings.Join(tags, ",")

	if sp.Query != "" {
		op.Title = fmt.Sprintf("Hacker News - %s submitted: \"%s\"", sp.ID, sp.Query)
	} else {
		op.Title = fmt.Sprintf("Hacker News: %s submitted", sp.ID)
	}
	op.Link = "https://news.ycombinator.com/submitted?id=" + sp.ID

	Generate(c, &sp, &op)
}

func Replies(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment"
	sp.SearchAttributes = "default"

	// If ID is a number, look for comments with a parent_id equal to the ID.
	// If ID is not a number, assume it is an author username and grab replies to their comments.
	_, err := strconv.Atoi(sp.ID)
	if err == nil {
		sp.Filters = "parent_id=" + sp.ID
		op.Title = "Hacker News: Replies to item #" + sp.ID
		op.Link = "https://news.ycombinator.com/item?id=" + sp.ID
	} else {
		values := make(url.Values)
		values.Set("tags", "comment,author_"+sp.ID)
		results, err := GetResults(values)
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadGateway, err.Error())
			return
		}

		var filters []string
		for _, hit := range results.Hits {
			filters = append(filters, "parent_id="+hit.ObjectID)
		}

		sp.Filters = strings.Join(filters, " OR ")
		op.Title = "Hacker News: Replies to " + sp.ID
		op.Link = "https://news.ycombinator.com/threads?id=" + sp.ID
	}

	Generate(c, &sp, &op)
}

func Item(c *gin.Context) {
	var sp SearchParams
	var op OutputParams
	ParseRequest(c, &sp, &op)

	sp.Tags = "comment,story_" + sp.ID
	sp.SearchAttributes = "default"

	// op.Title is set inside Generate to avoid the overhead of a
	// separate HTTP request to obtain the title.
	op.Link = "https://news.ycombinator.com/item?id=" + sp.ID

	Generate(c, &sp, &op)
}
