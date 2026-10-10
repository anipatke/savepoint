package run

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

// TestMain doubles as the fake provider: a test starts this binary with
// FAKE_AGENT set and the process behaves as the named script.
func TestMain(m *testing.M) {
	switch os.Getenv("FAKE_AGENT") {
	case "":
		os.Exit(m.Run())
	case "ok":
		fmt.Println(`{"type":"system"}`)
		fmt.Println("not json at all")
		fmt.Print(`{"type":"result"}`) // no trailing newline
		if os.Getenv("CI") != "1" || os.Getenv("NO_COLOR") != "1" {
			os.Exit(9)
		}
	case "fail":
		fmt.Println(`{"type":"system"}`)
		fmt.Fprintln(os.Stderr, "progress\x1b[31m")
		fmt.Fprintln(os.Stderr, "fatal: could not read https://user:secret@example.com/repo")
		os.Exit(3)
	case "replay": // write the file named by FAKE_REPLAY to stdout, as a provider would
		b, err := os.ReadFile(os.Getenv("FAKE_REPLAY"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(4)
		}
		os.Stdout.Write(b)
	case "idle", "tree", "stubborn", "sleeper":
		fakeStoppable(os.Getenv("FAKE_AGENT"))
	case "quit": // end at once, to leave behind a pid that is no longer running
	case "slow":
		fmt.Println("started")
		time.Sleep(600 * time.Millisecond)
	}
	os.Exit(0)
}

func project(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".savepoint"), 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func fakeSpec(t *testing.T, dir, mode string) Spec {
	t.Helper()
	t.Setenv("FAKE_AGENT", mode)
	return Spec{Dir: dir, Task: "T-001", Provider: "fake", Argv: []string{os.Args[0]}, NextLine: "Build T-001"}
}

func readRecord(t *testing.T, dir string) Record {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "run.yml"))
	if err != nil {
		t.Fatal(err)
	}
	var rec Record
	if err := yaml.Unmarshal(b, &rec); err != nil {
		t.Fatalf("run.yml does not parse: %v\n%s", err, b)
	}
	return rec
}

