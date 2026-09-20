---
type: audit-findings
audited: 2026-09-19
---

# Audit Findings: E49 Operate V2 through a clear terminal board

## Main Findings

### Verdict

**NEEDS WORK.** E49's board architecture, owner-write boundaries, reload behavior, evidence views, and board/resume parity are well supported, but one repository issue remains: authored terminal escape sequences can pass through the V2 record reader into supposedly ANSI-free plain output. No owner-run evidence or waiver is otherwise outstanding. **NOT READY TO COMMIT/PUSH.**

### What Needs Attention

#### 1. Authored control sequences escape into plain output

Task titles are accepted when non-empty, including YAML-escaped terminal control sequences, and the non-TTY renderer writes the decoded title directly. A title such as `"\u001b[31mred\u001b[0m"` therefore places ANSI bytes in redirected output, contradicting T010's requirement that non-TTY output contain no ANSI escapes. It also means project-authored text can change how a terminal renders piped output. Strip terminal sequences at the final plain-output boundary and add a regression using an on-disk V2 project whose selected Task title contains YAML-escaped ANSI.

Evidence: `internal/data/task_v2.go:89` validates only non-emptiness; `internal/board/v2/plain.go:57` interpolates `card.Task.Title` unchanged, while `internal/board/v2/run_test.go` tests ANSI absence only with ordinary titles.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| Authored control sequences escape into plain output | Low | Medium | Medium | Fix now at the plain-render boundary and add an on-disk regression covering ANSI/control input. |

### What Is Proven / Not Proven

Proven: schema dispatch and invalid-project diagnostics; empty V2 projects; three Task columns and Objective ownership filtering; all typed badge and clearance states; Task/Objective detail and complete Check history; the unified five-type Issues surface; owner-only acceptance, exception completion, and selection writes with migration/freshness guards; watcher exclusions, debounce, reload recovery, and identity preservation; narrow layouts and grapheme-safe CJK/emoji/combining-mark truncation; no-write browsing; and board/resume agreement across pending migration, replan, Task and Objective dependency, execute, check-needed, owner-validation, Objective-integration, ready, and plan rungs.

Not proven because it fails inspection: ANSI-free non-TTY output for the supported class of authored scalar text containing terminal escapes. No other acceptance criterion remains unverified.

### Audit Evidence

- Scope lock: all T001-T010 acceptance criteria; E49's dispatch, V2 TUI/plain surfaces, router/evidence/status writes, watcher/reload workflow, and shared Next wording; direct `data`, `migrate`, and `resume` dependencies; ARCH-01/02/04, DATA-01/02/03, FS-01/04/05/06, CFG-01/02, DEP-02, TEST-01/02/04/06/08, and applicable blocker policy.
- Coverage and workflow result: normal, empty, malformed, stale-selection, pending-migration, conflict, repeated-write, failed-reload/recovery, TTY/non-TTY, `NO_COLOR`, `TERM=dumb`, narrow/wide, Unicode-width, and control-text cells were classified. Selection, owner acceptance, exception completion, startup/reload, watcher debounce, and final-output side effects follow the intended order and preserve the primary failure. The control-text/non-TTY cell is the sole failure.
- File reality and drift: every path in the E49 delivery commit and task context logs exists. The `internal/resume` phrase exports and the third detail file are deliberate, pure-presentation splits reflected in AGENTS.md and consistent with the V2 design; no architecture blocker results.
- Gates: focused tests passed; `go test -race ./internal/board/v2` passed; uncached `go test -count=1 ./...` passed; `NO_COLOR=1` and `TERM=dumb` V2 suites passed; `go vet ./...`, `make build`, `make test`, `gofmt` scope check, and `git diff --check` passed.

### Guardrails Verification

- Rule IDs checked: ARCH-01, ARCH-02, ARCH-04, DATA-01, DATA-02, DATA-03, FS-01, FS-04, FS-05, FS-06, CFG-01, CFG-02, DEP-02, TEST-01, TEST-02, TEST-04, TEST-06, TEST-08, POL-01, and the STYLE guidelines. Template, distribution, dry-run, and network rules are not touched by E49.
- Health check mode: Full.
- Evidence: the focused, race, environment, full-suite, vet, build, formatting, and diff commands listed above.
- File reality evidence: all 69 delivery paths from `ae039d8..a01f31c` resolve in the current tree; the worktree was clean before this audit artifact.
- Waivers or unresolved findings: no waiver; Finding 1 remains open.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config.
- [x] STYLE-08 **Comments explain why** — not what the code already says.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic.
- [ ] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes. The single E49 delivery commit is approximately 12,600 added lines across 69 files; this is advisory and does not cause the verdict.

## Proposed Changes

### Target File
internal/board/v2/width.go

### Replace
```go
import (
	"strings"

	xansi "github.com/charmbracelet/x/ansi"
)
```

### With
```go
import (
	"strings"
	"unicode"

	xansi "github.com/charmbracelet/x/ansi"
)
```

### Target File
internal/board/v2/width.go

### Replace
```go
func fitLine(text string, width int) string {
	return truncateCells(strings.ReplaceAll(text, "\n", " "), terminalWidthOrOne(width))
}
```

### With
```go
func fitLine(text string, width int) string {
	return truncateCells(strings.ReplaceAll(text, "\n", " "), terminalWidthOrOne(width))
}

// stripTerminalControls makes redirected output plain even when authored
// record text contains YAML-escaped ANSI or other terminal control bytes.
// Newlines produced by the renderer remain structural; authored tabs become
// spaces and other controls are removed.
func stripTerminalControls(text string) string {
	text = xansi.Strip(text)
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n':
			return r
		case r == '\t':
			return ' '
		case unicode.IsControl(r):
			return -1
		default:
			return r
		}
	}, text)
}
```

### Target File
internal/board/v2/plain.go

### Replace
```go
	return b.String()
}
```

### With
```go
	return stripTerminalControls(b.String())
}
```

### Target File
internal/board/v2/run_test.go

### Replace
```go
func TestRunWithoutTTYIsDeterministic(t *testing.T) {
```

### With
```go
func TestRunWithoutTTYStripsAuthoredTerminalControls(t *testing.T) {
	root := writeValidProject(t)
	writeTask(t, root, "O001", "T001", `\u001b[31mred\u001b[0m\tname`, "status: planned\n")

	var stdout bytes.Buffer
	if err := Run(Options{Root: root, Stdout: &stdout, TTY: false}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if strings.Contains(stdout.String(), "\x1b") {
		t.Fatalf("plain output contains authored ANSI escapes: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "red name") {
		t.Errorf("plain output did not preserve readable title text: %q", stdout.String())
	}
}

func TestRunWithoutTTYIsDeterministic(t *testing.T) {
```
