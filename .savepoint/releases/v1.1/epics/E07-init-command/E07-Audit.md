---
type: audit-findings
audited: 2026-05-03
---

# Audit Findings: E07 Init Command

## Main Findings

E07 is implemented and the Go baseline passes: `go test ./internal/init ./cmd`, `go build ./...`, and `go test ./...` all succeeded. The command parser, target validation, scaffold copy, prompt rendering, clipboard best-effort path, and integration coverage are present.

Audit apply resolved the generated-template drift. New Savepoint projects now receive `stage` task-status instructions, `make build && make test` quality-gate guidance, and the current epic-local `E##-Audit.md` audit workflow instead of the archived `.savepoint/audit/.../proposals.md` pipeline.

Audit apply also tightened scaffold safety. Init validation now treats an existing root `agent-skills/` directory as a conflict unless `--force` is set, and `AtomicWrite` uses a portable replace helper that falls back to backup/restore when a platform cannot rename over an existing target.

The clipboard task record was reconciled: `T006-clipboard.md` now reflects the implemented `ClipboardResult` API, best-effort status handling, current-platform tests, and completed Go quality gates.

## Code Style Review

- [x] One job per file
- [x] One-sentence functions
- [x] Test branches
- [x] Types are documentation
- [x] Build, don't speculate
- [x] Errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

The audit-close patch brought generated workflow content back in line with the live root instructions, added tests that guard template freshness, and kept the scaffold conflict/overwrite fixes scoped to `internal/init`.

## Proposed Changes

### Target File
templates/project/AGENTS.md

### Replace
```
- `phase` (build/test/audit): only when `status: in_progress`
```

### With
```
- `stage` (build/test/audit): **required** when `status: in_progress` — omitting it is a parse error
```

### Target File
templates/project/AGENTS.md

### Replace
```
2. When starting implementation, set task frontmatter to `status: in_progress`
```

### With
```
2. When starting implementation, set task frontmatter to `status: in_progress` + `stage: build` (both required together)
```

### Target File
templates/project/AGENTS.md

### Replace
```
npm run build && npm run test
```

### With
```
make build && make test
```

### Target File
templates/project/.savepoint/router.md

### Replace
```
## Manual audit override

If the user explicitly asks you to audit an epic, perform the audit for that epic even if the router has not reached `state: audit-pending` yet.

Persist the audit artifacts before replying:

- Ensure `.savepoint/audit/{E##-epic}/snapshot.md` exists. Create a manual snapshot once if needed.
- Write the proposal bundle to `.savepoint/audit/{release}/{E##-epic}/proposals.md`.
- Do not stop at chat-only findings. The filesystem artifact is part of the task output.
```

### With
```
## Manual audit override

If the user explicitly asks you to audit an epic, perform the audit for that epic even if the router has not reached `state: audit-pending` yet.

Persist the audit artifact before replying:

- Write exactly one `.savepoint/releases/{release}/epics/{E##-epic}/E##-Audit.md`.
- Include `## Main Findings`, `## Code Style Review`, and `## Proposed Changes`.
- Keep file-specific `### Target File` / `### Replace` / `### With` blocks under `## Proposed Changes`.
- Do not apply proposals or mark the epic audited until the user says `apply audit`.
```

### Target File
templates/project/.savepoint/router.md

### Replace
```
**Next action (fresh session only):** Confirm `.savepoint/audit/{E##-epic}/snapshot.md` exists. If it is missing while the audit CLI is still unavailable, create one manual snapshot from the known epic scope once; do not search broadly for replacement inputs. Then read the snapshot, read the epic's `Design.md`, and read only the files listed as changed. Write one patch-shaped proposal bundle to `.savepoint/audit/{E##-epic}/proposals.md`:
```

### With
```
**Next action (fresh session only):** Read the epic's `E##-Detail.md`, task files, drift notes, `.savepoint/Design.md`, `AGENTS.md`, and scoped changed files. Write one epic-local audit file to `.savepoint/releases/{release}/epics/{E##-epic}/E##-Audit.md`:
```

### Target File
templates/project/.savepoint/router.md

### Replace
```
- `Design.md` section — merge only the epic delta into project architecture.
- `AGENTS.md` section — refresh Codebase Map entries from changed-module metadata; preserve existing rows.
- `epic-Design.md` section — add "Implemented as:" notes and deltas from the original plan.
- `Quality Review` section — semantic-review findings against the 10 Code Style rules.
```

### With
```
- `## Main Findings` — user-facing AC verification, important drift, and notable risks.
- `## Code Style Review` — checklist against the 10 AGENTS.md code style rules.
- `## Proposed Changes` — admin/apply metadata using `### Target File`, `### Replace`, and `### With`.
```

### Target File
templates/project/.savepoint/router.md

### Replace
```
After proposals are approved, apply approved proposals to live files, mark the epic `Design.md` as `status: audited`, update project `Design.md` `last_audited`, refresh `AGENTS.md` Codebase Map, and advance this router to the next epic state.
```

### With
```
After proposals are approved, apply approved proposals to live files, mark the epic `E##-Detail.md` as `status: audited`, update project `Design.md` `last_audited`, refresh `AGENTS.md` Codebase Map if needed, and advance this router to the next epic state.
```

### Target File
internal/init/validate.go

### Replace
```
var conflictingFiles = []string{
	"AGENTS.md",
}
```

### With
```
var conflictingFiles = []string{
	"AGENTS.md",
	"agent-skills",
}
```

### Target File
internal/init/scaffold_test.go

### Replace
```
func TestScaffold_preservesExistingWithoutForce(t *testing.T) {
```

### With
```
func TestScaffold_overwritesExistingAfterValidation(t *testing.T) {
```

### Target File
internal/init/validate_test.go

### Replace
```
func TestValidateTarget_conflictingFileWithForce(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, true)
	if err != nil {
		t.Fatalf("ValidateTarget() with --force error = %v, want nil", err)
	}
}
```

### With
```
func TestValidateTarget_conflictingFileWithForce(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, true)
	if err != nil {
		t.Fatalf("ValidateTarget() with --force error = %v, want nil", err)
	}
}

func TestValidateTarget_conflictingAgentSkillsDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "agent-skills"), 0755); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, false)
	if err == nil {
		t.Fatal("ValidateTarget() expected error for conflicting agent-skills directory")
	}
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("ValidateTarget() error type = %v, want ErrConflict", err)
	}
}
```

### Target File
internal/init/write.go

### Replace
```
	if err := os.Rename(tmpName, target); err != nil {
		return fmt.Errorf("rename temp to target: %w", err)
	}
