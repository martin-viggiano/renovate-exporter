package analyzer

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	reg := prometheus.NewRegistry()

	registry, err := newMetrics(reg)

	assert.NoError(t, err)
	assert.NotNil(t, registry)
}
