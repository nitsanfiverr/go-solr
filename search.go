package solr

import (
	"fmt"
	"net/url"
	"strings"
)


type Query struct {
	params url.Values
}

func NewQuery() *Query {
	q := new(Query)
	q.params = url.Values{}
	return q
}

func(q *Query) AddParam(k string, v string) {
	q.params.Add(k, v)
}

func(q *Query) RemoveParam(k string) {
	q.params.Del(k)
}

func (q *Query) String() string {
	return q.params.Encode()
}

type Search struct {
	queries []*Query
	conn    *Connection
	start   int
	rows    int
	debug   string
}

// NewSearch takes c and q as optional
func NewSearch(c *Connection, q *Query) * Search {
	s := new(Search)
	if q != nil {
		s.AddQuery(q)
	}
	if c != nil {
		s.conn = c
	}
	return s
}


func (s *Search) Query() *Query {
	q := NewQuery()
	s.AddQuery(q)
	return q
}

func (s *Search) AddQuery(q *Query) {
	s.queries = append(s.queries, q)
}

// QueryString return a query string of all queries, including start, rows debug and wt=json
// wt is always json
func (s *Search) QueryString() string {

	query := []string{"wt=json"}
	
	if s.start > 0 {
		query = append(query, fmt.Sprintf("start=%d", s.start))
	}

	if s.rows > 0 {
		query = append(query, fmt.Sprintf("rows=%d", s.rows))
	}

	if s.debug != "" {
		query = append(query, fmt.Sprintf("debug=%s&indent=true", s.debug))
	}

	if len(s.queries) > 0 {
		for _, v := range s.queries {
			query = append(query, v.String())
		}
	}

	return strings.Join(query, "&")
}

// Result will create a StandardResultParser if no parser specified
// parser must be an implement ResultParser interface
func (s *Search) Result(parser ResultParser) (*SolrResult, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("No connection found for making request to solr")
	}
	resp, err := s.conn.Select(s.QueryString())
	if err != nil {
		return nil, err
	}
	if parser == nil {
		parser = new(StandardResultParser)
	}
	return parser.Parse(resp)
}
