# Context Bench

Purpose: track whether workflow changes reduce agent context load between epics.

## Method

Use the same estimate for every entry:

```text
estimated_input_tokens = ceil(total_characters_read / 4)
```

This is not tokenizer-exact. It is intentionally simple and repeatable, so before/after comparisons are meaningful.

Count only files intentionally loaded for the active work. Do not count command output from `rg`, `git status`, formatting, or tests unless the output itself became task context.

For each phase, record:

| Field                    | Meaning                                                               |
| ------------------------ | --------------------------------------------------------------------- |
| `epic`                   | Active epic ID, e.g. `E02-data-model`                                 |
| `phase`                  | `design`, `task-breakdown`, `task-planning`, `task-building`, `audit` |
| `files_read`             | Number of files intentionally loaded                                  |
| `estimated_input_tokens` | `ceil(total_characters_read / 4)`                                     |
| `actual_input_tokens`    | API-reported input tokens, if available                               |
| `actual_output_tokens`   | API-reported output tokens, if available                              |
| `quality`                | `measured`, `reconstructed`, or `target`                              |
| `notes`                  | Why extra context was needed                                          |

## Baseline: E01-scaffolding

These values are reconstructed after the fact from the files known to have been loaded or required by the workflow. They are good enough for directional comparison, not precise billing analysis.

| Epic            | Phase                        | Files read | Estimated input tokens | Quality       | Notes                                                                                                       |
| --------------- | ---------------------------- | ---------- | ---------------------- | ------------- | ----------------------------------------------------------------------------------------------------------- |
| E01-scaffolding | implementation reconstructed | 20         | ~10,882                | reconstructed | Included router, project/release/epic docs, all task files, and scaffold source/config files.               |
| E01-scaffolding | audit closeout reconstructed | 17         | ~16,457                | reconstructed | Included router, project/release docs, live audit targets, proposals, snapshot, and touched scaffold files. |

## Target: E02-data-model

The process changes should make the next epic cheaper by default. The first comparison point is task breakdown, which should require only router state plus the active epic Design.

| Epic           | Phase          | Files read | Estimated input tokens | Quality | Notes                                                |
| -------------- | -------------- | ---------- | ---------------------- | ------- | ---------------------------------------------------- |
| E02-data-model | task-breakdown | 2          | ~2,295                 | target  | `.savepoint/router.md` + `E02-data-model/Design.md`. |

## Log: E02-data-model

Append measured entries here as the epic progresses.

| Date       | Phase                 | Files read | Estimated input tokens | Actual input tokens | Actual output tokens | Quality | Notes                                                    |
| ---------- | --------------------- | ---------- | ---------------------- | ------------------- | -------------------- | ------- | -------------------------------------------------------- |
| 2026-04-27 | task-breakdown target | 2          | ~2,295                 | n/a                 | n/a                  | target  | Expected minimum context before creating E02 task files. |

## Delta Log: E02-data-model

Use this to compare what the workflow expected against what the agent actually needed.

| Date       | Phase          | Target files | Measured files | Target tokens | Measured tokens | Delta tokens | Reason for delta               |
| ---------- | -------------- | ------------ | -------------- | ------------- | --------------- | ------------ | ------------------------------ |
| 2026-04-27 | task-breakdown | 2            | TBD            | ~2,295        | TBD             | TBD          | Fill after E02 task breakdown. |

## E02 Logging Instructions

At the end of each E02 phase, append one measured row to `Log: E02-data-model` and one comparison row to `Delta Log: E02-data-model`.

Use these phase targets unless the router or task explicitly changes scope:

| Phase          | Target files | Target tokens | Target context                                    |
| -------------- | ------------ | ------------- | ------------------------------------------------- |
| task-breakdown | 2            | ~2,295        | Router + `E02-data-model/Design.md`               |
| task-planning  | 3            | TBD           | Router + E02 Design + selected task               |
| task-building  | 4-8          | TBD           | Router + selected task + directly touched files   |
| audit          | TBD          | TBD           | Snapshot + changed files + patch-shaped proposals |

For every positive delta, write the concrete reason. Examples:

- Needed release PRD to resolve epic order.
- Needed project Design because architecture changed.
- Needed previous handoff because dependency crossed epic boundary.
- Tool output became task context.

## Future Target: E03-cli-foundation

`E03-cli-foundation` was prepped early by mistake while tightening the workflow. Keep its context budget and close criteria, but do not use it as the next measured comparison until E02 is audited.

## Comparison Rules

- Compare phase to phase. Do not compare one epic's audit closeout against another epic's implementation.
- Track both `files_read` and `estimated_input_tokens`; a lower file count can still be worse if one file is large.
- Mark after-the-fact estimates as `reconstructed`.
- Mark entries captured during the work as `measured`.
- Keep notes short and explain only context exceptions.
