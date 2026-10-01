package codehealth

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestHelperProcess is not a test: the tests below re-run the test binary as a
// stand-in analysis tool and select its behavior with the arguments after "--".
func TestHelperProcess(t *testing.T) {
	if os.Getenv("SAVEPOINT_HELPER_PROCESS") != "1" {
		return
	}
	args := os.Args
	for i, a := range args {
		if a == "--" {
			args = args[i+1:]
			break
		}
	}
	switch args[0] {
	case "ok":
		os.Stdout.WriteString("report-bytes")
	case "env":
		os.Stdout.WriteString(os.Getenv("CI") + "|" + os.Getenv("NO_COLOR") + "|" + os.Getenv("GIT_TERMINAL_PROMPT"))
	case "exit":
		code, _ := strconv.Atoi(args[1])
		os.Stderr.WriteString("bad\x1b[31m thing\nsecond line " + strings.Repeat("x", 500))
		os.Exit(code)
	case "huge":
		chunk := bytes.Repeat([]byte("a"), 1<<20)
		for i := 0; i < MaxReportBytes>>20+2; i++ {
			os.Stdout.Write(chunk)
		}
	case "hang":
		time.Sleep(time.Minute)
	case "child":
		// Start a grandchild that also hangs and record its pid.
		child := helperCommand("hang")
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		os.WriteFile(args[1], []byte(strconv.Itoa(child.Process.Pid)), 0o644)
		time.Sleep(time.Minute)
	case "report":
		os.WriteFile(args[1], []byte("file-report"), 0o644)
	case "report-then-exit":
		os.WriteFile(args[1], []byte("file-report"), 0o644)
		os.Exit(1)
	case "silent":
	}
	os.Exit(0)
}

func helperSpec(dir string, args ...string) ToolSpec {
	exe, _ := os.Executable()
	return ToolSpec{Dir: dir, Executable: exe, Args: append([]string{"-test.run=^TestHelperProcess$", "--"}, args...)}
}

func helperCommand(args ...string) *execCmd {
	s := helperSpec("", args...)
	return newExecCmd(s)
}

func helperEnv(t *testing.T) {
	t.Helper()
	t.Setenv("SAVEPOINT_HELPER_PROCESS", "1")
}

func TestExecRunnerResults(t *testing.T) {
	helperEnv(t)
	dir := t.TempDir()
	ctx := context.Background()

	t.Run("stdout", func(t *testing.T) {
		res, err := ExecRunner{}.Run(ctx, helperSpec(dir, "ok"))
		if err != nil || string(res.Stdout) != "report-bytes" || res.ExitCode != 0 {
			t.Fatalf("got %q exit %d err %v", res.Stdout, res.ExitCode, err)
		}
	})
	t.Run("environment is pinned", func(t *testing.T) {
		res, err := ExecRunner{}.Run(ctx, helperSpec(dir, "env"))
		if err != nil || string(res.Stdout) != "1|1|0" {
			t.Fatalf("got %q err %v", res.Stdout, err)
		}
	})
	t.Run("non-zero exit is a result with a bounded sanitized stderr", func(t *testing.T) {
		res, err := ExecRunner{}.Run(ctx, helperSpec(dir, "exit", "7"))
		if err != nil || res.ExitCode != 7 {
			t.Fatalf("got exit %d err %v", res.ExitCode, err)
		}
		if strings.ContainsAny(res.Stderr, "\x1b\n") || !strings.HasPrefix(res.Stderr, "bad [31m thing second line") || len(res.Stderr) > MaxReasonLen {
			t.Fatalf("stderr not sanitized and bounded: %q", res.Stderr)
		}
	})
	t.Run("missing executable is unavailable", func(t *testing.T) {
		_, err := ExecRunner{}.Run(ctx, ToolSpec{Dir: dir, Executable: filepath.Join(dir, "no-such-tool")})
		if !errors.Is(err, ErrToolUnavailable) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("not runnable is unavailable", func(t *testing.T) {
		notExec := filepath.Join(dir, "plain.txt")
		os.WriteFile(notExec, []byte("data"), 0o644)
		_, err := ExecRunner{}.Run(ctx, ToolSpec{Dir: dir, Executable: notExec})
		if !errors.Is(err, ErrToolUnavailable) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("oversized stdout is truncated, not failed", func(t *testing.T) {
		res, err := ExecRunner{}.Run(ctx, helperSpec(dir, "huge"))
		if err != nil || res.ExitCode != 0 || !res.Truncated || len(res.Stdout) != MaxReportBytes {
			t.Fatalf("got %d bytes truncated=%v exit %d err %v", len(res.Stdout), res.Truncated, res.ExitCode, err)
		}
	})
	t.Run("runs in the given directory without a shell", func(t *testing.T) {
		// A shell would expand this; an argument vector passes it through.
		res, err := ExecRunner{}.Run(ctx, helperSpec(dir, "ok", "$(echo hi)", "a;b"))
		if err != nil || res.ExitCode != 0 {
			t.Fatalf("got %v exit %d", err, res.ExitCode)
		}
	})
}

func TestExecRunnerStopsOnDeadlineAndCancel(t *testing.T) {
	helperEnv(t)
	dir := t.TempDir()

	t.Run("deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		start := time.Now()
		_, err := ExecRunner{}.Run(ctx, helperSpec(dir, "hang"))
		if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > 10*time.Second {
			t.Fatalf("got %v after %s", err, time.Since(start))
		}
	})
	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() { time.Sleep(200 * time.Millisecond); cancel() }()
		_, err := ExecRunner{}.Run(ctx, helperSpec(dir, "hang"))
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("already cancelled never starts", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if _, err := (ExecRunner{}).Run(ctx, helperSpec(dir, "hang")); !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v", err)
		}
	})
	t.Run("children die with the tool", func(t *testing.T) {
		pidFile := filepath.Join(dir, "child.pid")
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			for i := 0; i < 100; i++ {
				if _, err := os.Stat(pidFile); err == nil {
					break
				}
				time.Sleep(50 * time.Millisecond)
			}
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()
		ExecRunner{}.Run(ctx, helperSpec(dir, "child", pidFile))
		raw, err := os.ReadFile(pidFile)
		if err != nil {
			t.Fatalf("child never started: %v", err)
		}
		pid, _ := strconv.Atoi(string(raw))
		deadline := time.Now().Add(5 * time.Second)
		for processAlive(pid) {
			if time.Now().After(deadline) {
				t.Fatalf("child process %d survived cancellation", pid)
			}
			time.Sleep(50 * time.Millisecond)
		}
	})
}
