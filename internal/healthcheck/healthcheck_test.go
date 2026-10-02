package healthcheck

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/opencode/savepoint/internal/codehealth"
)

// recordingRunner is the fake process runner. It records every tool the
// collection tries to execute and answers lizard with a fixed CSV.
type recordingRunner struct {
	mu    sync.Mutex
	calls []codehealth.ToolSpec
	csv   []byte
}

func (r *recordingRunner) Run(_ context.Context, spec codehealth.ToolSpec) (codehealth.ToolResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.calls = append(r.calls, spec)
	if spec.Executable == "lizard" {
		return codehealth.ToolResult{Stdout: r.csv}, nil
	}
	return codehealth.ToolResult{}, codehealth.ErrToolUnavailable
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "codehealth", "testdata", "readers", name))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// newProject returns a temporary Git repository holding a minimal V2 project
// with Objective O-001.
func newProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed; collection needs repository state")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_SYSTEM", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	root := t.TempDir()
	writeFile(t, root, ".savepoint/config.yml", "schema_version: 2\nquality_gates:\n  lint: null\n  typecheck: null\n  build: null\n  test: null\ntheme: {}\n")
	writeFile(t, root, ".savepoint/router.md", "## Current state\n\n```yaml\nstate: task\nrelease: R-001\nobjective: O-001\ntask: none\nissue: none\n```\n")
	writeFile(t, root, ".savepoint/releases/R-001/Release.md", "---\nid: R-001\ntitle: \"Release\"\nstatus: planned\n---\n\n## Outcome\n\nx\n\n## Why\n\nx\n\n## Success Conditions\n\nx\n\n## Boundaries\n\nx\n")
	writeFile(t, root, ".savepoint/objectives/O-001-ship/Objective.md", "---\nid: O-001\ntitle: \"Ship\"\nstatus: planned\nrelease: R-001\n---\n\n# Ship\n")
	writeFile(t, root, "main.go", "package main\n")
	for _, args := range [][]string{
		{"init"}, {"add", "."},
		{"-c", "user.name=T", "-c", "user.email=t@example.com", "-c", "commit.gpgsign=false", "commit", "-m", "first"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	return root
}

// configure saves a configuration whose tests instance only reads go-test.json
// and whose complexity instance runs lizard through the fake runner.
func configure(t *testing.T, root, testReport string) {
	t.Helper()
	writeFile(t, root, "go-test.json", testReport)
	cfg := codehealth.Config{Version: codehealth.ConfigVersion, Capabilities: []codehealth.CapabilityConfig{
		{Capability: codehealth.CapabilityTests, Provider: codehealth.ProviderGoTestJSON, Report: "go-test.json"},
		{Capability: codehealth.CapabilityComplexity, Provider: codehealth.ProviderLizardCSV, Executable: "lizard", Args: []string{"--csv", "."}},
	}}
	if _, err := codehealth.NewStore(root).SaveConfig(cfg); err != nil {
		t.Fatal(err)
	}
}

func run(t *testing.T, root, objective string, runner *recordingRunner) (string, error) {
	t.Helper()
	var out bytes.Buffer
	err := Run(context.Background(), Request{Dir: root, Objective: objective, Runner: runner, Git: codehealth.GitRunner{}}, &out)
	return out.String(), err
}

func snapshots(t *testing.T, root string) []codehealth.Snapshot {
	t.Helper()
	got, err := codehealth.NewStore(root).LoadSnapshots()
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func assertOnlyLizardRan(t *testing.T, runner *recordingRunner) {
	t.Helper()
	if len(runner.calls) != 1 || runner.calls[0].Executable != "lizard" {
		t.Fatalf("executed tools = %+v, want only lizard; test and coverage instances read existing reports", runner.calls)
	}
}

func TestRun_notConfiguredSavesNothing(t *testing.T) {
	root := newProject(t)
	runner := &recordingRunner{}

	out, err := run(t, root, "O-001", runner)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
	if !strings.Contains(out, "not configured") {
		t.Errorf("output = %q, want it to say Code Health is not configured", out)
	}
	if len(runner.calls) != 0 || len(snapshots(t, root)) != 0 {
		t.Errorf("calls = %v, snapshots = %d; want nothing run or saved", runner.calls, len(snapshots(t, root)))
	}
	if _, err := os.Stat(filepath.Join(root, ".savepoint", "health")); !os.IsNotExist(err) {
		t.Errorf("health directory exists after an unconfigured run: %v", err)
	}
}

func TestRun_blockingVerdictStillSucceeds(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-fail.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}

	out, err := run(t, root, "O-001", runner)
	if err != nil {
		t.Fatalf("Run() error = %v, want nil for a blocking verdict", err)
	}
	assertOnlyLizardRan(t, runner)

	saved := snapshots(t, root)
	if len(saved) != 1 || saved[0].Origin != codehealth.OriginOfficial {
		t.Fatalf("snapshots = %+v, want exactly one official snapshot", saved)
	}
	cfg, _ := codehealth.NewStore(root).LoadConfig()
	verdict, err := codehealth.Evaluate(saved[0], cfg)
	if err != nil || !verdict.Blocks() {
		t.Fatalf("verdict blocks = %v, err = %v; want a blocking verdict from failing tests", verdict.Blocks(), err)
	}
	for _, want := range []string{"Objective: O-001", "Snapshot: " + saved[0].ID + " (created)", verdict.Render()} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestRun_nonBlockingVerdict(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}

	out, err := run(t, root, "O-001", runner)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	assertOnlyLizardRan(t, runner)
	saved := snapshots(t, root)
	if len(saved) != 1 {
		t.Fatalf("snapshots = %d, want 1", len(saved))
	}
	cfg, _ := codehealth.NewStore(root).LoadConfig()
	verdict, err := codehealth.Evaluate(saved[0], cfg)
	if err != nil || verdict.Blocks() {
		t.Fatalf("verdict blocks = %v, err = %v; want a non-blocking verdict", verdict.Blocks(), err)
	}
	if !strings.Contains(out, verdict.Render()) || !strings.Contains(out, saved[0].ID) {
		t.Errorf("output does not carry the snapshot ID and rendered verdict:\n%s", out)
	}
}

func TestRun_unknownObjectiveRunsNothing(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}

	_, err := run(t, root, "O-999", runner)
	if err == nil || !strings.Contains(err.Error(), "O-999") {
		t.Fatalf("Run() error = %v, want an error naming O-999", err)
	}
	if len(runner.calls) != 0 || len(snapshots(t, root)) != 0 {
		t.Errorf("calls = %v; want no tool run and no snapshot saved", runner.calls)
	}
}

