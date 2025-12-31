package analyzer

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Repositories            *prometheus.GaugeVec
	RepositoryDuration      *prometheus.GaugeVec
	RepositoryLibyears      *prometheus.GaugeVec
	PullRequests            *prometheus.GaugeVec
	DependenciesTotal       *prometheus.GaugeVec
	DependencyFilesTotal    *prometheus.GaugeVec
	DependencyOutdatedTotal *prometheus.GaugeVec
}

func newMetrics(reg *prometheus.Registry) (*Metrics, error) {
	m := &Metrics{
		Repositories: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_repositories",
				Help: "Current number of repositories being watched by Renovate.",
			},
			[]string{"repository", "status"},
		),
		RepositoryDuration: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_repositories_duration_seconds",
				Help: "Duration of the job for each of the repositories being watched by Renovate.",
			},
			[]string{"repository"},
		),
		RepositoryLibyears: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_repositories_libyears",
				Help: "Libyear of the repositories being watched by Renovate.",
			},
			[]string{"repository", "manager"},
		),
		PullRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_pull_requests_total",
				Help: "Current number of Pull Requests being managed by Renovate.",
			},
			[]string{"repository", "state"},
		),
		DependenciesTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_dependencies_total",
				Help: "Total number of dependencies per repository and manager discovered by Renovate.",
			},
			[]string{"repository", "manager"},
		),
		DependencyFilesTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_dependency_files_total",
				Help: "Number of dependency files per repository and manager discovered by Renovate.",
			},
			[]string{"repository", "manager"},
		),
		DependencyOutdatedTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_dependency_outdated_total",
				Help: "Total number of outdated dependencies per repository and manager discovered by Renovate.",
			},
			[]string{"repository", "manager"},
		),
	}

	reg.MustRegister(m.Repositories)
	reg.MustRegister(m.RepositoryDuration)
	reg.MustRegister(m.RepositoryLibyears)
	reg.MustRegister(m.PullRequests)
	reg.MustRegister(m.DependenciesTotal)
	reg.MustRegister(m.DependencyFilesTotal)
	reg.MustRegister(m.DependencyOutdatedTotal)

	return m, nil
}
