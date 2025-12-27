package analyzer

import (
	"log/slog"

	"github.com/martin-viggiano/renovate-exporter/internal/registry"
)

type Matcher struct {
	Name      string
	Predicate func(*LogEntry) bool
	Extract   func(*LogEntry, *registry.Registry)
}

type Skipper struct {
	Name      string
	Predicate func(*LogEntry) bool
}

type Engine struct {
	metrics  *registry.Registry
	skippers []Skipper
	matchers []Matcher
}

func NewEngine(reg *registry.Registry) *Engine {
	return &Engine{
		metrics:  reg,
		skippers: buildSkippers(reg),
		matchers: buildMatchers(reg),
	}
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
			slog.Debug("extracted metrics from log entry", slog.String("matcher", m.Name))
			m.Extract(entry, e.metrics)
			return nil
		}
	}

	return nil
}
