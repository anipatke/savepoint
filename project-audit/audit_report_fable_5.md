# Audit Report — Savepoint (Claude Fable 5, 2026-06-12)

Independent full-codebase audit per `project-audit/audit-prompt.md`. Findings were formed before reading prior reports; section 4 cross-checks against `consolidated-audit-report.md`.

## 1. Executive Summary

**Baseline:** `make build` and `make test` pass. `go vet ./...` is clean. 686 top-level test functions across 50 test files (~11,410 test lines) against ~7,250 production lines in 56 files — roughly 1.6× test-to-production ratio. Every package except `internal/testutil` (a test helper package) has tests. The codebase has grown ~30% since the prior audit and has visibly absorbed most of its recommendations.

### What is working well

- **The prior audit's two Critical findings are fixed.** All TUI file I/O now runs in `tea.Cmd` functions (`internal/board/io.go`) returning typed messages; cycle detection was rewritten with a path stack (`internal/doctor/checks.go:443-491`) and the comment even documents the old parent-map bug.
- **Mtime conflict detection with self-healing retry.** Task writes verify mtime, and on conflict re-read the file and retry only if the lifecycle base is unchanged (`internal/board/io.go:99-123`). This is careful, correct concurrent-edit handling.
- **Lifecycle rules live in one place.** `internal/data/lifecycle.go` owns status/stage validation, legacy alias resolution, and typed diagnostics; board and doctor both consume it.
- **Fuzz and policy tests.** `FuzzSplitFrontmatterBody` with a checked-in corpus (`internal/data/fuzz_test.go`), render-policy tests, and template-freshness tests (`internal/init/template_freshness_test.go`) that fail the build if `agent-skills/` and `templates/project/agent-skills/` drift apart.
- **Interfaces at I/O boundaries.** `internal/board/interfaces.go` and `internal/doctor/interfaces.go` define `taskDiscoverer`/`taskParser`/`routerReader`/config-reader seams with dependency injection, fixing the prior "concrete structs only" complaint.
- **Distribution hardened.** Windows targets cross-compiled and PE-header-verified (`internal/buildtool/main.go:157-172`), SHA-256 checksums generated, buildtool has tests, binaries are no longer committed (`.gitignore` covers `dist/`, `savepoint.exe`).

### Biggest risks

1. **Task/router/defect writes are not atomic** — `internal/data/write.go` uses plain `os.WriteFile` while a proper `AtomicWrite` already exists in `internal/init/write.go`. A crash mid-write corrupts the source of truth.
2. **File watching dies permanently after one failed reload** — a transient parse error (e.g. an agent mid-write) stops live updates until a manual `ctrl+r`, with no indication the watcher is dead.
3. **Doc/skill drift around `savepoint-create-plan`** — the scaffolded router tells agents to use a skill that the AGENTS.md activation table doesn't list.

### Extensibility

Good. New checks, overlays, key handlers, and commands all have obvious extension points with test seams. `update.go` is still the largest production file (730 lines) but dispatch is now split into per-overlay handlers, so it is no longer the monolith the prior audit described. The architecture (flat `internal/` packages, Elm-style TUI, embedded templates, markdown-as-database) remains well matched to the product with no over-engineering observed.

## 2. Severity-Ranked Findings

### Critical

None. Baseline is green and no data-corruption path is reachable without a crash or external interference.

### High

#### H1. Data-layer writes are not atomic; the atomic helper exists but only `init` uses it

- **Finding:** Every write in `internal/data/write.go` uses `os.WriteFile` directly: `ApplyProposal` (line 27), `updateFrontmatterField` (line 94), `WriteTaskStatus` (line 162), `WriteDefectStatus` (line 220), `WriteRouterState` (line 286). `os.WriteFile` truncates then writes, so a crash, kill, or full disk mid-write leaves a truncated or partial task/router/defect file. Meanwhile `internal/init/write.go:10-44` implements a correct `AtomicWrite` (temp file + fsync + rename, with a copy fallback for rename failures) that the data layer never uses.
- **Why it matters:** These files are the product's entire source of truth ("No database. Your filesystem is the source of truth"). The TUI also writes them while an fsnotify watcher and possibly an agent session read them; non-atomic writes mean readers can observe truncated content. The 100ms debounce usually hides this, but it is not a guarantee.
- **Evidence:** `internal/data/write.go:27,94,162,220,286`; `internal/init/write.go:10`.
- **Recommended fix:** Move `AtomicWrite` to a shared location (e.g. `internal/data`) and route all five write sites through it. This also resolves a one-source-of-truth violation (two write strategies for the same job).
- **Verified:** yes (code inspection; both implementations read in full).

