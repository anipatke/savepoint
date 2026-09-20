---
type: audit-findings
audited: 2026-05-10
---
# Audit Findings: E16 Pre-Production Refinement

## Main Findings

Applied E16 audit proposals on 2026-05-10. The `upgrade-assets --dry-run` path now resolves casing-variant agent guides through the same managed-block destination logic as write mode, with focused test coverage for `Agents.MD`.

The scaffolded router proposal-format example now matches the required `### Target File` / `### Replace` / `### With` audit apply contract, and the README duplicate status-line artifact was removed.

Architecture drift from T006 has been reconciled: `AGENTS.md`, `.savepoint/Design.md`, and `E16-Detail.md` now document `upgrade-assets`, managed AGENTS.md merge behavior, and the release skeleton refinement. E16 is closed as audited. Residual process note: several task files were already marked `done` with unticked checklist items; no code change was needed for that, but future build agents should keep task checkboxes synchronized before audit handoff.

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

## Proposed Changes

### Target File
internal/init/upgrade.go

### Replace
```go
		if dryRun {
			if _, err := os.Stat(targetPath); err != nil {
				if os.IsNotExist(err) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionSkipped})
				}
			} else {
				existingContent, err := os.ReadFile(targetPath)
				if err != nil {
					return fmt.Errorf("read existing %s: %w", path, err)
				}
				if isSkill {
					if string(existingContent) == string(content) {
						report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
					} else {
						report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
					}
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionMerged})
				}
			}
			return nil
		}
```

### With
```go
		if dryRun {
			if path == "AGENTS.md" {
				dest := FindAgentGuide(absTarget)
				if dest == "" {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
					return nil
				}

				existingContent, err := os.ReadFile(dest)
				if err != nil {
					return fmt.Errorf("read existing %s: %w", path, err)
				}
				block := managedBegin + "\n" + strings.TrimSpace(string(content)) + "\n" + managedEnd
				merged := replaceManagedBlock(string(existingContent), block)
				if merged == string(existingContent) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionMerged})
				}
				return nil
			}

			if _, err := os.Stat(targetPath); err != nil {
				if os.IsNotExist(err) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionSkipped})
				}
			} else {
				existingContent, err := os.ReadFile(targetPath)
				if err != nil {
					return fmt.Errorf("read existing %s: %w", path, err)
				}
				if string(existingContent) == string(content) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
				}
			}
			return nil
		}
```

### Target File
internal/init/upgrade_test.go

### Replace
```go
func TestUpgradeProjectAssets_multipleSkills(t *testing.T) {
```

### With
```go
func TestUpgradeProjectAssets_dryRunUsesCasingVariantAgentGuide(t *testing.T) {
	target := t.TempDir()
	savepointDir := filepath.Join(target, ".savepoint")
	if err := os.MkdirAll(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	variantPath := filepath.Join(target, "Agents.MD")
	testutil.WriteFile(t, variantPath, "# My Guide")

	templates := fstest.MapFS{
		"AGENTS.md": &fstest.MapFile{Data: []byte("# Savepoint Instructions")},
	}

	report, err := UpgradeProjectAssets(templates, target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	for _, e := range report.Actions {
		if e.Path == "AGENTS.md" {
			if e.Action != ActionMerged {
				t.Errorf("dry-run action = %v, want merged", e.Action)
			}
			return
		}
	}
	t.Fatal("AGENTS.md not in dry-run report")
}

func TestUpgradeProjectAssets_multipleSkills(t *testing.T) {
```

### Target File
templates/project/.savepoint/router.md

### Replace
```md
Proposal format:

```md
## Target File

`path/to/file.md`

## Replace

<exact old heading, marker, or section>

## With

<replacement text>
```
```

### With
````md
Proposal format:

```md
### Target File
path/to/file.md

### Replace
```
exact old text
```

### With
```
replacement text
```
```
````

### Target File
README.md

### Replace
```md
**License:** MIT  
**Status:** Recursive Construction (v1 MVP in progress)
us:** Recursive Construction (v1 MVP in progress)
```

### With
```md
**License:** MIT  
**Status:** Recursive Construction (v1 MVP in progress)
```

### Target File
AGENTS.md

### Replace
```md
| `main.go` | CLI entrypoint, --version |
| `cmd/` | CLI command arg parsing and dispatch for init, board, and doctor |
| `internal/init/` | Target validation, scaffold writing from templates |
```

### With
```md
| `main.go` | CLI entrypoint, --version, embedded template wiring for init and upgrade-assets |
| `cmd/` | CLI command arg parsing and dispatch for init, board, doctor, and upgrade-assets |
| `internal/init/` | Target validation, scaffold writing from templates, managed AGENTS.md merge behavior, and safe project asset refresh |
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Init command** (`savepoint init`) validates target directories, scaffolds rendered copies of `templates/project/`, prints the rendered magic prompt, attempts best-effort clipboard copy, and optionally runs `npm install` after scaffolding (v1.1 E07).
```

### With
```md
- **Init command** (`savepoint init`) validates target directories, scaffolds rendered copies of `templates/project/`, merges Savepoint instructions into an existing root agent guide using a managed block while preserving user content and casing variants, creates the initial `.savepoint/releases/v1/epics` skeleton plus release PRD, prints the rendered magic prompt, attempts best-effort clipboard copy, and optionally runs `npm install` after scaffolding (v1.1 E07, refined in E16).
- **Upgrade-assets command** (`savepoint upgrade-assets [dir] [--dry-run] [--force]`) refreshes package-owned `agent-skills/**/SKILL.md` files and the managed block in the root agent guide from embedded templates for existing Savepoint projects, while skipping `.savepoint/PRD.md`, `.savepoint/Design.md`, `.savepoint/releases/**`, and other project state.
```

### Target File
.savepoint/Design.md

### Replace
```md
| `savepoint init`       | Scaffold `.savepoint/`, print magic prompt to stdout + clipboard                  |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `--version` / `--help` | Standard global flags                                                             |
```

### With
```md
| `savepoint init`       | Scaffold `.savepoint/`, merge the managed agent guide block, print magic prompt to stdout + clipboard |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `savepoint upgrade-assets [dir] [--dry-run] [--force]` | Refresh package-owned agent skills and the managed agent-guide block without touching project state |
| `--version` / `--help` | Standard global flags                                                             |
```

### Target File
.savepoint/releases/v1.1/epics/E16-pre-prod-refinement/E16-Detail.md

### Replace
```md
## Boundaries
```

### With
```md
## Implemented as

- `internal/init/agents.go` owns agent-guide casing detection plus managed block insertion/replacement.
- `internal/init/scaffold.go` routes `AGENTS.md` scaffold writes through the managed merge path and now scaffolds the release skeleton from template assets.
- `cmd/upgrade-assets.go` adds the CLI parser for `upgrade-assets [dir] [--dry-run] [--force]`.
- `internal/init/upgrade.go` owns existing-project validation, package-owned asset allowlisting, dry-run reporting, idempotent skill refresh, and managed agent-guide block refresh.
- `templates/project/.savepoint/releases/v1/v1-PRD.md` seeds the release PRD referenced by the scaffolded router.
- `package.json` postinstall prints a notice only; project mutation remains explicit via `savepoint upgrade-assets`.

## Boundaries
```
