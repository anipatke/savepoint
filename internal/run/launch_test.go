package run

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/internal/testutil"
)

const (
	launchConfigOff = "schema_version: 2\nquality_gates:\n  test: go test ./...\n"
	launchConfigOn  = launchConfigOff + "agent_launcher:\n  enabled: true\n  builder:\n    command: FAKE\n    args: [\"{{prompt}}\", \"{{project_root}}\"]\n"
)

// launchProject is a small valid V2 project with one Task of the given status.
func launchProject(t *testing.T, config, taskStatus string) string {
	t.Helper()
	dir := t.TempDir()
	sp := filepath.Join(dir, ".savepoint")
	testutil.WriteFile(t, filepath.Join(sp, "config.yml"), strings.Replace(config, "FAKE", os.Args[0], 1))
	testutil.WriteFile(t, filepath.Join(sp, "router.md"),
		"## Current state\n\n```yaml\nstate: task\nrelease: G-001\nobjective: O-001\ntask: T-001\nissue: none\n```\n")
	testutil.WriteFile(t, filepath.Join(sp, "releases", "G-001", "Release.md"),
		"---\nid: G-001\ntitle: \"Goal\"\nstatus: planned\n---\n\n## Outcome\n\nx\n\n## Why\n\nx\n\n## Success Conditions\n\n- x\n\n## Boundaries\n\nx\n")
	testutil.WriteFile(t, filepath.Join(sp, "objectives", "O-001-ship", "Objective.md"),
		"---\nid: O-001\ntitle: \"Ship\"\nstatus: planned\nrelease: G-001\n---\n\n# Ship\n")
	task := "---\nid: T-001\ntitle: \"Write it\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-001}\nstatus: " + taskStatus + "\n"
	if taskStatus == "in_progress" {
		task += "stage: build\n"
	}
	testutil.WriteFile(t, filepath.Join(sp, "objectives", "O-001-ship", "tasks", "T-001-write.md"), task+"---\n\n# Write it\n")
	return dir
}

func newEnough(string) (string, error) { return "2.1.296 (Claude Code)", nil }

func launch(t *testing.T, dir string, version func(string) (string, error)) (int, error, string) {
	t.Helper()
	var out bytes.Buffer
	code, err := Launch(context.Background(), Request{Dir: dir, Task: "T-001", Version: version, Heartbeat: 50 * time.Millisecond}, &out)
	return code, err, out.String()
}

func noRunsDir(t *testing.T, dir string) {
	t.Helper()
	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "runs")); !os.IsNotExist(err) {
		t.Errorf("a runs directory exists (stat error %v); nothing should have been written", err)
	}
}

func TestLaunch_offWritesNothingAndSaysHowToTurnItOn(t *testing.T) {
	for name, config := range map[string]string{"absent": launchConfigOff, "false": launchConfigOff + "agent_launcher:\n  enabled: false\n"} {
		t.Run(name, func(t *testing.T) {
			dir := launchProject(t, config, "planned")
			code, err, _ := launch(t, dir, newEnough)
			if code != ExitRefused || err == nil || !strings.Contains(err.Error(), "agent_launcher.enabled: true") {
				t.Fatalf("code = %d, err = %v", code, err)
			}
			noRunsDir(t, dir)
		})
	}
}

func TestLaunch_enabledWithoutABuilderUsesTheDefaultClaudeCommand(t *testing.T) {
	dir := launchProject(t, launchConfigOff+"agent_launcher:\n  enabled: true\n", "planned")
	var asked string
	code, err, _ := launch(t, dir, func(command string) (string, error) {
		asked = command
		return "", errors.New("stop here")
	})
	if code != ExitFailed || err == nil || asked != "claude" {
		t.Fatalf("code = %d, err = %v, asked for the version of %q", code, err, asked)
	}
	noRunsDir(t, dir)
}

func TestLaunch_buildsTheNextLineAndRootIntoTheConfiguredBuilder(t *testing.T) {
	dir := launchProject(t, launchConfigOn, "in_progress")
	t.Setenv("FAKE_AGENT", "ok")
	if code, err, out := launch(t, dir, newEnough); code != ExitOK || err != nil {
		t.Fatalf("code = %d, err = %v, out = %s", code, err, out)
	}
	runs, _ := filepath.Glob(filepath.Join(dir, ".savepoint", "runs", "*-T-001"))
	if len(runs) != 1 {
		t.Fatalf("runs = %v", runs)
	}
	rec := readRecord(t, runs[0])
	if len(rec.Command) != 3 || rec.Command[1] != "Build T-001 — Write it (O-001)" || rec.Command[2] != dir || rec.NextLine != rec.Command[1] {
		t.Errorf("record = %+v", rec)
	}
}

func TestLaunch_refusesATaskThatCannotStartBeforeWriting(t *testing.T) {
	cases := map[string]struct {
		status, want string
		task         string
	}{
		"done":    {"done", "is done", "T-001"},
		"missing": {"planned", "does not exist", "T-099"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			dir := launchProject(t, launchConfigOn, c.status)
			var out bytes.Buffer
			code, err := Launch(context.Background(), Request{Dir: dir, Task: c.task, Version: newEnough}, &out)
			if code != ExitFailed || err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("code = %d, err = %v", code, err)
			}
			noRunsDir(t, dir)
		})
	}
}

