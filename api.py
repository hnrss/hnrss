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
        if query is not None:
            if ' OR ' in query:
                components = query.replace(' OR ', ' ')
                if '"' in components:
                    quoted_terms = components
                else:
                    quoted_terms = ' '.join('"%s"' % t for t in components.split())
                self.params['query'] = quoted_terms
                self.params['optionalWords'] = quoted_terms
            else:
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

    def jobs(self):
        return self._request('job')

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

    def who_is_hiring(self, include='all'):
        submitted = self.user('whoishiring', 'submitted')
        hits = submitted.get('hits', [])

        if include == 'all':
            thread_ids = [hit['objectID'] for hit in hits]
        elif include == 'jobs':
            thread_ids = [hit['objectID'] for hit in hits if 'Ask HN: Who is hiring?' in hit['title']]
        elif include == 'hired':
            thread_ids = [hit['objectID'] for hit in hits if 'Ask HN: Who wants to be hired?' in hit['title']]
        elif include == 'freelance':
            thread_ids = [hit['objectID'] for hit in hits if 'Ask HN: Freelancer? Seeking freelancer?' in hit['title']]

        thread_ids = map(int, thread_ids)
        story_ids = ['story_%d' % thread_id for thread_id in thread_ids]
        tags = 'comment,(%s)' % ','.join(story_ids)

        response = self._request(tags)
        slim_response = response.copy()
        slim_response['hits'] = []

        # Only include top-level comments
        for hit in response.get('hits', []):
            if hit['parent_id'] in thread_ids:
                slim_response['hits'].append(hit)

        return slim_response
