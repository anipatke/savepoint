## Target File

`.savepoint/releases/v1/epics/E03-cli-foundation/Design.md`

## Replace

## Architectural delta

Before this epic, the binary is a placeholder. After this epic, `savepoint` has a stable command boundary with stubbed behaviors behind it.

The command layer should not introduce data-model behavior. Commands that need project data stay as deterministic stubs until a later epic owns the data model.

## With

## Architectural delta

Before this epic, the binary is a placeholder. After this epic, `savepoint` has a stable command boundary with stubbed behaviors behind it.

The command layer should not introduce data-model behavior. Commands that need project data stay as deterministic stubs until a later epic owns the data model.

## Implemented As

- `src/cli.ts` now isolates process globals and delegates to `runCli()`.
- `src/cli/args.ts` implements the closed parser contract for bare invocation, global help/version flags, known commands, unknown top-level commands, and unknown top-level flags.
- `src/cli/run.ts` wires parser output to help, version output, command-level help, command stubs, and usage errors.
- `src/cli/help.ts` provides deterministic top-level and command-level help text.
- `src/cli/environment.ts` provides injectable TTY, color, and platform detection.
- `src/cli/exit-codes.ts` centralizes success, usage-error, and not-implemented exit codes.
- `src/commands/init.ts`, `src/commands/board.ts`, `src/commands/audit.ts`, and `src/commands/doctor.ts` return deterministic not-yet-implemented results with no project filesystem side effects.
- `src/commands/result.ts` was added as the shared command-result type used by the command stubs.
- Focused tests were added under `test/cli/` for parsing, help text, terminal environment detection, command stubs, and runner dispatch.
- The audit snapshot was generated manually because `savepoint audit` is still a stub.
- The router had to be corrected after task review because E03 implementation had completed while the router still pointed at T001.

Design delta notes:

- The planned `src/commands/*.ts` component was implemented as one module per command plus a small shared `result.ts`.
- The epic design table described command modules as the "five-command surface"; the implemented command surface has four commands plus top-level global help/version flags, matching the rest of the design.
- Command behavior remains stubbed by design. Later epics own real `init`, `board`, `audit`, and `doctor` behavior.
