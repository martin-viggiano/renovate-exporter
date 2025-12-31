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

func NewOutdatedDependencyMatcher() analyzer.Matcher {
	return &matcher{
		name: "package files",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "packageFiles with updates" && e.Repository != "" && e.PackageFilesConfig != nil
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			// Delete metrics for the same repository to be able to update the status.
			metrics.DependencyOutdatedTotal.DeletePartialMatch(prometheus.Labels{
				"repository": e.Repository,
			})
			for manager, configs := range e.PackageFilesConfig.Managers {
				for _, config := range configs {
					outdatedCount := 0
					for _, dep := range config.Deps {
						if len(dep.Updates) > 0 {
							outdatedCount++
						}
					}

					metrics.DependencyOutdatedTotal.WithLabelValues(e.Repository, manager).Set(float64(outdatedCount))
				}
			}
		},
	}
}
