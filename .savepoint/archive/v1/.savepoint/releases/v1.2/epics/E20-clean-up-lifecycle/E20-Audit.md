---
type: audit-findings
audited: 2026-05-16
---

# Audit Findings: E20 Clean-up Lifecycle

## Main Findings

E20 is applied and closed. `internal/data/lifecycle.go` is the sole owner of the task lifecycle contract: canonical statuses (`planned`/`in_progress`/`done`), canonical stages (`build`/`test`/`audit`), legacy aliases (`todo`, `implementation` stage, `phase` field), parse normalization, write validation, transition helpers (`AdvanceTaskLifecycle*`/`RetreatTaskLifecycle*`), diagnostics (`DiagnoseTaskLifecycle`), and transition-equality (`SameTaskLifecycleForTransition`). Parser, writer, doctor, and board transitions all delegate to that surface — no package re-implements status/stage decisions.

Consumers verified thin:
- `internal/board/transitions.go` — `Advance`/`Retreat` are 2-line delegates; `CanAdvance` reuses `AdvanceTaskLifecycleState` for next-target gating while keeping dependency/epic-audit checks board-local.
- `internal/doctor/checks.go:checkTaskLifecycle` consumes `DiagnoseTaskLifecycle` rather than branching on raw frontmatter.
- `internal/data/parser.go` routes load-time normalization through `ParseTaskLifecycle`; `internal/data/write.go` writes `status: <Task.Column>` and validates via `ValidateTaskLifecycle`.

T001's open decision was resolved during audit: the denormalized `Task.Status` string field has been removed. `Task.Column` and `Task.Stage` are the canonical in-memory lifecycle fields. The dead `card.go` glyph fallback that read `Task.Status` was deleted, the `applyTaskLifecycleState` mirror write was removed, and tests that asserted the mirror were updated to check `Task.Column`/`Task.Stage` only. Disk YAML still uses `status:` as the canonical key (written from `Task.Column` in `write.go:137`), so persistence and round-trip are unchanged.

Verification after apply:

- `go test ./...` passes (board, data, doctor, init, buildtool, root).
- `go build ./...` passes.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed (removed denormalized `Task.Status` field; closed T001 open decision)
- [x] Handle errors at boundaries
- [x] One source of truth (`internal/data/lifecycle.go`)
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

All applied during audit. No outstanding follow-ups.

### Target File
internal/data/task.go

### Replace
```go
ID               string         `yaml:"id"`
Title            string         `yaml:"title"`
Description      string         `yaml:"description,omitempty"`
Status           string         `yaml:"status,omitempty"`
Epic             string         `yaml:"epic"`
```

### With
```go
ID               string         `yaml:"id"`
Title            string         `yaml:"title"`
Description      string         `yaml:"description,omitempty"`
Epic             string         `yaml:"epic"`
```

### Target File
internal/data/lifecycle.go

### Replace
```go
func applyTaskLifecycleState(task *Task, state TaskLifecycleState) {
	task.Column = state.Status
	task.Stage = state.Stage
	task.Status = string(state.Status)
}
```

### With
```go
func applyTaskLifecycleState(task *Task, state TaskLifecycleState) {
	task.Column = state.Status
	task.Stage = state.Stage
}
```

### Target File
internal/board/card.go

### Replace
```go
if isRouterPriority(t, routerState) {
    return styles.TagDone.Render(glyphBuild)
}
if t.Status != "" {
    return statusGlyph(t.Status)
}
return stageGlyphStyled(t.Stage)
```

### With
```go
if isRouterPriority(t, routerState) {
    return styles.TagDone.Render(glyphBuild)
}
return stageGlyphStyled(t.Stage)
```

### Target File
.savepoint/Design.md

### Replace
Section 8 task-card status glyph description mentioning explicit `Task.Status` with legacy fallback.

### With
Description that task cards derive glyphs from canonical `Task.Column`/`Task.Stage` only, and Section 5 note clarifying no denormalized status mirror exists on the Task struct.
