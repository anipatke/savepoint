package codehealth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// tree writes files (slash-separated path to content) under a temp directory.
func tree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, content := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// allFound is a lookPath fake that finds every executable and records the calls.
type lookRecorder struct {
	missing map[string]bool
	calls   []string
}

func (l *lookRecorder) look(file string) (string, error) {
	l.calls = append(l.calls, file)
	if l.missing[file] {
		return "", errors.New("not found")
	}
	return "/bin/" + file, nil
}

func discover(t *testing.T, root string, missing ...string) []Proposal {
	t.Helper()
	rec := &lookRecorder{missing: map[string]bool{}}
	for _, m := range missing {
		rec.missing[m] = true
	}
	got, err := Discover(context.Background(), root, rec.look)
	if err != nil {
		t.Fatalf("Discover: %v", err)
	}
	for _, p := range got {
		if p.Config.Capability == "" {
			continue
		}
		cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{p.Config}}
		if err := cfg.Validate(); err != nil {
			t.Errorf("proposal %s/%s invalid: %v", p.Config.Capability, p.Config.Provider, err)
		}
	}
	return got
}

func find(t *testing.T, ps []Proposal, component string, p ProviderKey) Proposal {
	t.Helper()
	for _, x := range ps {
		if x.Component == component && x.Config.Provider == p {
			return x
		}
	}
	t.Fatalf("no proposal for %s in %q; got %s", p, component, summary(ps))
	return Proposal{}
}

func summary(ps []Proposal) string {
	var parts []string
	for _, p := range ps {
		parts = append(parts, fmt.Sprintf("%s:%s:%s gap=%q", p.Component, p.Config.Capability, p.Config.Provider, p.Gap))
	}
	return strings.Join(parts, "; ")
}

func gapsOf(ps []Proposal, gap Gap) []Proposal {
	var out []Proposal
	for _, p := range ps {
		if p.Gap == gap {
			out = append(out, p)
		}
	}
	return out
}

func TestDiscoverGoOnly(t *testing.T) {
	root := tree(t, map[string]string{"go.mod": "module x\n\ngo 1.22\n", "main.go": "package x\n"})
	got := discover(t, root)

	tests := find(t, got, ".", ProviderGoTestJSON)
	if tests.Gap != GapNoReportYet || tests.Config.Report != "go-test.json" || tests.GateFlag != "go test -json ./... > go-test.json" {
		t.Errorf("go tests proposal = %+v", tests)
	}
	if cov := find(t, got, ".", ProviderGoCoverProfile); cov.Config.Report != "coverage.out" {
		t.Errorf("go coverage report = %q", cov.Config.Report)
	}
	osv := find(t, got, ".", ProviderOSVScannerJSON)
	if osv.Gap != GapNone || osv.Config.Executable != "osv-scanner" {
		t.Errorf("osv proposal = %+v", osv)
	}
	if !strings.Contains(osv.Reason, "OSV.dev") {
		t.Errorf("osv reason does not state that package names go to OSV.dev: %q", osv.Reason)
	}
	for _, p := range got {
		if p.Config.Name != "" {
			t.Errorf("single-instance proposal %s has name %q", p.Config.Provider, p.Config.Name)
		}
	}
	if len(gapsOf(got, GapUnsupportedStack)) != 0 {
		t.Errorf("Go-only project reported unsupported stack: %s", summary(got))
	}
}

func TestDiscoverExistingReportClearsGap(t *testing.T) {
	root := tree(t, map[string]string{"go.mod": "module x\n", "go-test.json": "{}\n"})
	got := discover(t, root)
	if p := find(t, got, ".", ProviderGoTestJSON); p.Gap != GapNone || !strings.Contains(p.Reason, "already exists") {
		t.Errorf("existing report proposal = %+v", p)
	}
	if p := find(t, got, ".", ProviderGoCoverProfile); p.Gap != GapNoReportYet {
		t.Errorf("absent coverage report gap = %q", p.Gap)
	}
}

