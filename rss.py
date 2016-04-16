import re
import time
from xml.sax.saxutils import unescape as sax_unescape
from lxml import etree

def unescape(s):
    deref_ncr = lambda m: unichr(int(m.group(1), 16)) # '&#x2F;' -> '/'
    s = re.sub('&#[Xx]([A-Fa-f0-9]+);', deref_ncr, s)
    entities = {'&quot;': '"', '&apos;': "'"}
    return sax_unescape(s, entities)

class RSS(object):
    def __init__(self, api_response, title, link='https://news.ycombinator.com/'):
        self.api_response = api_response

        self.rss_root = etree.Element('rss', version='2.0')
        self.rss_channel = etree.SubElement(self.rss_root, 'channel')

        self.add_element(self.rss_channel, 'title', title)
        self.add_element(self.rss_channel, 'link', link)
        self.add_element(self.rss_channel, 'description', 'Hacker News RSS')
        self.add_element(self.rss_channel, 'docs', 'https://edavis.github.io/hnrss/')
        self.add_element(self.rss_channel, 'generator', 'https://github.com/edavis/hnrss')
        self.add_element(self.rss_channel, 'lastBuildDate', self.generate_rfc2822())

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
                elif self.api_response['description'] and self.api_response['link_to'] == 'url':
                    template = (
                        'Article URL: <a href="%(url)s">%(url)s</a><br>'
                        'Comments URL: <a href="%(hn_url)s">%(hn_url)s</a><br>'
                        'Points: %(points)s'
                    )
                    params = {
                        'url': hit.get('url') or hn_url,
                        'hn_url': hn_url,
                        'points': hit.get('points') or 0,
                    }
                    self.add_element(rss_item, 'description', template % params)
                elif self.api_response['description'] and self.api_response['link_to'] == 'comments':
                    self.add_element(rss_item, 'description', 'Article URL: <a href="%(url)s">%(url)s</a>' % {'url': hit.get('url') or hn_url})

            self.add_element(rss_item, 'pubDate', self.generate_rfc2822(hit.get('created_at_i')))

            if self.api_response['link_to'] == 'comments':
                self.add_element(rss_item, 'link', hn_url)
            else:
                self.add_element(rss_item, 'link', hit.get('url') or hn_url)

            self.add_element(rss_item, 'author', hit.get('author'))

            if ('story' in tags or 'poll' in tags):
                self.add_element(rss_item, 'comments', hn_url)

            self.add_element(rss_item, 'guid', hn_url, isPermaLink='false')

    def response(self):
        rss_xml = etree.tostring(
            self.rss_root, pretty_print=True, encoding='UTF-8', xml_declaration=True,
        )
        return (rss_xml, 200, {'Content-Type': 'text/xml'})

    def add_element(self, parent, tag, text, **attrs):
        el = etree.Element(tag, attrs)
        el.text = text
        parent.append(el)
        return el

    def generate_rfc2822(self, secs=None):
        t = time.gmtime(secs)
        return time.strftime('%a, %d %b %Y %H:%M:%S +0000', t)
