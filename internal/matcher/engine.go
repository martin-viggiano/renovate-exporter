package matcher

import (
	"log/slog"

	"github.com/martin-viggiano/renovate-exporter/internal/logentry"
	"github.com/martin-viggiano/renovate-exporter/internal/registry"
)

type Matcher struct {
	Name      string
	Predicate func(*logentry.LogEntry) bool
	Extract   func(*logentry.LogEntry, *registry.Registry)
}

type Skipper struct {
	Name      string
	Predicate func(*logentry.LogEntry) bool
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

func (e *Engine) Process(entry *logentry.LogEntry) {
	for _, s := range e.skippers {
		if s.Predicate(entry) {
			slog.Debug("skipped processing of log entry", slog.String("skipper", s.Name))
			return
		}
	}

	for _, m := range e.matchers {
		if m.Predicate(entry) {
			slog.Debug("extracted metrics from log entry", slog.String("matcher", m.Name))
			m.Extract(entry, e.metrics)
		}
	}
}