#### H2. File watching silently dies after a single failed reload

- **Finding:** The `watchFiles` command delivers exactly one message and terminates; it is re-armed only in two places — `Init()` (`internal/board/model.go:118-123`) and the `reloadMsg` branch of `Update()` (`internal/board/update.go:41-43`). When a file change triggers `reloadTasks` and `loadBoardData` fails — which any malformed frontmatter does, since `loadBoardData` aborts on the first parse error — the command returns `errorMsg` (`internal/board/watch.go:101-105`), and the `errorMsg` branch (`internal/board/update.go:93-94`) only sets `StatusMessage` without re-arming the watcher. From then on no file event is ever observed again.
- **Why it matters:** The most likely trigger is exactly the normal workflow: an agent or editor saves a task file in two steps and the watcher fires between them. The board shows "reload failed: …" once, then silently stops live-updating. The user sees stale state with no hint that watching is dead; the fix on their side is an undiscoverable manual `ctrl+r`.
- **Evidence:** `internal/board/update.go:20-24,41-43,93-94`; `internal/board/watch.go:94-110`; `internal/board/model.go:118-123`.
- **Recommended fix:** Make the reload failure path re-arm the watcher — e.g. have `reloadTasksWithMessage` return a dedicated `reloadFailedMsg`, and in its `Update` branch return `watchFiles(m.Watcher)` alongside setting the status message. Add a regression test: inject a failing loader, deliver `fileChangeMsg`, assert the returned command chain ends with a watch command.
- **Verified:** yes (full code trace of every `watchFiles` call site; not reproduced at runtime).

### Medium

#### M1. `ParseDefectFileFromDisk` is dead code carrying a latent mtime bug

- **Finding:** `internal/data/parser.go:224-243` defines `ParseDefectFileFromDisk`, which no production or test code calls (repo-wide search matches only the definition). It truncates the recorded mtime to whole seconds (`fi.ModTime().Truncate(time.Second)`), while `WriteDefectStatus` compares with `fi.ModTime().Equal(expectedMtime)` against an untruncated stat (`internal/data/write.go:175`). NTFS stores 100ns-precision mtimes, so any future caller wiring this function's `Mtime` into `WriteDefectStatus` would hit a spurious `ErrMtimeConflict` on virtually every write. The live defect-loading path (`internal/board/board.go:152-179`) correctly stores the untruncated mtime.
- **Why it matters:** Dead code violates "build only what is needed", and this particular dead code is a pre-planted bug for whoever discovers and reuses it.
- **Evidence:** `internal/data/parser.go:241`; `internal/data/write.go:175`; grep for `ParseDefectFileFromDisk` returns only the definition.
- **Recommended fix:** Delete the function. If it is kept for a planned feature, remove the `Truncate` call and add a test pairing it with `WriteDefectStatus`.
- **Verified:** yes (dead-code status verified by search; the truncation mismatch is verified by inspection, not runtime).

#### M2. Repair suggestions still matched by substring search on error text

- **Finding:** `internal/doctor/repairs.go:23-85` maps problems to suggestions via ~30 `strings.Contains` checks on error message text produced elsewhere (mostly `internal/data/lifecycle.go`). Typed errors exist and are checked first (`errors.Is`, lines 12-21), but only for four sentinel errors; everything else couples doctor behavior to exact wording in another package. Rewording `"task stage is required when status is in_progress"` in `lifecycle.go:273` would silently degrade the doctor's suggestion to the generic fallback — no test would fail at the data layer.
- **Why it matters:** One-source-of-truth violation with silent failure mode; this was prior finding H6 and is still present, only partially mitigated.
- **Evidence:** `internal/doctor/repairs.go:23-85`; message origins in `internal/data/lifecycle.go:241,249,264,273,280,289`.
- **Recommended fix:** Extend the existing pattern that already works: `DiagnoseTaskLifecycle` returns typed `TaskLifecycleDiagnosticCode` values — have doctor map codes (not message text) to suggestions, and add sentinel errors for the remaining cases.
- **Verified:** yes.

