---
id: v2/D007-absolute-v2-source-path-containment
release: v2
status: open
severity: high
title: "Absolute V2 source paths bypass project-root containment"
---

# D007: Absolute V2 source paths bypass project-root containment

## Symptom

`internal/data/write.go:resolveV2SourcePath()` checks a relative source path
with `filepath.Rel` and rejects paths that escape the project root, but returns
an absolute `source.Path` after `filepath.Clean` without applying the same
containment check.

## Expected Behavior

Every V2 source path, whether relative or absolute, must resolve inside the
record's project root before a write or read is attempted. Paths outside the
root must fail with the same explicit unsafe-path diagnostic.

## Reproduction

1. Construct a V2 source document with `ProjectRoot` set to a project root and
   `Path` set to an absolute path in a separate directory.
2. Pass the document through the V2 write-path source resolver.
3. Observe that the absolute path is accepted instead of being rejected as
   outside the project root.

## Impact

An in-memory or caller-supplied absolute source path can direct V2 record
writes outside the project being operated on, bypassing the relative-path
safety invariant.

## Fix Plan

Confirmed approach: normalize absolute paths and apply the same
root-containment check used for relative paths, including the platform
separator boundary, so an absolute path outside the root is rejected with
`ErrV2UnsafePath`. Keep the current pass-through for an absolute path with an
empty `ProjectRoot` (there is no root to contain against). Add tests for an
absolute path inside the root, an absolute path outside it, and equivalent
relative escape attempts.

## Acceptance Criteria

- [ ] Absolute paths inside the project root remain supported safely.
- [ ] Absolute paths outside the project root are rejected with
      `ErrV2UnsafePath`.
- [ ] Relative and absolute path tests enforce one containment invariant.

## Resolution Notes

Pending.
