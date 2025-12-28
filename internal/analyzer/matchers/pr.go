package matchers

import "github.com/martin-viggiano/renovate-exporter/internal/analyzer"

func NewPullRequestMatcher() analyzer.Matcher {
	return &matcher{
		name: "pull request",
		predicate: func(e *analyzer.Entry) bool {
			return e.Message == "Renovate repository PR statistics" && e.Repository != "" && e.PullRequestStatistics != nil
		},
		extract: func(e *analyzer.Entry, metrics analyzer.Metrics) {
			metrics.PullRequests.WithLabelValues(e.Repository, "open").Set(float64(e.PullRequestStatistics.Open))
			metrics.PullRequests.WithLabelValues(e.Repository, "merged").Set(float64(e.PullRequestStatistics.Merged))
			metrics.PullRequests.WithLabelValues(e.Repository, "closed").Set(float64(e.PullRequestStatistics.Closed))
		},
	}
}