#### M3. Skill-activation tables contradict each other about `savepoint-create-plan`

- **Finding:** The scaffolded router (`templates/project/.savepoint/router.md:29`) maps `pre-implementation → savepoint-draft-prd or savepoint-create-plan`, and the skill exists in both `agent-skills/savepoint-create-plan/` and the template copy, triggering on "router state is pre-implementation". But the Skill Activation tables in both the repo `AGENTS.md:15` and `templates/project/AGENTS.md:15` list only `savepoint-draft-prd` for `pre-implementation` and never mention `savepoint-create-plan` at all.
- **Why it matters:** These tables are the routing contract that every agent session reads first (AGENTS.md is step 2 of the workflow). An agent following AGENTS.md will never activate the planning skill; an agent following router.md will use a skill its guide says doesn't exist for that state. The project's own docs say "one source of truth — no duplicated rules".
- **Evidence:** `templates/project/.savepoint/router.md:29`; `AGENTS.md:15-20`; `templates/project/AGENTS.md:15-20`; `agent-skills/savepoint-create-plan/SKILL.md:13-14`.
- **Recommended fix:** Decide the canonical mapping (likely `pre-implementation → savepoint-draft-prd, then savepoint-create-plan once PRD/Design are approved`) and state it identically in both AGENTS.md files and the router template. Consider extending the template-freshness test to compare the activation tables.
- **Verified:** yes.

#### M4. Doctor is documented as read-only but executes arbitrary configured commands

