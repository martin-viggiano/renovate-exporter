package analyzer

import (
	"encoding/json"
	"time"
)

type Entry struct {
	V          int       `json:"v"`
	Time       time.Time `json:"time"`
	PID        int       `json:"pid"`
	Name       string    `json:"name"`
	Message    string    `json:"msg"`
	LogContext string    `json:"logContext"`
	Level      int       `json:"level"`
	Hostname   string    `json:"hostname"`

	Repository string `json:"repository,omitempty"`
	Cloned     bool   `json:"cloned,omitempty"`
	Duration   int    `json:"durationMs,omitempty"`
	Status     string `json:"status,omitempty"`
	Onboarded  bool   `json:"onboarded,omitempty"`

	RawStats json.RawMessage `json:"stats,omitempty"`

	PullRequestStatistics *PullRequestStatistics `json:"-"`
}

// PullRequestStatistics contains information about the merge requests.
type PullRequestStatistics struct {
	Total  int `json:"total"`
	Open   int `json:"open"`
	Closed int `json:"closed"`
	Merged int `json:"merged"`
}
