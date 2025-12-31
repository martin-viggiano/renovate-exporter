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

func TestDependencyCountMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewDependencyCountMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"247242df26cb","pid":96,"level":30,"logContext":"9c28e67e-ec9c-43dc-a985-2cf6b447aa6f","repository":"test","baseBranch":"master","stats":{"managers":{"gitlabci":{"fileCount":1,"depCount":1},"gomod":{"fileCount":1,"depCount":17},"renovate-config-presets":{"fileCount":1,"depCount":1}},"total":{"fileCount":3,"depCount":19}},"msg":"Dependency extraction complete","time":"2025-12-28T13:46:46.486Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependenciesTotal.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(17), testutil.ToFloat64(
		engine.Metrics().DependenciesTotal.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependenciesTotal.WithLabelValues("test", "renovate-config-presets"),
	))
}

func TestDependencyFilesMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewDependencyFilesMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"247242df26cb","pid":96,"level":30,"logContext":"9c28e67e-ec9c-43dc-a985-2cf6b447aa6f","repository":"test","baseBranch":"master","stats":{"managers":{"gitlabci":{"fileCount":1,"depCount":1},"gomod":{"fileCount":1,"depCount":17},"renovate-config-presets":{"fileCount":1,"depCount":1}},"total":{"fileCount":3,"depCount":19}},"msg":"Dependency extraction complete","time":"2025-12-28T13:46:46.486Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyFilesTotal.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyFilesTotal.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyFilesTotal.WithLabelValues("test", "renovate-config-presets"),
	))
}
