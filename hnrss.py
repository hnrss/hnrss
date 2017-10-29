#!/usr/bin/env python

from api import API
from flask import Flask, request, redirect
from flask_compress import Compress
from rss import RSS
try:
    from urllib import urlencode
except ImportError:
    from urllib.parse import urlencode

app = Flask(__name__)
Compress(app)

@app.route('/newest')
def newest():
    query = request.args.get('q')
    api = API.using_request(request)
    if query:
        rss_title = 'Hacker News: "%s"' % query
    else:
        rss_title = 'Hacker News: Newest'
    rss = RSS(api.newest(), rss_title, 'https://news.ycombinator.com/newest')
    return rss.response()

@app.route('/frontpage')
def frontpage():
    api = API.using_request(request)
    rss_title = 'Hacker News: Front Page'
    rss = RSS(api.frontpage(), rss_title, 'https://news.ycombinator.com/')
    return rss.response()

@app.route('/newcomments')
def new_comments():
    query = request.args.get('q')
    api = API.using_request(request)
    if query:
        del api.params['restrictSearchableAttributes']
        rss_title = 'Hacker News: "%s" comments' % query
    else:
        rss_title = 'Hacker News: New Comments'
    rss = RSS(api.comments(), rss_title, 'https://news.ycombinator.com/newcomments')
    return rss.response()

@app.route('/ask')
def ask():
    api = API.using_request(request)
    rss = RSS(api.ask_hn(), 'Hacker News: Ask HN', 'https://news.ycombinator.com/ask')
    return rss.response()

@app.route('/show')
def show():
    api = API.using_request(request)
    rss = RSS(api.show_hn(), 'Hacker News: Show HN')
    return rss.response()

@app.route('/polls')
def polls():
    api = API.using_request(request)
    rss = RSS(api.polls(), 'Hacker News: Polls')
    return rss.response()

@app.route('/jobs')
def jobs():
    api = API.using_request(request)
    rss = RSS(api.jobs(), 'Hacker News: Jobs')
    return rss.response()

@app.route('/item')
def story_comments():
    story_id = request.args.get('id')
    api = API.using_request(request)
    api_response = api.comments(story_id=story_id)

    if api_response['hits']:
        story = api_response['hits'][0]
        rss_title = 'Hacker News: New Comments on "%s"' % story['story_title']
    else:
        rss_title = 'Hacker News: New Comments'

    rss_link = 'https://news.ycombinator.com/item?id=%s' % story_id
    rss = RSS(api_response, rss_title, rss_link)
    return rss.response()

@app.route('/user')
def user():
    username = request.args.get('id')
    api = API.using_request(request)
    api_response = api.user(username)

    rss_title = 'Hacker News: %s RSS feed' % username
    rss_link = 'https://news.ycombinator.com/user?id=%s' % username
    rss = RSS(api_response, rss_title, rss_link)
    return rss.response()

@app.route('/submitted')
def user_submitted():
    username = request.args.get('id')
    api = API.using_request(request)
    api_response = api.user(username, 'submitted')

    rss_title = 'Hacker News: %s submitted RSS feed' % username
    rss_link = 'https://news.ycombinator.com/submitted?id=%s' % username
    rss = RSS(api_response, rss_title, rss_link)
    return rss.response()

@app.route('/threads')
def user_threads():
    username = request.args.get('id')
    api = API.using_request(request)
    api_response = api.user(username, 'threads')

    rss_title = 'Hacker News: %s threads RSS feed' % username
    rss_link = 'https://news.ycombinator.com/threads?id=%s' % username
    rss = RSS(api_response, rss_title, rss_link)
    return rss.response()

# Redirect the old RSS endpoints to the new ones. Keep the query
# string intact when doing so.
@app.route('/feeds/<location>')
def feeds_redirects(location):
    qs = '?' + urlencode(request.args)
    redirect_map = {
        'firehose.xml' : '/newest',
        'comments.xml' : '/newcomments',
        'askhn.xml'    : '/ask',
        'showhn.xml'   : '/show',
        'polls.xml'    : '/polls',
        'search.xml'   : '/newest',
    }
    if location in redirect_map:
        return redirect(redirect_map[location] + qs)
    else:
        return redirect('/newest')

@app.route('/feeds/author/<author>.xml')
def author_redirect(author):
    args = request.args.copy()
    args['id'] = author
    qs = '?' + urlencode(args)
    return redirect('/user' + qs)

@app.route('/feeds/')
@app.route('/')
def index():
    return redirect('https://edavis.github.io/hnrss/')

if __name__ == '__main__':
    app.run(debug=True)
