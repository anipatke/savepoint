---
id: v1.5/D001-in-progress-defect-cannot-resolve
release: v1.5
status: resolved
severity: medium
title: "In-progress defects cannot be resolved from the board"
---

# D001: In-progress defects cannot be resolved from the board

## Symptom

Pressing space on an `in_progress` defect shows "Defect in progress: resolve after lifecycle stage is closed" and leaves the defect unchanged.

## Expected Behavior

The board should transition an `in_progress` defect to `resolved` and remove its now-stale `stage` field.

## Reproduction

1. Open the defect overlay.
2. Select a defect with `status: in_progress` and a valid lifecycle stage.
3. Press space to resolve it.
4. Observe that the board displays a lifecycle notification instead of resolving the defect.

## Impact

Defects being actively repaired cannot complete their documented `open → in_progress → resolved` lifecycle through the board.

## Fix Plan

- Route `in_progress` defects through the board's existing resolution write path.
- Preserve the existing behavior that clears `stage` when a defect becomes `resolved`.
- Replace the test that expects an in-progress no-op with transition and persistence coverage.

## Acceptance Criteria

- [x] Pressing space on an `in_progress` defect persists `status: resolved`.
- [x] Resolving the defect removes its `stage` field.
- [x] The defect overlay model updates to the resolved status.
- [x] Existing open and already-resolved behavior remains covered.

## Context Files

- `internal/board/update.go`
- `internal/board/defect_overlay_test.go`
- `internal/board/io.go`
- `internal/data/write.go`
- `internal/data/write_test.go`

## Context Log

- Read: `internal/board/update.go`, `internal/board/defect_overlay_test.go`, `internal/board/io.go`, `internal/data/write.go`, `internal/data/write_test.go`.
- Edited: `internal/board/update.go`, `internal/board/defect_overlay_test.go`.
- Regression evidence: `go test ./internal/board ./internal/data` passed.
- Quality gate: `make build` passed; `make test` passed with clipboard/display access required by the existing clipboard integration test.
- Quick health check: skipped because `.savepoint/Health-Check.md` is absent.

## Resolution Notes

The board now routes `in_progress` defects through the existing resolution writer. Resolution persists `status: resolved`, clears `stage`, refreshes the overlay model, and retains the existing open and already-resolved behaviors.
