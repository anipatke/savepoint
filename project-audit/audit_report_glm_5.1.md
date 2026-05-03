# Codebase Audit Report

## 1. Executive Summary

**Savepoint** is a Go CLI/TUI tool (~5,600 lines of production code, ~5,600 lines of tests, 264 tests) that provides a kanban-style project workflow manager with a terminal UI. It uses the Bubble Tea framework and a custom "Atari-Noir" visual theme.

**What is working well:**
- Clean module separation: `cmd/` (CLI), `internal/data/` (persistence), `internal/board/` (TUI), `internal/doctor/` (diagnostics), `internal/init/` (scaffolding), `internal/styles/` (palette)
- Dependency injection in `cmd/` via function types (`InitRunner`, `BoardRunner`, `DoctorRunner`) makes commands trivially testable
- Good test coverage for the core modules (data, board, doctor, init, cmd) with 264 tests
- Atomic file writes in `data/write.go` and `init/write.go` prevent corruption
- Well-designed responsive TUI layout with breakpoint-based column rendering
- Policy tests (`render_policy_test.go`) enforce visual constraints cross-cutting

**Biggest risks:**
- Synchronous filesystem I/O inside the Bubble Tea `Update()` loop will freeze the TUI on slow disks or large files
- `update.go` (521 lines) is a monolithic switch that handles 5+ overlay types, keyboard input, and file writes in one function
- `doctor/checks.go` (585 lines) reimplements `strings.Contains`/`strings.Index` and the cycle detection has an accuracy issue
- No Windows build target in `buildtool`; no checksums in distribution archives
- `buildtool` and `styles` packages have zero tests

**Extensibility:** The architecture suits the current project size well. Adding new overlays, columns, or commands would be straightforward. The main risk to extensibility is the `update.go` monolith — every new keybind or overlay type adds more branches to an already complex function.

**Architecture fit:** Good for a small-to-medium project. The Bubble Tea pattern (Model-Update-View) is well-applied, and the rendering helpers are pure functions. The data layer correctly separates parsing, discovery, and writing. The main architectural debt is the I/O-in-update pattern and the size of `update.go`.

---

## 2. Severity-Ranked Recommendations

### Critical

#### C1. Synchronous file I/O in the TUI update loop

- **Finding:** `update.go` performs filesystem reads and writes directly inside `Update()`: `writeTaskStatus()`, `writeRouterTask()`, `writeRouterReleaseEpic()`, `readEpicDetailFile()`, `selectEpicPanelEpic()`. These block the TUI event loop.
- **Why it matters:** Any disk latency (network drives, slow SSDs, virus scanners) freezes the entire TUI. Bubble Tea's design intent is for `Update()` to be pure — I/O should happen in `tea.Cmd` functions that return messages.
- **Evidence:** `internal/board/update.go:285-330` (writeRouterTask, writeRouterReleaseEpic), `internal/board/update.go:217-250` (readEpicDetailFile), `internal/board/model.go:235-280` (writeRouterReleaseEpic, writeRouterTask)
- **Recommended fix:** Extract I/O operations into `tea.Cmd` functions. E.g., `writeTaskStatusCmd(task, path, mtime) tea.Cmd` returns a `tea.Msg` on completion. `Update()` dispatches the command and handles the result message. This is the standard Bubble Tea pattern.
- **Estimated effort:** Medium (requires refactoring `update.go` to emit commands for each I/O operation, add result message types, and handle them in separate `Update()` branches)

#### C2. Cycle detection produces inaccurate paths

- **Finding:** `detectCycles` in `checks.go` uses a `parent` map that gets overwritten when a node is visited from multiple paths. When a cycle is found, the path reconstructed via `parent` may not represent the actual cycle.
- **Why it matters:** Users could be shown a cycle path that doesn't actually exist, causing confusion or incorrect doctor reports.
- **Evidence:** `internal/doctor/checks.go` — the DFS `parent` map is a simple `map[string]string` that gets overwritten per-visit
- **Recommended fix:** Either use a stack-based cycle reconstruction (track the current DFS path as a slice) or validate the reconstructed path actually forms a cycle before reporting it.
- **Estimated effort:** Small

