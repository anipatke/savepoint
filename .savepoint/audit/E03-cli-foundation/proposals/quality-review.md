## Must Fix Before Close

None.

## Carry Forward

- `src/cli/help.ts` describes future command behavior with active wording such as "Initialize", "Display", "Run", and "Validate" while the command handlers are still stubs. This is acceptable as a contract preview, but it is slightly in tension with the E03 task constraint to avoid promising unimplemented behavior. Consider adding "not implemented yet" wording until each command lands, or update the text as each command becomes real.
- `src/cli/run.ts` calls `detectEnvironment(...)` and discards the result. That keeps the boundary ready for future color/TUI behavior, but it is not observable yet. Carry forward only if a near-term task will consume terminal capabilities.
- The audit snapshot omitted `src/commands/result.ts` from the "new files introduced" list even though the implementation uses it. The audit proposal for the epic design records the file.

## Already Fixed

- All E03 task files are now marked `done`.
- `.savepoint/router.md` now points to `state: audit-pending` for `E03-cli-foundation`.
- `src/cli/run.ts` now treats unknown command-level flags such as `savepoint init --bad` as usage errors, with focused runner coverage.
- The E03 snapshot records successful closeout gates: build, typecheck, lint, format check, and test.
- Formatting drift found during T006 was fixed before the quality gates were recorded.
