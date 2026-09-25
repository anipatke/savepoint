---
 type: epic-design
 status: planned
 ---

 # Epic E02: Domain Phase Model

 ## Purpose

 Extend the domain layer to support the build/test/audit phase model within the `in_progress` status. Simplify the config model to theme-only and collapse the router state machine from 6 to 3 states.

 ## What this epic adds

 - `TaskPhase` type with constants and transition rules.
 - Phase validation in task frontmatter parsing.
 - Phase-to-accent color mapping for TUI rendering.
 - Simplified config model (theme-only, no quality gates or audit config).
 - Collapsed router domain (3 states: planning, building, reviewing).

 ## Definition of Done

 - `status.ts` defines `TaskPhase`, phase transitions, and `PHASE_ACCENTS`.
 - `task.ts` validates `phase?: TaskPhase` in frontmatter.
 - `config.ts` contains only `ThemeConfig` and `SavepointConfig` (no quality_gates, audit, verify_strict).
 - `router.ts` defines 3 states (`planning`, `building`, `reviewing`).
 - All existing domain tests compile and pass (or are updated).
 - No references to old 6-state model remain in domain code.

 ## Components and files

 | Path | Purpose |
 |------|---------|
 | `src/domain/status.ts` | TaskPhase type, phase transitions, accent mapping |
 | `src/domain/task.ts` | Phase frontmatter field and validation |
 | `src/domain/config.ts` | Simplified theme-only config model |
 | `src/domain/router.ts` | 3-state router domain |
 | `test/domain/status.test.ts` | Phase transition tests |
 | `test/domain/task.test.ts` | Phase frontmatter validation tests |
 | `test/domain/config.test.ts` | Simplified config tests |
 | `test/domain/router.test.ts` | 3-state router tests |
