---
id: I-095
title: Windows health cancellation leaves child processes running
type: defect
status: open
source:
  kind: check
  check: C-941
  actor: {role: checker, session: o028-check-20261001}
  at: '2026-10-01T08:38:18Z'
tasks: [T-057]
checks: [C-941]
guardrail_ids: [CFG-02, CFG-03, TEST-01, TEST-02]
history:
  - at: '2026-10-01T08:50:00Z'
    actor: {role: executor, session: o028-repair-20261001}
    kind: repair_attempted
    note: "Windows cancel and deadline now run taskkill /T /F on the tool's tree before Process.Kill; real Windows liveness check added and the skip removed. Verified natively on Windows: children test passes with the fix and fails (child survives) with the old runner. Gates: make build, make test-fast, make test-full (go1.26.2 linux/amd64) pass; awaiting a checker."
    check: C-941
---

# I-095: Windows health cancellation leaves child processes running

## Summary

T-057 promises cancellation stops the process and its children. Windows cancellation kills only the immediate tool; WaitDelay bounds waiting for output pipes rather than killing descendants. The children test skips Windows even though this scenario can occur there, contrary to CFG-03.

## Evidence

C-941 frozen M7: native Windows/amd64 test binary with an independent overlay test. Use the existing helper to start a tool that spawns a sleeping child and records its PID; cancel the tool; wait for ExecRunner to return; open the child with a Windows process handle and WaitForSingleObject for one second. Actual: `runner returned context canceled; child pid=14132 wait=258 err=<nil>` (258 is WAIT_TIMEOUT: the child is still running). The harness terminates the child during cleanup. Expected: child has exited when cancellation cleanup completes.

`internal/codehealth/runner_windows.go:11` only calls `cmd.Process.Kill()`. `runner.go:67` configures WaitDelay. `runner_test.go:162` skips this supported Windows situation, and `runner_alive_windows_test.go:13` unconditionally returns false for process liveness. Unix's equivalent existing test passes, so this is a supported-platform behavior difference rather than a universal failure.

Native command: compile `GOOS=windows GOARCH=amd64 go test -overlay <overlay.json> -c -o <scratch>/codehealth.exe ./internal/codehealth`, then run the binary with `-test.run=^TestO028WindowsChildren$ -test.v -test.count=1` from PowerShell. Full test source is embedded in C-941. Native probe fails as described; cross-builds pass.

## Proof Needed

Implement and verify supported Windows process-tree cleanup on cancellation and deadline. Use real Windows liveness checks and remove the inappropriate skip. Recheck Unix and Windows cancel/deadline/pre-cancel/children cells and temp cleanup in frozen M7/M8, then run the full gate.
