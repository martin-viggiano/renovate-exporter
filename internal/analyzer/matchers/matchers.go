package matchers

import "github.com/martin-viggiano/renovate-exporter/internal/analyzer"

func DefaultMatchers() []analyzer.Matcher {
	return []analyzer.Matcher{
		NewRepositoryMatcher(),
	}
}

type matcher struct {
	name      string
	predicate func(e *analyzer.LogEntry) bool
	extract   func(e *analyzer.LogEntry, metrics analyzer.Metrics)
}

func (m *matcher) Name() string {
	return m.name
}

func (m *matcher) Predicate(e *analyzer.LogEntry) bool {
	return m.predicate(e)
}

func (m *matcher) Extract(e *analyzer.LogEntry, metrics analyzer.Metrics) {
	m.extract(e, metrics)
}