func TestStart_finishedRunKeepsStdoutVerbatim(t *testing.T) {
	dir := project(t)
	r, err := Start(fakeSpec(t, dir, "ok"))
	if err != nil {
		t.Fatal(err)
	}
	if base := filepath.Base(r.Dir); !strings.HasSuffix(base, "-T-001") || len(base) != len("20060102T150405Z-T-001") {
		t.Errorf("run directory name = %q", base)
	}
	if filepath.Dir(r.Dir) != filepath.Join(dir, ".savepoint", "runs") {
		t.Errorf("run directory = %q", r.Dir)
	}
	final := r.Wait()
	if final.State != StateFinished || final.ExitCode == nil || *final.ExitCode != 0 || final.StderrTail != "" {
		t.Fatalf("final = %+v", final)
	}
	if got := readRecord(t, r.Dir); got.State != StateFinished || got.PID == 0 || got.Provider != "fake" ||
		got.Cwd != dir || got.Task != "T-001" || got.NextLine != "Build T-001" || got.StartedAt == "" || got.UpdatedAt == "" {
		t.Errorf("run.yml = %+v", got)
	}
	events, err := os.ReadFile(filepath.Join(r.Dir, "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	want := "{\"type\":\"system\"}\nnot json at all\n{\"type\":\"result\"}\n"
	if string(events) != want {
		t.Errorf("events.jsonl = %q, want %q", events, want)
	}
}

func TestStart_recordsRunningBeforeAndWhileChildRuns(t *testing.T) {
	dir := project(t)
	r, err := Start(fakeSpec(t, dir, "slow"))
	if err != nil {
		t.Fatal(err)
	}
	rec := readRecord(t, r.Dir)
	if rec.State != StateRunning || rec.PID == 0 || rec.ExitCode != nil {
		t.Errorf("record while running = %+v", rec)
	}
	if final := r.Wait(); final.State != StateFinished {
		t.Errorf("final = %+v", final)
	}
}

func TestStart_failedRunRecordsExitCodeAndSanitisedStderrTail(t *testing.T) {
	r, err := Start(fakeSpec(t, project(t), "fail"))
	if err != nil {
		t.Fatal(err)
	}
	final := r.Wait()
	if final.State != StateFailed || final.ExitCode == nil || *final.ExitCode != 3 {
		t.Fatalf("final = %+v", final)
	}
	tail := final.StderrTail
	if strings.ContainsAny(tail, "\x1b\n") || strings.Contains(tail, "secret") || strings.Contains(tail, "user:") ||
		!strings.Contains(tail, "https://example.com/repo") {
		t.Errorf("stderr_tail = %q", tail)
	}
	if got := readRecord(t, r.Dir); got.State != StateFailed || got.StderrTail != tail {
		t.Errorf("run.yml = %+v", got)
	}
}

func TestStart_heartbeatRewritesUpdatedAtWhileRunning(t *testing.T) {
	dir := project(t)
	spec := fakeSpec(t, dir, "slow")
	spec.Heartbeat = 20 * time.Millisecond
	var tick atomic.Int64
	base := time.Date(2026, 10, 11, 9, 0, 0, 0, time.UTC)
	spec.Now = func() time.Time { return base.Add(time.Duration(tick.Add(1)) * time.Second) }
	r, err := Start(spec)
	if err != nil {
		t.Fatal(err)
	}
	initial := readRecord(t, r.Dir).UpdatedAt
	deadline := time.Now().Add(5 * time.Second)
	for {
		rec := readRecord(t, r.Dir) // never a partial file: it always parses
		if rec.State == StateRunning && rec.UpdatedAt != initial {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("updated_at never moved from %s while running", initial)
		}
		time.Sleep(10 * time.Millisecond)
	}
	r.Wait()
}

func TestStart_refusesWhatItCannotRecordSafely(t *testing.T) {
	cases := map[string]func(t *testing.T) Spec{
		"not a Savepoint project": func(t *testing.T) Spec {
			return fakeSpec(t, t.TempDir(), "ok")
		},
		"empty command": func(t *testing.T) Spec {
			s := fakeSpec(t, project(t), "ok")
			s.Argv = nil
			return s
		},
		"task escapes the runs directory": func(t *testing.T) Spec {
			s := fakeSpec(t, project(t), "ok")
			s.Task = "../T-001"
			return s
		},
	}
	for name, mk := range cases {
		t.Run(name, func(t *testing.T) {
			spec := mk(t)
			if _, err := Start(spec); err == nil {
				t.Fatal("Start succeeded")
			}
			if _, err := os.Stat(filepath.Join(spec.Dir, ".savepoint", "runs")); !os.IsNotExist(err) {
				t.Errorf("runs directory was written: %v", err)
			}
		})
	}
}

func TestStart_missingExecutableIsRecordedAsFailed(t *testing.T) {
	dir := project(t)
	spec := fakeSpec(t, dir, "ok")
	spec.Argv = []string{filepath.Join(dir, "no-such-agent")}
	if _, err := Start(spec); err == nil {
		t.Fatal("Start succeeded")
	}
	entries, err := os.ReadDir(filepath.Join(dir, ".savepoint", "runs"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("runs = %v, %v", entries, err)
	}
	rec := readRecord(t, filepath.Join(dir, ".savepoint", "runs", entries[0].Name()))
	if rec.State != StateFailed || rec.StderrTail == "" {
		t.Errorf("record = %+v", rec)
	}
}

func TestSanitizeTail_keepsEndOfOversizedText(t *testing.T) {
	got := sanitizeTail(strings.Repeat("x", 2000) + " the end")
	if !strings.HasPrefix(got, "…") || !strings.HasSuffix(got, "the end") || len(got) > maxTailLen+len("…") {
		t.Errorf("tail = %q", got)
	}
}
