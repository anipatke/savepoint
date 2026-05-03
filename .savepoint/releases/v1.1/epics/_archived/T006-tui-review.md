---
id: E06-audit-command/T006-tui-review
status: planned
objective: "Implement TUI review mode for approving/rejecting proposals"
depends_on: ["E06-audit-command/T005-proposals"]
---

# T006: TUI Review Mode

## Acceptance Criteria

- Shows each proposal as a side-by-side diff (old vs new)
- User can Approve, Reject, or Edit each proposal
- Edit opens inline editor to modify the "With" section
- High-divergence warning if proposal changes >50% of file
- Shows progress: X of Y proposals reviewed
- Keyboard: arrow nav, Enter approve, Backspace reject, e edit, q quit

## Implementation Plan

- [ ] Add `internal/audit/review.go` with Bubble Tea model
- [ ] Create `ReviewModel` with proposals, current index, decisions
- [ ] Implement `RenderProposal` showing old/new side-by-side
- [ ] Add keyboard handlers: Up/Down (nav), Enter (approve), Backspace (reject), e (edit)
- [ ] Implement edit mode: text area for "With" content
- [ ] Add high-divergence check (config threshold from config.yml)
- [ ] Show progress bar: "Proposal 1 of 5"
- [ ] Track decisions: approved/rejected/edited map
- [ ] Return decisions to caller for apply step
- [ ] Test review flow manually
- [ ] Run `make build && make test`