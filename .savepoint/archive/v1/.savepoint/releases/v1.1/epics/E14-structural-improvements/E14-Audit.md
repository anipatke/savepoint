---
type: audit-findings
audited: 2026-05-05
---

# Audit Findings: E14 Structural Improvements

## Main Findings

Audit proposals were applied and E14 is closed as audited. The board model is grouped into focused embedded state structs, dead board helpers were removed, board and doctor now use consumer-side data-access interfaces, doctor orphan traversal goes through `Discover.ListRootDirs`, audit tab hidden-section filtering uses exact heading matches, `TaskStatus` has been removed in favor of the canonical column/status values, and shared Go test fixtures now live in `internal/testutil`.

The shell tokenizer follow-up was applied: empty quoted arguments such as `-run ""` are preserved, trailing backslashes are retained, and the branch tests were added. Router state was advanced to `epic-design` for `E15-hardening`, `E14-Detail.md` now has `status: audited`, and `.savepoint/Design.md` has `last_audited: v1.1/E14-structural-improvements`.

All seven task files include `## Context Files`; no `## Drift Notes` sections were found. The remaining review note is that `package.json` contains a version bump outside the explicit E14 task scope; it was left untouched as apparent release housekeeping.

## Code Style Review

- [x] One job per file - E14 keeps production responsibilities in board/data/doctor and moves shared test fixtures to `internal/testutil`.
- [x] One job per function - the touched functions remain small and named around a single behavior.
- [x] Test branches - `splitCommand` now covers empty quoted arguments and trailing escapes.
- [x] Types are documentation - consumer-side interfaces and `ColumnType` document the intended contracts.
- [x] Build, don't speculate - the implementation maps to the listed audit findings.
- [x] Errors at boundaries - data discovery and quality-gate execution continue returning boundary errors.
- [x] One source of truth - task status constants are consolidated under `ColumnType`.
- [x] Comments explain WHY - new comments describe contracts or compatibility, not line-by-line mechanics.
- [x] Content in data files - no new product copy is embedded in logic.
- [ ] Small diffs - residual note: the `package.json` version bump remains outside the E14 task scope and was left untouched.

## Proposed Changes

### Target File
internal/doctor/gates.go

### Replace
```go
func splitCommand(command string) []string {
	var parts []string
	current := strings.Builder{}
	inDoubleQuote := false
	inSingleQuote := false

	for i := 0; i < len(command); i++ {
		c := command[i]

		if inSingleQuote {
			if c == '\'' {
				inSingleQuote = false
			} else {
				current.WriteByte(c)
			}
			continue
		}

		if inDoubleQuote {
			if c == '\\' && i+1 < len(command) {
				next := command[i+1]
				if next == '"' || next == '\\' || next == '$' || next == '`' {
					i++
					current.WriteByte(command[i])
					continue
				}
			}
			if c == '"' {
				inDoubleQuote = false
				continue
			}
			current.WriteByte(c)
			continue
		}

		if c == '\\' && i+1 < len(command) {
			i++
			current.WriteByte(command[i])
			continue
		}
		if c == '"' {
			inDoubleQuote = true
			continue
		}
		if c == '\'' {
			inSingleQuote = true
			continue
		}
		if c == ' ' {
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			continue
		}
		current.WriteByte(c)
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
```

### With
```go
func splitCommand(command string) []string {
	var parts []string
	current := strings.Builder{}
	inDoubleQuote := false
	inSingleQuote := false
	tokenStarted := false

	flush := func() {
		if tokenStarted || current.Len() > 0 {
			parts = append(parts, current.String())
			current.Reset()
			tokenStarted = false
		}
	}

	for i := 0; i < len(command); i++ {
		c := command[i]

		if inSingleQuote {
			tokenStarted = true
			if c == '\'' {
				inSingleQuote = false
			} else {
				current.WriteByte(c)
			}
			continue
		}

		if inDoubleQuote {
			tokenStarted = true
			if c == '\\' && i+1 < len(command) {
				next := command[i+1]
				if next == '"' || next == '\\' || next == '$' || next == '`' {
					i++
					current.WriteByte(command[i])
					continue
				}
			}
			if c == '"' {
				inDoubleQuote = false
				continue
			}
			current.WriteByte(c)
			continue
		}

		switch c {
		case '\\':
			tokenStarted = true
			if i+1 < len(command) {
				i++
				current.WriteByte(command[i])
			} else {
				current.WriteByte(c)
			}
		case '"':
			tokenStarted = true
			inDoubleQuote = true
		case '\'':
			tokenStarted = true
			inSingleQuote = true
		case ' ', '\t', '\n', '\r':
			flush()
		default:
			tokenStarted = true
			current.WriteByte(c)
		}
	}
	flush()
	return parts
}
```

### Target File
internal/doctor/gates_test.go

### Replace
```go
		{"echo 'hello world'", []string{"echo", "hello world"}},
		{"echo \"hello \\\"world\\\"\"", []string{"echo", "hello \"world\""}},
		{"echo hello\\ world", []string{"echo", "hello world"}},
		{"echo 'it''s'", []string{"echo", "its"}},
		{"", nil},
