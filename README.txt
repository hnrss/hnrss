hnrss
=====

hnrss generates RSS feeds for new items (posts and comments) as they
appear on Hacker News.

Examples
--------

- http://hnrss.org/feeds/firehose.xml
- http://hnrss.org/feeds/askhn.xml?comments=20
- http://hnrss.org/feeds/comments.xml?points=25
- http://hnrss.org/feeds/author/pg.xml
- http://hnrss.org/feeds/search.xml?query=Django

Feeds
-----

hnrss provides the following feeds:

* http://hnrss.org/feeds/firehose.xml

  Contains all 'stories' sorted by age descending. A 'story' is all
  regular articles, 'Ask HN' and 'Show HN' posts, and polls.

* http://hnrss.org/feeds/askhn.xml

  Like the firehose feed, but only with 'Ask HN' posts.

* http://hnrss.org/feeds/showhn.xml

  Like the firehose feed, but only with 'Show HN' posts.

* http://hnrss.org/feeds/polls.xml

  Like the firehose feed, but only with polls.

* http://hnrss.org/feeds/comments.xml

  Contains all comments posted throughout Hacker News sorted by age
  descending.

  To only show comments from a particular story, pass its story ID via
  the "id" parameter:

    http://hnrss.org/feeds/comments.xml?id=7763923

  The ID is the number after "?id=" on an article's comment page.

* http://hnrss.org/feeds/author/<username>.xml

  Returns all stories and comments by a given username.

  To limit the results to only stories or comments, use the "only"
  parameter:

    http://hnrss.org/feeds/author/edavis.xml?only=comments # or "stories"

* http://hnrss.org/feeds/search.xml?query=VALUE

  Full-text search with results sorted by age descending.

  By default, results only include stories. To include both stories
  and comments, pass "?all=1":

    http://hnrss.org/feeds/search.xml?query=Django&all=1
  
Filters
-------

To limit results, use the "points" and/or "comments" filter(s):

* Points

  Provide a "points" GET parameter to only include results containing
  more than N points:

    http://hnrss.org/feeds/firehose.xml?points=50 # All stories with > 50 points

* Comments

  Provide a "comments" GET parameter to only include results
  containing more than N comments:

    http://hnrss.org/feeds/askhn.xml?comments=25 # All 'Ask HN' posts with > 25 points

  This works for all feeds except /feeds/comments.xml.

You can provide one, or both, and they'll be ANDed together.

Links
-----

By default, stories that link to external websites have that external
URL in <link>. If you'd rather have <link> point to the story's
comment page, provide "link=comments" as a GET parameter:

  http://hnrss.org/feeds/firehose.xml?link=comments

Credits
-------

hnrss is powered by Algolia (https://hn.algolia.com/api), the official
Hacker News API provider [1].

[1] https://news.ycombinator.com/item?id=7547578
