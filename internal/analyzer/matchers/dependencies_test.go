package matchers_test

import (
	"os"
	"strings"
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

	data, err := os.ReadFile("testdata/dependency_extraction_complete.txt")
	require.NoError(t, err)

	err = engine.Process([]byte(strings.TrimSpace(string(data))))
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

	data, err := os.ReadFile("testdata/dependency_extraction_complete.txt")
	require.NoError(t, err)

	err = engine.Process([]byte(strings.TrimSpace(string(data))))
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

func TestOutdatedDependenciesMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewOutdatedDependencyMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	data, err := os.ReadFile("testdata/packagefiles_with_updates.txt")
	require.NoError(t, err)

	err = engine.Process([]byte(strings.TrimSpace(string(data))))
	assert.NoError(t, err)

	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().DependencyOutdatedTotal.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyOutdatedTotal.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().DependencyOutdatedTotal.WithLabelValues("test", "renovate-config-presets"),
	))
}
