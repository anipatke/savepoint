---
id: I-091
title: Health filename test fails on native Windows
type: verification
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-053]
checks: [C-939, C-940]
guardrail_ids: [CFG-02, CFG-03, TEST-08]
severity: high
resolution:
  disposition: verified
  check: C-940
  actor: {role: checker, session: o027-recheck-20261001}
  at: '2026-09-30T22:21:36Z'
  reason: Repair reproduced and proven by CLEAR Full Objective re-check C-940.
history:
  - at: '2026-09-30T21:34:30Z'
    actor: {role: checker, session: o027-full-check-20261001}
    kind: observed
    check: C-939
    note: Reproduced during the initial Full O-027 Check; see frozen matrix and embedded harness in C-939.
  - at: '2026-09-30T22:21:36Z'
    actor: {role: checker, session: o027-recheck-20261001}
    kind: rechecked
    check: C-940
    note: C-939 reproduction and independent re-check probes pass; see C-940 closure map.
---

# I-091: Health filename test fails on native Windows

## Summary

T-053 platform/test criteria 4 and 6 and CFG-03 require supported Windows behavior and passing native tests. Frozen M9/M5.

Native windows/amd64 execution of the unchanged Code Health test binary fails TestObserveUnusualFilenames at repository_test.go:268:
open ...\\quote'\"s.txt: The filename, directory name, or volume label syntax is incorrect.

The fixture list at repository_test.go:262 includes a double quote and tab on all platforms. The later Windows branch removes the tab only from the expected map, after write(t, ...) already tried creating it. Both names are illegal Windows filename cases; the double quote is the first observed failure.

Command: cross-compile the current package with GOOS=windows GOARCH=amd64 go test -c -o /tmp/o027-codehealth.test.exe ./internal/codehealth; invoke that binary natively using Windows PowerShell, with fixture cwd at the Code Health directory and "-test.v". Native package run exited 1; all other entered tests passed or explicitly skipped. Linux make test-full and cross-builds passed, demonstrating why those alone miss this.

Native symlink cases skipped because this Windows account lacks symlink privileges; Unix permission-bit cases skipped by design. A passing native full-suite CI result for this source is not established.

Existing I-030 concerned board/init/migration failures fixed and owner-resolved earlier; this is a distinct newly introduced Code Health fixture failure, so it is not a reopening of that repair.

## Proof Needed

Exclude Windows-impossible filename variants before creating files while retaining legal unusual-name coverage on Windows and all variants on POSIX. Run the focused test and full suite natively on Windows (CI evidence for CFG-03), plus the fresh full gate after repair.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

