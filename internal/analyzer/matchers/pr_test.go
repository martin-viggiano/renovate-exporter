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

func TestPullRequestMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewPullRequestMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"22111cc51078","pid":96,"level":20,"logContext":"4b569c43-f97c-43eb-ae32-443e49ca1b89","repository":"test","stats":{"total":7,"open":3,"closed":4,"merged":0},"msg":"Renovate repository PR statistics","time":"2025-12-23T22:23:16.987Z","v":0}`))
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
