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

func TestPullRequestMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewPullRequestMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	data, err := os.ReadFile("testdata/repository_pr_statistics.txt")
	require.NoError(t, err)

	err = engine.Process(data)
	assert.NoError(t, err)

	assert.Equal(t, float64(3), testutil.ToFloat64(
		engine.Metrics().PullRequests.WithLabelValues("test", "open"),
	))
	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().PullRequests.WithLabelValues("test", "merged"),
	))
	assert.Equal(t, float64(4), testutil.ToFloat64(
		engine.Metrics().PullRequests.WithLabelValues("test", "closed"),
	))
}
