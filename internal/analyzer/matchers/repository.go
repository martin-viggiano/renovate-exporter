package matchers

import (
	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/prometheus/client_golang/prometheus"
)

func NewRepositoryDurationMatcher() analyzer.Matcher {
	return &matcher{
		name: "repository",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Repository finished" && e.Repository != "" && e.Duration > 0
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			metrics.RepositoryDuration.WithLabelValues(e.Repository).Set(float64(e.Duration) / 1000)
		},
	}
}

func NewRepositoryStatusMatcher() analyzer.Matcher {
	return &matcher{
		name: "repository",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Repository finished" && e.Repository != "" && e.Status != ""
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			// Delete metrics for the same repository to be able to update the status.
			metrics.Repositories.DeletePartialMatch(prometheus.Labels{
				"repository": e.Repository,
			})
			metrics.Repositories.WithLabelValues(e.Repository, e.Status).Set(1)
		},
	}
}
