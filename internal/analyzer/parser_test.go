package analyzer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mustParseTime(value string) time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		panic(err)
	}

	return t
}

func TestParse(t *testing.T) {
	tt := []struct {
		name         string
		data         []byte
		wantLogEntry *LogEntry
		wantErr      error
	}{
		{
			name: "ok",
			data: []byte(`{"name":"renovate","hostname":"22111cc51078","pid":96,"level":20,"logContext":"ONCgiIpqZXDiRnCN_3XnS","msg":"Checking for config file in config.js","time":"2025-12-23T22:22:46.922Z","v":0}`),
			wantLogEntry: &LogEntry{
				V:          0,
				Time:       mustParseTime("2025-12-23T22:22:46.922Z"),
				PID:        96,
				Name:       "renovate",
				Message:    "Checking for config file in config.js",
				LogContext: "ONCgiIpqZXDiRnCN_3XnS",
				Level:      20,
				Hostname:   "22111cc51078",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Parse(tc.data)
			if tc.wantErr != nil {
				assert.Nil(t, result)
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.wantLogEntry, result)
			}
		})
	}
}
