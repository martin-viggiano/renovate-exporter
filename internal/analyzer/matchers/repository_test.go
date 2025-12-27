package matchers_test

import (
	"testing"

	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/martin-viggiano/renovate-exporter/internal/analyzer/matchers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositoryMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewRepositoryMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"22111cc51078","pid":96,"level":30,"logContext":"4b569c43-f97c-43eb-ae32-443e49ca1b89","repository":"test/repos","renovateVersion":"42.42.2","msg":"Repository started","time":"2025-12-23T22:22:48.654Z","v":0}`))
	assert.NoError(t, err)

	value := testutil.ToFloat64(
		engine.Metrics().Repositories.WithLabelValues("test/repos"),
	)

	assert.Equal(t, float64(1), value)
}