func TestDiscoverVitestWithV8(t *testing.T) {
	root := tree(t, map[string]string{
		"package.json":                 `{"devDependencies":{"vitest":"^3.2.0","@vitest/coverage-v8":"^3.2.0"}}`,
		"pnpm-lock.yaml":               "lockfileVersion: 9\n",
		"coverage/coverage-final.json": "{}",
	})
	got := discover(t, root)
	if p := find(t, got, ".", ProviderVitestJUnit); p.Gap != GapNoReportYet || p.Config.Report != "junit.xml" {
		t.Errorf("vitest tests proposal = %+v", p)
	}
	if p := find(t, got, ".", ProviderVitestV8); p.Gap != GapNone || p.Config.Report != "coverage/coverage-final.json" {
		t.Errorf("vitest coverage proposal = %+v", p)
	}
	if p := find(t, got, ".", ProviderOSVScannerJSON); p.Gap != GapNone {
		t.Errorf("pnpm lockfile not recognised: %+v", p)
	}
}

func TestDiscoverVitestWithoutV8CoverageIsUnsupported(t *testing.T) {
	root := tree(t, map[string]string{
		"package.json":     `{"scripts":{"test":"vitest run"}}`,
		"vitest.config.ts": "export default {}\n",
	})
	got := discover(t, root)
	p := find(t, got, ".", ProviderVitestV8)
	if p.Gap != GapUnsupportedStack || !strings.Contains(p.Reason, "@vitest/coverage-v8") {
		t.Errorf("coverage without v8 = %+v", p)
	}
	if osv := find(t, got, ".", ProviderOSVScannerJSON); osv.Gap != GapUnsupportedStack || !strings.Contains(osv.Reason, "lockfile") {
		t.Errorf("osv without lockfile = %+v", osv)
	}
}

func TestDiscoverPytestWithCoveragePy(t *testing.T) {
	root := tree(t, map[string]string{
		"pyproject.toml": "[tool.pytest.ini_options]\naddopts = \"-q\"\n\n[tool.coverage.run]\nbranch = true\n",
		"uv.lock":        "version = 1\n",
	})
	got := discover(t, root)
	if p := find(t, got, ".", ProviderPytestJUnit); p.Config.Report != "junit.xml" || p.GateFlag != "pytest --junitxml=junit.xml" {
		t.Errorf("pytest proposal = %+v", p)
	}
	if p := find(t, got, ".", ProviderCoveragePyJSON); p.Gap != GapNoReportYet || p.Config.Report != "coverage.json" {
		t.Errorf("coverage.py proposal = %+v", p)
	}
}

func TestDiscoverPythonWithoutPytestOrCoverage(t *testing.T) {
	root := tree(t, map[string]string{"pyproject.toml": "[project]\nname = \"x\"\n"})
	got := discover(t, root)
	notes := gapsOf(got, GapUnsupportedStack)
	if len(notes) == 0 || !strings.Contains(notes[0].Reason, "pytest") {
		t.Errorf("expected a pytest unsupported note, got %s", summary(got))
	}
	for _, p := range got {
		if p.Config.Provider == ProviderPytestJUnit || p.Config.Provider == ProviderCoveragePyJSON {
			t.Errorf("proposed %s without evidence", p.Config.Provider)
		}
	}
}

