package matchers_test

import (
	"testing"

	"github.com/martin-viggiano/renovate-exporter/internal/analyzer"
	"github.com/martin-viggiano/renovate-exporter/internal/analyzer/matchers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependencyCountMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewDependencyCountMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"247242df26cb","pid":96,"level":30,"logContext":"9c28e67e-ec9c-43dc-a985-2cf6b447aa6f","repository":"test","baseBranch":"master","stats":{"managers":{"gitlabci":{"fileCount":1,"depCount":1},"gomod":{"fileCount":1,"depCount":17},"renovate-config-presets":{"fileCount":1,"depCount":1}},"total":{"fileCount":3,"depCount":19}},"msg":"Dependency extraction complete","time":"2025-12-28T13:46:46.486Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependenciesTotal.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(17), testutil.ToFloat64(
		engine.Metrics().DependenciesTotal.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependenciesTotal.WithLabelValues("test", "renovate-config-presets"),
	))
}

func TestDependencyFilesMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewDependencyFilesMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"247242df26cb","pid":96,"level":30,"logContext":"9c28e67e-ec9c-43dc-a985-2cf6b447aa6f","repository":"test","baseBranch":"master","stats":{"managers":{"gitlabci":{"fileCount":1,"depCount":1},"gomod":{"fileCount":1,"depCount":17},"renovate-config-presets":{"fileCount":1,"depCount":1}},"total":{"fileCount":3,"depCount":19}},"msg":"Dependency extraction complete","time":"2025-12-28T13:46:46.486Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyFilesTotal.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyFilesTotal.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyFilesTotal.WithLabelValues("test", "renovate-config-presets"),
	))
}

func TestOutdatedDependenciesMatcher(t *testing.T) {
	reg := prometheus.NewRegistry()

	matchers := []analyzer.Matcher{
		matchers.NewOutdatedDependencyMatcher(),
	}

	engine, err := analyzer.NewEngine(reg, matchers)
	require.NoError(t, err)

	err = engine.Process([]byte(`{"name":"renovate","hostname":"247242df26cb","pid":805,"level":20,"logContext":"16eba67a-c5cf-48d1-abe6-680e16875f8f","repository":"test","baseBranch":"master","config":{"gitlabci":[{"packageFile":".gitlab-ci.yml","deps":[{"depName":"golang","packageName":"golang","currentValue":"1.25.5","replaceString":"golang:1.25.5","autoReplaceStringTemplate":"{{depName}}{{#if newValue}}:{{newValue}}{{/if}}{{#if newDigest}}@{{newDigest}}{{/if}}","datasource":"docker","depType":"image","updates":[],"versioning":"docker","warnings":[],"registryUrl":"https://index.docker.io","lookupName":"library/golang","currentVersion":"1.25.5","currentVersionTimestamp":"2025-12-30T09:23:31.782Z","currentVersionAgeInDays":0,"fixedVersion":"1.25.5"}]}],"gomod":[{"deps":[{"datasource":"golang-version","versioning":"go-mod-directive","depType":"golang","depName":"go","currentValue":"1.24.4","managerData":{"lineNumber":2},"updates":[],"packageName":"go","warnings":[],"sourceUrl":"https://github.com/golang/go","registryUrl":"https://raw.githubusercontent.com/golang/website","homepage":"https://go.dev/","mostRecentTimestamp":"2025-12-02T00:00:00.000Z","currentVersion":"1.24.4","currentVersionTimestamp":"2025-06-05T00:00:00.000Z","currentVersionAgeInDays":208,"fixedVersion":"1.24.4"},{"datasource":"golang-version","depType":"toolchain","depName":"go","currentValue":"1.25.5","managerData":{"lineNumber":4},"updates":[],"packageName":"go","versioning":"semver","warnings":[],"sourceUrl":"https://github.com/golang/go","registryUrl":"https://raw.githubusercontent.com/golang/website","homepage":"https://go.dev/","mostRecentTimestamp":"2025-12-02T00:00:00.000Z","currentVersion":"1.25.5","currentVersionTimestamp":"2025-12-02T00:00:00.000Z","currentVersionAgeInDays":28,"fixedVersion":"1.25.5"},{"datasource":"go","depType":"require","depName":"github.com/stretchr/testify","currentValue":"v1.7.5","managerData":{"multiLine":true,"lineNumber":12},"updates":[{"bucket":"non-major","newVersion":"v1.11.1","newValue":"v1.11.1","releaseTimestamp":"2025-08-27T10:46:31.000Z","newVersionAgeInDays":125,"newMajor":1,"newMinor":11,"newPatch":1,"updateType":"minor","isBreaking":false,"libYears":3.1765497209538305,"branchName":"renovate/github.com-stretchr-testify-1.x"}],"packageName":"github.com/stretchr/testify","versioning":"semver","warnings":[],"sourceUrl":"https://github.com/stretchr/testify","mostRecentTimestamp":"2025-08-27T10:46:31.000Z","currentVersion":"v1.7.5","currentVersionTimestamp":"2022-06-24T00:11:59.000Z","currentVersionAgeInDays":1285,"isSingleVersion":true,"fixedVersion":"v1.7.5"}],"packageFile":"go.mod"}],"renovate-config-presets":[{"deps":[{"depName":"mft/demo-poc/renovate-config","skipReason":"unsupported-datasource","updates":[],"packageName":"mft/demo-poc/renovate-config"}],"packageFile":"renovate.json"}]},"msg":"packageFiles with updates","time":"2025-12-30T11:56:08.800Z","v":0}`))
	assert.NoError(t, err)

	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().DependencyOutdatedTotal.WithLabelValues("test", "gitlabci"),
	))
	assert.Equal(t, float64(1), testutil.ToFloat64(
		engine.Metrics().DependencyOutdatedTotal.WithLabelValues("test", "gomod"),
	))
	assert.Equal(t, float64(0), testutil.ToFloat64(
		engine.Metrics().DependencyOutdatedTotal.WithLabelValues("test", "renovate-config-presets"),
	))
}
