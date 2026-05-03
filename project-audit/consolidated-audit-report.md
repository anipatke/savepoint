# Consolidated Audit Report — Savepoint

This report merges findings from two independent audits (Opus 4.6 and GLM 5.1) of the Savepoint codebase. Where both audits identified the same issue, the finding is merged with combined evidence. Where only one audit identified an issue, it is included with attribution. Severity disagreements are resolved in favour of the higher rating, per the principle that the more cautious assessment should prevail.

---

## 1. Executive Summary

**Savepoint** is a Go CLI/TUI tool (~5,600 lines of production code, ~5,600 lines of tests, 264 tests across 39 test files) that implements a file-based project state machine with a kanban-style terminal board. It uses the Charmbracelet stack (Bubble Tea, Lip Gloss) for TUI rendering, fsnotify for file watching, and YAML-frontmatter markdown files as its data layer.

### What is working well

- **File-per-responsibility** is followed diligently. Nearly every `.go` file does one job. The `board/` package splits cleanly into `model.go`, `view.go`, `update.go`, `card.go`, `column.go`, etc.
- **Test coverage is strong.** 264 tests across 39 files — roughly 1.6× the production code. All packages except `buildtool` and `styles` have tests, and integration tests exist for `board` and `init`.
- **Clean dependency tree.** Only 2 direct dependencies (`bubbletea`, `fsnotify`) plus the Charmbracelet ecosystem for rendering. No framework bloat.
- **All tests pass.** `go test ./...` reports zero failures.
- **Data model is honest.** The `data` package cleanly separates parsing, lifecycle validation, writing, and discovery.
- **Dependency injection in `cmd/`** via function types (`InitRunner`, `BoardRunner`, `DoctorRunner`) makes commands trivially testable.
- **Atomic file writes** in `data/write.go` and `init/write.go` prevent corruption.
- **Policy tests** (`render_policy_test.go`) enforce cross-cutting visual constraints.
- **Responsive TUI layout** with clean breakpoint-based column rendering.

### Biggest risks

1. **Synchronous file I/O inside the Bubble Tea `Update()` loop** freezes the TUI on slow disks. This is the highest-impact architectural issue.
2. **`update.go` is a 521-line monolith** with a deeply nested key-dispatch switch. Hard to extend or test in isolation.
3. **Duplicated YAML frontmatter read/write/parse logic** appears in `write.go`, `parser.go`, `board.go`, `checks.go`, and `epic_panel.go`.
4. **Cycle detection in `checks.go` produces inaccurate paths** — a correctness bug.
5. **No interfaces used for I/O boundaries** — all data-access types are concrete structs, making them impossible to mock without disk fixtures.

### Extensibility

The project is easy to extend for new checks, overlays, and commands. It is harder to extend for new data sources or rendering backends because I/O is baked into concrete functions. The main risk to extensibility is the `update.go` monolith.

### Architecture fit

The architecture (flat `internal/` packages, Elm-like TUI model, embedded templates) is well-suited for a small-to-medium CLI tool. No over-engineering is evident. The main architectural debts are the I/O-in-update anti-pattern and the `update.go` size.

---

## 2. Severity-Ranked Recommendations

### Critical

#### C1. Synchronous file I/O in the TUI update loop

- **Finding:** `update.go` performs filesystem reads and writes directly inside `Update()`: `writeTaskStatus()`, `writeRouterTask()`, `writeRouterReleaseEpic()`, `readEpicDetailFile()`, `selectEpicPanelEpic()`. These block the TUI event loop.
- **Why it matters:** Any disk latency (network drives, slow SSDs, virus scanners) freezes the entire TUI. Bubble Tea's design intent is for `Update()` to be pure — I/O should happen in `tea.Cmd` functions that return messages.
- **Evidence:** `internal/board/update.go` (writeRouterTask, writeRouterReleaseEpic, readEpicDetailFile), `internal/board/model.go:235-280` (writeRouterReleaseEpic, writeRouterTask)
- **Recommended fix:** Extract I/O operations into `tea.Cmd` functions. E.g., `writeTaskStatusCmd(task, path, mtime) tea.Cmd` returns a `tea.Msg` on completion. `Update()` dispatches the command and handles the result message. This is the standard Bubble Tea pattern.
- **Estimated effort:** Medium

#### C2. Cycle detection produces inaccurate paths