func TestRun_collectionErrorSavesNothing(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	// A damaged history fails collection before any tool runs.
	writeFile(t, root, ".savepoint/health/snapshots/junk.txt", "not a snapshot")
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}

	_, err := run(t, root, "O-001", runner)
	if err == nil || !strings.Contains(err.Error(), "no snapshot was saved") {
		t.Fatalf("Run() error = %v, want a no-snapshot error", err)
	}
	if len(runner.calls) != 0 {
		t.Errorf("calls = %v, want none after a history failure", runner.calls)
	}
}

func TestRun_neverProducesManualSnapshot(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	if _, err := run(t, root, "O-001", runner); err != nil {
		t.Fatal(err)
	}
	for _, s := range snapshots(t, root) {
		if s.Origin != codehealth.OriginOfficial || s.Retention != codehealth.RetentionPermanent {
			t.Errorf("snapshot %s origin %q retention %q, want official and permanent", s.ID, s.Origin, s.Retention)
		}
	}
}

func TestRun_cancelledContextSavesNothing(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var out bytes.Buffer
	err := Run(ctx, Request{Dir: root, Objective: "O-001", Runner: runner, Git: codehealth.GitRunner{}}, &out)
	if err == nil || !errors.Is(err, codehealth.ErrCollectionCancelled) || !strings.Contains(err.Error(), "cancelled; no snapshot was saved") {
		t.Fatalf("Run() error = %v, want a cancelled no-snapshot error", err)
	}
	if len(snapshots(t, root)) != 0 {
		t.Error("a cancelled run saved a snapshot")
	}
	for _, call := range runner.calls {
		t.Errorf("tool %q ran after cancellation", call.Executable)
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("output sink failed") }

// A saved snapshot is a successful collection even when the report cannot be
// written; the problem is noted on stderr with the saved identity.
func TestRun_outputFailureAfterSaveSucceeds(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	var stderr bytes.Buffer

	err := Run(context.Background(), Request{Dir: root, Objective: "O-001", Runner: runner, Git: codehealth.GitRunner{}, Stderr: &stderr}, failingWriter{})
	if err != nil {
		t.Fatalf("Run() error = %v, want nil after the snapshot was saved", err)
	}
	saved := snapshots(t, root)
	if len(saved) != 1 || saved[0].Origin != codehealth.OriginOfficial {
		t.Fatalf("snapshots = %+v, want one official snapshot", saved)
	}
	for _, want := range []string{saved[0].ID, "was saved", "output sink failed"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("stderr = %q, want it to contain %q", stderr.String(), want)
		}
	}
}

// Without a stderr writer the note is discarded, still not an error.
func TestRun_outputFailureWithoutStderr(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	if err := Run(context.Background(), Request{Dir: root, Objective: "O-001", Runner: runner, Git: codehealth.GitRunner{}}, failingWriter{}); err != nil {
		t.Fatalf("Run() error = %v, want nil", err)
	}
}

func reportPath(root string) string {
	return filepath.Join(root, ".savepoint", "health", "report.md")
}

func runReport(root string) (string, error) {
	var out bytes.Buffer
	err := RunReport(ReportRequest{Dir: root}, &out)
	return out.String(), err
}

