package matchers

import "github.com/martin-viggiano/renovate-exporter/internal/analyzer"

func NewRepositoryMatcher() analyzer.Matcher {
	return &matcher{
		name: "repository",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Repository started" && e.Repository != ""
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			metrics.Repositories.WithLabelValues(e.Repository).Set(1)
		},
	}
}
