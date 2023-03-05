package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	HitsPerPageLimit = 100
)

type OutputParams struct {
	Title       string
	Link        string
	Description string `form:"description"`
	LinkTo      string `form:"link"`
	Format      string
	SelfLink    string
}

type SearchParams struct {
	Tags             string
	Query            string `form:"q"`
	OptionalWords    string
	Filters          string
	Points           string `form:"points"`
	ID               string `form:"id"`
	Author           string `form:"author"`
	Comments         string `form:"comments"`
	SearchAttributes string `form:"search_attrs"`
	Count            string `form:"count"`
}

func (sp *SearchParams) numericFilters() string {
	var filters []string
	if sp.Points != "" {
		filters = append(filters, "points>="+sp.Points)
	}
	if sp.Comments != "" {
		filters = append(filters, "num_comments>="+sp.Comments)
	}
	if sp.Tags == "front_page" {
		// For /frontpage requests, limit to stories created within the past week
		//
		// This is a workaround to avoid showing ancient stories that never lost the `front_page` tag
		timestamp := time.Now().Unix() - 7*24*60*60
		createdAt := strconv.FormatInt(timestamp, 10)
		filters = append(filters, "created_at_i>="+createdAt)
	}
	return strings.Join(filters, ",")
}

// Encode transforms the search options into an Algolia search querystring
func (sp *SearchParams) Values() url.Values {
	params := make(url.Values)

	if sp.OptionalWords != "" {
		params.Set("query", sp.Query)
		params.Set("optionalWords", sp.OptionalWords)
	} else if sp.Query != "" {
		params.Set("query", fmt.Sprintf("\"%s\"", sp.Query))
	}

	if f := sp.numericFilters(); f != "" {
		params.Set("numericFilters", f)
	}

	searchAttrs := sp.SearchAttributes
	if searchAttrs == "" {
		searchAttrs = "title"
	}
	if searchAttrs != "default" {
		params.Set("restrictSearchableAttributes", searchAttrs)
	}

	if sp.Count != "" {
		c, err := strconv.Atoi(sp.Count)
		if err != nil {
			c = 20
		} else if c > HitsPerPageLimit {
			c = HitsPerPageLimit
		}
		params.Set("hitsPerPage", strconv.Itoa(c))
	}

	if sp.Filters != "" {
		params.Set("filters", sp.Filters)
	}

	if sp.Tags != "" {
		params.Set("tags", sp.Tags)
	}

	return params
}
