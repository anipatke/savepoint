---
type: audit-findings
audited: 2026-06-30
---

# Audit Findings: E31 Audit Register Data Model

## Main Findings

E31 is functionally close to complete: all four task files are marked `done`, and the implementation adds typed finding, prompt/register, run, discovery, and cross-record validation models under `internal/data`. The scoped and project gates pass: `go test ./internal/data` and `make build && make test`.

One unresolved T004 validation issue remains. `ValidateAuditFindings` currently checks exact task/epic IDs first, then falls back to short-ID matching for every reference. That means a full-looking but broken reference such as `E99-missing/T004-audit-register-validation` can pass if any discovered `T004` exists, and `E31-wrong-slug` can pass if an `E31-*` epic exists. This under-reports broken work-item links, so T004 AC5 is not fully satisfied until the fallback is limited to true shorthand references.

Process drift noted: the router still says `state: task-building` for E31/T004, and `E31-Detail.md` still has `status: planned`, even though all E31 tasks are `done`. This audit was written from the explicit user request, not from an `audit-pending` router handoff. Apply/close should normalize the epic/router state after the proposed fix is approved.

Known template/model drift from T001 remains non-blocking for this epic: the E30 finding template documents singular `work_item`, while E31 supports typed plural `releases`, `epics`, `tasks`, and `defects` plus compatibility parsing for `work_item`.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [ ] Test branches - missing coverage for full-looking task/epic refs that should not resolve through shorthand fallback
- [x] Types document intent
- [x] Build only what is needed
- [ ] Handle errors at boundaries - broken full task/epic references can be accepted as valid
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File
internal/data/audit_validate.go

### Replace
```go
		for _, ref := range f.Epics {
			if !epicIDs[ref] && !epicShortIDs[epicShortID(ref)] {
				add(f, AuditEpicRefMissing, fmt.Sprintf("references unknown epic %q", ref))
			}
		}
		for _, ref := range f.Tasks {
			if !taskIDs[ref] && !taskShortIDs[taskShortID(ref)] {
				add(f, AuditTaskRefMissing, fmt.Sprintf("references unknown task %q", ref))
			}
		}
```

### With
```go
		for _, ref := range f.Epics {
			if !resolvesEpicRef(ref, epicIDs, epicShortIDs) {
				add(f, AuditEpicRefMissing, fmt.Sprintf("references unknown epic %q", ref))
			}
		}
		for _, ref := range f.Tasks {
			if !resolvesTaskRef(ref, taskIDs, taskShortIDs) {
				add(f, AuditTaskRefMissing, fmt.Sprintf("references unknown task %q", ref))
			}
		}
```

### Target File
internal/data/audit_validate.go

### Replace
```go
func shortIDSet(ids []string, short func(string) string) map[string]bool {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		if s := short(id); s != "" {
			set[s] = true
		}
	}
	return set
}
```

### With
```go
func shortIDSet(ids []string, short func(string) string) map[string]bool {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		if s := short(id); s != "" {
			set[s] = true
		}
	}
	return set
}

func resolvesEpicRef(ref string, fullIDs, shortIDs map[string]bool) bool {
	if fullIDs[ref] {
		return true
	}
	return isShortEpicRef(ref) && shortIDs[ref]
}

func resolvesTaskRef(ref string, fullIDs, shortIDs map[string]bool) bool {
	if fullIDs[ref] {
		return true
	}
	if !isShortTaskRef(ref) {
		return false
	}
	return shortIDs[taskShortID(ref)]
}
```

### Target File
internal/data/audit_validate_test.go

### Replace
```go
		{
			name:     "unknown epic",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E99-missing"}},
			wantCode: AuditEpicRefMissing,
		},
		{
			name:     "epic by short id",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E31"}},
			wantNone: true,
		},
		{
			name:     "unknown task",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"E31-audit-register-data-model/T099-ghost"}},
			wantCode: AuditTaskRefMissing,
		},
		{
			name:     "task by short id",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"T004"}},
			wantNone: true,
		},
```

### With
```go
		{
			name:     "unknown epic",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E99-missing"}},
			wantCode: AuditEpicRefMissing,
		},
		{
			name:     "wrong epic slug does not resolve by short id",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E31-wrong-slug"}},
			wantCode: AuditEpicRefMissing,
		},
		{
			name:     "epic by short id",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E31"}},
			wantNone: true,
		},
		{
			name:     "unknown task",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"E31-audit-register-data-model/T099-ghost"}},
			wantCode: AuditTaskRefMissing,
		},
		{
			name:     "wrong full task path does not resolve by short id",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"E99-missing/T004-audit-register-validation"}},
			wantCode: AuditTaskRefMissing,
		},
		{
			name:     "task by short id",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"T004"}},
			wantNone: true,
		},
```
