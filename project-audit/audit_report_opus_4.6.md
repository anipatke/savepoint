# Codebase Audit Report — Savepoint

## 1. Executive Summary

**Savepoint** is a small, well-structured Go CLI + TUI tool (~133 KB of production code across 42 source files) that implements a file-based project state machine with a kanban-style terminal board. It uses the Charmbracelet stack (Bubble Tea, Lip Gloss) for TUI rendering and YAML frontmatter-based markdown files as its data layer.

### What is working well

- **File-per-responsibility is followed diligently.** Nearly every `.go` file does one job. The `board/` package is split into `model.go`, `view.go`, `update.go`, `card.go`, `column.go`, etc. This maps to the project's own code-style rules and makes individual files easy to reason about.
- **Test coverage is strong.** 38 test files totalling ~213 KB — roughly 1.6× the production code. All packages except `buildtool` and `styles` have tests, and integration tests exist for `board` and `init`.
- **Clean dependency tree.** Only 2 direct dependencies (`bubbletea`, `fsnotify`) plus the Charmbracelet ecosystem for rendering. No framework bloat.
- **All tests pass.** `go test ./...` reports zero failures.
- **Data model is honest.** The `data` package cleanly separates parsing, lifecycle validation, writing, and discovery. The frontmatter-based approach is simple and appropriate for the project's scale.
- **The doctor/diagnostics subsystem is thorough.** Checks config, router, structure, dependencies, orphans, audit state, and quality gates — with actionable repair suggestions.

### Biggest risks

1. **`update.go` is a 522-line monolith** with the `Update()` method containing a deeply nested key-dispatch switch. It is the hardest file to extend or test in isolation.
2. **Duplicated YAML-frontmatter read/write/parse logic** appears in `write.go`, `parser.go`, `board.go`, and `checks.go` — each creating `NewParser()` independently and re-implementing body extraction.
3. **No interfaces used for I/O boundaries.** The `Discover`, `Parser`, `ConfigReader`, `RouterReader` types are all concrete structs with no interfaces, making them impossible to mock without test fixtures on disk.
4. **Committed binaries** (`savepoint`, `savepoint.exe`, `dist/`) inflate the repository and cause spurious diffs.

### Extensibility

The project is easy to extend for new checks, overlays, and commands, but harder to extend for new data sources or rendering backends because I/O is baked into concrete functions.

### Architecture fit

The architecture (flat `internal/` packages, Elm-like TUI model, embedded templates) is well-suited for a small-to-medium CLI tool. No over-engineering is evident.

---

## 2. Severity-Ranked Recommendations

### Critical

No critical issues were found. The application builds, tests pass, and the data model is consistent.

---

### High

#### H1 — Committed binaries in the repository

- **Finding:** `savepoint` (5.5 MB), `savepoint.exe` (6.0 MB), and `ink-cli-ui-design.zip` (15 KB) are tracked in Git.
- **Why it matters:** Bloats clone size, causes merge noise, and risks accidentally shipping stale binaries.
- **Evidence:** [.gitignore](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/.gitignore) does not exclude the root binaries. `dist/` is also checked in.
- **Recommended fix:** Add `savepoint`, `savepoint.exe`, `dist/`, and `*.zip` to `.gitignore`. Run `git rm --cached savepoint savepoint.exe dist/ ink-cli-ui-design.zip`. Build in CI only.
- **Estimated effort:** Small

#### H2 — `update.go` complexity: God-method Update()

- **Finding:** [update.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/update.go) is 522 lines. The `Update()` method alone spans lines 19–192 with 4 levels of nesting inside `case tea.KeyMsg`.
- **Why it matters:** Adding a new keybinding or overlay requires editing a deeply nested switch. Bug surface area grows with each addition.
- **Evidence:** The space-bar handler (lines 132–158) and backspace handler (lines 159–181) have nearly identical structure — find-task-by-ID, mutate, write, refresh.
- **Recommended fix:** Extract key handlers into named methods: `handleAdvanceTask()`, `handleRetreatTask()`, `handleSetPriority()`. Extract `updateBoardKeys()` and `updateOverlayKeys()` from the top-level switch.
- **Estimated effort:** Medium

