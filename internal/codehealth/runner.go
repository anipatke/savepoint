package codehealth

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// ErrToolUnavailable means the executable was not found or cannot be started.
var ErrToolUnavailable = errors.New("tool is not available")

// Limits on what one executed tool may hand back. Stdout is the report when no
// {report} file is used, so it gets the report cap; stderr only feeds a short
// reason.
const (
	MaxReportBytes = 32 << 20
	maxStderrBytes = 64 << 10
	// killWait bounds how long a cancelled tool's output pipes may stay open,
	// for example when a child process outlives its parent.
	killWait = 2 * time.Second
)

// ToolSpec is one process to run: an executable and an argument vector, never a
// shell string. Dir is the working directory.
type ToolSpec struct {
	Dir        string
	Executable string
	Args       []string
}

// ToolResult is a tool that ran to completion. A non-zero ExitCode is a result,
// not an error. Truncated reports that stdout passed MaxReportBytes and the
// excess was dropped. Stderr is a bounded, sanitized single line that keeps its end.
type ToolResult struct {
	Stdout    []byte
	Truncated bool
	Stderr    string
	ExitCode  int
}

// ToolRunner runs one analysis tool. It returns ErrToolUnavailable when the
// process cannot be started and the context's error when it is stopped.
type ToolRunner interface {
	Run(ctx context.Context, spec ToolSpec) (ToolResult, error)
}

// ExecRunner runs tools directly with os/exec. The environment is the caller's
// plus settings that keep tools from prompting or coloring their output.
type ExecRunner struct{}

// Run implements ToolRunner. On cancellation the whole process group (Unix) is
// killed, and output pipes are force-closed after killWait everywhere.
func (ExecRunner) Run(ctx context.Context, spec ToolSpec) (ToolResult, error) {
	if err := ctx.Err(); err != nil {
		return ToolResult{}, err
	}
	cmd := exec.CommandContext(ctx, spec.Executable, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = append(os.Environ(), "CI=1", "NO_COLOR=1", "GIT_TERMINAL_PROMPT=0")
	cmd.WaitDelay = killWait
	configureProcess(cmd)
	out := &limitedBuffer{limit: MaxReportBytes}
	errOut := &limitedBuffer{limit: maxStderrBytes}
	cmd.Stdout, cmd.Stderr = out, errOut

	if err := cmd.Start(); err != nil {
		if ctx.Err() != nil {
			return ToolResult{}, ctx.Err()
		}
		// Not found, not executable, and any other failure to start all mean
		// the project's tool cannot be used.
		return ToolResult{}, ErrToolUnavailable
	}
	err := cmd.Wait()
	if ctx.Err() != nil {
		return ToolResult{}, ctx.Err()
	}
	res := ToolResult{Stdout: out.Bytes(), Truncated: out.truncated, Stderr: sanitizeTail(errOut.String())}
	var exit *exec.ExitError
	switch {
	case err == nil:
	case errors.As(err, &exit):
		res.ExitCode = exit.ExitCode()
	case errors.Is(err, exec.ErrWaitDelay):
		// The tool finished but a child kept the pipes open; its output is
		// whatever arrived before they were closed.
	default:
		return ToolResult{}, err
	}
	return res, nil
}

// limitedBuffer keeps the first limit bytes and drops the rest while still
// accepting writes, so a chatty tool is never blocked or failed by the cap.
//
// It wraps a Buffer rather than embedding it: an embedded Buffer's ReadFrom
// would let io.Copy bypass Write and the limit.
type limitedBuffer struct {
	buf       bytes.Buffer
	limit     int
	truncated bool
}

func (b *limitedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if room := b.limit - b.buf.Len(); n > room {
		b.truncated = true
		p = p[:room]
	}
	b.buf.Write(p)
	return n, nil
}

func (b *limitedBuffer) Bytes() []byte  { return b.buf.Bytes() }
func (b *limitedBuffer) String() string { return b.buf.String() }

// cleanLine reduces tool text to one printable line, without bounding it.
func cleanLine(s string) string {
	s = strings.ToValidUTF8(s, "")
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
	return strings.Join(strings.Fields(s), " ")
}

// sanitizeLine reduces tool text to one bounded, printable line for a reason.
func sanitizeLine(s string) string {
	s = cleanLine(s)
	if len(s) > MaxReasonLen/2 {
		cut := MaxReasonLen / 2
		for cut > 0 && !utf8.RuneStart(s[cut]) {
			cut--
		}
		s = s[:cut] + "…"
	}
	return s
}

// sanitizeTail is sanitizeLine for a tool's stderr: tools print the error last,
// after any progress output, so an oversized line keeps its end, not its start.
func sanitizeTail(s string) string {
	s = cleanLine(s)
	if len(s) > MaxReasonLen/2 {
		cut := len(s) - MaxReasonLen/2
		for cut < len(s) && !utf8.RuneStart(s[cut]) {
			cut++
		}
		s = "…" + s[cut:]
	}
	return s
}
