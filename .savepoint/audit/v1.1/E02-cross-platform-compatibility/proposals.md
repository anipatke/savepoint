# Audit Proposals: v1.1 E02 Cross-Platform Compatibility

## Target File

`.savepoint/Design.md`

## Replace

```md
## 12. Distribution & build

> Audit note: the live repository is transitioning from the documented TypeScript/Node implementation to a Go module (`github.com/opencode/savepoint`). The architecture document still contains substantial TypeScript-era implementation detail and should be reconciled as Go epics are audited.

- **License:** MIT.
- **Install:** primary `npx savepoint init`, persistent `npm i -g savepoint` → `savepoint`.
- **Runtime:** Node 20.10+ LTS, ESM-only, no native deps. macOS / Linux / Windows-Terminal.
- **Repo:** single package. TypeScript strict. `tsup` build → `dist/`. Bin `dist/cli.js` shebanged.
- **No telemetry.** Ever.
```

## With

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

## Target File

`AGENTS.md`

## Replace

```md
| `main.go`                            | CLI Entrypoint and root command wiring                                                               |
```

## With

```md
| `main.go`                            | CLI entrypoint, root command wiring, and `--version` handling via build-time version injection        |
```

## Target File

`.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`

## Insert After

```md
## Architectural notes

- All Go code uses `path/filepath` — no syscalls are platform-specific
- Bubble Tea handles terminal differences across platforms
- No Go build tags or platform-specific files needed
- This is purely a build-tooling and release-workflow epic
```

## With

```md
## Implemented as

- `Makefile` now has `build-linux`, `build-darwin`, `build-all`, `dist`, and `smoke-test` targets.
- `VERSION` defaults from `git describe --tags --abbrev=0` with a `v0.0.0` fallback and is injected into the binary with `-ldflags "-X main.version=$(VERSION)"`.
- `main.go` handles `--version` before launching the TUI, which gives smoke tests a non-interactive success path.
- `dist/` remains ignored by git and contains both raw cross-compiled binaries and versioned `.tar.gz` archives for Linux and Darwin builds.
- Windows zip packaging was not implemented in this epic; the final implementation covers the Linux and Darwin targets named in the epic's build-all definition of done.
```

## Target File

`.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`

## Replace

```md
- [x] Makefile replaced with Go-native commands (no `rm`, `cp`, `mkdir`)
```

## With

```md
- [ ] Makefile replaced with Go-native commands (no `rm`, `cp`, `mkdir`)
```

## Quality Review

## Must Fix Before Close

- `Makefile`: `clean` still runs `rm -rf dist/`, which violates T001 acceptance criteria: "`make clean` removes the binary without using `rm -f` or other Unix-only commands" and "No shell commands in Makefile are platform-dependent."
- `Makefile`: `build-linux` and `build-darwin` still use `mkdir -p` and inline `GOOS=... GOARCH=...` environment assignment. That is fine for Linux, macOS, Git Bash, and WSL, but it is not a shell-agnostic Makefile implementation.
- `Makefile`: `dist` relies on `tar`, so packaging is not yet toolchain-native or Windows-shell portable.
- `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/tasks/T004-smoke-tests-and-artifacts.md`: the task says `go test ./...` was not run even though the task completion protocol requires the full quality-gate suite before `status: done`.

## Carry Forward

- Windows `.zip` artifacts were deferred in T004. Keep that as a known scope gap unless the project narrows E02 explicitly to Linux and Darwin distribution.
- Consider moving cross-platform build and archive logic into a tiny Go helper if the project wants true shell portability while keeping `make` as the human entrypoint.

## Already Fixed

- `main.go` has a headless `--version` path, so smoke testing no longer needs to start the TUI.
- Linux and Darwin raw binary targets exist for amd64 and arm64.
- `AGENTS.md` already documents `build-all`, `dist`, and `smoke-test`.