### High

#### H1. `update.go` is a monolithic 521-line file handling all input, overlays, and I/O

- **Finding:** The `Update()` method alone is ~190 lines with deeply nested switches. Overlay handling mixes 5 overlay types in one `updateOverlay()` function.
- **Why it matters:** Every new feature (new keybind, new overlay) increases the branching complexity. It's hard to reason about which keys do what in which state.
- **Evidence:** `internal/board/update.go` — single file handling quit, navigation, task transitions, overlays for help/epic/release/detail/epic-detail, file watching, phase management
- **Recommended fix:** Extract `updateOverlay()` into per-type handlers (`updateHelpOverlay`, `updateEpicOverlay`, etc.). Extract board key handling into `handleBoardKey()`. Consider a key-to-handler dispatch map rather than a switch chain.
- **Estimated effort:** Medium

#### H2. Stdlib reimplementation in `repairs.go` and `buildtool/main.go`

- **Finding:** `contains()` and `indexOf()` in `repairs.go` reimplement `strings.Contains` and `strings.Index`. `trimSpace()` in `buildtool/main.go` reimplements `bytes.TrimSpace`.
- **Why it matters:** Makes code harder to read for Go developers expecting standard library calls. The custom implementations may have subtle differences from the stdlib versions (e.g., Unicode whitespace handling).
- **Evidence:** `internal/doctor/repairs.go:67-79`, `internal/buildtool/main.go:211-219`
- **Recommended fix:** Replace with `strings.Contains`, `strings.Index`, and `bytes.TrimSpace` respectively. Import the relevant packages.
- **Estimated effort:** Small

#### H3. Fragile repair suggestion matching via substring search on error messages

- **Finding:** `SuggestRepair()` pattern-matches against error message substrings to suggest fixes. If error messages change format, repairs silently break.
- **Evidence:** `internal/doctor/repairs.go:8-66` — hard-coded substring matching against messages like `"not found"`, `"invalid"`, `"missing"`
- **Recommended fix:** Define typed error codes or sentinel errors in `data/` and `doctor/`, and match on error type rather than substring. E.g., `func SuggestRepair(err error) string` switches on `errors.Is(err, ErrNoFrontmatter)` instead of `contains(err.Error(), "no frontmatter")`.
- **Estimated effort:** Medium

#### H4. Quality gate command execution has no timeout

- **Finding:** `RunQualityGates` executes commands from `config.yml` with no timeout. A hung command blocks the doctor indefinitely.
- **Evidence:** `internal/doctor/gates.go:32-55` — `exec.Command` with `Run()` and no `Context` timeout
- **Recommended fix:** Use `exec.CommandContext(ctx, ...)` with a configurable timeout (default: 60s). Add a `gate_timeout` config option.
- **Estimated effort:** Small

#### H5. `\r\n` normalization is scattered across files instead of centralized

- **Finding:** `strings.ReplaceAll(content, "\r\n", "\n")` appears in `parser.go`, `write.go` (3 times), `epic_panel.go`, and likely elsewhere. This is a cross-cutting concern that should be handled once.
- **Why it matters:** If the normalization logic changes (e.g., also handling `\r` alone), every call site must be found and updated. Missing a call site causes subtle cross-platform bugs.
- **Evidence:** `internal/data/parser.go:92`, `internal/data/write.go:22,46,104`, `internal/board/epic_panel.go`
- **Recommended fix:** Create a `normalizeLineEndings(s string) string` function in `internal/data/` and use it everywhere. Or read files with a `bufio.Scanner` which handles line endings, or use `strings.ReplaceAll` once at the point of file reading.
- **Estimated effort:** Small

### Medium

#### M1. Dead code in `column.go` and `report.go`