func TestDiscoverMonorepoNamesAndScopes(t *testing.T) {
	root := tree(t, map[string]string{
		"services/api/go.mod":             "module api\n",
		"web/package.json":                `{"devDependencies":{"vitest":"3"}}`,
		"web/package-lock.json":           "{}",
		"tools/py/pyproject.toml":         "[tool.pytest.ini_options]\n",
		"services/billing/go.mod":         "module billing\n",
		"services/billing/vendor/x":       "",
		"README.md":                       "# repo\n",
		"web/node_modules/a/package.json": `{"devDependencies":{"vitest":"3"}}`,
	})
	got := discover(t, root)

	api := find(t, got, "services/api", ProviderLizardCSV)
	billing := find(t, got, "services/billing", ProviderLizardCSV)
	if api.Config.Name != "services-api" || billing.Config.Name != "services-billing" {
		t.Errorf("lizard names = %q, %q", api.Config.Name, billing.Config.Name)
	}
	if !reflect.DeepEqual(api.Config.Scope, []string{"services/api/**"}) {
		t.Errorf("api scope = %v", api.Config.Scope)
	}
	if api.Config.Report == billing.Config.Report {
		t.Errorf("components share report %q", api.Config.Report)
	}
	if p := find(t, got, "services/api", ProviderGoTestJSON); p.GateFlag != "cd services/api && go test -json ./... > go-test.json" || p.Config.Report != "services/api/go-test.json" {
		t.Errorf("api go tests = %+v", p)
	}
	if !slices.Contains(billing.Config.Exclusions, "services/billing/vendor/**") {
		t.Errorf("billing exclusions = %v", billing.Config.Exclusions)
	}
	// node_modules is never searched, so web/node_modules adds no component.
	for _, p := range got {
		if strings.Contains(p.Component, "node_modules") {
			t.Errorf("searched node_modules: %s", p.Component)
		}
		if p.Component == "." {
			t.Errorf("monorepo root produced an aggregate proposal: %+v", p)
		}
	}
	// The same Config the proposals form must pass whole-config validation, which
	// requires distinct names and scopes.
	var cfg Config
	cfg.Version = ConfigVersion
	for _, p := range got {
		if p.Config.Capability != "" && p.Gap != GapUnsupportedStack {
			cfg.Capabilities = append(cfg.Capabilities, p.Config)
		}
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("combined proposals invalid: %v", err)
	}
}

func TestDiscoverDefaultExclusions(t *testing.T) {
	root := tree(t, map[string]string{
		"go.mod":            "module x\n",
		"vendor/a/a.go":     "package a\n",
		"third_party/b/b.c": "",
	})
	got := discover(t, root)
	want := []string{"vendor/**", "third_party/**", "**/*.pb.go", "**/*_generated.*", "**/*.min.js"}
	if ex := find(t, got, ".", ProviderJscpdJSON).Config.Exclusions; !reflect.DeepEqual(ex, want) {
		t.Errorf("exclusions = %v, want %v", ex, want)
	}
}

func TestDiscoverExecutedToolArgsTranslateExclusions(t *testing.T) {
	root := tree(t, map[string]string{"go.mod": "module x\n", "vendor/a.go": ""})
	got := discover(t, root)

	liz := find(t, got, ".", ProviderLizardCSV).Config
	wantLiz := []string{"--csv", "-o", ".savepoint/health/reports/root-lizard.csv",
		"-x", "*/vendor/*", "-x", "*.pb.go", "-x", "*_generated.*", "-x", "*.min.js", "."}
	if !reflect.DeepEqual(liz.Args, wantLiz) || liz.Report != ".savepoint/health/reports/root-lizard.csv" {
		t.Errorf("lizard = %v report %q", liz.Args, liz.Report)
	}
	jscpd := find(t, got, ".", ProviderJscpdJSON).Config
	if jscpd.Report != ".savepoint/health/reports/root/jscpd-report.json" ||
		!slices.Contains(jscpd.Args, "vendor/**,**/*.pb.go,**/*_generated.*,**/*.min.js,**/.git/**") ||
		!slices.Contains(jscpd.Args, ".savepoint/health/reports/root") {
		t.Errorf("jscpd = %v report %q", jscpd.Args, jscpd.Report)
	}
	osv := find(t, got, ".", ProviderOSVScannerJSON).Config
	if !slices.Contains(osv.Args, "g:vendor/**") || !slices.Contains(osv.Args, "--output-file") {
		t.Errorf("osv args = %v", osv.Args)
	}
	for _, c := range []CapabilityConfig{liz, jscpd, osv} {
		if _, ok := DefaultTimeoutSeconds(c.Provider); !ok {
			t.Errorf("%s has no default timeout", c.Provider)
		}
	}
}

