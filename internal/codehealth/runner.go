package codehealth

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// ErrToolUnavailable means the executable was not found or cannot be started.
var ErrToolUnavailable = errors.New("tool is not available")

// Limits on what one executed tool may hand back. Stdout is the report when no
// {report} file is used, so it gets the report cap and keeps its start; stderr
// only feeds a short reason and keeps its end.
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
	errOut := &tailBuffer{limit: maxStderrBytes}
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
	res := ToolResult{Stdout: out.Bytes(), Truncated: out.truncated, Stderr: sanitizeTail(errOut.Text())}
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

// tailBuffer keeps the last limit bytes and drops the earlier ones, so the
// error a tool prints last survives however much progress came first. Like
// limitedBuffer it always accepts writes.
type tailBuffer struct {
	buf     []byte
	limit   int
	dropped bool
}

func (b *tailBuffer) Write(p []byte) (int, error) {
	b.buf = append(b.buf, p...)
	// Trim in batches so a chatty tool costs amortized constant time per byte.
	if len(b.buf) > 2*b.limit {
		b.dropped = true
		b.buf = append(b.buf[:0], b.buf[len(b.buf)-b.limit:]...)
	}
	return len(p), nil
}

// cutUserinfo matches the start of text that may be the end of a URL whose
// "scheme://" was cut away: some "user:password@".
var cutUserinfo = regexp.MustCompile(`^[^\s/?#@]*[:/]*[^\s/?#@]*@`)

// Text is the kept tail as stderr text. When earlier bytes were dropped it
// starts mid-word, so a leading "user:password@" fragment is removed.
func (b *tailBuffer) Text() string {
	s := b.String()
	if b.dropped || len(b.buf) > b.limit {
		s = cutUserinfo.ReplaceAllString(s, "")
	}
	return s
}

func (b *tailBuffer) String() string {
	if len(b.buf) > b.limit {
		return string(b.buf[len(b.buf)-b.limit:])
	}
	return string(b.buf)
}

func (b *limitedBuffer) Bytes() []byte  { return b.buf.Bytes() }
func (b *limitedBuffer) String() string { return b.buf.String() }

// urlUserinfo matches the credentials of a URL, "user:password@" after "://".
var urlUserinfo = regexp.MustCompile(`://[^\s/?#@]*@`)

// cleanLine reduces tool text to one printable line, without bounding it. A
// credential inside a URL is dropped, since the text is saved in permanent
// evidence.
func cleanLine(s string) string {
	s = strings.ToValidUTF8(s, "")
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	return urlUserinfo.ReplaceAllString(s, "://")
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