- **Finding:** `taskLabel()` in `column.go` is defined but never called. `CheckResult` in `report.go` is defined but never used. `GateSuggestion` has an unused `exitCode` parameter.
- **Evidence:** `internal/board/column.go`, `internal/doctor/report.go`
- **Recommended fix:** Remove `taskLabel`. Remove `CheckResult` type. Remove or use the `exitCode` parameter in `GateSuggestion`.
- **Estimated effort:** Small

#### M2. Hardcoded state constants in `checks.go` duplicate `data/` definitions

- **Finding:** `validStates` map in `checks.go` duplicates the canonical state names defined in `data/lifecycle.go`. If a state is added to `lifecycle.go`, `checks.go` must be manually updated.
- **Evidence:** `internal/doctor/checks.go` — `validStates` map vs `data.IsCanonicalColumn()` / `data.IsCanonicalStage()`
- **Recommended fix:** Remove `validStates` and use `data.IsCanonicalColumn()` / `data.IsCanonicalStage()` for validation.
- **Estimated effort:** Small

#### M3. Inconsistent directory abstraction: `CheckOrphans` bypasses `data.Discover`

- **Finding:** Most of `checks.go` uses `data.Discover` for filesystem traversal, but `CheckOrphans` uses `os.ReadDir` directly.
- **Evidence:** `internal/doctor/checks.go` — `CheckOrphans` function
- **Recommended fix:** Add a `ListRootDirs()` or similar method to `data.Discover` and use it in `CheckOrphans`.
- **Estimated effort:** Small

#### M4. `shortID()` and `WrapText()`/`SplitLongWord()` should be in a shared utilities file

- **Finding:** `shortID()` is defined in `card.go` but used in `board.go`, `view.go`, `update.go`, and `detail.go`. `WrapText()` and `SplitLongWord()` in `detail.go` are general-purpose text utilities.
- **Why it matters:** Scattered utility functions make it hard to find shared logic and risk duplication.
- **Recommended fix:** Create `internal/board/util.go` (or `internal/text/util.go`) for `shortID`, `WrapText`, `SplitLongWord`, and `truncate`. Move `epicIndex`/`releaseIndex` to a generic `indexOrZero` helper.
- **Estimated effort:** Small

#### M5. Layout constants split between `layout.go` and `column.go`

- **Finding:** `colOverhead` is defined in `column.go` (value 4) and used in both `column.go` and `layout.go`. `minColWidth` is in `layout.go`. These layout-related constants should be co-located.
- **Evidence:** `internal/board/column.go:17` (`const colOverhead = 4`), `internal/board/layout.go:8-15`
- **Recommended fix:** Move all layout constants to `layout.go`.
- **Estimated effort:** Small

#### M6. `AtomicWrite` cross-device rename fallback doesn't actually solve cross-device moves

- **Finding:** `replaceFile()` in `init/write.go` tries `os.Rename` first, then falls back to creating a backup and renaming, but `os.Rename` will still fail on cross-device moves. The fallback should use `os.Link` + `os.Remove` or a plain copy.
- **Evidence:** `internal/init/write.go:52-63`
- **Recommended fix:** In the fallback path, use `os.Open` + `io.Copy` + `os.Remove` instead of `os.Rename` for cross-filesystem cases.
- **Estimated effort:** Small

#### M7. Ad-hoc Markdown parsing in `epic_panel.go` is fragile

- **Finding:** `epicDetailBody()` and `epicAuditBody()` contain hand-rolled markdown section extraction that skips headings containing "component" or "files" via substring matching. Any heading with those substrings will be silently hidden.
- **Evidence:** `internal/board/epic_panel.go` — text manipulation functions
- **Recommended fix:** Either use a proper markdown parser (e.g., `goldmark`) for section extraction, or use exact heading matches with a configurable allowlist/blocklist rather than substring matching.
- **Estimated effort:** Medium (if switching to a markdown parser) / Small (if switching substring to exact match)

#### M8. `Config.Theme` defaults not filling individual accent colors

