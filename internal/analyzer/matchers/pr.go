package matchers

import "github.com/martin-viggiano/renovate-exporter/internal/analyzer"

func NewPullRequestMatcher() analyzer.Matcher {
	// {"name":"renovate","hostname":"22111cc51078","pid":96,"level":20,"logContext":"4b569c43-f97c-43eb-ae32-443e49ca1b89","repository":"test","stats":{"total":7,"open":3,"closed":4,"merged":0},"msg":"Renovate repository PR statistics","time":"2025-12-23T22:23:16.987Z","v":0}
	return &matcher{
		name: "pull request",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Renovate repository PR statistics" && e.Repository != "" && e.Stats != nil
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			metrics.PullRequests.WithLabelValues(e.Repository, "open").Set(float64(e.Stats.Open))
			metrics.PullRequests.WithLabelValues(e.Repository, "merged").Set(float64(e.Stats.Merged))
			metrics.PullRequests.WithLabelValues(e.Repository, "closed").Set(float64(e.Stats.Closed))
		},
	}
}
