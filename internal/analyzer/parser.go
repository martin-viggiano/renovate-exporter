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

	Stats *PullRequestStatistics `json:"stats,omitempty"`
}

// PullRequestStatistics contains information about the merge requests.
type PullRequestStatistics struct {
	Total  int `json:"total"`
	Open   int `json:"open"`
	Closed int `json:"closed"`
	Merged int `json:"merged"`
}

func Parse(data []byte) (*Entry, error) {
	le := Entry{}

	if err := json.Unmarshal(data, &le); err != nil {
		return nil, err
	}

	return &le, nil
}
