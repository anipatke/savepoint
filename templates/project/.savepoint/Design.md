---
type: project-design
status: active
last_audited: never
---

# {{PROJECT_NAME}} — System Architecture

> Project-level architecture. Audit-kept fresh: every epic's audit step merges its delta into this document.
>
> **Visual identity** lives separately in `.savepoint/visual-identity.md` and is loaded only for TUI/theme/visual tasks.

## 1. Architecture model

<!-- High-level architecture pattern. -->

## 2. Directory layout

<!-- Expected file structure. -->

## 3. Hierarchy semantics

<!-- Definitions of release, epic, task, sub-task. -->

## 4. Status model & gates

<!-- How work moves through states. -->

## 5. Dependencies

<!-- How tasks depend on each other. -->

## 6. CLI surface

<!-- Commands and flags if applicable. -->

## 7. Agent audit workflow

<!-- Savepoint audit is agent-led and skill-driven, not a CLI pipeline. At epic close, a fresh audit agent writes one epic-local `E##-Audit.md` with exactly these user-facing sections: `## Main Findings` and `## Code Style Review`. File-specific `### Target File` / `### Replace` / `### With` blocks belong under a separate `## Proposed Changes` admin section so the TUI Audit tab can omit them. -->

## 8. Testing strategy

<!-- How the project is tested. -->

## 9. Release versioning

<!-- Version scheme. -->
