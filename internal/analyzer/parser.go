package analyzer

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	V          int       `json:"v"`
	Time       time.Time `json:"time"`
	PID        int       `json:"pid"`
	Name       string    `json:"name"`
	Message    string    `json:"msg"`
	LogContext string    `json:"logContext"`
	Level      int       `json:"level"`
	Hostname   string    `json:"hostname"`
}

func Parse(data []byte) (*LogEntry, error) {
	le := LogEntry{}

	if err := json.Unmarshal(data, &le); err != nil {
		return nil, err
	}

	return &le, nil
}
