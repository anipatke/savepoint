---
type: epic-design
status: active
---

# Epic E10: release-validation

## Purpose

Run the final validation matrix for `0.1.0`, resolve release-blocking issues, and prepare the public release cut.

## What this epic adds

- Manual e2e test matrix across Claude, Cursor, Gemini, and Aider where available.
- Cross-platform checks for macOS, Linux, and Windows Terminal where available.
- Fresh-project `npx savepoint init` validation.
- Full command walkthrough: init, board, audit, and doctor.
- Final version check for `0.1.0`.
- Release notes.
- Publish-readiness checklist.

## Components and files

Expected files introduced or extended by this epic:

| Path                        | Purpose                                           |
| --------------------------- | ------------------------------------------------- |
| `docs/release-checklist.md` | Final release checklist and results.              |
| `docs/e2e-matrix.md`        | Manual agent/platform matrix results.             |
| `CHANGELOG.md`              | Initial release notes.                            |
| `package.json`              | Final version metadata if needed.                 |
| `.savepoint/audit-log.md`   | Audit history if audits were skipped or recorded. |

## Architectural delta

Before this epic, the package should be feature-complete for v0.1.0. After this epic, the release has documented validation evidence and can be published with known residual risks.

This epic should prefer fixing release blockers in the smallest responsible place over broad refactors.

## Boundaries

In scope:

- Manual matrix execution and recording.
- Cross-platform smoke validation.
- Final bug fixes required to pass release criteria.
- Release notes and version confirmation.

Out of scope:

- New features.
- v0.2.0 semantic review.
- File watching.
- Search.
- MCP server work.

## Quality gates

- All automated tests pass.
- Packed install smoke passes.
- Manual e2e matrix passes against at least three of four target agents.
- Known limitations are documented.

## Design constraints

- Keep late fixes narrow and auditable.
- Do not relax quality gates to ship.
- Record evidence, not vibes.
