---
id: v1.1/D001-shared
release: v1.1
status: open
severity: low
title: "Follow-up drafted before the v1.1 baseline is confirmed stable"
reference: E01-example/T002-follow-up
---

# D001: Follow-up Drafted Before Baseline Confirmed

## Symptom

`T002-follow-up` was drafted and scoped while `T001-shared` (v1.1) was still in
progress, before the redone baseline was confirmed stable.

## Expected Behavior

A dependent follow-up Task should not be started until its same-epic dependency Task
reaches done.

## Reproduction

1. Open the v1.1 release, epic `E01-example`.
2. Observe `T002-follow-up` exists with `depends_on: [T001-shared]` while T001-shared is
   still `in_progress`.

## Impact

Low: the follow-up remains `planned`, not started, so no premature work has happened.
Left open in this frozen fixture as release-scoped migration source evidence, distinct
from the resolved v1 record sharing the same short ID.
