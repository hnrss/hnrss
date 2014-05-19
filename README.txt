hnrss
=====

hnrss generates RSS feeds for new posts (and comments) as they appear
on Hacker News.

It can be found at http://hnrss.org/

Examples
--------

- http://hnrss.org/feeds/firehose.xml
- http://hnrss.org/feeds/askhn.xml?comments=20
- http://hnrss.org/feeds/comments?points=25
- http://hnrss.org/author/pg.xml
- http://hnrss.org/search.xml?query=Django

Feeds
-----

hnrss provides the following feeds:

* /feeds/firehose.xml

  Contains all 'stories' sorted by age descending. A 'story' is all
  regular articles, 'Ask HN' and 'Show HN' posts, and polls.

* /feeds/askhn.xml

  Like the firehose feed, but only with 'Ask HN' posts.

* /feeds/showhn.xml

  Like the firehose feed, but only with 'Show HN' posts.

* /feeds/polls.xml

  Like the firehose feed, but only with polls.

* /feeds/comments.xml

  Contains all comments posted throughout Hacker News sorted by age
  descending.

  To only show comments from a particular story, pass its story ID via
  the "id" parameter:

    /feeds/comments.xml?id=7763923

  The ID is the number after "?id=" on an article's comment page.

* /feeds/author/<username>.xml

  Returns all stories and comments by a given username.

  To limit the results to only stories or comments, use the "only"
  parameter:

    /feeds/author/edavis.xml?only=comments # or "stories"

* /feeds/search.xml?query=TERM

  Full-text search with results sorted by age descending.

  By default, results only include stories. To include both stories
  and comments, pass "?all=1":

    /feeds/search.xml?query=Django&all=1
  
Filters
-------

To limit results, use the "points" and/or "comments" filter(s):

* Points

  Provide a "points" GET parameter to only include results containing
  more than N points:

  /feeds/firehose.xml?points=50 # All stories with > 50 points

* Comments

  Provide a "comments" GET parameter to only include results
  containing more than N comments:

  /feeds/askhn.xml?comments=25 # All 'Ask HN' posts with > 25 points

  This works for all feeds except /feeds/comments.xml.

You can provide one, or both, and they'll be ANDed together.

Links
-----

By default, stories that link to external websites have that external
URL in <link>. If you'd rather have <link> point to the story's
comment page, provide "link=comments" as a GET parameter:

  /feeds/firehose.xml?link=comments

Titles
------

Comments and Ask/Show HN posts without body text don't have "titles"
in the traditional sense.

hnrss handles this by including this information in a single
<description> element without a <title>. This adheres to the RSS
specification, but can cause problems with some feed readers.

To include comments and Ask/Show HN headlines in a <title> element,
pass "usetitles=1" in your request:

  /feeds/comments.xml?points=10&usetitles=1

Credits
-------

hnrss is powered by Algolia (https://hn.algolia.com/api), the official
Hacker News API provider [1].

[1] https://news.ycombinator.com/item?id=7547578
