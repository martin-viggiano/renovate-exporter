package analyzer

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Repositories *prometheus.GaugeVec
	PullRequests *prometheus.GaugeVec
}

func newMetrics(reg *prometheus.Registry) (*Metrics, error) {
	m := &Metrics{
		Repositories: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_repositories",
				Help: "Current number of repositories being watched by Renovate.",
			},
			[]string{"repository"},
		),
		PullRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_pull_requests_total",
				Help: "Current number of Pull Requests being managed by Renovate.",
			},
			[]string{"repository", "state"},
		),
	}

	var err error
	err = errors.Join(err, reg.Register(m.Repositories))
	err = errors.Join(err, reg.Register(m.PullRequests))

	return m, err
}