- **Finding:** `detectCycles` in `checks.go` uses a `parent` map that gets overwritten when a node is visited from multiple paths. When a cycle is found, the path reconstructed via `parent` may not represent the actual cycle.
- **Why it matters:** Users could be shown a cycle path that doesn't actually exist, causing confusion or incorrect doctor reports.
- **Evidence:** `internal/doctor/checks.go` — the DFS `parent` map is a simple `map[string]string` that gets overwritten per-visit
- **Recommended fix:** Use a stack-based cycle reconstruction (track the current DFS path as a slice) or validate the reconstructed path actually forms a cycle before reporting it.
- **Estimated effort:** Small

---

### High

#### H1. Committed binaries in the repository

- **Finding:** `savepoint` (5.5 MB), `savepoint.exe` (6.0 MB), `dist/`, and `ink-cli-ui-design.zip` are tracked in Git.
- **Why it matters:** Bloats clone size (12+ MB), causes merge noise, and risks accidentally shipping stale binaries.
- **Evidence:** `.gitignore` does not exclude root binaries. `dist/` is also checked in.
- **Recommended fix:** Add `savepoint`, `savepoint.exe`, `dist/`, and `*.zip` to `.gitignore`. Run `git rm --cached` on tracked binaries. Build in CI only.
- **Estimated effort:** Small

#### H2. `update.go` complexity: God-method `Update()`

- **Finding:** `update.go` is 521 lines. The `Update()` method alone spans ~190 lines with 4+ levels of nesting inside `case tea.KeyMsg`. Overlay handling mixes 5 overlay types in one `updateOverlay()` function.
- **Why it matters:** Adding a new keybinding or overlay requires editing a deeply nested switch. Bug surface area grows with each addition. The space-bar and backspace handlers have nearly identical structure — find-task-by-ID, mutate, write, refresh.
- **Evidence:** `internal/board/update.go` — lines 19–192 for `Update()`, lines 132–158 and 159–181 for near-duplicate handler structure
- **Recommended fix:** Extract key handlers into named methods: `handleAdvanceTask()`, `handleRetreatTask()`, `handleSetPriority()`. Extract `updateBoardKeys()` and `updateOverlayKeys()` from the top-level switch. Split the overlay update into per-type handlers.
- **Estimated effort:** Medium

#### H3. Duplicated frontmatter body-extraction logic

- **Finding:** The pattern "extract frontmatter → unmarshal YAML → compute body start offset → reconstruct file" appears in:
  - `write.go:updateFrontmatterField` (lines 40–83)
  - `write.go:WriteTaskStatus` (lines 85–150)
  - `board.go` (router reading)
  - `checks.go` (task validation)
  - `epic_panel.go:epicDetailBody` and `epicAuditBody` (frontmatter stripping, lines 50–60 and 125–134)
  - The magic `delimLen := 4; bodyStart := delimLen + len(raw) + delimLen` appears in two places
- **Why it matters:** A change to frontmatter format must be patched in 4+ places. The body-offset calculation is fragile.
- **Evidence:** `internal/data/write.go:74-79` and `internal/data/write.go:140-145` are identical; `internal/board/epic_panel.go:51-59` and `internal/board/epic_panel.go:126-133` are identical
- **Recommended fix:** Extract `SplitFrontmatterBody(content string) (yaml string, body string, err error)` in the `data` package. Extract `stripFrontmatter(content string) string` for `epic_panel.go`.
- **Estimated effort:** Small

#### H4. No interfaces for data-access types

- **Finding:** `Discover`, `Parser`, `ConfigReader`, `RouterReader` are all concrete structs. Every consumer calls `data.NewDiscover()`, `data.NewParser()`, etc. directly.
- **Why it matters:** Test helpers in `board` and `doctor` must create real filesystem fixtures to test business logic. This is expensive and brittle.
- **Evidence:** `internal/board/board.go` lines 37, 85, 130 all call `data.NewDiscover()` with no injection point. `internal/doctor/checks.go` lines 116, 293, 454, 525 do the same.
- **Recommended fix:** Define interfaces at the consumer side (e.g., `type taskDiscoverer interface { ListReleases(root string) ([]data.ReleaseInfo, error) ... }`). Keep the existing structs as production implementations. Accept the interface in board/doctor constructors.
- **Estimated effort:** Medium

#### H5. Stdlib reimplementation in `repairs.go` and `buildtool/main.go`