#### H3 — Duplicated frontmatter body-extraction logic

- **Finding:** The pattern "extract frontmatter → unmarshal YAML → compute body start offset → reconstruct file" appears in:
  - [write.go:updateFrontmatterField](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/write.go#L40-L83) (lines 40–83)
  - [write.go:WriteTaskStatus](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/write.go#L85-L150) (lines 85–150)
  - The magic `delimLen := 4; bodyStart := delimLen + len(raw) + delimLen` appears in both
- **Why it matters:** A change to frontmatter format (e.g. supporting `---` on the same line as metadata) must be patched in multiple places. The body-offset calculation is fragile.
- **Evidence:** Lines 74–79 and 140–145 of `write.go` are identical.
- **Recommended fix:** Extract a `SplitFrontmatterBody(content string) (yaml string, body string, err error)` function and use it in all write paths.
- **Estimated effort:** Small

#### H4 — No interfaces for data-access types

- **Finding:** `Discover`, `Parser`, `ConfigReader`, `RouterReader` are all concrete structs. Every consumer calls `data.NewDiscover()`, `data.NewParser()`, etc. directly.
- **Why it matters:** Test helpers in `board` and `doctor` must create real filesystem fixtures to test business logic. This is expensive and brittle.
- **Evidence:** [board.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/board.go) lines 37, 85, 130 all call `data.NewDiscover()` with no injection point. [checks.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/doctor/checks.go) does the same (lines 116, 293, 454, 525).
- **Recommended fix:** Define interfaces at the consumer side (e.g. `type taskDiscoverer interface { ListReleases(root string) ([]data.ReleaseInfo, error) ... }`). Keep the existing structs as the production implementations. Accept the interface in board/doctor constructors.
- **Estimated effort:** Medium

---

### Medium

#### M1 — `Discover` is instantiated repeatedly with no state

- **Finding:** `data.NewDiscover()` returns `&Discover{}` — a zero-value struct. It is re-created in `board.go` (twice), `checks.go` (4 times), `main.go` (once).
- **Why it matters:** Unnecessary allocation noise and a missed opportunity for future caching (e.g., directory memoization).
- **Evidence:** [discover.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/discover.go) — `Discover` has zero fields.
- **Recommended fix:** Either make `Discover` methods package-level functions (since there's no state) or make it a singleton / accept it as a dependency.
- **Estimated effort:** Small

#### M2 — `Parser` and `ConfigReader` are similarly stateless singletons

- **Finding:** Same pattern as M1. `NewParser()` returns `&Parser{}` — no configuration.
- **Why it matters:** Code reads as if these types might have configuration, but they don't. This is an accidental abstraction that adds indirection without benefit.
- **Evidence:** [parser.go:12-14](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/parser.go#L12-L14), [config.go:52-54](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/config.go#L52-L54)
- **Recommended fix:** Convert to package-level functions unless interfaces are introduced (per H4).
- **Estimated effort:** Small

#### M3 — `newProgramModel()` contains a hardcoded epic slug

- **Finding:** [board.go:33](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/board.go#L32-L34) — `NewModel(nil, "v1", "E03-board-tui-core")` is a leftover from development.
- **Why it matters:** If `newProgramModel()` is ever called without the project-loading path, it shows stale default data.
- **Evidence:** The function is declared but never called in production (the real path is `newProjectModel`). However it is exported-adjacent and could confuse contributors.
- **Recommended fix:** Delete `newProgramModel()` or replace the hardcoded values with empty strings.
- **Estimated effort:** Small

#### M4 — `repairs.go` reimplements `strings.Contains`

- **Finding:** [repairs.go:50-61](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/doctor/repairs.go#L50-L61) defines custom `contains()` and `indexOf()` functions that do exactly what `strings.Contains` does.
- **Why it matters:** Standard library duplication is a maintainability red flag and confuses readers.
- **Evidence:** `func contains(s, substr string) bool` vs `strings.Contains`.
- **Recommended fix:** Replace with `strings.Contains`.
- **Estimated effort:** Small

#### M5 — `buildtool/main.go` reimplements `strings.TrimSpace`

- **Finding:** [buildtool/main.go:211-219](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/buildtool/main.go#L211-L219) defines a custom `trimSpace()` function that trims `\n\r\t `.
- **Why it matters:** `strings.TrimSpace` exists in the standard library with identical behaviour.
- **Evidence:** The function is only called once, on line 199.
- **Recommended fix:** Replace with `strings.TrimSpace(string(output))`.
- **Estimated effort:** Small

#### M6 — `buildtool` has no tests

- **Finding:** `internal/buildtool` is the only production package with `[no test files]`.
- **Why it matters:** Build/dist logic is the kind of code most likely to break silently across platforms.
- **Evidence:** `go test ./...` output shows `? github.com/opencode/savepoint/internal/buildtool [no test files]`.
- **Recommended fix:** Add tests for `run()`, `version()`, `splitCommand()`, and `writeTarGz()` at minimum. The `run()` function is already well-structured for testing.
- **Estimated effort:** Medium

#### M7 — `splitCommand` in `gates.go` is a naïve shell tokeniser

- **Finding:** [gates.go:101-124](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/doctor/gates.go#L101-L124) handles only double-quote grouping — no escaping, no single quotes, no backslash-escapes.
- **Why it matters:** Users configuring quality gates with complex commands (e.g., `go test -run "Test Foo"`) may get unexpected tokenisation.
- **Evidence:** `case c == '"': inQuote = !inQuote` — no escape handling.
- **Recommended fix:** Document the limitation or use `shellwords` parsing. For this project size, documenting the limitation is sufficient.
- **Estimated effort:** Small

#### M8 — `epicDetailBody` and `epicAuditBody` duplicate frontmatter stripping

- **Finding:** Both [epicDetailBody](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/epic_panel.go#L43-L91) and [epicAuditBody](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/epic_panel.go#L118-L168) contain identical YAML frontmatter stripping logic (lines 50–60 and 125–134).
- **Why it matters:** Duplicated logic that must be changed in lockstep.
- **Evidence:** Compare lines 51–59 with 126–133 — identical block.
- **Recommended fix:** Extract a `stripFrontmatter(content string) []string` helper.
- **Estimated effort:** Small

---

### Low

#### L1 — `ColumnType` and `TaskStatus` are parallel enumerations for the same concept

- **Finding:** [task.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/task.go) defines both `ColumnType` (`planned`, `in_progress`, `done`) and `TaskStatus` (`planned`, `in_progress`, `done`, `audited`). The `Task` struct has both `Column ColumnType` and `Status string`.
- **Why it matters:** Two representations of the same state create confusion. `syncTaskStatus` in `transitions.go` manually keeps them in sync.
- **Evidence:** [transitions.go:57-59](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/transitions.go#L57-L59) — `t.Status = string(t.Column)`.
- **Recommended fix:** Consider unifying into a single `TaskStatus` type and deriving the column from it. Low priority since both are kept in sync.
- **Estimated effort:** Medium

#### L2 — `package.json` exists for npm distribution but `scripts.test` is misleading

- **Finding:** [package.json:24](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/package.json#L24) — `"test": "savepoint init"` — this is not a real test.
- **Why it matters:** Running `npm test` will scaffold a project in the current directory rather than running tests.
- **Recommended fix:** Change to `"test": "echo \"Run 'make test' for Go tests\""` or remove the scripts section.
- **Estimated effort:** Small

#### L3 — `shortID` and `shortRouterID` are near-duplicates

- **Finding:** [card.go:shortID](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/card.go#L119-L127) and [view.go:shortRouterID](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/view.go#L160-L169) extract the same pattern (last segment after `/`, prefix before `-`).
- **Why it matters:** Minor duplication — both are small helpers but do the same thing.
- **Recommended fix:** Consolidate into one exported `ShortID(full string) string` function.
- **Estimated effort:** Small

#### L4 — `epicIndex` and `releaseIndex` are identical functions

- **Finding:** [epic_panel.go:epicIndex](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/epic_panel.go#L249-L256) and [release.go:releaseIndex](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/release.go#L35-L42) have identical logic — find a string in a slice, return 0 if not found.
- **Why it matters:** Two functions doing the same thing.
- **Recommended fix:** Extract `sliceIndex(items []string, target string) int`.
- **Estimated effort:** Small

#### L5 — Unused function `taskLabel` in `column.go`

- **Finding:** [column.go:91-96](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/column.go#L91-L96) defines `taskLabel` but it is never called.
- **Why it matters:** Dead code.
- **Recommended fix:** Delete it.
- **Estimated effort:** Small

#### L6 — `loadAllTasks` in `board.go` wraps `loadBoardData` but discards 4 return values

- **Finding:** [board.go:125-128](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/board.go#L125-L128) — `func loadAllTasks(root string) ([]data.Task, error)` is defined but never called.
- **Why it matters:** Dead code.
- **Recommended fix:** Delete it.
- **Estimated effort:** Small

#### L7 — No linter configured

- **Finding:** No `.golangci.yml` or linter configuration in the repository.
- **Why it matters:** Lint could catch many of the issues in this audit automatically (dead code, duplicate functions, etc.).
- **Recommended fix:** Add a basic `.golangci.yml` with `unused`, `errcheck`, `staticcheck`, and `govet`.
- **Estimated effort:** Small

---

## 3. Complexity & Modularity Review

### Overly large files or functions

| File | Lines | Concern |
|------|-------|---------|
| [update.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/update.go) | 522 | `Update()` is 170 lines; `updateOverlay()` is 100 lines |
| [checks.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/doctor/checks.go) | 586 | Reasonable for a checks aggregator, but could be split by check type |
| [write.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/data/write.go) | 217 | Two large functions with duplicated body-offset arithmetic |

### Tight coupling

- **Board → data:** The board package imports `data` for models, parsing, and writing. This is appropriate for a small project.
- **Board → os:** `model.go` and `update.go` call `os.ReadFile`, `os.Stat`, and `os.WriteFile` directly. This couples business logic to the filesystem. Not a problem today but will make unit tests slower if the board grows.

### Repeated logic

1. Frontmatter stripping (3 places)
2. `shortID` extraction (2 places)
3. `indexOf` / `epicIndex` / `releaseIndex` pattern (3 places)
4. Space-bar and Backspace handlers in `update.go` (nearly identical structure)
5. `Discover` instantiation (7 places)

### Poor naming

None found. Naming is consistently clear and idiomatic Go.

### Unclear data flow

The `Task.Status` vs `Task.Column` duality is the main source of confusion. `Status` is the string from the YAML, `Column` is the normalised enum, and `syncTaskStatus` bridges them. A newcomer would need to trace this chain.

### Components doing too much

`update.go`'s `Update()` handles: quit, overlays (epic, release, help, detail, epic-detail), navigation (left, right, up, down, pgup, pgdown), task transitions (space, backspace), router priority (p), epic panel focus switching, and window resize. This should be 3–4 functions.

### Missing abstraction

A `FrontmatterDocument` type that encapsulates "parse YAML → modify → reserialise with body preserved" would eliminate the duplicated write logic.

### Excessive abstraction

The `Discover`, `Parser`, `ConfigReader`, and `RouterReader` structs are over-abstracted — they hold no state and exist only as method namespaces. They should be either interfaces (H4) or package-level functions (M2).

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

This structure is clean and appropriate. `cmd/` handles arg parsing, `internal/` contains all logic, templates are embedded. No changes recommended.

### Domain boundaries ✅

Boundaries are sensible. `data` owns models and persistence, `board` owns presentation, `doctor` owns diagnostics. No circular dependencies.

### State management

The Bubble Tea `Model` struct has 27 fields. This is on the upper bound of manageable. If the TUI grows, consider grouping related fields into sub-structs (e.g., `EpicPanelState`, `OverlayState`).

### API / data layer

The data layer is file-based (markdown + YAML frontmatter). This is the right choice for the project's stated design philosophy of "slow down, write things down." No database or network calls exist.

### Configuration approach ✅

`config.yml` with defaults baked into Go code. Clean and simple.

### Error handling strategy

Error handling is generally good — errors are wrapped with `fmt.Errorf("context: %w", err)` consistently. One concern: `reloadTasks` in [watch.go:58](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/watch.go#L56-L64) silently swallows all errors by returning `nil`:

```go
tasks, releases, releaseEpics, epicStatuses, err := loadBoardData(root)
if err != nil {
    return nil  // error silently dropped
}
```

This means a malformed task file will cause the board to silently stop refreshing. Consider emitting an `errorMsg` type instead.

---

## 5. Best-Practice Review

### Framework conventions ✅

The Bubble Tea Model/Update/View pattern is followed correctly. Messages are properly typed. Commands are returned where appropriate.

### Type safety ✅

Custom types (`ColumnType`, `ProgressStage`, `OverlayType`, `TaskStatus`) are used appropriately instead of raw strings. The `ColumnType` → `ColumnType` validation in `lifecycle.go` is solid.

### Linting/formatting ⚠️

No linter is configured. `go fmt` appears to have been run (code is consistently formatted), but static analysis tools like `staticcheck` or `golangci-lint` would catch the dead code and duplication issues found in this audit.

### Testing approach ✅

Strong test coverage. Tests use filesystem fixtures (temp directories) and verify behaviour end-to-end. The `board` package has both unit tests and integration tests. The test-to-code ratio of ~1.6:1 is excellent.

### Dependency management ✅

`go.mod` is clean. Only 2 direct dependencies. All indirect deps are from the Charmbracelet ecosystem. `go.sum` is checked in (correct for Go modules).

### Build/deployment setup ⚠️

The `buildtool` is a self-contained Go binary invoked via `go run ./internal/buildtool`. This is creative but unconventional — most Go projects use `goreleaser` or direct `go build` invocations. The approach works but:
- No CI configuration is visible (no `.github/workflows/`, no `Makefile` CI target)
- Binaries are committed to the repo (H1)

### Environment variable handling ✅

`NO_COLOR` and `COLORTERM` are respected in [theme.go](file:///c:/Users/User/Branding/03-VIBE-LAB/savepoint/internal/board/theme.go). No secrets or sensitive env vars are used.

### Logging/debugging practices ⚠️

No logging at all. The TUI uses `StatusMessage` for user feedback, but there is no debug log for filesystem operations, parser errors, or watcher events. For a TUI tool, this is borderline acceptable, but a `--verbose` or `--debug` flag that writes to a log file would help troubleshoot issues in the field.

---

## 6. Refactor Roadmap

### Phase 1 — Safe cleanup

**Objective:** Remove dead code, fix obvious duplication, improve hygiene.

| Task | Risk |
|------|------|
| Delete `taskLabel()` from `column.go` | None |
| Delete `loadAllTasks()` from `board.go` | None |
| Delete or update `newProgramModel()` in `board.go` | None |
| Replace custom `contains`/`indexOf` in `repairs.go` with `strings.Contains` | None |
| Replace custom `trimSpace` in `buildtool/main.go` with `strings.TrimSpace` | None |
| Add `savepoint`, `savepoint.exe`, `dist/`, `*.zip` to `.gitignore` and remove from tracking | None |
| Fix `package.json` test script | None |
| Add `.golangci.yml` with basic linters | None |

**Expected benefit:** Cleaner codebase, smaller repo, automated lint catches.
**Risk level:** Low.

---

### Phase 2 — Structural improvements

**Objective:** Reduce duplication, improve modularity, make key files easier to extend.

| Task | Risk |
|------|------|
| Extract `SplitFrontmatterBody()` in `data` package; refactor `write.go` | Low — well-tested |
| Extract `stripFrontmatter()` helper for `epic_panel.go` | Low |
| Consolidate `shortID` / `shortRouterID` into one function | Low |
| Consolidate `epicIndex` / `releaseIndex` into `sliceIndex` | Low |
| Split `Update()` into `handleBoardKey()`, `handleOverlayKey()`, and named action methods | Medium — many tests depend on exact Update behaviour |
| Group `Model` fields into sub-structs (`EpicPanelState`, `OverlayState`) | Medium |
| Convert stateless data types to package-level functions or introduce consumer-defined interfaces | Medium |

**Expected benefit:** Easier to add features, reduced duplication, more testable.
**Risk level:** Medium.

---

### Phase 3 — Hardening

**Objective:** Improve error handling, add diagnostic capabilities, fill test gaps.

| Task | Risk |
|------|------|
| Handle `reloadTasks` errors by emitting an `errorMsg` instead of returning `nil` | Low |
| Add `buildtool` tests | Low |
| Add `--verbose` debug logging to a log file | Low |
| Document `splitCommand` tokenisation limitations | None |
| Add CI workflow (GitHub Actions) with `make build && make test` | Low |
| Unify `ColumnType` / `TaskStatus` into a single status model | Medium — data format change |

**Expected benefit:** Fewer silent failures, better debugging, CI safety net.
**Risk level:** Low–Medium.

---

## 7. Top 10 Action List

- [ ] **1. Add binaries to `.gitignore` and remove from Git tracking** — Severity: High — Files: `.gitignore`, `savepoint`, `savepoint.exe`, `dist/`, `ink-cli-ui-design.zip` — Benefit: Repo shrinks by 12 MB+, eliminates merge noise
- [ ] **2. Extract `SplitFrontmatterBody()` to deduplicate write logic** — Severity: High — Files: `internal/data/write.go` — Benefit: Single source of truth for frontmatter reconstruction
- [ ] **3. Split `Update()` into named key-handler methods** — Severity: High — Files: `internal/board/update.go` — Benefit: 170-line method becomes 5–6 focused methods; easier to extend
- [ ] **4. Delete dead code: `taskLabel()`, `loadAllTasks()`, `newProgramModel()`** — Severity: Low — Files: `internal/board/column.go`, `internal/board/board.go` — Benefit: Less confusion for contributors
- [ ] **5. Replace custom `contains`/`indexOf`/`trimSpace` with stdlib** — Severity: Medium — Files: `internal/doctor/repairs.go`, `internal/buildtool/main.go` — Benefit: Idiomatic Go, less code to maintain
- [ ] **6. Extract `stripFrontmatter()` helper for `epic_panel.go`** — Severity: Medium — Files: `internal/board/epic_panel.go` — Benefit: DRY, consistent behaviour
- [ ] **7. Add `.golangci.yml` linter configuration** — Severity: Medium — Files: project root — Benefit: Automated detection of dead code, unchecked errors, and style issues
- [ ] **8. Handle `reloadTasks` errors instead of silently swallowing** — Severity: Medium — Files: `internal/board/watch.go` — Benefit: Board surfaces parse errors to the user instead of silently freezing
- [ ] **9. Add tests for `buildtool`** — Severity: Medium — Files: `internal/buildtool/` — Benefit: Cross-platform build confidence
- [ ] **10. Consolidate `shortID`/`shortRouterID` and `epicIndex`/`releaseIndex` duplicates** — Severity: Low — Files: `internal/board/card.go`, `internal/board/view.go`, `internal/board/epic_panel.go`, `internal/board/release.go` — Benefit: Single source of truth for slug extraction and list search
