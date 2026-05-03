---
id: E08-board-command/T006-status-transitions
status: done
objective: "Implement status transitions with gate enforcement"
depends_on: ["E08-board-command/T005-detail-pane"]
---

# T006: Status Transitions

## Acceptance Criteria

- Space advances: planned → in_progress → done
- Backspace retreats: done → in_progress → planned
- Phase gating: can only reach done after audit phase
- Dependency gating: cannot advance if depends_on not done
- mtime conflict detection: warn if task modified since load
- Write status back to task frontmatter

## Implementation Plan

- [x] Add `internal/board/transitions.go` (extend existing or create new)
- [x] Implement phase advancement (space key)
- [x] Implement phase retreat (backspace key)
- [x] Check phase gating (must reach audit phase before done)
- [x] Check dependency gating (all depends_on status: done)
- [x] Implement mtime conflict detection before write
- [x] Write new status to task frontmatter via data.WriteTaskStatus
- [x] Test transition behavior in TUI
- [x] Run `make build && make test`