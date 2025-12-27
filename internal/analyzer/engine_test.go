package analyzer_test

import (
	"testing"

	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngine(t *testing.T) {
	reg := prometheus.NewRegistry()

	matched := make(map[string]struct{})
	extracted := make(map[string]struct{})

	engine, err := analyzer.NewEngine(reg, []analyzer.Matcher{
		&matcherMock{
			name: "all",
			predicate: func(e *analyzer.LogEntry) bool {
				matched[e.Message] = struct{}{}
				return true
			},
			extract: func(e *analyzer.LogEntry, metrics analyzer.Metrics) {
				extracted[e.Message] = struct{}{}
			},
		},
	})
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"22111cc51078","pid":96,"level":30,"logContext":"4b569c43-f97c-43eb-ae32-443e49ca1b89","repository":"test/repos","renovateVersion":"42.42.2","msg":"Repository started","time":"2025-12-23T22:22:48.654Z","v":0}`))
	assert.NoError(t, err)

	assert.Contains(t, matched, `Repository started`)
	assert.Contains(t, extracted, `Repository started`)
}

type matcherMock struct {
	name      string
	predicate func(e *analyzer.LogEntry) bool
	extract   func(e *analyzer.LogEntry, metrics analyzer.Metrics)
}

func (m *matcherMock) Name() string {
	return m.name
}

func (m *matcherMock) Predicate(e *analyzer.LogEntry) bool {
	return m.predicate(e)
}

func (m *matcherMock) Extract(e *analyzer.LogEntry, metrics analyzer.Metrics) {
	m.extract(e, metrics)
}
