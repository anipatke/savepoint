---
type: audit-findings
audited: 2026-05-06
---
# Audit Findings: E15 Phase 3 - Hardening

## Main Findings

Applied the E15 audit proposals and closed the epic as audited.

- Fixed the Windows archive packaging issue: Windows tarballs now preserve `savepoint.exe` as the archive member name instead of extracting as `savepoint`.
- Added a regression test covering the Windows archive member name.
- Reconciled E15 documentation drift in `Design.md`, `AGENTS.md`, and `E15-Detail.md`.
- Marked `E15-Detail.md` as `status: audited`.
- Updated `Design.md` `last_audited` to `v1.1/E15-hardening`.
- Advanced router state to `epic-design` with no active epic because there are no later unarchived epics in v1.1.

Verification after applying proposals:

- `go test ./internal/buildtool`
- `make build`
- `make test`
- `make ci`

No remaining audit blockers are known. The `## Proposed Changes` section remains below as the apply trace.

## Code Style Review

- [x] One job per file - E15 changes stayed within board rendering/debug, data parsing/fuzzing, buildtool packaging, CI, and task metadata.
- [x] One job per function - new helpers such as `executableName`, `writeChecksums`, `stripDebugFlag`, `debugf`, and abbreviation checks are narrowly scoped.
- [x] Test branches - Windows build output and Windows archive member naming are now covered.
- [x] Types document intent - existing Go types and small helper APIs remain explicit enough for this scope.
- [x] Build only what is needed - no speculative runtime features were added beyond E15 scope.
- [x] Handle errors at boundaries - buildtool, filesystem, checksum, fuzz, and parser errors are propagated with context.
- [x] One source of truth - target lists and executable naming are centralized in buildtool; the proposed archive fix reuses that source.
- [x] Comments explain why - new comments document abbreviation suppression and audit-tab hidden sections rather than restating obvious code.
- [x] Content lives in data - workflow and architecture copy remain in markdown; no new product copy was embedded in logic.
- [x] Small diffs - most changes are localized, test-backed additions.

## Proposed Changes

### Target File
internal/buildtool/main.go

### Replace
```go
		if err := writeTarGz(archive, source, "savepoint"); err != nil {
```

### With
```go
		if err := writeTarGz(archive, source, executableName(target.os)); err != nil {
```

### Target File
internal/buildtool/main_test.go

### Replace
```go
import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)
```

### With
```go
import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)
```

### Target File
internal/buildtool/main_test.go

### Replace
```go
func TestExecutableName(t *testing.T) {
	if got := executableName("windows"); got != "savepoint.exe" {
		t.Errorf("executableName(windows) = %q, want savepoint.exe", got)
	}
	if got := executableName("linux"); got != "savepoint" {
		t.Errorf("executableName(linux) = %q, want savepoint", got)
	}
	if got := executableName("darwin"); got != "savepoint" {
		t.Errorf("executableName(darwin) = %q, want savepoint", got)
	}
}
```

### With
```go
func TestExecutableName(t *testing.T) {
	if got := executableName("windows"); got != "savepoint.exe" {
		t.Errorf("executableName(windows) = %q, want savepoint.exe", got)
	}
	if got := executableName("linux"); got != "savepoint" {
		t.Errorf("executableName(linux) = %q, want savepoint", got)
	}
	if got := executableName("darwin"); got != "savepoint" {
		t.Errorf("executableName(darwin) = %q, want savepoint", got)
	}
}

func TestWriteTarGzPreservesWindowsExecutableName(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "savepoint.exe")
	if err := os.WriteFile(source, []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(dir, "savepoint-windows-amd64.tar.gz")
	if err := writeTarGz(archive, source, executableName("windows")); err != nil {
		t.Fatalf("writeTarGz: %v", err)
	}

	f, err := os.Open(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	header, err := tr.Next()
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != "savepoint.exe" {
		t.Fatalf("archive member = %q, want savepoint.exe", header.Name)
	}

	content, err := io.ReadAll(tr)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "binary" {
		t.Fatalf("archive content = %q, want binary", content)
	}
}
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Structural improvement baseline** (v1.1 E14) groups board `Model` fields into focused embedded state structs, defines consumer-side board/doctor data-access interfaces, routes doctor orphan discovery through `Discover.ListRootDirs`, renders audit-tab hidden sections via exact heading matches, improves quality-gate shell tokenization for quoted and escaped arguments, removes the separate `TaskStatus` enum in favor of `ColumnType`, and adds `internal/testutil` for shared Go test fixtures.
- **Agent audit workflow** is skill-driven, not a CLI pipeline. At `audit-pending`, a fresh audit agent writes one epic-local `E##-Audit.md`; the user reviews its Audit tab, then asks an agent to apply the admin proposal blocks, update the visible audit findings to reflect the applied outcome, and close the epic.
```

### With
```md
- **Structural improvement baseline** (v1.1 E14) groups board `Model` fields into focused embedded state structs, defines consumer-side board/doctor data-access interfaces, routes doctor orphan discovery through `Discover.ListRootDirs`, renders audit-tab hidden sections via exact heading matches, improves quality-gate shell tokenization for quoted and escaped arguments, removes the separate `TaskStatus` enum in favor of `ColumnType`, and adds `internal/testutil` for shared Go test fixtures.
- **Hardening baseline** (v1.1 E15) adds board render/layout benchmarks, data frontmatter fuzz targets, debug logging via CLI `--debug` or `SAVEPOINT_DEBUG`, abbreviation-aware task checklist sentence splitting, root test package isolation, documented audit-tab hidden-section allowlisting, repo-local CI, `make ci`, distribution SHA256 checksums, and Windows amd64/arm64 build outputs.
- **Agent audit workflow** is skill-driven, not a CLI pipeline. At `audit-pending`, a fresh audit agent writes one epic-local `E##-Audit.md`; the user reviews its Audit tab, then asks an agent to apply the admin proposal blocks, update the visible audit findings to reflect the applied outcome, and close the epic.
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Cross-platform builds:** `make build-all` cross-compiles linux-amd64, linux-arm64, darwin-amd64, and darwin-arm64 raw binaries into `dist/{platform}-{arch}/savepoint`.
- **Artifacts:** `make dist` creates versioned `.tar.gz` archives in `dist/` for the Linux and Darwin targets using Go archive APIs, not shell `tar`.
```

### With
```md
- **Cross-platform builds:** `make build-all` cross-compiles linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64, and windows-arm64 raw binaries into `dist/{platform}-{arch}/savepoint` or `savepoint.exe` for Windows. `make ci` runs the repo-local verification sequence used by CI.
- **Artifacts:** `make dist` creates versioned `.tar.gz` archives in `dist/` for Linux, Darwin, and Windows targets using Go archive APIs, not shell `tar`, and writes SHA256 hashes to `dist/checksums.txt`.
```

### Target File
AGENTS.md

### Replace
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, async update I/O commands, shared board utilities |
```

