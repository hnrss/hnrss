import re
import time
import hashlib
from xml.sax.saxutils import unescape as sax_unescape
from flask import request
from lxml import etree

try:
    unichr(0)
except NameError:
    unichr = chr
try:
    xrange(0)
except NameError:
    xrange = range

def unescape(s):
    deref_ncr = lambda m: unichr(int(m.group(1), 16)) # '&#x2F;' -> '/'
    s = re.sub('&#[Xx]([A-Fa-f0-9]+);', deref_ncr, s)
    entities = {'&quot;': '"', '&apos;': "'"}
    return sax_unescape(s, entities)

def insert_donation_request(guid):
    h = hashlib.sha1(guid).hexdigest()
    if h.startswith(('0', '1', '2', '3')):
        return '''
<hr><p>hnrss is a labor of love, but if the project has made your job
or hobby project easier and you want to show some gratitude, <a
href="https://donate.hnrss.org/">donations are very much
appreciated</a>. PayPal and Bitcoin both accepted. Thanks!</p>
        '''
    else:
        return ''

class RSS(object):
    def __init__(self, api_response, title, link='https://news.ycombinator.com/'):
        self.api_response = api_response

        nsmap = {
            'dc': 'http://purl.org/dc/elements/1.1/',
            'atom': 'http://www.w3.org/2005/Atom',
        }
        self.rss_root = etree.Element('rss', version='2.0', nsmap=nsmap)
        self.rss_channel = etree.SubElement(self.rss_root, 'channel')

        self.add_element(self.rss_channel, 'title', title)
        self.add_element(self.rss_channel, 'link', link)
        self.add_element(self.rss_channel, 'description', 'Hacker News RSS')
        self.add_element(self.rss_channel, 'docs', 'https://edavis.github.io/hnrss/')
        self.add_element(self.rss_channel, 'generator', 'https://github.com/edavis/hnrss')
        self.add_element(self.rss_channel, 'lastBuildDate', self.generate_rfc2822())

        # FIXME: Is there a way to tell Flask or nginx we're running under HTTPS so this is correct off the bat?
        atom_link = request.url.replace('http://', 'https://')
        atom_link = atom_link.replace('"', '%22').replace(' ', '%20')
        self.add_element(self.rss_channel, '{http://www.w3.org/2005/Atom}link', text='', rel='self', type='application/rss+xml', href=atom_link)

        if 'hits' in api_response:
            self.generate_body()

    def generate_body(self):
        for hit in self.api_response['hits']:
            rss_item = etree.SubElement(self.rss_channel, 'item')
            hn_url = 'https://news.ycombinator.com/item?id=%s' % hit['objectID']
            tags = hit.get('_tags', [])

            if 'comment' in tags:
                if hit.get('story_title') and hit.get('comment_text'):
                    self.add_element(rss_item, 'title', 'New comment by %s in "%s"' % (
                        hit.get('author'), hit.get('story_title')))
                    self.add_element(rss_item, 'description', unescape(hit.get('comment_text')))
            else:
                if hit.get('title'):
                    self.add_element(rss_item, 'title', hit.get('title'))
                if hit.get('story_text'):
                    self.add_element(rss_item, 'description', unescape(hit.get('story_text')))
                elif self.api_response['description']:
                    template = (
                        '<p>Article URL: <a href="%(url)s">%(url)s</a></p>'
                        '<p>Comments URL: <a href="%(hn_url)s">%(hn_url)s</a></p>'
                        '<p>Points: %(points)s</p>'
                        '<p># Comments: %(comments)s</p>'
                    ) + insert_donation_request(hn_url)
                    params = {
                        'url': hit.get('url') or hn_url,
                        'hn_url': hn_url,
                        'points': hit.get('points', 0) or 0,
                        'comments': hit.get('num_comments', 0) or 0,
                    }
                    self.add_element(rss_item, 'description', template % params)

            self.add_element(rss_item, 'pubDate', self.generate_rfc2822(hit.get('created_at_i')))

            if self.api_response['link_to'] == 'comments':
                self.add_element(rss_item, 'link', hn_url)
            else:
                self.add_element(rss_item, 'link', hit.get('url') or hn_url)

            self.add_element(rss_item, '{http://purl.org/dc/elements/1.1/}creator', hit.get('author'))

            if ('story' in tags or 'poll' in tags):
                self.add_element(rss_item, 'comments', hn_url)

            self.add_element(rss_item, 'guid', hn_url, isPermaLink='false')

    def response(self):
        rss_xml = etree.tostring(
            self.rss_root, pretty_print=True, encoding='UTF-8', xml_declaration=True,
        )

        if self.api_response.get('hits'):
            latest = max(hit['created_at_i'] for hit in self.api_response['hits'])
            last_modified = self.generate_rfc2822(latest)

            # Set max-age=N to the average number of seconds between new items
            timestamps = sorted(map(lambda h: h['created_at_i'], self.api_response['hits']), reverse=True)
            seconds = sum(timestamps[idx] - timestamps[idx+1] for idx in xrange(0, len(timestamps) - 1)) / float(len(timestamps))
        else:
            last_modified = self.generate_rfc2822()
            seconds = 5 * 60

        # Cap between 5 minutes and 1 hour
        if seconds < (5 * 60):
            seconds = (5 * 60)
        elif seconds > (60 * 60):
            seconds = (60 * 60)

        headers = {
            'Content-Type': 'text/xml; charset=utf-8',
            'Last-Modified': last_modified.replace('+0000', 'GMT'),
            'Cache-Control': 'max-age=%d' % int(seconds),
            'Expires': self.generate_rfc2822(int(time.time() + seconds)).replace('+0000', 'GMT'),
        }

        return (rss_xml, 200, headers)

    def add_element(self, parent, tag, text, **attrs):
        el = etree.Element(tag, attrs)
        el.text = text
        parent.append(el)
        return el

    def generate_rfc2822(self, secs=None):
        t = time.gmtime(secs)
        return time.strftime('%a, %d %b %Y %H:%M:%S +0000', t)
