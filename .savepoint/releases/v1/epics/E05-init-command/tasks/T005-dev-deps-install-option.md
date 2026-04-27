---
id: E05-init-command/T005-dev-deps-install-option
status: planned
objective: "Support the optional one-shot dev-dependency installation path promised by the init command contract."
depends_on:
  - E05-init-command/T001-init-cli-contract
---

# T005: dev-deps-install-option

## Implementation Plan

- [ ] Create `src/init/dev-deps.ts` to build the package-manager-specific command needed for the optional Savepoint dev-dependency install.
- [ ] Detect the package manager from lockfiles or package metadata without modifying unrelated project files.
- [ ] Respect init flags that force install, skip install, or leave install as an interactive prompt decision.
- [ ] Run the install command only after scaffold writing succeeds and surface install failures separately from scaffold failures.
- [ ] Add tests for command construction, package manager detection, flag precedence, skipped install, and failed install reporting.