- **Finding:** `AGENTS.md:105` describes `internal/doctor/` as "Read-only project diagnostics", and the README says doctor "checks that the project state is still coherent". But `RunQualityGates` (`internal/doctor/gates.go:23-65`) executes whatever `lint`/`typecheck`/`test` commands appear in `config.yml` via `exec.CommandContext` with full environment inheritance. Additionally, `splitCommand` (`gates.go:123-193`) is a quote-aware tokenizer with no shell semantics, so a gate like `test: "go vet ./... && go test ./..."` passes `&&` as a literal argument and fails confusingly.
- **Why it matters:** Running configured gates is a legitimate, intentional feature — the issue is the documentation mismatch (a user or agent told "doctor is read-only" may run it in a state they didn't want mutated by test/build side effects) and the silent operator misbehavior, which produces a failing gate with a misleading error.
- **Evidence:** `internal/doctor/gates.go:81-85,123-193`; `AGENTS.md:105`.
- **Recommended fix:** Reword the AGENTS.md module description ("diagnostics plus configured quality-gate execution"); in `runGate`, detect shell operators (`&&`, `||`, `|`, `;`, `>`) in the command string and return a clear "use a script or Makefile target; gates do not run a shell" error.
- **Verified:** yes (code inspection; operator behavior follows directly from the tokenizer).

#### M5. Residual duplicated frontmatter/CRLF handling outside `internal/data`

- **Finding:** Two remnants of the previously-fixed duplication: (1) `stripFrontmatter` in `internal/board/epic_panel.go:43-55` re-implements `---` delimiter scanning line-by-line instead of using `data.SplitFrontmatterBody`; its acceptance of indented/whitespace-padded `---` lines also differs from the canonical parser's rules, so the overlay can render a file differently than the data layer parses it. (2) `hasAcceptanceCriteria` in `internal/doctor/checks.go:276` does its own `\r\n` replacement, omitting the bare-`\r` (legacy Mac) case the canonical `normalizeLineEndings` handles (`internal/data/parser.go:99-102`).
- **Why it matters:** This is the exact failure class the prior audit flagged (H3/H8): parallel parsing rules drift independently.
- **Evidence:** `internal/board/epic_panel.go:43-55`; `internal/doctor/checks.go:276`; `internal/data/parser.go:99-102`.
- **Recommended fix:** Export a `data.StripFrontmatter` (or reuse `SplitFrontmatterBody` and tolerate its error by falling back to full content), and export `normalizeLineEndings` for doctor's use.
- **Verified:** yes.

#### M6. Nothing ties the npm package version to the binary's `--version`

- **Finding:** `package.json` declares `1.2.12`; the binary's version is injected at build time from, in priority order, `-version` flag → `VERSION` env → `git describe --tags --abbrev=0` → `"v0.0.0"` (`internal/buildtool/main.go:267-281`). `npm publish` triggers `prepublishOnly → build-npm` with no version flag, so the binary silently gets the latest git tag. Publishing before tagging (or from a machine without tags) ships a package whose `savepoint --version` disagrees with the installed npm version — with no error. The formats also differ (`v1.2.12` vs `1.2.12`).
- **Why it matters:** Version skew between `npm ls savepoint` and `savepoint --version` makes bug reports and upgrade verification unreliable, and `upgrade-assets` workflows depend on users knowing which binary they actually have.
- **Evidence:** `internal/buildtool/main.go:267-281`; `package.json:3,29-34`. Currently consistent (tag `v1.2.12` = package 1.2.12), so the risk is procedural, not active.
- **Recommended fix:** In `buildNPM`, read the version from `package.json` (or compare it against `git describe` and fail on mismatch) instead of trusting tag state at publish time.
- **Verified:** yes (mechanism verified; skew not currently present).

### Low

#### L1. `Scaffold` and `UpgradeProjectAssets` accept `force` parameters they never use

- **Finding:** `Scaffold(templates, targetDir, projectName, force)` (`internal/init/scaffold.go:13-53`) and `UpgradeProjectAssets(templates, targetDir, dryRun, force)` (`internal/init/upgrade.go:71`) both take `force` and never reference it; callers pass `opts.Force` (`main.go:106,126`). Force semantics are actually enforced only in `ValidateTarget` — once validation passes, `Scaffold` unconditionally overwrites every template-managed file.
- **Why it matters:** Misleading API: a reader assumes per-file force behavior that doesn't exist. Dead parameters are "build only what is needed" violations.
- **Evidence:** `internal/init/scaffold.go:13`; `internal/init/upgrade.go:71`; `main.go:106,126`.
- **Recommended fix:** Drop both parameters, or implement the implied behavior (skip existing files unless force).
- **Verified:** yes.

#### L2. Dead artifact: `scripts/vitest-preload.cjs`

- **Finding:** The tracked file stubs `net use` for esbuild/vitest — machinery from the project's earlier TypeScript/Ink incarnation. Nothing in the current Go project references vitest except archived v1 epic documents; there is no `vitest` config, no `node_modules` usage, and `package.json` has no vitest dependency.
- **Evidence:** `scripts/vitest-preload.cjs`; references only under `.savepoint/releases/v1/epics/_archived/`.
- **Recommended fix:** Delete the file.
- **Verified:** yes.

#### L3. Makefile `build-all` omits Windows and there is no lint/vet target

- **Finding:** `Makefile:23` defines `build-all: build-linux build-darwin`, while `go run ./internal/buildtool build-all` builds all six targets including Windows — same name, different meaning. Separately, `go vet` is clean but nothing runs it: no `lint`/`vet` target exists and `ci: test build` skips it (prior finding L6 still open).
- **Evidence:** `Makefile:23,31`; `internal/buildtool/main.go:24-31,67-68`.
- **Recommended fix:** `build-all: build-linux build-darwin build-windows` (adding the missing make target), and add `vet: go vet ./...` to `ci`.
- **Verified:** yes.

#### L4. Detail overlay shows "Stage: build" for tasks that have no stage

- **Finding:** `stageLabel` returns `"build"` as the default case (`internal/board/detail.go:214-223`), so `planned` and `done` tasks — which by the project's own lifecycle rules must have no stage — render `Stage: build` in the task detail overlay.
- **Why it matters:** The TUI is the workflow's authority display; showing a stage on a non-in_progress task contradicts the lifecycle rules the rest of the system enforces strictly.
- **Evidence:** `internal/board/detail.go:32,214-223`; lifecycle rule at `internal/data/lifecycle.go:426-428`.
- **Recommended fix:** Return `"—"` (or empty) for an empty stage; keep `"build"` only as the in_progress default if desired.
- **Verified:** yes.

#### L5. `WriteRouterState` index math breaks if non-ASCII text precedes the state block

- **Finding:** `internal/data/write.go:263` finds the state block with `strings.Index(strings.ToLower(normalized), strings.ToLower(stateBlockStart))` and uses the resulting index into the **original** string. `strings.ToLower` can change byte length for non-ASCII runes (e.g. `İ` → `i̇`), so a router.md containing such characters before the "Current state" heading would misalign the indices and corrupt the file on write.
- **Why it matters:** Real-world router files are ASCII templates, so exposure is low — but the failure mode is file corruption, and the fix is trivial.
- **Evidence:** `internal/data/write.go:263-265`.
- **Recommended fix:** Use a case-preserving search (e.g. find via `strings.Index` on both common casings, or scan line-by-line with `strings.EqualFold`).
- **Verified:** yes (mechanism verified by inspection; not reproduced).

#### L6. Watcher error channel is drained but errors are discarded silently

- **Finding:** `internal/board/watch.go:85-88` reads `w.Errors` and does nothing with the value — not even the package's own `debugf`. fsnotify errors (e.g. watch overflow on Windows) vanish.
- **Recommended fix:** `debugf("watcher: error %v", err)` at minimum.
- **Verified:** yes.

## 3. Code Style Scorecard

| Rule | Verdict | Evidence |
|------|---------|----------|
| One job per file | **Partial** | Generally excellent (43 focused files in `board/` alone), but `update.go` is 730 lines mixing message dispatch, key handling, scroll math, and file reading (`internal/board/update.go:488-507` is I/O in a dispatch file) |
| One job per function | **Pass** | Functions are small and named for one action throughout; e.g. `internal/data/lifecycle.go` is 30+ single-purpose functions |
| Test branches | **Pass** | 686 tests incl. fuzz corpus (`internal/data/fuzz_test.go`), render-policy tests, conflict-retry tests; meaningful conditionals covered |
| Types document intent | **Partial** | Strong typed enums and diagnostic codes, but `Task.Column data.ColumnType` still names the *status* concept after a UI column (`internal/data/task.go`), requiring the `IsCanonicalColumn = IsCanonicalTaskStatus` shim (`internal/data/lifecycle.go:389-391`) |
| Build only what is needed | **Partial** | Dead `ParseDefectFileFromDisk` (M1), unused `force` params (L1), dead `scripts/vitest-preload.cjs` (L2) |
| Handle errors at boundaries | **Pass** | Validation at parse and write, mtime conflict handling, typed init validation errors; minor gap: discarded watcher errors (L6) |
| One source of truth | **Partial** | Lifecycle rules centralized well, but two write strategies (H1), duplicated frontmatter stripping (M5), suggestion-by-message-text (M2), divergent `build-all` meanings (L3) |
| Comments explain why | **Pass** | e.g. `conservativeColumnPageSize` indicator-budget rationale (`internal/board/update.go:616-618`), `detectCycles` referencing the bug it prevents (`internal/doctor/checks.go:443-445`) |
| Content lives in data | **Pass** | Templates embedded as files, gates configured in `config.yml`, palette isolated in `internal/styles` |
| Small diffs | **Pass** | Recent git history shows focused, reviewable commits |

## 4. Prior-Audit Cross-Check

Status of every Critical/High finding from `consolidated-audit-report.md`, plus notable Medium/Low items:

| Prior ID | Status | Current evidence |
|----------|--------|------------------|
| C1 sync I/O in Update loop | **Fixed** | All I/O in `tea.Cmd`s (`internal/board/io.go`); `Update()` is pure dispatch |
| C2 inaccurate cycle paths | **Fixed** | Path-stack DFS, `internal/doctor/checks.go:443-491` |
| H1 committed binaries | **Fixed** | `.gitignore` covers `dist/`, `savepoint`, `savepoint.exe`; `git ls-files` shows none tracked |
| H2 `update.go` god-method | **Largely fixed** | Dispatch split into per-overlay handlers; file still largest at 730 lines (see scorecard) |
| H3 duplicated frontmatter logic | **Mostly fixed** | Centralized in `data`; residue in `epic_panel.go:43` and `checks.go:276` (new M5) |
| H4 no I/O interfaces | **Fixed** | `internal/board/interfaces.go`, `internal/doctor/interfaces.go` with DI |
| H5 stdlib reimplementation | **Fixed** | buildtool uses `archive/tar`, `crypto/sha256`; repairs.go is a plain mapping |
| H6 fragile repair matching | **Still present** | `internal/doctor/repairs.go:23-85` (carried forward as M2) |
| H7 gates without timeout | **Fixed** | `context.WithTimeout`, configurable via `config.yml` (`internal/doctor/gates.go:36-41`) |
| H8 scattered CRLF handling | **Mostly fixed** | `normalizeLineEndings` canonical in `data`; one local copy in doctor (M5) |
| M6 broken AtomicWrite fallback | **Fixed** | Copy-based `replaceFile` fallback (`internal/init/write.go:46-73`) |
| M7 ad-hoc markdown parsing in epic_panel | **Partially fixed** | Rendering simplified but `stripFrontmatter` remains ad hoc (M5) |
| M9 naïve `splitCommand` | **Partially fixed** | Now quote/escape-aware; still no shell operators (folded into M4) |
| M10 no Windows build target | **Fixed** | Six targets incl. windows/arm64 + PE verification |
| M11 buildtool untested | **Fixed** | `internal/buildtool/main_test.go`; package passes in CI run |
| L6 no linter | **Still present** | No lint/vet target in Makefile (L3) |
| L7 no checksums | **Fixed** | `writeChecksums` (`internal/buildtool/main.go:191-213`) |
| L8 hardcoded skill count | **Fixed** | `agent_skills_test.go` enumerates directories |
| L10 abbreviation sentence-splitting | **Fixed** | `knownAbbreviations` (`internal/board/detail.go:97-113`) |
| L13 reload errors swallowed | **Fixed, with regression** | Errors now surface as `errorMsg` — but that path forgets to re-arm the watcher (new H2) |

**New findings prior audits missed:** H2 (watcher death — introduced by the L13 fix), M1 (dead defect parser with mtime bug), M3 (create-plan table drift), M6 (version skew mechanism), L1, L4, L5.

## 5. Quick Wins

1. Re-arm `watchFiles` on the reload-failure path (fixes H2; ~10 lines + test).
2. Route `internal/data/write.go` writes through `AtomicWrite` (fixes H1 once the helper is moved/shared).
3. Delete `ParseDefectFileFromDisk` (M1).
4. Delete `scripts/vitest-preload.cjs` (L2).
5. Add `savepoint-create-plan` to both AGENTS.md activation tables, matching router.md (M3).
6. Drop the unused `force` parameters from `Scaffold`/`UpgradeProjectAssets` (L1).
7. Return "—" from `stageLabel` for empty stage (L4).
8. Log watcher errors via `debugf` (L6).
9. Add `vet`/`lint` to the Makefile and `ci` target (L3).
10. Make Makefile `build-all` include `build-windows` or rename it (L3).

## 6. Do Not Change

Things that look unusual but are correct or intentional:

- **Legacy status aliases (`todo`, `complete`, `completed`, `phase`, `implementation`) accepted on load** (`internal/data/lifecycle.go:8-16`). AGENTS.md's "Never: todo, doing, …" governs *writing*; read-side tolerance is deliberate parse compatibility, and writes normalize.
- **Doctor executing quality gates** is a feature, not a violation — only its documentation needs the M4 wording fix.
- **Full reload on any file change** (`reloadTasksWithMessage` re-reads everything). At this data scale, simple-and-correct beats incremental cache invalidation; don't "optimize" this.
- **Skill content duplicated into `templates/project/agent-skills/`** is by design, and `template_freshness_test.go` enforces byte-identical sync — the duplication is mechanically guarded.
- **`go build main.go` (by filename) in buildtool** is fine: the root `main` package is intentionally a single file plus tests.
- **Forced color profile** (`applyColorProfile`) is the deliberate Atari-Noir rendering policy, backed by `render_policy_test.go`.
- **`checkWritable` probe file** (`internal/init/validate.go:74-84`) looks hacky but is the portable way to test writability on Windows.

---

*Audit performed read-only; no production code was modified. Build/test/vet were run for baseline verification only.*
