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

func TestLibYearMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewLibYearMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	// Set metric for same repository and different status.
	engine.Metrics().RepositoryLibyears.WithLabelValues("test", "unknown").Set(1)

	require.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().RepositoryLibyears.WithLabelValues("test", "unknown"),
	))

	err = engine.Process([]byte(`{"name":"renovate","hostname":"5d4f86fd4030","pid":173,"level":20,"logContext":"f63eb89c-1099-4bc3-b7b4-745eb6f6f3ba","repository":"test","libYears":{"managers":{"gitlabci":0,"gomod":3.347419393708777,"renovate-config-presets":0},"total":3.347419393708777},"dependencyStatus":{"outdated":5,"total":19},"msg":"Repository libYears","time":"2025-12-28T00:56:16.223Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().RepositoryLibyears.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(3.347419393708777), testutil.ToFloat64(
		engine.Metrics().RepositoryLibyears.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().RepositoryLibyears.WithLabelValues("test", "renovate-config-presets"),
	))
}
