package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	tempDir := t.TempDir()

	testLogPath := filepath.Join(tempDir, "renovate.log")
	sourceData, err := os.ReadFile("testdata/complete_logs.txt")
	require.NoError(t, err, "failed to read test log file")

	testPort := "19090"
	metricsURL := fmt.Sprintf("http://localhost:%s/metrics", testPort)

	rootCmd.SetArgs([]string{
		"--path", tempDir,
		"--address", fmt.Sprintf(":%s", testPort),
	})

	// Run command in background
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cmdErrCh := make(chan error, 1)
	go func() {
		cmdErrCh <- rootCmd.ExecuteContext(ctx)
	}()

	// Wait for program to start
	assert.Eventually(t, func() bool {
		_, err := net.Dial("tcp", "localhost:19090")
		return err == nil
	}, 10*time.Second, time.Second)

	err = os.WriteFile(testLogPath, sourceData, 0o644)
	require.NoError(t, err, "failed to write test log to temp directory")

	<-time.After(5 * time.Second)

	var metricsBody string
	assert.Eventually(t, func() bool {
		resp, err := http.Get(metricsURL)
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return false
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return false
		}

		metricsBody = string(body)
		// Check if we have some metrics (at least one renovate metric)
		return strings.Contains(metricsBody, "renovate_")
	}, 5*time.Second, 100*time.Millisecond, "metrics endpoint did not become available with metrics")

	t.Log("Successfully connected to metrics endpoint")

	t.Log(metricsBody)

	// Parse metrics
	parser := expfmt.NewTextParser(model.UTF8Validation)
	metricFamilies, err := parser.TextToMetricFamilies(strings.NewReader(metricsBody))
	require.NoError(t, err, "failed to parse metrics")

	// Helper function to get metric value
	getMetricValue := func(name string, labels map[string]string) (float64, bool) {
		family, ok := metricFamilies[name]
		if !ok {
			return 0, false
		}

		for _, metric := range family.GetMetric() {
			if matchesLabels(metric, labels) {
				if metric.Gauge != nil {
					return metric.Gauge.GetValue(), true
				}
				if metric.Counter != nil {
					return metric.Counter.GetValue(), true
				}
			}
		}
		return 0, false
	}

	t.Log("Collected metrics:")
	for name, family := range metricFamilies {
		t.Logf("  - %s: %d samples", name, len(family.GetMetric()))
	}

	// Verify that core metrics exist
	require.NotEmpty(t, metricFamilies, "no metrics were collected")
	assert.Contains(t, metricFamilies, "renovate_repositories", "missing renovate_repositories metric")

	tt := []struct {
		name      string
		labels    map[string]string
		wantValue float64
	}{
		{
			name: "renovate_repositories",
			labels: map[string]string{
				"repository": "test/repo",
				"status":     "activated",
			},
			wantValue: 1,
		},
		{
			name: "renovate_repositories_libyears",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "gitlabci",
			},
			wantValue: 0,
		},
		{
			name: "renovate_repositories_libyears",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "gomod",
			},
			wantValue: 6.5239691146626075,
		},
		{
			name: "renovate_repositories_libyears",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "renovate-config-presets",
			},
			wantValue: 0,
		},
		{
			name: "renovate_pull_requests_total",
			labels: map[string]string{
				"repository": "test/repo",
				"state":      "open",
			},
			wantValue: 5,
		},
		{
			name: "renovate_pull_requests_total",
			labels: map[string]string{
				"repository": "test/repo",
				"state":      "merged",
			},
			wantValue: 4,
		},
		{
			name: "renovate_pull_requests_total",
			labels: map[string]string{
				"repository": "test/repo",
				"state":      "closed",
			},
			wantValue: 6,
		},
		{
			name: "renovate_dependency_outdated_total",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "gitlabci",
			},
			wantValue: 0,
		},
		{
			name: "renovate_dependency_outdated_total",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "gomod",
			},
			wantValue: 6,
		},
		{
			name: "renovate_dependency_outdated_total",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "renovate-config-presets",
			},
			wantValue: 0,
		},
		{
			name: "renovate_dependencies_total",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "gitlabci",
			},
			wantValue: 1,
		},
		{
			name: "renovate_dependencies_total",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "gomod",
			},
			wantValue: 18,
		},
		{
			name: "renovate_dependencies_total",
			labels: map[string]string{
				"repository": "test/repo",
				"manager":    "renovate-config-presets",
			},
			wantValue: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			value, ok := getMetricValue(tc.name, tc.labels)
			assert.True(t, ok, "metric not found")
			assert.Equal(t, tc.wantValue, value)
		})
	}

	// Cancel context to stop the command
	cancel()

	select {
	case err := <-cmdErrCh:
		assert.NoError(t, err, "command should gracefully shut down")
	case <-time.After(2 * time.Second):
		t.Fatal("command did not stop in time")
	}
}

// matchesLabels checks if a metric's labels match the expected labels
func matchesLabels(metric *dto.Metric, expected map[string]string) bool {
	if len(expected) == 0 {
		return len(metric.Label) == 0
	}

	for _, label := range metric.Label {
		expectedValue, ok := expected[label.GetName()]
		if !ok {
			continue
		}
		if label.GetValue() != expectedValue {
			return false
		}
	}

	return true
}