// blockReport makes the report path a directory so writing it fails while the
// rest of the health directory stays writable.
func blockReport(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(reportPath(root), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestRun_rewritesReportForTheSavedSnapshot(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-fail.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}

	if _, err := run(t, root, "O-001", runner); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(reportPath(root))
	if err != nil {
		t.Fatalf("no report after a check: %v", err)
	}
	if !strings.Contains(string(got), "savepoint health check O-001") {
		t.Errorf("report does not name the checked Objective:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, ".savepoint", "health", ".gitignore")); err != nil {
		t.Errorf("no ignore rule beside the report: %v", err)
	}
}

func TestRun_reportWriteFailureIsOnlyAWarning(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	blockReport(t, root)
	var out, stderr bytes.Buffer
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}

	err := Run(context.Background(), Request{Dir: root, Objective: "O-001", Runner: runner, Git: codehealth.GitRunner{}, Stderr: &stderr}, &out)
	if err != nil {
		t.Fatalf("Run() error = %v, want the check to succeed", err)
	}
	if len(snapshots(t, root)) != 1 || !strings.Contains(out.String(), "Snapshot: ") {
		t.Errorf("snapshot or verdict lost: snapshots = %d, output = %q", len(snapshots(t, root)), out.String())
	}
	if !strings.Contains(stderr.String(), "warning") || !strings.Contains(stderr.String(), "report") {
		t.Errorf("stderr = %q, want a report warning", stderr.String())
	}
}

func TestRunReport_rewritesDeletedReportWithoutNewSnapshot(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-fail.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	if _, err := run(t, root, "O-001", runner); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(reportPath(root))
	if err := os.Remove(reportPath(root)); err != nil {
		t.Fatal(err)
	}
	calls := len(runner.calls)

	out, err := runReport(root)
	if err != nil {
		t.Fatalf("RunReport() error = %v", err)
	}
	after, _ := os.ReadFile(reportPath(root))
	if string(before) != string(after) {
		t.Errorf("rewritten report differs:\n--- before\n%s\n--- after\n%s", before, after)
	}
	if !strings.Contains(out, reportPath(root)) {
		t.Errorf("output = %q, want the report path", out)
	}
	if len(snapshots(t, root)) != 1 || len(runner.calls) != calls {
		t.Errorf("snapshots = %d, tool calls = %d; want no new snapshot and no tool run", len(snapshots(t, root)), len(runner.calls)-calls)
	}
}

func TestRunReport_nothingToReportWritesNothing(t *testing.T) {
	t.Run("not set up", func(t *testing.T) {
		root := newProject(t)
		_, err := runReport(root)
		if !errors.Is(err, codehealth.ErrReportNotConfigured) {
			t.Fatalf("RunReport() error = %v, want ErrReportNotConfigured", err)
		}
		if _, err := os.Stat(filepath.Join(root, ".savepoint", "health")); !os.IsNotExist(err) {
			t.Errorf("health directory exists: %v", err)
		}
	})
	t.Run("no snapshot", func(t *testing.T) {
		root := newProject(t)
		configure(t, root, fixture(t, "tests/go-pass.jsonl"))
		_, err := runReport(root)
		if !errors.Is(err, codehealth.ErrReportNoSnapshot) {
			t.Fatalf("RunReport() error = %v, want ErrReportNoSnapshot", err)
		}
		if _, err := os.Stat(reportPath(root)); !os.IsNotExist(err) {
			t.Errorf("report exists without a snapshot: %v", err)
		}
		if len(snapshots(t, root)) != 0 {
			t.Error("a snapshot was created")
		}
	})
}

func TestRunReport_unwritableReportIsAnError(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	if _, err := run(t, root, "O-001", runner); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(reportPath(root)); err != nil {
		t.Fatal(err)
	}
	blockReport(t, root)

	if _, err := runReport(root); err == nil {
		t.Fatal("RunReport() error = nil, want a write error")
	}
}

func TestRunReport_namesTheRoutersObjective(t *testing.T) {
	root := newProject(t)
	configure(t, root, fixture(t, "tests/go-pass.jsonl"))
	runner := &recordingRunner{csv: []byte(fixture(t, "complexity/mixed.csv"))}
	if _, err := run(t, root, "O-001", runner); err != nil {
		t.Fatal(err)
	}
	if _, err := runReport(root); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(reportPath(root))
	if !strings.Contains(string(got), "savepoint health check O-001") {
		t.Errorf("report does not name the router's Objective:\n%s", got)
	}

	writeFile(t, root, ".savepoint/router.md", "## Current state\n\n```yaml\nstate: design\nrelease: R-001\nobjective: none\ntask: none\nissue: none\n```\n")
	if _, err := runReport(root); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(reportPath(root))
	if !strings.Contains(string(got), "savepoint health check O-###") {
		t.Errorf("report does not use the O-### placeholder:\n%s", got)
	}
}
