---
id: v1/D001-shared
release: v1
status: resolved
severity: medium
title: "Baseline output missing trailing newline"
reference: E01-example/T001-shared
---

# D001: Baseline Output Missing Trailing Newline

## Symptom

The v1 baseline output file was written without a trailing newline, tripping the
release's byte-comparison check.

## Expected Behavior

The baseline output should end with a single trailing newline, matching every other
example artifact in the release.

## Resolution Notes

Baseline rewritten with the trailing newline before T001-shared was marked done. Frozen
here as resolved, release-scoped migration source evidence only.