- **Finding:** `contains()` and `indexOf()` in `repairs.go` reimplement `strings.Contains` and `strings.Index`. `trimSpace()` in `buildtool/main.go` reimplements `bytes.TrimSpace`.
- **Why it matters:** Makes code harder to read for Go developers expecting standard library calls. The custom implementations may have subtle differences from the stdlib versions (e.g., Unicode whitespace handling).
- **Evidence:** `internal/doctor/repairs.go:50-61`, `internal/buildtool/main.go:211-219`
- **Recommended fix:** Replace with `strings.Contains`, `strings.Index`, and `strings.TrimSpace` respectively.
- **Estimated effort:** Small

#### H6. Fragile repair suggestion matching via substring search on error messages

- **Finding:** `SuggestRepair()` pattern-matches against error message substrings to suggest fixes. If error messages change format, repairs silently break.
- **Evidence:** `internal/doctor/repairs.go:8-66` — hard-coded substring matching against messages like `"not found"`, `"invalid"`, `"missing"`
- **Recommended fix:** Define typed error codes or sentinel errors in `data/` and `doctor/`, and match on error type rather than substring.
- **Estimated effort:** Medium

#### H7. Quality gate command execution has no timeout

- **Finding:** `RunQualityGates` executes commands from `config.yml` with no timeout. A hung command blocks the doctor indefinitely.
- **Evidence:** `internal/doctor/gates.go:32-55` — `exec.Command` with `Run()` and no `Context` timeout
- **Recommended fix:** Use `exec.CommandContext(ctx, ...)` with a configurable timeout (default: 60s). Add a `gate_timeout` config option.
- **Estimated effort:** Small

#### H8. `\r\n` normalization is scattered across files

- **Finding:** `strings.ReplaceAll(content, "\r\n", "\n")` appears in `parser.go`, `write.go` (3 times), `epic_panel.go`, and likely elsewhere.
- **Why it matters:** If the normalization logic changes (e.g., also handling `\r` alone), every call site must be found and updated. Missing a call site causes subtle cross-platform bugs.
- **Evidence:** `internal/data/parser.go:92`, `internal/data/write.go:22,46,104`, `internal/board/epic_panel.go`
- **Recommended fix:** Create a `normalizeLineEndings(s string) string` function in `internal/data/` and use it everywhere.
- **Estimated effort:** Small

---

### Medium

#### M1. `Discover`, `Parser`, `ConfigReader` are stateless singletons instantiated repeatedly

