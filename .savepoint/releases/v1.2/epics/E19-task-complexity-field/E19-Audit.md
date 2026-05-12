---
type: audit-findings
audited: 2026-05-12
---

# Audit Findings: E19 Task Complexity Field

## Main Findings

E19 is applied and closed. The data model and parser hold `complexity_tier` and `complexity_reason`, parser and doctor validation reject malformed metadata, task status writes preserve existing complexity fields, board cards show the compact tier, detail overlays show the tier plus reason with wrapping, and the live/scaffolded create-task skills are textually aligned with the shared rubric.

The audit fixes were applied:

1. `ValidateTaskLifecycle` now validates complexity metadata before any lifecycle early return, so `in_progress` task writes reject invalid tiers and malformed reason pairs.
2. A write-path regression test covers invalid complexity on an `in_progress` task.
3. The v1.2 task set now has complexity metadata on E17, E18, and E19 task files.

Verification after apply:

- `go test ./internal/data ./internal/board ./internal/doctor ./internal/init` passes.
- `make build` passes.
- Full `make test` still has the known unrelated root failure: `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` reports `agent-skills\savepoint-audit\SKILL.md missing YAML frontmatter`.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File
internal/data/lifecycle.go

### Replace
```go
func ValidateTaskLifecycle(task *Task) error {
	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", task.Column)
	}

	if task.Column == ColumnInProgress {
		if task.Stage == "" {
			task.Stage = StageBuild
			return nil
		}
		if !IsCanonicalStage(task.Stage) {
			return fmt.Errorf("invalid phase %q: use build, test, or audit. Add 'phase: build' to task frontmatter", task.Stage)
		}
		return nil
	}

	if task.Stage != "" {
		return fmt.Errorf("phase field %q is only valid when status is in_progress. Remove 'phase' or change status to in_progress", task.Stage)
	}

	if err := ValidateComplexity(task.ComplexityTier, task.ComplexityReason); err != nil {
		return err
	}

	return nil
}
```

### With
```go
func ValidateTaskLifecycle(task *Task) error {
	if err := ValidateComplexity(task.ComplexityTier, task.ComplexityReason); err != nil {
		return err
	}

	if !IsCanonicalColumn(task.Column) {
		return fmt.Errorf("invalid status %q: use planned, in_progress, or done. Add 'status: planned' or 'status: in_progress' to task frontmatter", task.Column)
	}

	if task.Column == ColumnInProgress {
		if task.Stage == "" {
			task.Stage = StageBuild
			return nil
		}
		if !IsCanonicalStage(task.Stage) {
			return fmt.Errorf("invalid phase %q: use build, test, or audit. Add 'phase: build' to task frontmatter", task.Stage)
		}
		return nil
	}

	if task.Stage != "" {
		return fmt.Errorf("phase field %q is only valid when status is in_progress. Remove 'phase' or change status to in_progress", task.Stage)
	}

	return nil
}
```

### Target File
internal/data/write_test.go

### Replace
```go
func TestWriteTaskStatus_preservesComplexityFields(t *testing.T) {
```

### With
```go
func TestWriteTaskStatus_rejectsInvalidComplexityOnInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "task.md")
	content := `---
id: E19/T011
status: planned
objective: "Invalid complexity"
---`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	fi, _ := os.Stat(path)
	task := &Task{
		ID:               "E19/T011",
		Column:           ColumnInProgress,
		Stage:            StageBuild,
		ComplexityTier:   ComplexityTier("extreme"),
		ComplexityReason: "Invalid tier should be rejected before writing.",
	}

	err := WriteTaskStatus(path, task, fi.ModTime())
	if err == nil {
		t.Fatal("WriteTaskStatus() expected invalid complexity error")
	}
	if !strings.Contains(err.Error(), "invalid complexity_tier") {
		t.Fatalf("WriteTaskStatus() error = %v, want invalid complexity_tier", err)
	}
}

func TestWriteTaskStatus_preservesComplexityFields(t *testing.T) {
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T007-create-defect-skill.md

### Replace
```md
objective: Add a create-defect skill and scaffold guidance for capturing release-level defects before repair work starts
depends_on: [E17-defect-workflow-tui/T006-doctor-docs-and-templates]
```

### With
```md
objective: Add a create-defect skill and scaffold guidance for capturing release-level defects before repair work starts
depends_on: [E17-defect-workflow-tui/T006-doctor-docs-and-templates]
complexity_tier: medium
complexity_reason: "Adds a new workflow skill plus scaffold copy and guidance across template and live docs"
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T008-footer-defects-shortcut.md

### Replace
```md
objective: Add the defects overlay shortcut to the board footer helper line
depends_on: [E17-defect-workflow-tui/T004-defects-overlay]
```

### With
```md
objective: Add the defects overlay shortcut to the board footer helper line
depends_on: [E17-defect-workflow-tui/T004-defects-overlay]
complexity_tier: low
complexity_reason: "Small board footer copy change advertising an existing shortcut"
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T009-defect-warning-glyph.md

### Replace
```md
objective: Standardize defect indicators on a terminal-safe warning glyph
depends_on: [E17-defect-workflow-tui/T005-defect-detail-and-markers]
```

### With
```md
objective: Standardize defect indicators on a terminal-safe warning glyph
depends_on: [E17-defect-workflow-tui/T005-defect-detail-and-markers]
complexity_tier: low
complexity_reason: "Glyph standardization is a narrow rendering constant and expectation update"
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T010-defect-resolve-hotkey.md

### Replace
```md
objective: Allow a user to resolve an open defect directly from the Defects overlay by pressing `space` on the focused defect row.
depends_on: [T004-defects-overlay, T005-defect-detail-and-markers]
```

### With
```md
objective: Allow a user to resolve an open defect directly from the Defects overlay by pressing `space` on the focused defect row.
depends_on: [T004-defects-overlay, T005-defect-detail-and-markers]
complexity_tier: medium
complexity_reason: "Adds mutation behavior from the overlay and must preserve defect lifecycle rules"
```

### Target File
.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T011-task-defect-count-marker.md

### Replace
```md
objective: Show related defect counts on task cards instead of a single defect number
depends_on: [E17-defect-workflow-tui/T005-defect-detail-and-markers, E17-defect-workflow-tui/T009-defect-warning-glyph]
```

### With
```md
objective: Show related defect counts on task cards instead of a single defect number
depends_on: [E17-defect-workflow-tui/T005-defect-detail-and-markers, E17-defect-workflow-tui/T009-defect-warning-glyph]
complexity_tier: medium
complexity_reason: "Changes card marker aggregation and rendering for related defect counts"
```

### Target File
.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/tasks/T004-artifact-templates.md

### Replace
```md
objective: Add explicit artifact template blocks to savepoint-audit, savepoint-create-task, and savepoint-system-design skills so any agent produces consistent E##-Audit.md, T###-slug.md, and E##-Detail.md files across all projects
depends_on: [E18-template-skill-optimisation/T001-canonical-guides]
```

### With
```md
objective: Add explicit artifact template blocks to savepoint-audit, savepoint-create-task, and savepoint-system-design skills so any agent produces consistent E##-Audit.md, T###-slug.md, and E##-Detail.md files across all projects
depends_on: [E18-template-skill-optimisation/T001-canonical-guides]
complexity_tier: medium
complexity_reason: "Adds canonical artifact templates across multiple skills and scaffolded copies"
```
