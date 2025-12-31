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

	LibYearStatistics *LibYearStatistics `json:"libYears,omitempty"`

	PullRequestStatistics *PullRequestStatistics `json:"-"` // Present only in the "Renovate repository PR statistics" log message
	ManagerStatistics     *ManagerStatistics     `json:"-"` // Present only in the "Dependency extraction complete" log message
}

// PullRequestStatistics contains information about the merge requests.
type PullRequestStatistics struct {
	Total  int `json:"total"`
	Open   int `json:"open"`
	Closed int `json:"closed"`
	Merged int `json:"merged"`
}

// LibYearStatistics contains information about the lib-year statistics.
type LibYearStatistics struct {
	Total    float64            `json:"total"`
	Managers map[string]float64 `json:"managers"`
}

// ManagerData contrains infromation about a manager's total dependencies
// and dependency file count.
type ManagerData struct {
	FileCount       int `json:"fileCount"`
	DependencyCount int `json:"depCount"`
}

// ManagerStatistics contrains infromation about all the manager's dependencies
// and dependency file count.
type ManagerStatistics struct {
	Managers map[string]ManagerData `json:"managers"`
}
