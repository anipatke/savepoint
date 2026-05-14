---
id: E20-strategic-npm-packaging/T002-binary-format-validation
status: done
objective: Validate generated Windows, Linux, and macOS binaries before npm packaging or publish
depends_on: [T001-platform-package-architecture]
complexity_tier: medium
complexity_reason: "Extends build validation across binary formats without changing runtime CLI behavior"
---

# T002: Binary Format Validation

## Problem

D009 was caught only after Windows tried to execute a Linux ELF binary named `savepoint.exe`. The strategic package flow needs explicit validation for every supported binary format before artifacts can be packed or published.

## Context Files

- `internal/buildtool/main.go`
- `internal/buildtool/main_test.go`
- `package.json`

## Acceptance Criteria

- [x] Windows binaries are rejected unless they start with a PE `MZ` signature
- [x] Linux binaries are rejected unless they start with an ELF signature
- [x] macOS binaries are rejected unless they start with a Mach-O signature for the target architecture
- [x] Validation runs during npm artifact generation
- [x] Tests cover accepted and rejected headers for PE, ELF, and Mach-O
- [x] Validation errors name the target and artifact path clearly

## Implementation Plan

- [x] Generalize the current Windows header guard into target-aware binary validation
- [x] Add ELF and Mach-O header detection helpers
- [x] Call validation for every npm build target
- [x] Add table-driven tests for valid and invalid format headers
- [x] Ensure failure messages are actionable for CI logs

## Context Log

- Added `internal/buildtool/binary_format.go` with `validateBinaryFormat` dispatching to PE, ELF, and Mach-O header checks; ELF and Mach-O also verify the arch (e_machine / cpu_type) so a cross-target mix is caught.
- `buildNPM` now calls `validateBinaryFormat` for every target, replacing the windows-only `requireWindowsExecutable` guard (deleted).
- New `internal/buildtool/binary_format_test.go` covers accept/reject across all six targets, format swaps, arch swaps, unknown OS, missing file, and truncated headers.
- `go test ./internal/buildtool/` passes; pre-existing `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` failure in root package is unrelated to this task.
