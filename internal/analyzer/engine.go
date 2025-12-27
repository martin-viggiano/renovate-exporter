package analyzer

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

type Matcher interface {
	Name() string
	Predicate(e *LogEntry) bool
	Extract(e *LogEntry, m Metrics)
}

type Skipper struct {
	Name      string
	Predicate func(*LogEntry) bool
}

type Engine struct {
	metrics  *Metrics
	skippers []Skipper
	matchers []Matcher
}

func NewEngine(reg *prometheus.Registry, matchers []Matcher) (*Engine, error) {
	m, err := newMetrics(reg)
	if err != nil {
		return nil, err
	}

	return &Engine{
		metrics:  m,
		skippers: buildSkippers(m),
		matchers: matchers,
	}, nil
}

// Process receives a byte slice, parses it into a log entry and evaluates the skippers and matchers.
func (e *Engine) Process(data []byte) error {
	entry, err := Parse(data)
	if err != nil {
		return err
	}

	for _, s := range e.skippers {
		if s.Predicate(entry) {
			slog.Debug("skipped processing of log entry", slog.String("skipper", s.Name))
			return nil
		}
	}

	for _, m := range e.matchers {
		if m.Predicate(entry) {
			slog.Debug("extracted metrics from log entry", slog.String("matcher", m.Name()))
			m.Extract(entry, *e.metrics)
			return nil
		}
	}

	return nil
}

func (e *Engine) Metrics() *Metrics {
	return e.metrics
}
