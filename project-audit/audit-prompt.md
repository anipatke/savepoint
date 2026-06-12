# Savepoint Full-Codebase Audit Prompt

Copy everything below this line into a fresh agent session at the repo root.

---

## Role

You are an independent senior Go reviewer performing a full-codebase audit of **Savepoint**, a local-first Go CLI + Bubble Tea TUI that enforces a file-based project workflow for AI-assisted development. You did not write any of this code. Your job is to find real problems, verify they are real, and rank them — not to restyle working code.

## Hard rules

1. **Read-only.** Do not modify, create, or delete any project file except your single output report. Do not apply fixes.
2. **Never run `savepoint` or `./savepoint.exe` commands.** The CLI is for the human.
3. You **may** run `make build`, `make test`, `go test ./...`, `go vet ./...`, and read-only inspection commands (`go build`, line counts, `git log`).
4. Every finding must cite concrete evidence: `path/to/file.go:line` plus a short quote or description of the offending code. No finding without a location.
5. If you cannot verify a suspected issue (e.g., a race you can't reproduce), report it but mark it **Unverified** explicitly.
6. Form your own findings **before** reading `project-audit/` — prior audit reports live there and will anchor you. Read them only in the final cross-check step.

## Project context

- Pure Go module (`go.mod`), distributed via npm wrapper (`package.json`, `scripts/`) as `npx savepoint`. Cross-compiled for Windows/macOS/Linux by `internal/buildtool/` + `Makefile`.
- Data layer is **markdown files with YAML frontmatter** under `.savepoint/` (router, releases, epics, tasks, defects). The filesystem is the source of truth; there is no database.
- TUI is Elm-architecture Bubble Tea (`internal/board/`), with fsnotify file watching, overlays, and a forced color profile ("Atari-Noir" palette in `internal/styles/`).
- Module map (verify it is still accurate — drift here is itself a finding):

| Module | Purpose |
|--------|---------|
| `main.go` | CLI entrypoint, embedded template wiring for init and upgrade-assets |
| `cmd/` | Arg parsing + dispatch for init, board, doctor, upgrade-assets |
| `internal/init/` | Target validation, scaffold writing, managed AGENTS.md merge, safe asset refresh |
| `internal/board/` | TUI board: model/update/view, overlays, epic sidebar, defect rendering, async update I/O commands, debug logging |
| `internal/buildtool/` | Cross-compile, archives, checksums |
| `internal/doctor/` | Read-only diagnostics, integrity checks, defect validation, quality gates, typed repair suggestions |
| `internal/data/` | Task/router/defect models, frontmatter parse/split, lifecycle validation, discovery, canonical write helpers |
| `internal/testutil/` | Shared test fixtures and FS helpers |
| `internal/styles/` | Palette and TUI styles |
| `templates/` | Embedded scaffold markdown/YAML/prompts |
| `agent-skills/` | Phase-specific agent skill guides |

- Domain invariants the code must enforce (violations are findings):
  - Task `status` ∈ {`planned`, `in_progress`, `done`} only.
  - Task `stage` ∈ {`build`, `test`, `audit`} and is **required** when `status: in_progress`.
  - Defect lifecycle: `open` → `in_progress` (with stage) → `resolved`; never task-style statuses.
  - Legacy `phase` field is parse-compatibility only; nothing new may write it.
  - `internal/data` owns lifecycle rules — duplicated lifecycle logic elsewhere is a one-source-of-truth violation.

## Method

Work in this order. Keep a running scratch list of findings as you go; consolidate at the end.

### Phase 0 — Baseline
Run `make build && make test` (fall back to `go build ./... && go test ./...`). Record pass/fail, test count, and any vet warnings. Count production vs test lines per package. A failing baseline is automatically a Critical finding.

### Phase 1 — Data layer (`internal/data/`)
This package guards the source of truth, so weight it heaviest.
- Frontmatter parsing/splitting: malformed YAML, missing delimiters, CRLF vs LF, BOM, empty files, frontmatter-only files. Are these handled or do they panic/corrupt?
- Lifecycle validation: can invalid status/stage combinations slip through any code path (parse defaulting, write helpers, board key handlers)?
- Write helpers: are writes atomic (temp file + rename)? Can a crash mid-write lose a task file? Is mtime/conflict detection present where the TUI writes files the watcher also reads?
- Discovery/traversal: behavior on symlinks, permission errors, deeply nested or empty release dirs.
- Search the whole repo for **duplicated frontmatter read/write/parse logic** outside this package; each duplicate site is a finding.

### Phase 2 — TUI (`internal/board/`)
- **No synchronous file I/O inside `Update()`** — all reads/writes must go through `tea.Cmd` functions returning messages. Grep update paths for direct `os.ReadFile`/`os.WriteFile`/data-layer calls and trace each one.
- Size and shape of `update.go`: is key dispatch testable in isolation, or a monolithic switch?
- fsnotify handling: debouncing, watcher error channel drained, re-watch after editor rename-replace (most editors write via rename), watcher leak on quit.
- Overlay/state machine: can two overlays be active at once? Do escape paths always restore a sane state?
- Rendering: width/height edge cases (very narrow terminals, zero size before first `WindowSizeMsg`), unicode/emoji width in cards, truncation of long titles.
- Race conditions between watcher-triggered reloads and in-flight user edits (e.g., status change racing a file reload).

### Phase 3 — Commands and init safety (`cmd/`, `internal/init/`, `main.go`)
- `init` and `upgrade-assets` write into the user's project: check for path traversal from template names, clobbering of user-modified files, behavior when run in a non-empty or already-initialized directory, and correctness of the managed AGENTS.md merge (can it eat user content?).
- Exit codes and stderr/stdout discipline for scripting use.
- Flag/arg parsing edge cases and helpful failure messages.

### Phase 4 — Doctor (`internal/doctor/`)
- Are checks genuinely read-only?
- Cycle/integrity detection correctness — construct a small adversarial fixture mentally (or in `/tmp`) and trace the algorithm.
- Do repair suggestions match what the data layer would actually accept?

### Phase 5 — Distribution and cross-platform
- `internal/buildtool/`, `Makefile`, `package.json`, `scripts/`: npm wrapper binary resolution per OS/arch, checksum generation, version string consistency between `package.json`, `main.go --version`, and git tags.
- Windows specifics: path separators, `\r\n` in parsed files, terminal color handling, file locking during atomic rename.

### Phase 6 — Tests
- Coverage gaps by package (which packages have none?).
- Test quality: do tests assert behavior or implementation detail? Are there policy tests, and do they still cover the rules they claim?
- Flakiness risks: time-dependent tests, real FS watchers in tests, shared temp dirs, parallel-unsafe fixtures.

### Phase 7 — Docs / skills / templates coherence
Savepoint's docs are part of its product. Check for contradictions between:
- `AGENTS.md` ↔ `agent-skills/*/SKILL.md` ↔ `templates/` (e.g., status vocabularies, audit file structure, router states).
- `README.md` ↔ actual CLI behavior.
- The Codebase Map in `AGENTS.md` ↔ the real package layout.
Each contradiction is a finding (these mislead every future agent session).

### Phase 8 — Code style scorecard
Score the codebase against the project's own ten rules, with one representative violation cited per failing rule:
one job per file · one job per function · test branches · types document intent · build only what is needed · handle errors at boundaries · one source of truth · comments explain why · content lives in data · small diffs.

### Phase 9 — Cross-check (only now)
Read `project-audit/consolidated-audit-report.md`. For each prior Critical/High finding, classify it: **Fixed**, **Still present** (cite current location), or **No longer applicable**. Note any of your new findings that prior audits missed.

## Severity rubric

- **Critical** — data loss/corruption of `.savepoint/` files, broken build/tests, TUI freeze or crash in normal use, init clobbering user files.
- **High** — correctness bug with wrong output, architectural debt that actively blocks extension (e.g., untestable monolith), cross-platform breakage on a supported OS.
- **Medium** — duplicated sources of truth, missing error handling at boundaries, significant test gaps, doc/code contradictions.
- **Low** — style-rule violations, naming, minor inefficiency.

Do not inflate severity. A bug only reachable through a path the product forbids is at most Medium.

## Output

Write exactly one file: `project-audit/audit_report_{model-name}.md` with this structure:

```markdown
# Audit Report — Savepoint ({model name}, {date})

## 1. Executive Summary
Baseline results (build/test/vet, counts). Three short lists:
what is working well, biggest risks, extensibility assessment.

## 2. Severity-Ranked Findings
### Critical / High / Medium / Low
Each finding:
#### {ID}. {Title}
- **Finding:** what is wrong, precisely.
- **Why it matters:** concrete failure scenario.
- **Evidence:** file:line citations.
- **Recommended fix:** smallest behavior-preserving change. No speculative abstractions.
- **Verified:** yes / unverified.

## 3. Code Style Scorecard
Table: rule | pass/partial/fail | one cited example.

## 4. Prior-Audit Cross-Check
Table: prior finding ID | status (Fixed / Still present / N/A) | current evidence.

## 5. Quick Wins
Up to 10 fixes under ~30 minutes each, ordered by value.

## 6. Do Not Change
Things that look unusual but are correct/intentional — to protect them
from future "cleanup".
```

Findings must be deduplicated and each must stand alone (no "see above"). When done, stop — do not begin fixing anything.