### With
```md
| `internal/board/` | TUI board, overlays, epic sidebar, Next Activity line, router priority key, detail checklist rendering, status glyphs, forced color profile, debug logging hooks, async update I/O commands, shared board utilities |
```

### Target File
AGENTS.md

### Replace
```md
| `internal/buildtool/` | Makefile helper, cross-compile, archives |
```

### With
```md
| `internal/buildtool/` | Makefile helper, cross-compile including Windows targets, archives, distribution checksums |
```

### Target File
.savepoint/releases/v1.1/epics/E15-hardening/E15-Detail.md

### Replace
```md
## Components

| Module | Purpose |
|--------|---------|
| `internal/board/view_test.go` | Add render benchmarks |
| `internal/data/parser_test.go` | Add fuzz targets |
| `cmd/main.go` | Add --debug / SAVEPOINT_DEBUG |
| `internal/buildtool/main.go` | Add checksums, Windows targets |
| `internal/board/detail.go` | Fix abbreviation handling |
| `internal/board/epic_panel.go` | Extract allowlist constant |
| `agent_skills_test.go` | Move to cmd_test package |
| `Makefile` | Add `ci` target for local verification |
| `.github/workflows/ci.yml` | Add repo CI workflow |

## Boundaries
```

### With
```md
## Components

| Module | Purpose |
|--------|---------|
| `internal/board/view_test.go` | Add render and layout benchmarks |
| `internal/board/card_test.go` | Add card render benchmarks |
| `internal/board/column_test.go` | Add column render benchmarks |
| `internal/data/fuzz_test.go` | Add frontmatter and split-body fuzz targets |
| `main.go` | Add global --debug parsing and SAVEPOINT_DEBUG activation |
| `internal/board/debug.go` | Add board debug logging helper |
| `internal/board/board.go` | Log board initialization under debug mode |
| `internal/board/watch.go` | Log watcher and reload activity under debug mode |
| `internal/board/update.go` | Log update dispatch and preserve focused-card visibility |
| `internal/buildtool/main.go` | Add checksums and Windows targets |
| `internal/board/detail.go` | Fix abbreviation handling |
| `internal/board/epic_panel.go` | Document audit hidden-section allowlist |
| `agent_skills_test.go` | Move to external root test package |
| `Makefile` | Add `ci` target for local verification |
| `.github/workflows/ci.yml` | Add repo CI workflow |

## Implemented As

- Debug activation is in root `main.go` because global flags are parsed before dispatching into `cmd/`.
- Fuzz targets live in `internal/data/fuzz_test.go` rather than `internal/data/parser_test.go`.
- Fuzzing found and fixed `SplitFrontmatterBody` delimiter-offset handling and CR line-ending normalization.
- Root `agent_skills_test.go` remains at the repository root but now uses external package `main_test`, which removes the root package coupling without moving fixture-relative paths.
- CI was included as repo-local guardrail scope for E15 via `T008-ci-and-release-automation`.

## Boundaries
```
