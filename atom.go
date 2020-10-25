package main

// https://validator.w3.org/feed/docs/atom.html
type Atom struct {
	XMLName string      `xml:"feed"`
	NS      string      `xml:"xmlns,attr"`
	ID      string      `xml:"id"`
	Title   string      `xml:"title"`
	Updated string      `xml:"updated"`
	Links   []AtomLink  `xml:"link"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomEntry struct {
	Title     CDATA       `xml:"title"`
	Links     []AtomLink  `xml:"link"`
	Author    string      `xml:"author>name"`
	Content   AtomContent `xml:"content"`
	Updated   string      `xml:"updated"`
	Published string      `xml:"published"`
	ID        string      `xml:"id"`
}

type AtomContent struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",cdata"`
}

type AtomLink struct {
	Reference    string `xml:"href,attr"`
	Relationship string `xml:"rel,attr,omitempty"`
	Type         string `xml:"type,attr,omitempty"`
}

func NewAtom(results *AlgoliaSearchResponse, op *OutputParams) *Atom {
	atom := Atom{
		NS:      NSAtom,
		ID:      op.SelfLink,
		Title:   op.Title,
		Updated: Timestamp("atom", UTCNow()),
		Links: []AtomLink{
			AtomLink{op.SelfLink, "self", "application/atom+xml"},
		},
	}

	for _, hit := range results.Hits {
		entry := AtomEntry{
			ID:        hit.GetPermalink(),
			Title:     CDATA{hit.GetTitle()},
			Updated:   Timestamp("atom", hit.GetCreatedAt()),
			Published: Timestamp("atom", hit.GetCreatedAt()),
			Links: []AtomLink{
				AtomLink{hit.GetURL(op.LinkTo), "alternate", ""},
			},
			Author:  hit.Author,
			Content: AtomContent{"html", hit.GetDescription()},
		}
		atom.Entries = append(atom.Entries, entry)
	}

	return &atom
}
