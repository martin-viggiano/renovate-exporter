package matchers

import (
	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/prometheus/client_golang/prometheus"
)

func NewLibYearMatcher() analyzer.Matcher {
	return &matcher{
		name: "libyear",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Repository libYears" && e.Repository != "" && e.LibYearStatistics != nil
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			// Delete metrics for the same repository to be able to update the status.
			metrics.RepositoryLibyears.DeletePartialMatch(prometheus.Labels{
				"repository": e.Repository,
			})
			for manager, value := range e.LibYearStatistics.Managers {
				metrics.RepositoryLibyears.WithLabelValues(e.Repository, manager).Set(value)
			}
		},
	}
}
