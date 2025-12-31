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

	RawStats  json.RawMessage `json:"stats,omitempty"`
	RawConfig json.RawMessage `json:"config,omitempty"`

	LibYearStatistics *LibYearStatistics `json:"libYears,omitempty"`

	PullRequestStatistics *PullRequestStatistics `json:"-"` // Present only in the "Renovate repository PR statistics" log message
	ManagerStatistics     *ManagerStatistics     `json:"-"` // Present only in the "Dependency extraction complete" log message
	PackageFilesConfig    *PackageFilesConfig    `json:"-"` // Present only in the "packageFiles with updates" log message
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
	// Map of manager name (e.g., "gomod", "gitlabci") to manager data
	Managers map[string]ManagerData `json:"managers"`
}

// PackageFilesConfig represents the config field in "packageFiles with updates" log
type PackageFilesConfig struct {
	// Map of manager name (e.g., "gomod", "gitlabci") to array of package files
	Managers map[string][]PackageFile
}

// PackageFile represents a single package file with its dependencies
type PackageFile struct {
	PackageFile string       `json:"packageFile,omitempty"`
	Deps        []Dependency `json:"deps"`
}

// Dependency represents a single dependency
type Dependency struct {
	DepName                 string   `json:"depName"`
	PackageName             string   `json:"packageName,omitempty"`
	CurrentValue            string   `json:"currentValue,omitempty"`
	CurrentVersion          string   `json:"currentVersion,omitempty"`
	CurrentVersionAgeInDays *int     `json:"currentVersionAgeInDays,omitempty"`
	DepType                 string   `json:"depType,omitempty"`
	Datasource              string   `json:"datasource,omitempty"`
	Versioning              string   `json:"versioning,omitempty"`
	Enabled                 *bool    `json:"enabled,omitempty"`
	SkipReason              string   `json:"skipReason,omitempty"`
	Updates                 []Update `json:"updates"`
	FixedVersion            string   `json:"fixedVersion,omitempty"`
	IsSingleVersion         bool     `json:"isSingleVersion,omitempty"`

	// Timestamps
	CurrentVersionTimestamp string `json:"currentVersionTimestamp,omitempty"`
	MostRecentTimestamp     string `json:"mostRecentTimestamp,omitempty"`
}

// Update represents an available update for a dependency
type Update struct {
	Bucket              string   `json:"bucket,omitempty"`
	NewVersion          string   `json:"newVersion"`
	NewValue            string   `json:"newValue,omitempty"`
	NewVersionAgeInDays *int     `json:"newVersionAgeInDays,omitempty"`
	ReleaseTimestamp    string   `json:"releaseTimestamp,omitempty"`
	UpdateType          string   `json:"updateType,omitempty"` // "patch", "minor", "major"
	IsBreaking          bool     `json:"isBreaking"`
	LibYears            *float64 `json:"libYears,omitempty"`
	BranchName          string   `json:"branchName,omitempty"`

	// Semantic versioning fields
	NewMajor int `json:"newMajor,omitempty"`
	NewMinor int `json:"newMinor,omitempty"`
	NewPatch int `json:"newPatch,omitempty"`
}