- **Finding:** `fillThemeDefaults()` fills base theme colors when empty, but for accents, it's all-or-nothing: if the user specifies any accent, they must specify all accents, since the function only uses defaults when `len(theme.Accents) == 0`.
- **Evidence:** `internal/data/config.go:91-94`
- **Recommended fix:** Fill missing accent keys individually from `defaultTheme.Accents` instead of replacing the entire map.
- **Estimated effort:** Small

#### M9. No Windows build target in `buildtool`

- **Finding:** The `targets` list only includes linux/amd64, linux/arm64, darwin/amd64, darwin/arm64. No Windows target despite the project running on Windows (there's a `localExecutable()` Windows branch).
- **Evidence:** `internal/buildtool/main.go:21-26`
- **Recommended fix:** Add `{os: "windows", arch: "amd64"}` and `{os: "windows", arch: "arm64"}` to the targets list. Add `.exe` suffix handling in `writeTarGz`.
- **Estimated effort:** Small

### Low

#### L1. Audit section allowlist in `epic_panel.go` should be configurable or data-driven

- **Finding:** `allowedSections` is a hard-coded map that determines which audit sections render in the TUI.
- **Recommended fix:** Either extract to a constant with a comment explaining why these sections are allowed, or make it configurable via config.
- **Estimated effort:** Small

#### L2. No distribution checksums

- **Finding:** `dist()` creates tar.gz archives but no SHA256 checksums file.
- **Recommended fix:** Generate a `checksums.txt` file during `dist` with SHA256 hashes of each archive.
- **Estimated effort:** Small

#### L3. `GateSuggestion` has unused `exitCode` parameter

- **Finding:** The second parameter is never referenced in the function body.
- **Evidence:** `internal/doctor/repairs.go:67`
- **Recommended fix:** Remove the parameter or use it to provide exit-code-specific suggestions.
- **Estimated effort:** Small

#### L4. `agent_skills_test.go` hardcodes expected skill count of 6

- **Finding:** The test asserts exactly 6 skills exist. Adding or removing a skill requires updating this test.
- **Recommended fix:** Remove the count assertion and just verify each skill has valid frontmatter, or derive the expected count from the directory listing.
- **Estimated effort:** Small

#### L5. Test helper duplication across packages

- **Finding:** File-writing helpers (`writeFile`, `writeTaskFixture`, etc.) are reimplemented in `board_test.go`, `checks_test.go`, and `write_test.go`.
- **Recommended fix:** Create an `internal/testutil` package with shared fixtures.
- **Estimated effort:** Medium

#### L6. `splitChecklistSentences` doesn't handle abbreviations

- **Finding:** Text splitting on `.`, `!`, `?` doesn't account for abbreviations like "e.g." or "i.e.", which could cause incorrect sentence breaks.
- **Recommended fix:** Use a more careful sentence boundary detector, or at minimum skip periods preceded by known abbreviations.
- **Estimated effort:** Small

#### L7. `Package main` test file at root level

- **Finding:** `agent_skills_test.go` is in `package main` but tests structural properties of the `agent-skills/` directory, not main package functionality. This makes it invisible to per-package test runs.
- **Recommended fix:** Move it to an appropriate subdirectory or create a `build/` or `meta/` test package.
- **Estimated effort:** Small

---

## 3. Complexity & Modularity Review

### Overly large files

| File | Lines | Concern |
|------|-------|---------|
| `internal/board/update.go` | 521 | Monolithic switch handling all input, overlays, and I/O dispatch |
| `internal/doctor/checks.go` | 585 | "Check everything" monolith; each check function is OK but they share no common validation abstractions |
| `internal/data/write.go` | 216 | Two write paths (`WriteTaskStatus` and `WriteRouterState`) share YAML manipulation boilerplate |
| `internal/board/epic_panel.go` | 256 | Sidebar, detail, audit, and dropdown rendering mixed |

### Tight coupling

- `internal/board/model.go` is a god struct (25+ fields) that every file in the board package mutates. This is idiomatic for Bubble Tea but creates a wide coupling surface.
- `update.go` directly calls I/O functions that modify the filesystem. The TUI update loop should not know about filesystem paths.

### Mixed responsibilities

- `internal/board/board.go`'s `newProjectModel()` does discovery, file loading, watcher creation, and router reading. It's an initialization bottleneck mixing concerns.
- `model.go` contains both UI state and filesystem mutation functions (`writeRouterReleaseEpic`, `writeRouterTask`). Disk writes should be `tea.Cmd`s.

### Repeated logic

- `\r\n` normalization appears in 4+ call sites across 2 packages
- `shortID()` is defined in `card.go` but consumed by 6+ files
- `epicIndex` and `releaseIndex` are structurally identical
- YAML frontmatter extraction + reconstruction appears in `write.go` twice with significant duplication between `updateFrontmatterField()` and `WriteTaskStatus()`

### Unclear data flow

- `syncTaskStatus` in `transitions.go` sets `Status = string(Column)`, creating a dual-field synchronization issue. The `Task` struct has both `Status` and `Column`, and it's unclear which is authoritative.
- The `watch.go` reload silently swallows errors (returns `nil` on error).

### Excessive abstraction (minor)

- `data.NewParser()`, `data.NewDiscover()`, `data.NewConfigReader()`, `data.NewRouterReader()` all return pointer-to-empty-struct. These could be package-level functions instead of methods on empty structs, removing unnecessary allocation and indirection.

---

## 4. Architecture Review

### Folder organisation

```
.
├── main.go                  # CLI entrypoint, --version, command dispatch
├── cmd/                     # CLI arg parsing (3 commands: init, board, doctor)
├── internal/
│   ├── board/               # TUI (17 source files, ~2300 lines)
│   ├── data/                # Models, parsing, discovery, writing (8 files)
│   ├── doctor/              # Diagnostics, checks, gates, repairs (4 files)
│   ├── init/                # Scaffolding, validation, clipboard (7 files)
│   ├── buildtool/           # Build helper (1 file, 219 lines)
│   └── styles/              # Palette + lipgloss styles (2 files)
├── templates/               # Embedded project templates
└── agent-skills/            # Skill guides for AI agents
```

This is well-organised. The `cmd/` → `internal/` boundary is clean: `cmd/` handles arg parsing and delegates to `internal/` via dependency-injected function types. The `internal/` packages have clear responsibilities.

### Domain boundaries

- **Data layer** (`internal/data/`): Task lifecycle, parsing, discovery, config, writing. This is the right boundary.
- **Board/TUI** (`internal/board/`): Rendering, interaction, file watching. This package is the largest and could benefit from sub-packages, but Go's convention is to keep it flat.
- **Doctor** (`internal/doctor/`): Checks, gates, repairs, report. Each file has one job. Clean.
- **Init** (`internal/init/`): Scaffold, validate, write, clipboard, prompt. Clean.

### State management

State management follows Bubble Tea's Elm Architecture: `Model` holds state, `Update()` processes messages and returns commands, `View()` renders. This is solid. The issue is I/O leaking into `Update()` rather than being emitted as `tea.Cmd`.

### API/data layer structure

The data layer has no interfaces — all types are concrete structs. For this project size, this is appropriate. The `NewParser()` / `NewDiscover()` pattern returns `*Parser` / `*Discover` which are empty structs — these could just be package-level functions.

### Configuration approach

Config is read from `.savepoint/config.yml` with sensible defaults. The `QualityGates` struct allows project-specific lint/test/typecheck commands. Theme customization supports three color tiers (truecolor, ANSI256, ANSI16). This is well-designed.

### Error handling strategy

- **Data layer:** Uses sentinel errors (`ErrNoFrontmatter`, `ErrSavepointDirectoryMissing`, `ErrMtimeConflict`, `ErrProposalNotFound`) — good.
- **Board:** Errors from I/O inside `Update()` are displayed as status messages to the user — acceptable for a TUI.
- **Doctor:** Errors are collected into a structured `DiagnosticReport` with repair suggestions — well-designed.
- **Init:** Errors propagate up to `main.go` which prints to stderr and exits — simple and correct.

One gap: `watch.go`'s `reloadTasks` silently returns `nil` on filesystem errors, meaning a watcher failure silently blanks the board.

---

## 5. Best-Practice Review

### Framework conventions

- Bubble Tea conventions are generally followed: `Init()`, `Update()`, `View()` pattern. However, I/O in `Update()` violates the framework's core principle that `Update()` should be pure. This is the biggest convention deviation.

### Type safety

- `ColumnType`, `ProgressStage`, and `TaskStatus` are all `string` types. `TaskStatus` is defined but `Task.Column` uses `ColumnType` while `Task.Status` is a plain `string`. The `Status` field is redundant with `Column` and creates a dual-source-of-truth risk.
- `RouterState` uses `string` fields for `State`, `Release`, `Epic`, `Task` — these could benefit from typed constants similar to `ColumnType`.

### Linting/formatting

- No linter configuration found (no `.golangci.yml`, no `golint` config). Assuming `go fmt` and `go vet` are used via `make test`.
- **Recommended:** Add `golangci-lint` with at least `errcheck`, `gosimple`, `staticcheck`, `unused`, and `ineffassign` linters.

### Testing approach

- 264 tests across 39 test files — good coverage for core logic.
- Missing: `internal/buildtool/` (0 tests), `internal/styles/` (0 tests).
- Test helpers are duplicated across packages (write fixtures, temp directories).
- No benchmarks, no fuzz tests for the YAML parser.
- No `testdata/` directory — all fixtures are constructed inline.

### Dependency management

- Dependencies are minimal and appropriate: `bubbletea`, `lipgloss`, `fsnotify`, `yaml.v3`, `clipboard`.
- No dependency is used for a single call that could be replaced by stdlib (except clipboard, which is inherently platform-specific).
- `gopkg.in/yaml.v3` is the only YAML parser — used consistently.

### Build/deployment setup

- `Makefile` + `buildtool/main.go` provides `build`, `test`, `clean`, `dist`, `smoke-test`.
- The `go:embed` directive properly bundles templates.
- Missing: Windows cross-compilation, checksum generation, version injection via git tags (present but falls back to `v0.0.0`).

### Environment variable handling

- `buildtool/version()` checks `VERSION` env var — appropriate.
- Clipboard detection uses `runtime.GOOS` — correct.

### Logging/debugging practices

- No structured logging. Errors are printed to stderr in `main.go`.
- The TUI shows error messages as status bars. This is acceptable for a TUI app.
- **Gap:** No verbose/debug flag for TUI troubleshooting. Adding `--debug` or `SAVEPOINT_DEBUG` env var would help.

---

## 6. Refactor Roadmap

### Phase 1 — Safe cleanup

**Objective:** Remove dead code, fix easy bugs, reduce cognitive load without behavior changes.

| Task | Risk |
|------|------|
| Remove `taskLabel()` from `column.go` (dead code) | Very Low |
| Remove `CheckResult` type from `report.go` (dead code) | Very Low |
| Remove unused `exitCode` parameter from `GateSuggestion` | Very Low |
| Replace `contains()`/`indexOf()` with `strings.Contains`/`strings.Index` in `repairs.go` | Very Low |
| Replace `trimSpace()` with `bytes.TrimSpace` in `buildtool/main.go` | Very Low |
| Replace `validStates` map with `data.IsCanonicalColumn()`/`data.IsCanonicalStage()` | Low |
| Co-locate layout constants (`colOverhead`) into `layout.go` | Low |
| Add `\r\n` normalization helper to `data/` package | Low |
| Centralize duplicated test helpers into `internal/testutil/` | Low |
| Fix cycle detection path reconstruction in `checks.go` | Low |

**Expected benefit:** Cleaner codebase, fewer surprises, eliminates one actual bug (cycle detection).

### Phase 2 — Structural improvements

**Objective:** Decompose the update monolith, improve data flow, extract shared utilities.

| Task | Risk |
|------|------|
| Extract overlay handlers from `update.go` into per-overlay functions | Medium |
| Extract `handleBoardKey()` from `update.go` | Medium |
| Move `shortID`, `WrapText`, `SplitLongWord`, `truncate` to a utility file | Medium |
| Generalize `epicIndex`/`releaseIndex` into `indexOrZero` | Low |
| Make `data.NewParser/Discover/ConfigReader/RouterReader` package-level functions | Low |
| Fix `AtomicWrite` cross-device fallback to use `io.Copy` | Medium |
| Make `Config.Theme.Accents` fill missing keys individually | Low |
| Add Windows build targets to `buildtool` | Low |
| Add checksum generation to `dist` | Low |

**Expected benefit:** Reduced file complexity, better separation of concerns, cross-platform correctness.

### Phase 3 — Hardening

**Objective:** Fix the I/O-in-update anti-pattern, add missing tests, improve security.

| Task | Risk |
|------|------|
| Refactor all filesystem I/O in `update.go` to `tea.Cmd` functions | Medium-High |
| Add `tea.Cmd`-based router writing for priority key and epic selection | Medium |
| Add timeout to quality gate execution (`exec.CommandContext`) | Low |
| Convert `SuggestRepair` to typed error matching | Medium |
| Add tests for `buildtool/` | Low |
| Add tests for `styles/` | Low |
| Add benchmark tests for render functions (`RenderColumn`, `RenderCard`, `RenderDetail`) | Low |
| Add fuzz targets for YAML frontmatter parsing | Low |
| Fix `watch.go` error handling (don't silently swallow errors) | Medium |
| Add `--debug`/`SAVEPOINT_DEBUG` flag for TUI troubleshooting | Low |
| Improve `splitChecklistSentences` for abbreviation handling | Low |

**Expected benefit:** TUI responsiveness, correct error propagation, test coverage completeness, security.

---

## 7. Top 10 Action List

- [ ] **Extract I/O from `update.go` into `tea.Cmd` functions** — Critical — `internal/board/update.go`, `internal/board/model.go` — Prevents TUI freezes on slow disk; follows Bubble Tea conventions; highest architectural impact
- [ ] **Fix cycle detection path reconstruction in `checks.go`** — High (bug) — `internal/doctor/checks.go` — Produces inaccurate error messages today
- [ ] **Decompose `update.go` overlay handling into per-type handlers** — High — `internal/board/update.go` — Reduces 521-line file complexity; makes adding new overlays safe
- [ ] **Replace stdlib reimplmentations (`contains`, `indexOf`, `trimSpace`)** — High — `internal/doctor/repairs.go`, `internal/buildtool/main.go` — Eliminates confusing custom code
- [ ] **Add timeout to quality gate execution** — High — `internal/doctor/gates.go` — Prevents indefinite blocking
- [ ] **Centralize `\r\n` normalization into a shared function** — Medium — `internal/data/parser.go`, `internal/data/write.go`, `internal/board/epic_panel.go` — Single source of truth for cross-platform line endings
- [ ] **Consolidate layout constants and shared utilities** — Medium — `internal/board/column.go`, `internal/board/layout.go`, `internal/board/card.go`, `internal/board/detail.go` — Reduces scatter and makes layout logic findable
- [ ] **Fix `AtomicWrite` cross-device rename fallback** — Medium — `internal/init/write.go` — Prevents silent data loss on cross-filesystem moves
- [ ] **Remove dead code (`taskLabel`, `CheckResult`, unused params)** — Low — `internal/board/column.go`, `internal/doctor/report.go`, `internal/doctor/repairs.go` — Reduces clutter
- [ ] **Add tests for `buildtool/` and `styles/`** — Medium — `internal/buildtool/`, `internal/styles/` — Achieves baseline coverage for untested packages