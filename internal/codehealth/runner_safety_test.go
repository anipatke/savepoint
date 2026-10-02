package codehealth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// waitGone fails the test when pid is still running after a short wait.
func waitGone(t *testing.T, pid int) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for processAlive(pid) {
		if time.Now().After(deadline) {
			t.Fatalf("process %d is still running", pid)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func readPid(t *testing.T, file string) int {
	t.Helper()
	raw, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("helper never recorded a pid: %v", err)
	}
	pid, err := strconv.Atoi(string(raw))
	if err != nil {
		t.Fatalf("bad pid %q: %v", raw, err)
	}
	return pid
}

func TestExecRunnerOutputCapsAtTheBoundary(t *testing.T) {
	helperEnv(t)
	dir := t.TempDir()

	for _, c := range []struct {
		name      string
		size      int
		truncated bool
	}{
		{"below the report limit", MaxReportBytes - 1, false},
		{"exactly the report limit", MaxReportBytes, false},
		{"above the report limit", MaxReportBytes + 1, true},
	} {
		t.Run("stdout "+c.name, func(t *testing.T) {
			res, err := ExecRunner{}.Run(context.Background(), helperSpec(dir, "stdout-size", strconv.Itoa(c.size)))
			want := min(c.size, MaxReportBytes)
			if err != nil || res.ExitCode != 0 || res.Truncated != c.truncated || len(res.Stdout) != want {
				t.Fatalf("got %d bytes truncated=%v exit %d err %v", len(res.Stdout), res.Truncated, res.ExitCode, err)
			}
		})
	}

	// The end of stderr is the useful part, so it must survive however much
	// progress text came before it, and a chatty tool must never be blocked.
	for _, c := range []struct {
		name string
		size int
	}{
		{"below the stderr cap", maxStderrBytes - 1},
		{"exactly the stderr cap", maxStderrBytes},
		{"above the stderr cap", maxStderrBytes + 1},
		{"far above the stderr cap", 8 * maxStderrBytes},
	} {
		t.Run("stderr "+c.name, func(t *testing.T) {
			res, err := ExecRunner{}.Run(context.Background(), helperSpec(dir, "stderr-tail", strconv.Itoa(c.size)))
			if err != nil || res.ExitCode != 127 {
				t.Fatalf("got exit %d err %v", res.ExitCode, err)
			}
			if !strings.HasSuffix(res.Stderr, "END") || len(res.Stderr) > MaxReasonLen {
				t.Fatalf("stderr lost its end or is unbounded (%d bytes): %q", len(res.Stderr), res.Stderr)
			}
		})
	}
}

func TestExecRunnerPassesArgumentsAndDirectoryUnchanged(t *testing.T) {
	helperEnv(t)
	dir := filepath.Join(t.TempDir(), "dir with spaces")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	// None of these may be expanded, split, or interpreted.
	args := []string{"$(echo hi)", "a;b", "`id`", "two words", "*", "--flag=$HOME", "", "'quoted'", "line\nbreak"}
	res, err := ExecRunner{}.Run(context.Background(), helperSpec(dir, append([]string{"args"}, args...)...))
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("got exit %d err %v", res.ExitCode, err)
	}
	wd, got, _ := strings.Cut(string(res.Stdout), "\x00")
	if want, _ := filepath.EvalSymlinks(dir); !strings.EqualFold(filepath.Clean(wd), want) {
		if gotDir, _ := filepath.EvalSymlinks(wd); !strings.EqualFold(gotDir, want) {
			t.Errorf("working directory = %q, want %q", wd, dir)
		}
	}
	if got != strings.Join(args, "\x00") {
		t.Errorf("arguments changed: %q", got)
	}
}

func TestExecRunnerStopsTheWholeTreeOnDeadlineAndRunsAgain(t *testing.T) {
	helperEnv(t)
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "child.pid")

	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err := ExecRunner{}.Run(ctx, helperSpec(dir, "child", pidFile))
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > 10*time.Second {
		t.Fatalf("got %v after %s", err, time.Since(start))
	}
	waitGone(t, readPid(t, pidFile))

	// A failed run leaves nothing behind that stops the next one.
	res, err := ExecRunner{}.Run(context.Background(), helperSpec(dir, "ok"))
	if err != nil || string(res.Stdout) != "report-bytes" {
		t.Fatalf("run after a timeout: %q %v", res.Stdout, err)
	}
}

func TestExecRunnerBoundsPipesHeldByAnOrphan(t *testing.T) {
	helperEnv(t)
	dir := t.TempDir()
	pidFile := filepath.Join(dir, "orphan.pid")

	start := time.Now()
	res, err := ExecRunner{}.Run(context.Background(), helperSpec(dir, "orphan", pidFile))
	elapsed := time.Since(start)
	if pid := readPid(t, pidFile); pid > 0 {
		// The tool finished normally, so nothing cancelled its child: end it here.
		defer func() {
			if p, err := os.FindProcess(pid); err == nil {
				p.Kill()
			}
			waitGone(t, pid)
		}()
	}
	if err != nil || res.ExitCode != 0 {
		t.Fatalf("a finished tool with a lingering child is still a result: exit %d err %v", res.ExitCode, err)
	}
	if elapsed > killWait+5*time.Second {
		t.Fatalf("pipes were held for %s, want about %s", elapsed, killWait)
	}
}

