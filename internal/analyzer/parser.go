package analyzer

import (
	"encoding/json"
)

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
		case "Dependency extraction complete":
			var s ManagerStatistics
			if err := json.Unmarshal(e.RawStats, &s); err != nil {
				return nil, err
			}

			e.ManagerStatistics = &s
		}
	}

	if e.RawConfig != nil && e.Message == "packageFiles with updates" {
		var config PackageFilesConfig
		if err := json.Unmarshal(e.RawConfig, &config.Managers); err != nil {
			return nil, err
		}

		e.PackageFilesConfig = &config
	}

	return &e, nil
}
