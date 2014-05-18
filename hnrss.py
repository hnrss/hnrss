#!/usr/bin/env python

import time
import requests
from flask import Flask, request, make_response
from lxml import etree

API_BASE_URL = 'https://hn.algolia.com/api/v1'

app = Flask(__name__)

def do_search(request, endpoint, tags):
    params = {}
    filters = []

    if tags:
        params['tags'] = tags

    if request.args.get('query'):
        params['query'] = '"%s"' % request.args.get('query')

    if request.args.get('points'):
        filters.append('points>%s' % request.args.get('points'))
    if request.args.get('comments'):
        filters.append('num_comments>%s' % request.args.get('comments'))
    if filters:
        params['numericFilters'] = ','.join(filters)

    return requests.get('%s/%s' % (API_BASE_URL, endpoint),
                        params=params, verify=False)

def add_element(parent, tag, value, **kwargs):
    el = etree.Element(tag, kwargs)
    el.text = value
    parent.append(el)
    return el

def generate_rfc2822(secs=None):
    t = time.gmtime(secs)
    return time.strftime('%a, %d %b %Y %H:%M:%S GMT', t)
    
def generate_rss(request, response, title):
    """
    Generate a RSS document from a search API response.
    """
    rss_root = etree.Element('rss', version='2.0')
    rss_channel = etree.SubElement(rss_root, 'channel')

    add_element(rss_channel, 'title', title)
    add_element(rss_channel, 'link', 'https://news.ycombinator.com/')
    add_element(rss_channel, 'description', 'Hacker News RSS')
    add_element(rss_channel, 'docs', 'http://cyber.law.harvard.edu/rss/rss.html')
    add_element(rss_channel, 'generator', 'hnrss.py (https://github.com/edavis/hnrss)')
    add_element(rss_channel, 'lastBuildDate', generate_rfc2822())

    for hit in response.json()['hits']:
        hn_url = 'https://news.ycombinator.com/item?id=%s' % hit['objectID']
        tags = hit.get('_tags', [])

        rss_item = etree.SubElement(rss_channel, 'item')

        if hit.get('title') and hit.get('story_text'):
            add_element(rss_item, 'title', hit.get('title'))
            add_element(rss_item, 'description', hit.get('story_text'))
        elif hit.get('title') and not hit.get('story_text'):
            add_element(rss_item, 'description', hit.get('title'))
        elif hit.get('comment_text'):
            add_element(rss_item, 'description', hit.get('comment_text'))

        add_element(rss_item, 'pubDate', generate_rfc2822(hit.get('created_at_i')))

        if ('ask_hn' in tags or 'poll' in tags or 'comment' in tags):
            add_element(rss_item, 'link', hn_url)
        elif 'story' in tags:
            if request.args.get('link', '') == 'comments':
                add_element(rss_item, 'link', hn_url)
            elif hit.get('url'):
                add_element(rss_item, 'link', hit.get('url'))
        
        add_element(rss_item, 'author', hit.get('author'))

        if 'story' in tags:
            add_element(rss_item, 'comments', hn_url)

        add_element(rss_item, 'guid', hn_url)

    return rss_root

def make_rss_response(rss_doc):
    rss_xml = etree.tostring(rss_doc, pretty_print=True, encoding='UTF-8', xml_declaration=True)
    return (rss_xml, 200, {'Content-Type': 'text/xml'})

@app.route('/feeds/firehose.xml')
def stories():
    response = do_search(request, 'search_by_date', 'story')
    rss = generate_rss(request, response, 'Hacker News: Firehose')
    return make_rss_response(rss)

@app.route('/feeds/askhn.xml')
def askhn():
    response = do_search(request, 'search_by_date', 'ask_hn')
    rss = generate_rss(request, response, 'Hacker News: Ask HN')
    return make_rss_response(rss)

@app.route('/feeds/showhn.xml')
def showhn():
    response = do_search(request, 'search_by_date', 'show_hn')
    rss = generate_rss(request, response, 'Hacker News: Show HN')
    return make_rss_response(rss)

@app.route('/feeds/polls.xml')
def polls():
    response = do_search(request, 'search_by_date', 'poll')
    rss = generate_rss(request, response, 'Hacker News: Polls')
    return make_rss_response(rss)

@app.route('/feeds/comments.xml')
def comments():
    tags = ['comment']
    if request.args.get('id'):
        tags.append('story_%s' % request.args.get('id'))
        response = do_search(request, 'search_by_date', ','.join(tags))
        rss_title = '"%s" comments' % response.json()['hits'][0]['story_title']
    else:
        response = do_search(request, 'search_by_date', ','.join(tags))
        rss_title = 'New Comments'
    rss = generate_rss(request, response, 'Hacker News: %s' % rss_title)
    return make_rss_response(rss)

@app.route('/feeds/author/<username>.xml')
def author(username):
    tags = ['author_%s' % username]
    if request.args.get('only'):
        only = request.args.get('only')
        if only == 'stories':
            tags.append('story')
        elif only == 'comments':
            tags.append('comment')
    response = do_search(request, 'search_by_date', ','.join(tags))
    rss = generate_rss(request, response, 'Hacker News: %s RSS feed' % username)
    return make_rss_response(rss)

@app.route('/feeds/search.xml')
def search():
    query = request.args.get('query')
    tags = 'story'
    if request.args.get('all'):
        tags = '(story,comment)'
    response = do_search(request, 'search_by_date', tags)
    rss = generate_rss(request, response, 'Hacker News: Search for "%s"' % query)
    return make_rss_response(rss)

if __name__ == '__main__':
    app.run(debug=True)
