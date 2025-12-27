package analyzer

import (
	"errors"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Repositories *prometheus.GaugeVec
}

func newMetrics(reg *prometheus.Registry) (*Metrics, error) {
	m := &Metrics{
		prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "renovate_repositories",
				Help: "Current number of repositories being watched by Renovate.",
			},
			[]string{"repository"},
		),
	}

	err := errors.Join(reg.Register(m.Repositories))

	return m, err
}
