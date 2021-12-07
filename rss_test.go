package main_test

import (
	"testing"

	hnrss "github.com/hnrss/hnrss"
	"github.com/stretchr/testify/require"
)

func TestNewRSSOutputParamsDescription(t *testing.T) {
	t.Parallel()

	results := &hnrss.AlgoliaSearchResponse{
		Hits: []hnrss.AlgoliaSearchHit{
			{
				Tags:      []string{"story", "author_dragonsh", "story_29367715"},
				ObjectID:  "29367715",
				Title:     "Mercurial Release 6.0",
				URL:       "https://www.mercurial-scm.org/wiki/Release6.0",
				Author:    "dragonsh",
				CreatedAt: "2021-11-28T10:11:15.000Z",
				Points:    1,
			},
		},
	}

	tests := map[string]struct {
		outputParamsFunc func(*hnrss.OutputParams)
		want             require.ValueAssertionFunc
	}{
		"Blank": {
			want: require.NotNil,
		},
		"Enabled": {
			outputParamsFunc: func(op *hnrss.OutputParams) {
				op.Description = "1"
			},
			want: require.NotNil,
		},
		"Disabled": {
			outputParamsFunc: func(op *hnrss.OutputParams) {
				op.Description = "0"
			},
			want: require.Nil,
		},
		"RandomValue": {
			outputParamsFunc: func(op *hnrss.OutputParams) {
				op.Description = "foo"
			},
			want: require.NotNil,
		},
	}

	for n, testCase := range tests {
		tc := testCase

		t.Run(n, func(t *testing.T) {
			t.Parallel()

			outputParams := &hnrss.OutputParams{
				Title:    "Hacker News: Newest",
				Link:     "https://news.ycombinator.com/newest",
				Format:   "rss",
				SelfLink: "https://hnrss.org/newest?description=0",
			}

			if tc.outputParamsFunc != nil {
				tc.outputParamsFunc(outputParams)
			}

			feed := hnrss.NewRSS(results, outputParams)

			require.Len(t, feed.Items, len(results.Hits))
			tc.want(t, feed.Items[0].Description)
		})
	}
}
