package matchers

import (
	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/prometheus/client_golang/prometheus"
)

func NewDependencyCountMatcher() analyzer.Matcher {
	return &matcher{
		name: "total dependencies",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Dependency extraction complete" && e.Repository != "" && e.ManagerStatistics != nil
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			// Delete metrics for the same repository to be able to update the status.
			metrics.DependenciesTotal.DeletePartialMatch(prometheus.Labels{
				"repository": e.Repository,
			})
			for manager, value := range e.ManagerStatistics.Managers {
				metrics.DependenciesTotal.WithLabelValues(e.Repository, manager).Set(float64(value.DependencyCount))
			}
		},
	}
}

func NewDependencyFilesMatcher() analyzer.Matcher {
	return &matcher{
		name: "dependency files",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Dependency extraction complete" && e.Repository != "" && e.ManagerStatistics != nil
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			// Delete metrics for the same repository to be able to update the status.
			metrics.DependencyFilesTotal.DeletePartialMatch(prometheus.Labels{
				"repository": e.Repository,
			})
			for manager, value := range e.ManagerStatistics.Managers {
				metrics.DependencyFilesTotal.WithLabelValues(e.Repository, manager).Set(float64(value.FileCount))
			}
		},
	}
}
