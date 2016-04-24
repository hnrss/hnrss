import requests

class API(object):
    base_url = 'https://hn.algolia.com/api/v1'
    count_limit = 100

    def __init__(self, points=None, comments=None, link_to='url',
                 query=None, search_attrs='title', description=True, count=None):
        self.endpoint = 'search_by_date'
        self.params = {}
        if points or comments:
            numeric_filters = []
            if points: numeric_filters.append('points>%s' % points)
            if comments: numeric_filters.append('num_comments>%s' % comments)
            self.params['numericFilters'] = ','.join(numeric_filters)
        if query:
            self.params['query'] = '"%s"' % query
            if search_attrs != 'default':
                self.params['restrictSearchableAttributes'] = search_attrs
        self.link_to = link_to
        self.description = description
        if count is not None:
            try:
                self.count = min(self.count_limit, int(count))
            except ValueError:
                pass
            else:
                self.params['hitsPerPage'] = self.count

    @classmethod
    def using_request(cls, request):
        return cls(
            points = request.args.get('points'),
            comments = request.args.get('comments'),
            link_to = request.args.get('link', 'url'),
            query = request.args.get('q'),
            search_attrs = request.args.get('search_attrs', 'title'),
            description = bool(int(request.args.get('description', 1))),
            count = request.args.get('count'),
        )

    def _request(self, tags):
        params = self.params.copy()
        params['tags'] = tags
        resp = requests.get(
            '%s/%s' % (self.base_url, self.endpoint),
            params = params,
        )
        resp.raise_for_status()
        obj = resp.json().copy()
        obj.update({
            'link_to': self.link_to,
            'description': self.description,
        })
        return obj

    def newest(self):
        return self._request('(story,poll)')

    def frontpage(self):
        return self._request('front_page')

    def ask_hn(self):
        return self._request('ask_hn')

    def show_hn(self):
        return self._request('show_hn')

    def polls(self):
        return self._request('poll')

    def comments(self, story_id=None):
        tags = ['comment']
        if story_id is not None:
            tags.append('story_%s' % story_id)
        return self._request(','.join(tags))

    def user(self, username, include='all'):
        tags = ['author_%s' % username]
        if include == 'all':
            tags.append('(story,poll,comment)')
        elif include == 'submitted':
            tags.append('(story,poll)')
        elif include == 'threads':
            tags.append('comment')
        return self._request(','.join(tags))