```

### With
```
	if err := replaceFile(tmpName, target); err != nil {
		return fmt.Errorf("replace target with temp file: %w", err)
	}
```

### Target File
internal/init/write.go

### Replace
```
	return nil
}
```

### With
```
	return nil
}

func replaceFile(tmpName, target string) error {
	if err := os.Rename(tmpName, target); err == nil {
		return nil
	}

	if _, statErr := os.Stat(target); statErr != nil {
		return os.Rename(tmpName, target)
	}

	backup := target + ".savepoint-bak"
	_ = os.Remove(backup)
	if err := os.Rename(target, backup); err != nil {
		return err
	}
	if err := os.Rename(tmpName, target); err != nil {
		_ = os.Rename(backup, target)
		return err
	}
	return os.Remove(backup)
}
```

### Target File
.savepoint/Design.md

### Replace
```
- **Init command** (`savepoint init`) validates, scaffolds, prints prompt, clipboard, optional install (epic E05).
```

### With
```
- **Init command** (`savepoint init`) validates target directories, scaffolds rendered copies of `templates/project/`, prints the rendered magic prompt, attempts best-effort clipboard copy, and optionally runs `npm install` after scaffolding (v1.1 E07).
```

### Target File
.savepoint/releases/v1.1/epics/E07-init-command/E07-Detail.md

### Replace
```
## Boundaries
```

### With
```
## Implemented As

- `cmd/init.go` owns init argument parsing and delegates execution through `InitRunner`.
- `main.go` embeds `templates/project/` and `templates/prompts/`, then wires validation, scaffold, prompt rendering, clipboard copy, and optional install in sequence.
- `internal/init/` owns validation, scaffold interpolation, atomic writes, prompt rendering, clipboard copy, dependency install, and integration tests.
- Audit found that generated workflow templates need reconciliation with the current epic-local audit workflow before closing.

## Boundaries
```

### Target File
.savepoint/releases/v1.1/epics/E07-init-command/tasks/T006-clipboard.md

### Replace
```
- [ ] Add `internal/init/clipboard.go`
- [ ] Implement `CopyToClipboard(text) error`
- [ ] Detect platform: Windows (clip.exe), macOS (pbcopy), Linux (xclip/xsel)
- [ ] Fall back gracefully if clipboard tools unavailable
- [ ] Log skip or failure without failing init
- [ ] Test clipboard on each platform
- [ ] Run `make build && make test`
```

### With
```
- [x] Add `internal/init/clipboard.go`
- [x] Implement `CopyToClipboard(text) ClipboardResult`
- [x] Detect platform: Windows (clip.exe), macOS (pbcopy), Linux (xclip/xsel)
- [x] Fall back gracefully if clipboard tools unavailable
- [x] Log skip or failure without failing init
- [x] Test clipboard status handling on the current platform
- [x] Run quality gates (`go test ./internal/init ./cmd`, `go build ./...`, `go test ./...`)
```