- **Finding:** `data.NewDiscover()`, `data.NewParser()`, `data.NewConfigReader()`, `data.NewRouterReader()` all return pointers to zero-value structs. `Discover` is re-created 7 times across `board.go`, `checks.go`, and `main.go`.
- **Recommended fix:** Either convert to package-level functions (since there's no state) or introduce consumer-side interfaces per H4.
- **Estimated effort:** Small

#### M2. `newProgramModel()` contains hardcoded epic slug

- **Finding:** `board.go:33` — `NewModel(nil, "v1", "E03-board-tui-core")` is a development leftover.
- **Recommended fix:** Delete `newProgramModel()` or replace hardcoded values with empty strings.
- **Estimated effort:** Small

#### M3. Hardcoded state constants in `checks.go` duplicate `data/` definitions

- **Finding:** `validStates` map in `checks.go` duplicates state names defined in `data/lifecycle.go`.
- **Recommended fix:** Remove `validStates` and use `data.IsCanonicalColumn()` / `data.IsCanonicalStage()`.
- **Estimated effort:** Small

#### M4. Inconsistent directory abstraction: `CheckOrphans` bypasses `data.Discover`

- **Finding:** Most of `checks.go` uses `data.Discover` for filesystem traversal, but `CheckOrphans` uses `os.ReadDir` directly.
- **Recommended fix:** Add a `ListRootDirs()` method to `data.Discover` and use it in `CheckOrphans`.
- **Estimated effort:** Small

#### M5. Layout constants split between `layout.go` and `column.go`

- **Finding:** `colOverhead` is defined in `column.go` (value 4) and used in both `column.go` and `layout.go`.
- **Recommended fix:** Move all layout constants to `layout.go`.
- **Estimated effort:** Small

#### M6. `AtomicWrite` cross-device rename fallback is broken

- **Finding:** `replaceFile()` tries `os.Rename` first, then falls back to creating a backup and renaming, but `os.Rename` will still fail on cross-device moves in the fallback path.
- **Evidence:** `internal/init/write.go:52-63`
- **Recommended fix:** Use `os.Open` + `io.Copy` + `os.Remove` for cross-filesystem fallback.
- **Estimated effort:** Small

#### M7. Ad-hoc Markdown parsing in `epic_panel.go` is fragile

- **Finding:** `epicDetailBody()` and `epicAuditBody()` skip headings containing "component" or "files" via substring matching. Any heading with those substrings will be silently hidden.
- **Recommended fix:** Use exact heading matches with a configurable allowlist/blocklist rather than substring matching.
- **Estimated effort:** Small (exact match) / Medium (markdown parser)

#### M8. `Config.Theme` defaults not filling individual accent colors

- **Finding:** `fillThemeDefaults()` fills base theme colors when empty, but for accents it's all-or-nothing: `len(theme.Accents) == 0` triggers the default.
- **Recommended fix:** Fill missing accent keys individually from `defaultTheme.Accents`.
- **Estimated effort:** Small

#### M9. `splitCommand` in `gates.go` is a naïve shell tokenizer

- **Finding:** Only handles double-quote grouping — no escaping, no single quotes, no backslash-escapes.
- **Recommended fix:** Document the limitation or use `shellwords` parsing.
- **Estimated effort:** Small

#### M10. No Windows build target in `buildtool`

- **Finding:** `targets` list only includes Linux and Darwin. No Windows target despite a `localExecutable()` Windows branch.
- **Recommended fix:** Add Windows amd64 and arm64 targets. Add `.exe` suffix handling.
- **Estimated effort:** Small

#### M11. `buildtool` has no tests

- **Finding:** Only production package with `[no test files]`.
- **Recommended fix:** Add tests for `run()`, `version()`, and `writeTarGz()`.

---

### Low

#### L1. `ColumnType` and `TaskStatus` are parallel enumerations for the same concept

- **Recommended fix:** Consider unifying into a single status type. Low priority.
- **Estimated effort:** Medium

#### L2. `package.json` test script is misleading

- **Finding:** `"test": "savepoint init"` — `npm test` scaffolds a project instead of running tests.
- **Recommended fix:** Change to `"test": "echo \"Run 'make test' for Go tests\""`.

#### L3. Dead code: `taskLabel()`, `loadAllTasks()`, `newProgramModel()`, `CheckResult`

- **Recommended fix:** Delete all four.

#### L4. `shortID` and `shortRouterID` are near-duplicates

- **Recommended fix:** Consolidate into one `ShortID(full string) string` function.

#### L5. `epicIndex` and `releaseIndex` are identical functions

- **Recommended fix:** Extract `sliceIndex(items []string, target string) int`.

#### L6. No linter configured

- **Recommended fix:** Add `.golangci.yml` with `unused`, `errcheck`, `staticcheck`, `govet`, `ineffassign`.

#### L7. No distribution checksums

- **Finding:** `dist()` creates tar.gz archives but no SHA256 checksums file.
- **Recommended fix:** Generate `checksums.txt` during `dist`.

#### L8. `agent_skills_test.go` hardcodes expected skill count of 6

- **Recommended fix:** Remove count assertion or derive from directory listing.

#### L9. Test helper duplication across packages

- **Recommended fix:** Create `internal/testutil` package with shared fixtures.

#### L10. `splitChecklistSentences` doesn't handle abbreviations

- **Recommended fix:** Skip periods preceded by known abbreviations (e.g., "e.g.", "i.e.").

#### L11. `package main` test file at root level

- **Recommended fix:** Move to a dedicated test package.

#### L12. Audit section allowlist should be configurable

- **Recommended fix:** Extract `allowedSections` to a named constant with documentation.

#### L13. `reloadTasks` silently swallows errors

- **Finding:** `watch.go` returns `nil` on error, causing the board to silently stop refreshing.
- **Recommended fix:** Return an `errorMsg` so the TUI can surface it.

---

## 3. Complexity & Modularity Review

### Overly large files

| File | Lines | Concern |
|------|-------|---------|
| `internal/board/update.go` | 521 | `Update()` is 190 lines; `updateOverlay()` is 100 lines |
| `internal/doctor/checks.go` | 585 | Single-file check aggregator; could be split by check type |
| `internal/data/write.go` | 216 | Two large functions with duplicated body-offset arithmetic |
| `internal/board/epic_panel.go` | 256 | Sidebar, detail, audit, and dropdown rendering mixed |

### Tight coupling

- **Board → data:** Appropriate for project size. Clean boundary.
- **Board → os:** `model.go` and `update.go` call `os.ReadFile`, `os.Stat`, `os.WriteFile` directly. Couples UI logic to filesystem.

### Repeated logic

1. Frontmatter stripping (3+ places)
2. `\r\n` normalization (4+ call sites across 2 packages)
3. `shortID` extraction (2 near-duplicate implementations)
4. `indexOf`/`epicIndex`/`releaseIndex` pattern (3 implementations)
5. Space-bar and Backspace handlers in `update.go` (nearly identical structure)
6. `Discover` instantiation (7 separate `NewDiscover()` calls)

### Unclear data flow

- **`Task.Status` vs `Task.Column`:** Two representations of the same state. `syncTaskStatus` manually keeps them in sync.
- **`watch.go` silent error swallowing:** `reloadTasks` returns `nil` on error.

### Excessive abstraction

- `Discover`, `Parser`, `ConfigReader`, `RouterReader` are empty structs used as method namespaces. Should be package-level functions or interfaces.

---

## 4. Architecture Review

### Folder organisation ✅

```
savepoint/
├── cmd/           # CLI arg parsing — clean separation from execution
├── internal/
│   ├── board/     # TUI model/view/update — Elm architecture
│   ├── buildtool/ # Standalone Go binary for build automation
│   ├── data/      # Models, parsing, writing, discovery
│   ├── doctor/    # Read-only diagnostics
│   ├── init/      # Scaffolding
│   └── styles/    # Centralised palette + styles
├── templates/     # Embedded scaffold templates
└── agent-skills/  # Prompt documents for AI agents
```

### Domain boundaries ✅

- **Data layer** (`internal/data/`): Task lifecycle, parsing, discovery, config, writing. Correct boundary.
- **Board/TUI** (`internal/board/`): Rendering, interaction, file watching. Largest package.
- **Doctor** (`internal/doctor/`): Checks, gates, repairs, report. Each file has one job. Clean.
- **Init** (`internal/init/`): Scaffold, validate, write, clipboard, prompt. Clean.

### State management

Bubble Tea `Model` with 27 fields. Upper bound of manageable. Consider grouping related fields into sub-structs if TUI grows. I/O-in-Update is the main concern (C1).

### Configuration approach ✅

`config.yml` with defaults baked into Go code. Three-tier color support. `QualityGates` for project-specific commands. Clean.

### Error handling

Generally good — errors wrapped with `fmt.Errorf("context: %w", err)`. Key exceptions:
- `reloadTasks` silently swallows errors (L13)
- `SuggestRepair` relies on substring matching (H6)
- Quality gates have no timeout (H7)

---

## 5. Best-Practice Review

### Framework conventions ⚠️

Bubble Tea conventions generally followed. I/O in `Update()` violates the framework's core principle — see C1.

### Type security ✅

Custom types (`ColumnType`, `ProgressStage`, `OverlayType`) used appropriately. `ColumnType`/`TaskStatus` duality is a minor concern (L1).

### Linting/formatting ⚠️

No linter configured. Add `.golangci.yml` with `unused`, `errcheck`, `staticcheck`, `govet`, `ineffassign`.

### Testing ✅

Strong coverage (264 tests, ~1.6:1 test-to-code ratio). Missing: `buildtool`, `styles`, benchmarks, fuzz tests.

### Dependency management ✅

`go.mod` is clean. 2 direct dependencies. All indirect deps from Charmbracelet ecosystem.

### Build/deployment ⚠️

- No CI configuration visible
- Binaries committed to repo (H1)
- No Windows build target (M10)
- No distribution checksums (L7)

### Logging/debugging ⚠️

No logging. TUI uses `StatusMessage` for user feedback. No `--debug` flag. Recommend adding `SAVEPOINT_DEBUG` env var.

---

## 6. Refactor Roadmap

### Phase 1 — Safe cleanup

**Objective:** Remove dead code, fix obvious duplication, improve hygiene.

| Task | Risk |
|------|------|
| Add binaries to `.gitignore` and remove from Git tracking | None |
| Delete `taskLabel()`, `loadAllTasks()`, `newProgramModel()`, `CheckResult` | None |
| Remove unused `exitCode` param from `GateSuggestion` | None |
| Replace `contains()`/`indexOf()` with `strings.Contains`/`strings.Index` | None |
| Replace `trimSpace()` with `strings.TrimSpace` | None |
| Replace `validStates` map with `data.IsCanonicalColumn()`/`data.IsCanonicalStage()` | Low |
| Fix `package.json` test script | None |
| Add `.golangci.yml` with basic linters | None |
| Co-locate layout constants (`colOverhead`) into `layout.go` | Low |
| Consolidate `shortID`/`shortRouterID` and `epicIndex`/`releaseIndex` | Low |

**Expected benefit:** Cleaner codebase, smaller repo, automated lint catches.
**Risk level:** Very low.

### Phase 2 — Structural improvements

**Objective:** Reduce duplication, improve modularity, make key files easier to extend.

| Task | Risk |
|------|------|
| Extract `SplitFrontmatterBody()` in `data` package | Low — well-tested |
| Extract `stripFrontmatter()` for `epic_panel.go` | Low |
| Add `normalizeLineEndings()` to `data` package; use everywhere | Low |
| Split `Update()` into `handleBoardKey()`, `handleOverlayKey()`, named methods | Medium |
| Group `Model` fields into sub-structs | Medium |
| Convert stateless data types to package-level functions or consumer-defined interfaces | Medium |
| Fix `AtomicWrite` cross-device fallback | Low |
| Make `Config.Theme.Accents` fill missing keys individually | Low |

**Expected benefit:** Easier to add features, reduced duplication, more testable.
**Risk level:** Medium.

### Phase 3 — Hardening

**Objective:** Fix the I/O-in-update anti-pattern, add missing tests, improve error handling.

| Task | Risk |
|------|------|
| Extract all filesystem I/O from `update.go` into `tea.Cmd` functions | Medium-High |
| Add `tea.Cmd`-based router writing for priority key and epic selection | Medium |
| Handle `reloadTasks` errors by emitting `errorMsg` | Low |
| Add timeout to quality gate execution | Low |
| Convert `SuggestRepair` to typed error matching | Medium |
| Add tests for `buildtool/` | Low |
| Add tests for `styles/` | Low |
| Add benchmark tests for render functions | Low |
| Add fuzz targets for YAML frontmatter parsing | Low |
| Add `--debug`/`SAVEPOINT_DEBUG` flag | Low |
| Fix cycle detection path reconstruction | Low |

**Expected benefit:** TUI responsiveness, correct error propagation, test coverage completeness.
**Risk level:** Low–Medium.

---

## 7. Top 10 Action List

- [ ] **1. Extract I/O from `update.go` into `tea.Cmd` functions** — Critical — `internal/board/update.go`, `internal/board/model.go` — Prevents TUI freezes on slow disk; follows Bubble Tea conventions; highest architectural impact
- [ ] **2. Remove committed binaries from the repository** — High — `.gitignore`, `savepoint`, `savepoint.exe`, `dist/`, `ink-cli-ui-design.zip` — Repo shrinks by 12 MB+, eliminates merge noise
- [ ] **3. Fix cycle detection path reconstruction in `checks.go`** — High (bug) — `internal/doctor/checks.go` — Produces inaccurate error messages today
- [ ] **4. Extract `SplitFrontmatterBody()` to deduplicate write logic** — High — `internal/data/write.go`, `internal/board/epic_panel.go` — Single source of truth for frontmatter reconstruction
- [ ] **5. Split `Update()` into named key-handler methods** — High — `internal/board/update.go` — 190-line method becomes 5–6 focused methods; easier to extend
- [ ] **6. Replace stdlib reimplementations (`contains`, `indexOf`, `trimSpace`)** — High — `internal/doctor/repairs.go`, `internal/buildtool/main.go` — Eliminates confusing custom code
- [ ] **7. Add timeout to quality gate execution** — High — `internal/doctor/gates.go` — Prevents indefinite blocking
- [ ] **8. Centralize `\r\n` normalization and frontmatter stripping** — Medium — `internal/data/parser.go`, `internal/data/write.go`, `internal/board/epic_panel.go` — Single source of truth for cross-platform line endings and body extraction
- [ ] **9. Consolidate layout constants, shared utilities, and duplicate functions** — Medium — `internal/board/column.go`, `internal/board/layout.go`, `internal/board/card.go`, `internal/board/view.go`, `internal/board/detail.go`, `internal/board/epic_panel.go`, `internal/board/release.go` — Reduces scatter and makes logic findable
- [ ] **10. Fix `AtomicWrite` cross-device rename fallback** — Medium — `internal/init/write.go` — Prevents silent data loss on cross-filesystem moves

---

*Consolidated from audits by Opus 4.6 and GLM 5.1 on 2026-05-03.*