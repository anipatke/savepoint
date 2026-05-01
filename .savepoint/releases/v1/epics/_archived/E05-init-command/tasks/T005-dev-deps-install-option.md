---
id: E05-init-command/T005-dev-deps-install-option
status: done
objective: "Support the optional one-shot dev-dependency installation path promised by the init command contract."
depends_on:
  - E05-init-command/T001-init-cli-contract
---

# T005: dev-deps-install-option

## Implementation Plan

- [x] Create `src/init/dev-deps.ts` to build the package-manager-specific command needed for the optional Savepoint dev-dependency install.
- [x] Detect the package manager from lockfiles or package metadata without modifying unrelated project files.
- [x] Respect init flags that force install, skip install, or leave install as an interactive prompt decision.
- [x] Run the install command only after scaffold writing succeeds and surface install failures separately from scaffold failures.
- [x] Add tests for command construction, package manager detection, flag precedence, skipped install, and failed install reporting.
- [x] Update `.savepoint/metrics/context-bench.md` with this task's measured context before setting status to review.

## Context Log

- Files read: 12
- Estimated input tokens: ~6,269
- Notes: Read T001 and two existing init modules + tests to align with codebase patterns. runInstallCommand is exported separately so tests can drive it with node without a real package manager.
