package registry

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestRegistry(t *testing.T) {
	reg := prometheus.NewRegistry()

	registry, err := New(reg)

	assert.NoError(t, err)
	assert.NotNil(t, registry)
}