func TestLizardExclude(t *testing.T) {
	for in, want := range map[string]string{
		"vendor/**":        "*/vendor/*",
		"**/*.pb.go":       "*.pb.go",
		"**/*_generated.*": "*_generated.*",
		"a/b/vendor/**":    "*/a/b/vendor/*",
	} {
		if got := lizardExclude(in); got != want {
			t.Errorf("lizardExclude(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDiscoverMissingExecutableIsProposedNotInstalled(t *testing.T) {
	root := tree(t, map[string]string{"go.mod": "module x\n"})
	got := discover(t, root, "lizard", "jscpd")
	for _, p := range []ProviderKey{ProviderLizardCSV, ProviderJscpdJSON} {
		x := find(t, got, ".", p)
		if x.Gap != GapMissingExecutable || !strings.Contains(x.Reason, "is not on PATH") {
			t.Errorf("%s = %+v", p, x)
		}
	}
	if x := find(t, got, ".", ProviderOSVScannerJSON); x.Gap != GapNone {
		t.Errorf("found tool got gap %q", x.Gap)
	}
}

func TestDiscoverUnsupportedStack(t *testing.T) {
	root := tree(t, map[string]string{"Cargo.toml": "[package]\n", "src/main.rs": ""})
	got := discover(t, root)
	if len(got) != 1 || got[0].Gap != GapUnsupportedStack || got[0].Config.Capability != "" {
		t.Fatalf("got %s", summary(got))
	}
}

func TestDiscoverJSWithoutVitestIsExplicitGap(t *testing.T) {
	root := tree(t, map[string]string{"package.json": `{"devDependencies":{"jest":"29"}}`})
	got := discover(t, root)
	notes := gapsOf(got, GapUnsupportedStack)
	if len(notes) == 0 || !strings.Contains(notes[0].Reason, "Vitest") {
		t.Errorf("got %s", summary(got))
	}
}

func TestDiscoverOversizedManifestIsAGap(t *testing.T) {
	root := tree(t, map[string]string{
		"package.json": `{"devDependencies":{"vitest":"3"},"pad":"` + strings.Repeat("x", MaxDiscoveryFileBytes) + `"}`,
	})
	got := discover(t, root)
	notes := gapsOf(got, GapOversizedFile)
	if len(notes) != 1 || !strings.Contains(notes[0].Reason, "package.json") {
		t.Fatalf("got %s", summary(got))
	}
	for _, p := range got {
		if p.Config.Provider == ProviderVitestJUnit {
			t.Error("proposed Vitest from a file that was not read")
		}
	}
}

func TestDiscoverMalformedManifestIsAGap(t *testing.T) {
	root := tree(t, map[string]string{"package.json": `{"devDependencies": `})
	got := discover(t, root)
	if notes := gapsOf(got, GapUnreadableFile); len(notes) != 1 {
		t.Fatalf("got %s", summary(got))
	}
}

func TestDiscoverDoesNotFollowEscapingSymlink(t *testing.T) {
	outside := tree(t, map[string]string{"package.json": `{"devDependencies":{"vitest":"3"}}`})
	root := tree(t, map[string]string{"go.mod": "module x\n"})
	link := filepath.Join(root, "package.json")
	if err := os.Symlink(filepath.Join(outside, "package.json"), link); err != nil {
		t.Skipf("cannot create symlinks on this platform or account: %v", err)
	}
	got := discover(t, root)
	notes := gapsOf(got, GapUnreadableFile)
	if len(notes) != 1 || !strings.Contains(notes[0].Reason, "leaves the project") {
		t.Fatalf("got %s", summary(got))
	}
	for _, p := range got {
		if p.Config.Provider == ProviderVitestJUnit {
			t.Error("followed a symlink out of the project")
		}
	}
}

func TestDiscoverSearchDepthIsBounded(t *testing.T) {
	root := tree(t, map[string]string{
		"a/b/c/go.mod":   "module ok\n",
		"a/b/c/d/go.mod": "module toodeep\n",
	})
	got := discover(t, root)
	for _, p := range got {
		if p.Component == "a/b/c/d" {
			t.Errorf("searched past depth %d: %s", MaxDiscoveryDepth, summary(got))
		}
	}
	find(t, got, "a/b/c", ProviderGoTestJSON)
}

func TestDiscoverReadLimitIsAGap(t *testing.T) {
	files := map[string]string{}
	for i := range MaxDiscoveryFiles + 2 {
		files[fmt.Sprintf("p%03d/package.json", i)] = `{}`
	}
	got := discover(t, tree(t, files))
	if len(gapsOf(got, GapReadLimit)) == 0 {
		t.Errorf("no read-limit gap in %d proposals", len(got))
	}
}

func TestDiscoverIsDeterministic(t *testing.T) {
	root := tree(t, map[string]string{
		"b/go.mod": "module b\n", "a/go.mod": "module a\n", "c/package.json": `{"devDependencies":{"vitest":"3"}}`,
	})
	first := discover(t, root)
	for range 5 {
		if again := discover(t, root); !reflect.DeepEqual(first, again) {
			t.Fatalf("results differ between runs:\n%s\n%s", summary(first), summary(again))
		}
	}
}

func TestDiscoverIsReadOnly(t *testing.T) {
	root := tree(t, map[string]string{
		"go.mod": "module x\n", "web/package.json": `{"devDependencies":{"vitest":"3"}}`, "vendor/a.go": "",
	})
	type stamp struct {
		size int64
		mod  time.Time
		body string
	}
	snapshot := func() map[string]stamp {
		out := map[string]stamp{}
		err := filepath.WalkDir(root, func(p string, e os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			info, err := e.Info()
			if err != nil {
				return err
			}
			s := stamp{size: info.Size(), mod: info.ModTime()}
			if !e.IsDir() {
				b, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				s.body = string(b)
			}
			out[p] = s
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
		return out
	}
	before := snapshot()
	rec := &lookRecorder{}
	if _, err := Discover(context.Background(), root, rec.look); err != nil {
		t.Fatal(err)
	}
	if after := snapshot(); !reflect.DeepEqual(before, after) {
		t.Error("Discover changed the project tree")
	}
	// The injected lookPath is the only call that leaves the tree, and it is
	// asked only about the three executed tools.
	for _, c := range rec.calls {
		if !slices.Contains([]string{"lizard", "jscpd", "osv-scanner"}, c) {
			t.Errorf("lookPath asked about %q", c)
		}
	}
	if len(rec.calls) == 0 {
		t.Error("lookPath was never consulted")
	}
}

func TestDiscoverRejectsBadRoot(t *testing.T) {
	rec := &lookRecorder{}
	if _, err := Discover(context.Background(), filepath.Join(t.TempDir(), "missing"), rec.look); err == nil {
		t.Error("missing root accepted")
	}
	file := filepath.Join(t.TempDir(), "f")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Discover(context.Background(), file, rec.look); err == nil {
		t.Error("file root accepted")
	}
}

func TestDiscoverHonoursCancellation(t *testing.T) {
	root := tree(t, map[string]string{"go.mod": "module x\n"})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rec := &lookRecorder{}
	if _, err := Discover(ctx, root, rec.look); !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestDiscoverUnrepresentableNameBecomesNote(t *testing.T) {
	long := strings.Repeat("n", MaxTokenLen)
	root := tree(t, map[string]string{
		"a/" + long + "/go.mod": "module a\n",
		"b/go.mod":              "module b\n",
	})
	got := discover(t, root)
	var sawNote bool
	for _, p := range got {
		if p.Component == "a/"+long && p.Config.Capability == "" && strings.Contains(p.Reason, "cannot be configured") {
			sawNote = true
		}
	}
	if !sawNote {
		t.Errorf("expected a note for the over-long instance name, got %s", summary(got))
	}
}
