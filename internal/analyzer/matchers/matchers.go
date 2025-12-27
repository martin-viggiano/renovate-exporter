package matchers

import "github.com/martin-viggiano/renovate-exporter/internal/analyzer"

func DefaultMatchers() []analyzer.Matcher {
	return []analyzer.Matcher{
		NewRepositoryMatcher(),
		NewPullRequestMatcher(),
	}
}

type matcher struct {
	name      string
	predicate func(e *analyzer.Entry) bool
	extract   func(e *analyzer.Entry, metrics analyzer.Metrics)
}

func (m *matcher) Name() string {
	return m.name
}

func (m *matcher) Predicate(e *analyzer.Entry) bool {
	return m.predicate(e)
}

func (m *matcher) Extract(e *analyzer.Entry, metrics analyzer.Metrics) {
	m.extract(e, metrics)
}