func TestCollectCancelWinsOverTheInstanceDeadline(t *testing.T) {
	dir := project(t)
	before := storeFiles(t, dir)
	ctx, cancel := context.WithCancel(context.Background())
	// The tool reports the instance deadline, but the whole run was already
	// cancelled: cancellation must win and nothing is saved.
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		"slow": func(context.Context, ToolSpec) (ToolResult, error) {
			cancel()
			return ToolResult{}, context.DeadlineExceeded
		},
	}}
	cfg := cfgOf(lizardInstance("slow", "slow", "a/**"))
	_, err := Collect(ctx, CollectRequest{Root: dir, Origin: OriginManual, Config: cfg, Readers: Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, Runner: tools, Clock: testClock})
	if !errors.Is(err, ErrCollectionCancelled) {
		t.Fatalf("err = %v", err)
	}
	if after := storeFiles(t, dir); len(after) != len(before) {
		t.Errorf("store changed: %v -> %v", before, after)
	}
}

func TestCollectFailureReasonsLeakNeitherReportNorEnvironment(t *testing.T) {
	dir := project(t)
	t.Setenv("SAVEPOINT_TEST_SECRET", "hunter2-token")
	const source = "func secretSourceBody() {}"
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){
		// A failing tool's stdout is its (possibly source-bearing) report; only
		// the bounded stderr line may reach the saved reason.
		"fails": func(context.Context, ToolSpec) (ToolResult, error) {
			return ToolResult{Stdout: []byte(source), ExitCode: 3, Stderr: "boom"}, nil
		},
		"unreadable": stdout(source),
	}}
	readers := Readers{ProviderLizardCSV: readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
		return Reading{}, errors.New("bad csv near " + strings.Repeat("z", 5*MaxReasonLen))
	})}
	got := collect(t, dir, cfgOf(lizardInstance("fails", "fails", "a/**"), lizardInstance("unreadable", "unreadable", "b/**")), readers, tools)

	for _, name := range []string{"fails", "unreadable"} {
		r := collectedFor(t, got, CapabilityComplexity, name).Result
		if r.Outcome != OutcomeFailed || len(r.Reason) > MaxReasonLen {
			t.Errorf("%s = %s with a %d-byte reason", name, r.Outcome, len(r.Reason))
		}
		if strings.Contains(r.Reason, "secretSourceBody") || strings.Contains(r.Reason, "hunter2") {
			t.Errorf("%s reason leaks: %q", name, r.Reason)
		}
	}
}

func TestCollectReportFileSizeBoundary(t *testing.T) {
	for _, c := range []struct {
		name    string
		size    int64
		outcome Outcome
	}{
		{"below the limit", MaxReportBytes - 1, OutcomeAvailable},
		{"exactly the limit", MaxReportBytes, OutcomeAvailable},
		{"above the limit", MaxReportBytes + 1, OutcomePartial},
	} {
		t.Run(c.name, func(t *testing.T) {
			dir := project(t)
			f, err := os.Create(filepath.Join(dir, "r.out"))
			if err != nil {
				t.Fatal(err)
			}
			f.Truncate(c.size)
			f.Close()
			var read int
			reader := readerFunc(func(_ context.Context, in ReportInput) (Reading, error) {
				read = len(in.Data)
				return goodReading(1, UnitPercent), nil
			})
			cfg := cfgOf(CapabilityConfig{Capability: CapabilityCoverage, Provider: ProviderGoCoverProfile, Report: "r.out"})
			r := collectedFor(t, collect(t, dir, cfg, Readers{ProviderGoCoverProfile: reader}, &fakeTools{t: t}), CapabilityCoverage, "").Result
			if r.Outcome != c.outcome {
				t.Fatalf("outcome = %s (%s), want %s", r.Outcome, r.Reason, c.outcome)
			}
			if c.outcome == OutcomeAvailable && int64(read) != c.size {
				t.Errorf("reader got %d bytes, want %d", read, c.size)
			}
			if c.outcome == OutcomePartial && read != 0 {
				t.Errorf("an oversized report reached the reader (%d bytes)", read)
			}
		})
	}
}

func TestCollectRealProcessTimeoutEndsDescendantsAndLaterInstanceRuns(t *testing.T) {
	helperEnv(t)
	dir := project(t)
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	exe, _ := os.Executable()
	base := []string{"-test.run=^TestHelperProcess$", "--"}
	tool := func(name string, timeout int, args ...string) CapabilityConfig {
		cc := lizardInstance(name, exe, name+"/**")
		cc.Args, cc.TimeoutSeconds = append(append([]string{}, base...), args...), timeout
		return cc
	}
	got := collect(t, dir, cfgOf(tool("tree", 1, "child", pidFile), tool("after", 0, "ok")), Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, ExecRunner{})

	if r := collectedFor(t, got, CapabilityComplexity, "tree").Result; r.Outcome != OutcomeTimedOut {
		t.Errorf("tree = %s (%s), want timed_out", r.Outcome, r.Reason)
	}
	if r := collectedFor(t, got, CapabilityComplexity, "after").Result; r.Outcome != OutcomeAvailable {
		t.Errorf("after = %s (%s); a timed-out neighbour must not affect it", r.Outcome, r.Reason)
	}
	waitGone(t, readPid(t, pidFile))
}
