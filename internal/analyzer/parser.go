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
		}
	}

	return &e, nil
}