func TestLaunch_refusesATaskWhoseNextStepIsNotAStartReportingTheVerb(t *testing.T) {
	dir := launchProject(t, launchConfigOn, "in_progress")
	task := filepath.Join(dir, ".savepoint", "objectives", "O-001-ship", "tasks", "T-001-write.md")
	b, _ := os.ReadFile(task)
	if err := os.WriteFile(task, []byte(strings.Replace(string(b), "stage: build", "stage: audit", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	code, err, _ := launch(t, dir, newEnough)
	if code != ExitFailed || err == nil || !strings.Contains(err.Error(), "not ready to start") {
		t.Fatalf("code = %d, err = %v", code, err)
	}
	noRunsDir(t, dir)
}

func TestLaunch_refusesAnOldOrUnreadableAgentBeforeWriting(t *testing.T) {
	cases := map[string]func(string) (string, error){
		"old":       func(string) (string, error) { return "2.1.258 (Claude Code)", nil },
		"old major": func(string) (string, error) { return "1.9.999", nil },
		"no number": func(string) (string, error) { return "hello", nil },
		"not found": func(string) (string, error) { return "", errors.New("not found") },
	}
	for name, version := range cases {
		t.Run(name, func(t *testing.T) {
			dir := launchProject(t, launchConfigOn, "planned")
			code, err, _ := launch(t, dir, version)
			if code != ExitFailed || err == nil || !strings.Contains(err.Error(), MinClaudeVersion) {
				t.Fatalf("code = %d, err = %v", code, err)
			}
			noRunsDir(t, dir)
		})
	}
	for _, v := range []string{"2.1.259", "2.2.0", "10.0.0"} {
		if err := checkAgentVersion("claude", func(string) (string, error) { return v, nil }); err != nil {
			t.Errorf("version %s refused: %v", v, err)
		}
	}
}

func TestLaunch_aFailingAgentExitsNonZeroAndRecordsFailed(t *testing.T) {
	dir := launchProject(t, launchConfigOn, "planned")
	t.Setenv("FAKE_AGENT", "fail")
	code, err, out := launch(t, dir, newEnough)
	if code != ExitFailed || err != nil || !strings.Contains(out, "run: failed") {
		t.Fatalf("code = %d, err = %v, out = %s", code, err, out)
	}
}

func TestLaunch_endingTheContextCancelsTheRun(t *testing.T) {
	dir := launchProject(t, launchConfigOn, "planned")
	t.Setenv("FAKE_AGENT", "idle")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	var out bytes.Buffer
	code, err := Launch(ctx, Request{Dir: dir, Task: "T-001", Version: newEnough, Heartbeat: 50 * time.Millisecond, Grace: time.Second}, &out)
	if code != ExitFailed || err != nil || !strings.Contains(out.String(), "run: cancelled") {
		t.Fatalf("code = %d, err = %v, out = %s", code, err, out.String())
	}
}

func TestLaunch_neverWritesAProjectRecord(t *testing.T) {
	dir := launchProject(t, launchConfigOn, "planned")
	before := snapshotTree(t, filepath.Join(dir, ".savepoint"))
	t.Setenv("FAKE_AGENT", "ok")
	if code, err, _ := launch(t, dir, newEnough); code != ExitOK || err != nil {
		t.Fatalf("code = %d, err = %v", code, err)
	}
	after := snapshotTree(t, filepath.Join(dir, ".savepoint"))
	for path, content := range before {
		if after[path] != content {
			t.Errorf("%s changed", path)
		}
	}
	for path := range after {
		if _, ok := before[path]; !ok && !strings.HasPrefix(filepath.ToSlash(path), "runs/") {
			t.Errorf("%s was created outside runs/", path)
		}
	}
}

func snapshotTree(t *testing.T, root string) map[string]string {
	t.Helper()
	files := map[string]string{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			b, _ := os.ReadFile(path)
			rel, _ := filepath.Rel(root, path)
			files[rel] = string(b)
		}
		return nil
	})
	return files
}

func TestLostRuns_namesOnlyRunsWhoseRunnerIsGone(t *testing.T) {
	dir := project(t)
	now := time.Date(2026, 10, 11, 12, 0, 0, 0, time.UTC)
	old := now.Add(-time.Hour).Format(time.RFC3339)
	write := func(name string, rec Record) {
		t.Helper()
		d := filepath.Join(dir, ".savepoint", "runs", name)
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := writeRecord(filepath.Join(d, "run.yml"), rec); err != nil {
			t.Fatal(err)
		}
	}
	if got := LostRuns(dir, now); len(got) != 0 {
		t.Fatalf("no runs directory content: %v", got)
	}
	write("a-lost", Record{State: StateRunning, PID: 0x7ffffff0, UpdatedAt: old})
	write("b-finished", Record{State: StateFinished, PID: 0x7ffffff0, UpdatedAt: old})
	write("c-fresh", Record{State: StateRunning, PID: 0x7ffffff0, UpdatedAt: now.Format(time.RFC3339)})
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint", "runs", "d-empty"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got := LostRuns(dir, now); len(got) != 1 || got[0] != "a-lost" {
		t.Errorf("LostRuns = %v", got)
	}
}
