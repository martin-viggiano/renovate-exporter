package matchers_test

import (
	"os"
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

	data, err := os.ReadFile("testdata/repository_libyears.txt")
	require.NoError(t, err)

	// Set metric for same repository and different status.
	engine.Metrics().RepositoryLibyears.WithLabelValues("test", "unknown").Set(1)

	require.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().RepositoryLibyears.WithLabelValues("test", "unknown"),
	))

	err = engine.Process(data)
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
