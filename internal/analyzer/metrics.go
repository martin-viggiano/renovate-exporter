package analyzer

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Repositories       *prometheus.GaugeVec
	RepositoryDuration *prometheus.GaugeVec
	RepositoryLibyears *prometheus.GaugeVec
	PullRequests       *prometheus.GaugeVec
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
	}

	reg.MustRegister(m.Repositories)
	reg.MustRegister(m.RepositoryDuration)
	reg.MustRegister(m.RepositoryLibyears)
	reg.MustRegister(m.PullRequests)

	return m, nil
}
