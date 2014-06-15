hnrss — Hacker News RSS
========================

hnrss provides custom, realtime RSS feeds for Hacker News.

The [project page](http://hnrss.org/) explains all available RSS feeds.

Overview
--------

The following feeds are available:

- **Firehose** — Every new [post](http://hnrss.org/newest) and [comment](http://hnrss.org/newcomments) as it arrives
- **Points** — [Posts](http://hnrss.org/newest?points=300) and [comments](http://hnrss.org/newcomments?points=25) with more than N points
- **Activity** — [Posts](http://hnrss.org/newest?comments=250) with more than N comments
- **Self-posts** — All "[Ask HN](http://hnrss.org/ask)" and "[Show HN](http://hnrss.org/show)" posts, along with [polls](http://hnrss.org/polls)
- **Users** — New [posts](http://hnrss.org/submitted?id=tokenadult) and [comments](http://hnrss.org/threads?id=tptacek) made by a given user
- **Threads** — Each new comment made [in a given thread](http://hnrss.org/item?id=7864813)
- **Searches** — New [posts](http://hnrss.org/newest?q=git) and [comments](http://hnrss.org/newcomments?q=django) matching a given search term

Credits
-------

hnrss is powered by [Algolia](https://hn.algolia.com/api), the
[official Hacker News API provider](https://news.ycombinator.com/item?id=7547578).
