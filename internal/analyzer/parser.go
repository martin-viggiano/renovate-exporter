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

func Parse(data []byte) (*Entry, error) {
	e := Entry{}

	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}

	if e.RawStats != nil {
		switch e.Message {
		case "Renovate repository PR statistics":
			var s PullRequestStatistics
			if err := json.Unmarshal(e.RawStats, &s); err != nil {
				return nil, err
			}

			e.PullRequestStatistics = &s
		}
	}

	return &e, nil
}
