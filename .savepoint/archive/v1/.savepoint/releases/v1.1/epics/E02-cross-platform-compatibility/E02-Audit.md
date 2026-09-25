---
type: audit-findings
audited: 2026-05-02
---

# Audit Findings: E02 Cross-Platform Compatibility

## Main Findings

E02 found that the cross-platform build work moved Savepoint further into its Go-based distribution path, including version handling, Linux/Darwin build targets, archive generation, and smoke-test coverage.

The audit also identified that the original Makefile still retained shell-portability gaps at the time of review, including Unix-specific cleanup/archive commands and incomplete Windows packaging. Those details are retained under Proposed Changes for apply/close testing, while this visible section stays focused on the user-facing audit outcome.

## Code Style Review


- [ ] One job per file
- [ ] One-sentence functions
- [ ] Test branches
- [ ] Types are documentation
- [ ] Build, don't speculate
- [ ] Errors at boundaries
- [ ] One source of truth
- [ ] Comments explain WHY
- [ ] Content in data files
- [ ] Small diffs

## Proposed Changes

### Target File

`.savepoint/Design.md`

### Replace

```md
## 12. Distribution & build

> Audit note: the live repository is transitioning from the documented TypeScript/Node implementation to a Go module (`github.com/opencode/savepoint`). The architecture document still contains substantial TypeScript-era implementation detail and should be reconciled as Go epics are audited.

- **License:** MIT.
- **Install:** primary `npx savepoint init`, persistent `npm i -g savepoint` → `savepoint`.
- **Runtime:** Node 20.10+ LTS, ESM-only, no native deps. macOS / Linux / Windows-Terminal.
- **Repo:** single package. TypeScript strict. `tsup` build → `dist/`. Bin `dist/cli.js` shebanged.
- **No telemetry.** Ever.
```

### With

```md
## 12. Distribution & build

> Audit note: the live repository is now a Go module (`github.com/opencode/savepoint`). Remaining TypeScript-era distribution details should be removed as Go epics are audited.

- **License:** MIT.
- **Runtime:** Go CLI binary. Source builds with `go build`; tests run with `go test ./...`.
- **Local build:** `make build` builds `savepoint` and injects `main.version` from `VERSION`.
- **Cross-platform builds:** `make build-all` cross-compiles linux-amd64, linux-arm64, darwin-amd64, and darwin-arm64 raw binaries into `dist/{platform}-{arch}/savepoint`.
- **Artifacts:** `make dist` creates versioned `.tar.gz` archives in `dist/` for the Linux and Darwin targets.
- **Smoke validation:** `make smoke-test` builds the local binary and runs `./savepoint --version` as a headless exit-0 check.
- **No telemetry.** Ever.
```

---

### Target File

`AGENTS.md`

### Replace

```md
| `main.go`                            | CLI Entrypoint and root command wiring                                                               |
```

### With

```md
| `main.go`                            | CLI entrypoint, root command wiring, and `--version` handling via build-time version injection        |
```

---

### Target File

`.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`

### Insert After

```md
## Architectural notes
```

### With

```md
## Implemented as

- `Makefile` now has `build-linux`, `build-darwin`, `build-all`, `dist`, and `smoke-test` targets.
- `VERSION` defaults from `git describe --tags --abbrev=0` with a `v0.0.0` fallback and is injected into the binary with `-ldflags "-X main.version=$(VERSION)"`.
- `main.go` handles `--version` before launching the TUI, which gives smoke tests a non-interactive success path.
- `dist/` remains ignored by git and contains both raw cross-compiled binaries and versioned `.tar.gz` archives for Linux and Darwin builds.
- Windows zip packaging was not implemented in this epic; the final implementation covers the Linux and Darwin targets named in the epic's build-all definition of done.
```

---

### Target File

`.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`

### Replace

```md
- [x] Makefile replaced with Go-native commands (no `rm`, `cp`, `mkdir`)
```

### With

```md
- [ ] Makefile replaced with Go-native commands (no `rm`, `cp`, `mkdir`)
```

---
