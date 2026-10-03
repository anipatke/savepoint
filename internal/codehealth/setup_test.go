package codehealth

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func goProject(t *testing.T) string {
	t.Helper()
	return tree(t, map[string]string{
		"go.mod":               "module example\n\ngo 1.22\n",
		"go.sum":               "",
		".savepoint/router.md": "state: idea\n",
	})
}

// fileState records every file's bytes and modification time.
func fileState(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Lstat, not d.Info(): on NTFS a directory entry's times can lag behind the
		// file, which would read as a change the code never made.
		info, err := os.Lstat(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, p)
		if d.IsDir() {
			out[rel] = "dir " + info.ModTime().String()
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		out[rel] = info.ModTime().String() + " " + string(b)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func planFor(t *testing.T, root string, missing ...string) SetupPlan {
	t.Helper()
	rec := &lookRecorder{missing: map[string]bool{}}
	for _, m := range missing {
		rec.missing[m] = true
	}
	plan, err := PlanProject(context.Background(), root, rec.look)
	if err != nil {
		t.Fatalf("PlanProject: %v", err)
	}
	return plan
}

func TestPreviewWritesNothing(t *testing.T) {
	root := goProject(t)
	before := fileState(t, root)
	time.Sleep(10 * time.Millisecond)

	plan := planFor(t, root)
	out := plan.Preview()

	if !reflect.DeepEqual(before, fileState(t, root)) {
		t.Fatal("preview changed the project")
	}
	for _, want := range []string{"New:", "complexity via", "savepoint health setup --apply", "Default exclusions", "OSV-Scanner contacts OSV.dev"} {
		if !strings.Contains(out, want) {
			t.Errorf("preview missing %q:\n%s", want, out)
		}
	}
	if len(plan.New) == 0 {
		t.Fatal("no new proposals for a Go project")
	}
}

func TestApplyThenReapplyIsUnchanged(t *testing.T) {
	root := goProject(t)
	store := NewStore(root)

	first := planFor(t, root)
	changed, err := first.Apply(store)
	if err != nil || !changed {
		t.Fatalf("first Apply = %v, %v; want changed", changed, err)
	}
	if _, err := store.LoadConfig(); err != nil {
		t.Fatalf("saved config does not load: %v", err)
	}
	state := fileState(t, root)

	second := planFor(t, root)
	if len(second.New) != 0 {
		t.Fatalf("second plan has new entries: %+v", second.New)
	}
	changed, err = second.Apply(store)
	if err != nil || changed {
		t.Fatalf("second Apply = %v, %v; want unchanged", changed, err)
	}
	if !reflect.DeepEqual(state, fileState(t, root)) {
		t.Fatal("re-apply changed the project")
	}
	if !strings.Contains(second.Applied(false), "No changes") || !strings.Contains(second.Preview(), "No changes") {
		t.Fatal("unchanged run does not say so")
	}
}

func TestReconciliationKeepsOwnerEdits(t *testing.T) {
	root := goProject(t)
	store := NewStore(root)
	if _, err := planFor(t, root).Apply(store); err != nil {
		t.Fatal(err)
	}
	cfg, err := store.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	// Keep only the complexity entry and edit it, as an owner would.
	var edited CapabilityConfig
	for _, cc := range cfg.Capabilities {
		if cc.Capability == CapabilityComplexity {
			edited = cc
		}
	}
	edited.Required = true
	edited.Args = append(append([]string{}, edited.Args...), "--extra")
	edited.Exclusions = []string{"only-this/**"}
	edited.TimeoutSeconds = 90
	edited.Thresholds = &Threshold{Good: 8, Watch: 12}
	if _, err := store.SaveConfig(Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{edited}}); err != nil {
		t.Fatal(err)
	}

	plan := planFor(t, root)
	if len(plan.New) == 0 {
		t.Fatal("expected the other tools to be proposed again")
	}
	if _, err := plan.Apply(store); err != nil {
		t.Fatal(err)
	}
	got, err := store.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Capabilities[0], edited) {
		t.Fatalf("confirmed entry changed:\n got %+v\nwant %+v", got.Capabilities[0], edited)
	}
	if len(got.Capabilities) != 1+len(plan.New) {
		t.Fatalf("entries = %d, want %d", len(got.Capabilities), 1+len(plan.New))
	}
}