```

### With
```go
		{"echo 'hello world'", []string{"echo", "hello world"}},
		{"echo \"hello \\\"world\\\"\"", []string{"echo", "hello \"world\""}},
		{"echo hello\\ world", []string{"echo", "hello world"}},
		{"echo 'it''s'", []string{"echo", "its"}},
		{"go test -run \"\" ./...", []string{"go", "test", "-run", "", "./..."}},
		{"printf ''", []string{"printf", ""}},
		{"echo trailing\\", []string{"echo", "trailing\\"}},
		{"", nil},
```

### Target File
.savepoint/router.md

### Replace
```yaml
epic: E14
```

### With
```yaml
epic: E14-structural-improvements
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Audit remediation baseline** (v1.1 E13) centralizes frontmatter/body splitting and line-ending normalization in `internal/data`, uses typed sentinel errors for doctor repair suggestions, applies a configurable `quality_gates.gate_timeout`, removes tracked build artifacts from source control, adds `.golangci.yml`, and moves board filesystem reads/writes behind Bubble Tea command messages while preserving direct file I/O inside command helpers.
- **Agent audit workflow** is skill-driven, not a CLI pipeline. At `audit-pending`, a fresh audit agent writes one epic-local `E##-Audit.md`; the user reviews its Audit tab, then asks an agent to apply the admin proposal blocks, update the visible audit findings to reflect the applied outcome, and close the epic.
```

### With
```md
- **Audit remediation baseline** (v1.1 E13) centralizes frontmatter/body splitting and line-ending normalization in `internal/data`, uses typed sentinel errors for doctor repair suggestions, applies a configurable `quality_gates.gate_timeout`, removes tracked build artifacts from source control, adds `.golangci.yml`, and moves board filesystem reads/writes behind Bubble Tea command messages while preserving direct file I/O inside command helpers.
- **Structural improvement baseline** (v1.1 E14) groups board `Model` fields into focused embedded state structs, defines consumer-side board/doctor data-access interfaces, routes doctor orphan discovery through `Discover.ListRootDirs`, renders audit-tab hidden sections via exact heading matches, improves quality-gate shell tokenization for quoted and escaped arguments, removes the separate `TaskStatus` enum in favor of `ColumnType`, and adds `internal/testutil` for shared Go test fixtures.
- **Agent audit workflow** is skill-driven, not a CLI pipeline. At `audit-pending`, a fresh audit agent writes one epic-local `E##-Audit.md`; the user reviews its Audit tab, then asks an agent to apply the admin proposal blocks, update the visible audit findings to reflect the applied outcome, and close the epic.
```

### Target File
AGENTS.md

### Replace
```md
| `internal/data/` | Task/router models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery, canonical write helpers |
| `internal/styles/` | Atari-Noir palette, TUI styles |
```

### With
```md
| `internal/data/` | Task/router models, frontmatter parsing/splitting, lifecycle validation/defaulting, discovery including root-dir traversal, unified task status constants, canonical write helpers |
| `internal/testutil/` | Shared Go test fixtures and filesystem helpers for internal package tests |
| `internal/styles/` | Atari-Noir palette, TUI styles |
```

### Target File
.savepoint/releases/v1.1/epics/E14-structural-improvements/E14-Detail.md

### Replace
```md
**Out of scope:**
- New UI features or commands
- Audit Phase 3 hardening items (separate epic E15)
```

### With
```md
**Out of scope:**
- New UI features or commands
- Audit Phase 3 hardening items (separate epic E15)

## Implemented As

- `Model` grouping landed as embedded state structs in `internal/board/model.go`, preserving existing call sites while separating view, data, navigation, epic, release, and data-access state.
- Consumer-side interfaces landed in `internal/board/interfaces.go` and `internal/doctor/interfaces.go`; concrete `internal/data` types satisfy them without data-layer changes.
- `TaskStatus` was removed. `ColumnType` is the canonical task lifecycle type, while `data.Task.Status` remains a string mirror for parsed frontmatter and board glyph compatibility.
- `internal/testutil` was added for shared test filesystem and fixture helpers across board, doctor, data, and init tests.
- Audit follow-up before close: preserve empty quoted shell arguments in `splitCommand` and normalize router epic from `E14` to `E14-structural-improvements`.
```
