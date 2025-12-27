package registry

import "github.com/prometheus/client_golang/prometheus"

type Registry struct {
	Metrics Metrics
}

type Metrics struct{}

func New(reg *prometheus.Registry) (*Registry, error) {
	return &Registry{
		Metrics: Metrics{
			// metrics would be registered here
		},
	}, nil
}