func TestMissingToolIsReportedNotRemoved(t *testing.T) {
	root := goProject(t)
	store := NewStore(root)
	if _, err := planFor(t, root).Apply(store); err != nil {
		t.Fatal(err)
	}
	before, _ := store.LoadConfig()

	plan := planFor(t, root, "lizard")
	if len(plan.New) != 0 {
		t.Fatalf("new = %+v", plan.New)
	}
	found := false
	for _, a := range plan.Attention {
		if a.Config.Executable == "lizard" && strings.Contains(a.Problem, "not on PATH") {
			found = true
		}
	}
	if !found {
		t.Fatalf("lizard not reported: %+v", plan.Attention)
	}
	if !strings.Contains(plan.Preview(), "Needs attention") {
		t.Fatal("preview has no attention section")
	}
	if changed, err := plan.Apply(store); err != nil || changed {
		t.Fatalf("Apply = %v, %v", changed, err)
	}
	after, _ := store.LoadConfig()
	if !reflect.DeepEqual(before, after) {
		t.Fatal("configuration changed")
	}
}

func TestMissingScopeIsReported(t *testing.T) {
	cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Executable: "lizard", Scope: []string{"gone/**"}},
	}}
	probe := Probe{LookPath: (&lookRecorder{}).look, Exists: func(string) bool { return false }}
	plan := Plan(&cfg, nil, probe)
	if len(plan.Attention) != 1 || !strings.Contains(plan.Attention[0].Problem, "gone") {
		t.Fatalf("attention = %+v", plan.Attention)
	}
}

func TestApplyWithNothingNewKeepsHandEditedFileBytes(t *testing.T) {
	root := goProject(t)
	store := NewStore(root)
	if _, err := planFor(t, root).Apply(store); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, ".savepoint", "health", "config.json")
	raw, _ := os.ReadFile(path)
	handFormatted := []byte(strings.ReplaceAll(string(raw), "\n  ", "\n      "))
	if err := os.WriteFile(path, handFormatted, 0o644); err != nil {
		t.Fatal(err)
	}
	if changed, err := planFor(t, root).Apply(store); err != nil || changed {
		t.Fatalf("Apply = %v, %v", changed, err)
	}
	after, _ := os.ReadFile(path)
	if string(after) != string(handFormatted) {
		t.Fatal("file was rewritten")
	}
}

func TestApplyRefusesNonProjectWithoutWriting(t *testing.T) {
	root := tree(t, map[string]string{"go.mod": "module example\n\ngo 1.22\n"})
	plan := planFor(t, root)
	if len(plan.New) == 0 {
		t.Fatal("expected proposals")
	}
	_, err := plan.Apply(NewStore(root))
	if !errors.Is(err, ErrNotProject) {
		t.Fatalf("Apply error = %v, want ErrNotProject", err)
	}
	if _, statErr := os.Stat(filepath.Join(root, ".savepoint")); statErr == nil {
		t.Fatal(".savepoint was created")
	}
}

func TestUnusableExistingConfigFailsClearly(t *testing.T) {
	root := goProject(t)
	dir := filepath.Join(root, ".savepoint", "health")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte("{nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := PlanProject(context.Background(), root, (&lookRecorder{}).look)
	if err == nil || !strings.Contains(err.Error(), "cannot be used") {
		t.Fatalf("error = %v", err)
	}
}

func TestConflictingProposalIsNotSavedAndIsExplained(t *testing.T) {
	// An existing unnamed instance with a different scope, plus a new proposal
	// of the same capability and provider, would need names; Apply refuses.
	existing := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Executable: "lizard", Scope: []string{"a/**"}},
	}}
	proposal := Proposal{Component: "b", Config: CapabilityConfig{
		Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Executable: "lizard", Scope: []string{"b/**"}}}
	probe := Probe{LookPath: (&lookRecorder{}).look, Exists: func(string) bool { return true }}
	plan := Plan(&existing, []Proposal{proposal}, probe)

	root := goProject(t)
	if _, err := plan.Apply(NewStore(root)); err == nil || !strings.Contains(err.Error(), "nothing was saved") {
		t.Fatalf("Apply error = %v", err)
	}
	if !strings.Contains(plan.Preview(), "conflict") {
		t.Fatal("preview does not explain the conflict")
	}
	if _, err := os.Stat(filepath.Join(root, ".savepoint", "health")); err == nil {
		t.Fatal("health directory was created")
	}
}

