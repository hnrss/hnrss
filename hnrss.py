#!/usr/bin/env python

import time
import requests
from flask import Flask, request, redirect
from lxml import etree

app = Flask(__name__)

class API(object):
    base_url = 'https://hn.algolia.com/api/v1'

    def __init__(self, endpoint='search_by_date', points=None, comments=None):
        self.endpoint = endpoint
        self.params = {}
        if points or comments:
            numeric_filters = []
            if points: numeric_filters.append('points>%s' % points)
            if comments: numeric_filters.append('num_comments>%s' % comments)
            self.params['numericFilters'] = ','.join(numeric_filters)

    @classmethod
    def using_request(cls, request):
        return cls(
            points = request.args.get('points'),
            comments = request.args.get('comments'),
        )

    def _request(self, tags, query=None):
        params = self.params.copy()
        params['tags'] = tags
        if query:
            params['query'] = query
        resp = requests.get(
            '%s/%s' % (self.base_url, self.endpoint),
            params = params,
            verify = False,
        )
        resp.raise_for_status()
        return resp.json()

    def firehose(self):
        return self._request('(story,poll)')

    def ask_hn(self):
        return self._request('ask_hn')

    def show_hn(self):
        return self._request('show_hn')

    def polls(self):
        return self._request('poll')

class RSS(object):
    def __init__(self, api_response, title, link='https://news.ycombinator.com/'):
        self.api_response = api_response

        self.rss_root = etree.Element('rss', version='2.0')
        self.rss_channel = etree.SubElement(self.rss_root, 'channel')

        self.add_element(self.rss_channel, 'title', title)
        self.add_element(self.rss_channel, 'link', link)
        self.add_element(self.rss_channel, 'description', 'Hacker News RSS')
        self.add_element(self.rss_channel, 'docs', 'http://cyber.law.harvard.edu/rss/rss.html')
        self.add_element(self.rss_channel, 'generator', 'https://github.com/edavis/hnrss')
        self.add_element(self.rss_channel, 'lastBuildDate', self.generate_rfc2822())

        self.generate_body()

    def generate_body(self):
        for hit in self.api_response['hits']:
            rss_item = etree.SubElement(self.rss_channel, 'item')
            hn_url = 'https://news.ycombinator.com/item?id=%s' % hit['objectID']
            tags = hit.get('_tags', [])

            if 'comment' in tags:
                self.add_element(rss_item, 'title', 'New comment by %s in "%s"' % (
                    hit.get('author'), hit.get('story_title')))
                self.add_element(rss_item, 'description', hit.get('comment_text'))
            else:
                if hit.get('title'):
                    self.add_element(rss_item, 'title', hit.get('title'))
                if hit.get('story_text'):
                    self.add_element(rss_item, 'description', hit.get('story_text'))

            self.add_element(rss_item, 'pubDate', self.generate_rfc2822(hit.get('created_at_i')))

            if ('ask_hn' in tags or 'poll' in tags or 'comment' in tags):
                self.add_element(rss_item, 'link', hn_url)
            elif 'story' in tags:
                self.add_element(rss_item, 'link', hit.get('url') or hn_url)

            self.add_element(rss_item, 'author', hit.get('author'))

            if ('story' in tags or 'poll' in tags):
                self.add_element(rss_item, 'comments', hn_url)

            self.add_element(rss_item, 'guid', hn_url)

    def add_element(self, parent, tag, text):
        el = etree.Element(tag)
        el.text = text
        parent.append(el)
        return el

    def generate_rfc2822(self, secs=None):
        t = time.gmtime(secs)
        return time.strftime('%a, %d %b %Y %H:%M:%S GMT', t)

    def response(self):
        rss_xml = etree.tostring(
            self.rss_root, pretty_print=True, encoding='UTF-8', xml_declaration=True,
        )
        return (rss_xml, 200, {'Content-Type': 'text/xml'})

@app.route('/feeds/firehose.xml')
def firehose():
    api = API.using_request(request)
    rss = RSS(api.firehose(), 'Hacker News: Firehose', 'https://news.ycombinator.com/newest')
    return rss.response()

@app.route('/feeds/askhn.xml')
def askhn():
    api = API.using_request(request)
    rss = RSS(api.ask_hn(), 'Hacker News: Ask HN', 'https://news.ycombinator.com/ask')
    return rss.response()

@app.route('/feeds/showhn.xml')
def showhn():
    api = API.using_request(request)
    rss = RSS(api.show_hn(), 'Hacker News: Show HN')
    return rss.response()

@app.route('/feeds/polls.xml')
def polls():
    api = API.using_request(request)
    rss = RSS(api.polls(), 'Hacker News: Polls')
    return rss.response()

if __name__ == '__main__':
    app.run(debug=True)
