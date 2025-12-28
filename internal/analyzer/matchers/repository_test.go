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

// func TestRepositoryMatcher(t *testing.T) {
// 	reg := prometheus.NewRegistry()

// 	matchers := []analyzer.Matcher{
// 		matchers.NewRepositoryMatcher(),
// 	}

// 	engine, err := analyzer.NewEngine(reg, matchers)
// 	require.NoError(t, err)

// 	err = engine.Process([]byte(`{"name":"renovate","hostname":"22111cc51078","pid":96,"level":30,"logContext":"4b569c43-f97c-43eb-ae32-443e49ca1b89","repository":"test/repos","renovateVersion":"42.42.2","msg":"Repository started","time":"2025-12-23T22:22:48.654Z","v":0}`))
// 	assert.NoError(t, err)

// 	value := testutil.ToFloat64(
// 		engine.Metrics().Repositories.WithLabelValues("test/repos"),
// 	)

// 	assert.Equal(t, float64(1), value)
// }

func TestRepositoryStatusMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewRepositoryStatusMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	// Set metric for same repository and different status.
	engine.Metrics().Repositories.WithLabelValues("test/repos", "unknown").Set(1)

	require.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().Repositories.WithLabelValues("test/repos", "unknown"),
	))

	err = engine.Process([]byte(`{"name":"renovate","hostname":"5d4f86fd4030","pid":20,"level":30,"logContext":"da05cd78-e34e-4bf2-b1a9-c9b6aae13710","repository":"test/repos","cloned":true,"durationMs":19174,"status":"onboarding","enabled":true,"onboarded":false,"msg":"Repository finished","time":"2025-12-28T00:15:36.757Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().Repositories.WithLabelValues("test/repos", "onboarding"),
	))
	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().Repositories.WithLabelValues("test/repos", "unknown"),
	))
}