func TestUnsupportedProposalsAreListedNotAdded(t *testing.T) {
	root := tree(t, map[string]string{".savepoint/router.md": "x"})
	plan := planFor(t, root)
	if len(plan.New) != 0 || len(plan.NotSuggested) == 0 {
		t.Fatalf("plan = %+v", plan)
	}
	if !strings.Contains(plan.Preview(), "Not suggested") {
		t.Fatal("preview does not list what was skipped")
	}
}

func TestMissingToolAttentionSaysHowToInstall(t *testing.T) {
	root := goProject(t)
	store := NewStore(root)
	if _, err := planFor(t, root).Apply(store); err != nil {
		t.Fatal(err)
	}
	out := planFor(t, root, "lizard").Preview()
	if !strings.Contains(out, "lizard is not on PATH; install it yourself, Savepoint does not (pip install lizard)") {
		t.Errorf("attention line has no install hint:\n%s", out)
	}
}

func TestAbsentReportAttentionNamesTheCommandThatWritesIt(t *testing.T) {
	cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Report: "junit.xml"},
		{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Report: "pytest-junit.xml"},
	}}
	plan := Plan(&cfg, nil, Probe{LookPath: func(string) (string, error) { return "", nil }, Exists: func(string) bool { return false }})
	out := plan.Preview()
	for _, want := range []string{"vitest run --reporter=junit --outputFile=junit.xml", "pytest --junitxml=pytest-junit.xml"} {
		if !strings.Contains(out, want) {
			t.Errorf("preview does not give the command %q:\n%s", want, out)
		}
	}
}

func TestSetupWarnsWhenTwoInstancesShareAReportFile(t *testing.T) {
	cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Report: "junit.xml"},
		{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Report: "junit.xml"},
	}}
	plan := Plan(&cfg, nil, Probe{LookPath: func(string) (string, error) { return "", nil }, Exists: func(string) bool { return true }})
	if len(plan.Warnings) != 1 {
		t.Fatalf("warnings = %v, want one", plan.Warnings)
	}
	for _, want := range []string{"pytest-junit", "vitest-junit", "junit.xml", "own report file", "vitest-junit.xml"} {
		if !strings.Contains(plan.Warnings[0], want) {
			t.Errorf("warning %q is missing %q", plan.Warnings[0], want)
		}
	}
	if !strings.Contains(plan.Preview(), "Warnings:") {
		t.Error("preview has no warnings section")
	}
	if len(plan.New) != 0 || len(plan.Attention) != 0 {
		t.Errorf("a warning must not edit or drop configured entries: %+v", plan)
	}
}

func TestSharedReportsIgnoresDistinctAndEmptyPaths(t *testing.T) {
	cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Report: "junit.xml"},
		{Capability: CapabilityTests, Provider: ProviderPytestJUnit, Report: "pytest-junit.xml"},
		{Capability: CapabilityComplexity, Provider: ProviderLizardCSV, Executable: "lizard"},
	}}
	if got := SharedReports(cfg); len(got) != 0 {
		t.Errorf("SharedReports() = %+v, want none", got)
	}
}

func TestSetupSaysWhichGeneratedReportsGitDoesNotIgnore(t *testing.T) {
	root := t.TempDir()
	cfg := Config{Version: ConfigVersion, Capabilities: []CapabilityConfig{
		{Capability: CapabilityTests, Provider: ProviderVitestJUnit, Report: "junit.xml"},
		{Capability: CapabilityCoverage, Provider: ProviderVitestV8, Report: "coverage/coverage-final.json"},
		{Capability: CapabilityTests, Provider: ProviderGoTestJSON, Report: "go-test.json"},
	}}
	if got := unignoredReports(root, cfg); !reflect.DeepEqual(got, []string{"junit.xml", "coverage/", "go-test.json"}) {
		t.Errorf("without a .gitignore: %v", got)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("# build\nnode_modules\n/coverage/\n*.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := unignoredReports(root, cfg); !reflect.DeepEqual(got, []string{"junit.xml"}) {
		t.Errorf("with coverage/ and *.json ignored: %v, want only junit.xml", got)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte("junit.xml\ncoverage\ngo-test.json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := unignoredReports(root, cfg); len(got) != 0 {
		t.Errorf("with everything ignored: %v", got)
	}
}
