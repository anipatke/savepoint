---
id: E13-audit-remediation/T004-infrastructure
status: done
objective: Remove committed binaries from git, add gitignore rules, add golangci-lint config, fix package.json test script
depends_on: []
---

# T004: Infrastructure — Gitignore, Linter Config, Binary Cleanup

## Context Files

- `.gitignore` — missing entries for `savepoint`, `savepoint.exe`, `dist/`, `*.zip`
- `package.json` — `scripts.test` is `"savepoint init"` (misleading)
- Root directory — `savepoint`, `savepoint.exe` binaries tracked in git
- No `.golangci.yml` exists

## Acceptance Criteria

- [ ] `.gitignore` updated with entries for `savepoint`, `savepoint.exe`, `dist/`, `*.zip`
- [ ] `git rm --cached` run on `savepoint`, `savepoint.exe`, `dist/`, `ink-cli-ui-design.zip` (if tracked)
- [ ] `package.json` `scripts.test` changed to `echo "Run 'make test' for Go tests"` or removed
- [ ] `.golangci.yml` added with linters: `unused`, `errcheck`, `staticcheck`, `govet`, `ineffassign`, `gosimple`
- [ ] `golangci-lint run` passes (or only pre-existing issues are flagged)
- [ ] `make build && make test` still passes

## Implementation Plan

- [ ] Add `savepoint`, `savepoint.exe`, `dist/`, `*.zip` to `.gitignore`
- [ ] Run `git rm --cached savepoint savepoint.exe` and remove `dist/` and `ink-cli-ui-design.zip` from tracking
- [ ] Update `package.json` test script
- [ ] Create `.golangci.yml` with recommended linters
- [ ] Run `golangci-lint run` and review output
- [ ] Run `make build && make test`