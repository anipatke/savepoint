---
id: v2/D002-apply-writeability-probe-destroys-sentinel
release: v2
status: resolved
severity: blocker
title: "Apply writeability probe can destroy a pre-existing sentinel"
reference: E45-safe-migration/T011-run-migration-from-the-command-line
---

# D002: Apply writeability probe can destroy a pre-existing sentinel

## Symptom

The apply-only writeability probe writes to the fixed
`.savepoint-migrate-write-test` path. If that path already exists, the probe
truncates it and then removes it.

## Expected Behavior

The probe must create an exclusive, unique temporary file and remove only the
file created by that probe. Existing user files, including names sharing the
probe prefix, must remain byte-identical.

## Reproduction

1. Create a project with a file named `.savepoint-migrate-write-test` containing
   user data.
2. Run migration in apply mode.
3. Observe that the sentinel is empty or missing after the writeability probe.

## Impact

This violates FS-01 at the apply preflight boundary and can destroy unrelated
user content before migration starts.

## Fix Plan

Use `os.CreateTemp` with an exclusive unique prefix, clean up only the returned
path, and add an apply-mode regression test for pre-existing sentinel and
prefix-collision files.

## Acceptance Criteria

- [x] The apply-only probe never truncates or removes a pre-existing sentinel.
- [x] A pre-existing file sharing the probe prefix remains byte-identical.
- [x] Focused and repository quality gates pass.

## Resolution Notes

Replaced the fixed-path `os.WriteFile` probe with an exclusive `os.CreateTemp`
probe using the `.savepoint-migrate-write-test-*` prefix. Cleanup removes only
the unique path returned by that creation. Added
`TestRunCommand_applyWriteabilityProbePreservesPreExistingSentinelAndPrefix`,
which exercises apply mode with both collision names, verifies their exact
bytes, and detects leaked transient probe files.

Verification passed: focused and full `internal/migrate` tests, `go vet ./...`,
`GOOS=windows GOARCH=amd64 go build ./...`, `make build`, `make test`, and
`git diff --check`.
