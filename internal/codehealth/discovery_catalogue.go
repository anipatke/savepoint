package codehealth

import "strings"

// This file is the data behind Discover: what to look for, where each provider's
// report conventionally lives, and how to translate exclusions into each tool's
// own flags. Logic in discovery.go reads these tables and holds no copy of them.

// Bounds on discovery, so inspection of a large tree stays cheap and predictable.
const (
	MaxDiscoveryDepth     = 3
	MaxDiscoveryFiles     = 64
	MaxDiscoveryFileBytes = 256 << 10
	MaxDiscoveryDirs      = 2000
)

// defaultExclusions are proposed for confirmation. A pattern ending in "/**"
// names a directory and is proposed only when it exists; any other pattern is a
// generated-file suffix and is always proposed.
var defaultExclusions = []string{
	"vendor/**",
	"node_modules/**",
	"third_party/**",
	"dist/**",
	"build/**",
	"**/*.pb.go",
	"**/*_generated.*",
	"**/*.min.js",
}

// reportsDir holds the reports of tools Savepoint runs itself.
const reportsDir = ".savepoint/health/reports"

// stack is the language ecosystem of one manifest.
type stack string

const (
	stackGo     stack = "go"
	stackJS     stack = "javascript"
	stackPython stack = "python"
)

// manifestStacks maps a manifest file name to the stack it marks. A directory
// holding one is a component.
var manifestStacks = map[string]stack{
	"go.mod":         stackGo,
	"package.json":   stackJS,
	"pyproject.toml": stackPython,
	"pytest.ini":     stackPython,
	"setup.cfg":      stackPython,
}

// lockfiles lists the resolved-dependency files OSV-Scanner reads, per stack.
var lockfiles = map[stack][]string{
	stackGo: {"go.mod"},
	stackJS: {"package-lock.json", "npm-shrinkwrap.json", "pnpm-lock.yaml", "yarn.lock", "bun.lock", "bun.lockb"},
	stackPython: {"uv.lock", "poetry.lock", "Pipfile.lock", "pdm.lock", "pylock.toml",
		"requirements.txt"},
}

// vitestConfigs name the Vitest configuration files that mark Vitest as the runner.
var vitestConfigs = []string{
	"vitest.config.ts", "vitest.config.mts", "vitest.config.cts",
	"vitest.config.js", "vitest.config.mjs", "vitest.config.cjs",
}

// providerConfigs are existing provider settings worth mentioning in a reason.
var providerConfigs = map[ProviderKey][]string{
	ProviderJscpdJSON:      {".jscpd.json"},
	ProviderLizardCSV:      {".lizardrc"},
	ProviderCoveragePyJSON: {".coveragerc"},
}

// reportOnly describes a provider whose report the project's own gate produces.
// Report is the conventional path relative to the component directory; Gate is
// the command that produces it, run from that directory.
type reportOnly struct {
	Report string
	Gate   string
}

var reportOnlyProviders = map[ProviderKey]reportOnly{
	ProviderGoTestJSON:     {"go-test.json", "go test -json ./... > go-test.json"},
	ProviderVitestJUnit:    {"junit.xml", "vitest run --reporter=junit --outputFile=junit.xml"},
	ProviderPytestJUnit:    {"junit.xml", "pytest --junitxml=junit.xml"},
	ProviderGoCoverProfile: {"coverage.out", "go test -coverprofile=coverage.out ./..."},
	ProviderVitestV8:       {"coverage/coverage-final.json", "vitest run --coverage --coverage.provider=v8 --coverage.reporter=json"},
	ProviderCoveragePyJSON: {"coverage.json", "coverage run -m pytest && coverage json -o coverage.json"},
}

// executed describes a provider Savepoint runs. Args builds the argument vector
// from the target directory, the report path, and the translated exclusions.
type executed struct {
	Executable string
	Reason     string
	Report     func(instance string) string
	Args       func(target, report string, exclusions []string) []string
}

var executedProviders = map[ProviderKey]executed{
	ProviderLizardCSV: {
		Executable: "lizard",
		Reason:     "Lizard measures the highest cyclomatic complexity across Go, JavaScript, TypeScript, and Python.",
		Report:     func(i string) string { return reportsDir + "/" + i + "-lizard.csv" },
		Args: func(target, report string, ex []string) []string {
			args := []string{"--csv", "-o", report}
			for _, p := range ex {
				args = append(args, "-x", lizardExclude(p))
			}
			return append(args, target)
		},
	},
	ProviderJscpdJSON: {
		Executable: "jscpd",
		Reason:     "jscpd measures the share of duplicated lines across Go, JavaScript, TypeScript, and Python.",
		Report:     func(i string) string { return reportsDir + "/" + i + "/jscpd-report.json" },
		Args: func(target, report string, ex []string) []string {
			args := []string{"--reporters", "json", "--output", report[:strings.LastIndex(report, "/")], "--workers", "1"}
			if len(ex) > 0 {
				args = append(args, "--ignore", strings.Join(ex, ","))
			}
			return append(args, target)
		},
	},
	ProviderOSVScannerJSON: {
		Executable: "osv-scanner",
		Reason:     "OSV-Scanner lists known vulnerabilities in resolved dependencies. It runs in its normal online mode: the scanner, not Savepoint, sends package names and versions to OSV.dev.",
		Report:     func(i string) string { return reportsDir + "/" + i + "-osv.json" },
		Args: func(target, report string, ex []string) []string {
			args := []string{"scan", "source", "-r", "--format", "json", "--output-file", report}
			for _, p := range ex {
				args = append(args, "--experimental-exclude", "g:"+p)
			}
			return append(args, target)
		},
	},
}

// lizardExclude turns a repository glob into Lizard's fnmatch form, where "*"
// also crosses directory separators.
func lizardExclude(pattern string) string {
	p := strings.TrimPrefix(pattern, "**/")
	if dir, ok := strings.CutSuffix(p, "/**"); ok {
		p = dir + "/*"
	}
	if !strings.HasPrefix(p, "*") {
		p = "*/" + p
	}
	return p
}
