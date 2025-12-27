package matcher

import (
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
