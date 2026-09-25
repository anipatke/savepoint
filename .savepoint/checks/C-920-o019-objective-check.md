---
id: C-920
scope: {kind: objective, id: O-019}
result: NEEDS WORK
checked_by: {role: checker, session: o019-objective-check-20260925}
executed_session: unrecorded-multiple-executor-sessions
checked_at: '2026-09-24T21:22:30Z'
reviewed:
  base_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  head_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  files:
    - internal/data/task_ids.go
    - internal/data/task_create.go
    - cmd/create_task.go
    - main.go
    - .github/workflows/ci.yml
    - AGENTS.md
    - agent-skills/savepoint-design/SKILL.md
    - .savepoint/Design.md
  dependencies: []
issues: [I-049]
supersedes: null
---

# C-920: O-019 Full Objective Check

## Independence and Scope Lock

This session did not execute T-010, T-011, or T-012. All three owned Tasks are `done` with owner Task-check waivers; this Full Check evaluates each directly. I-024 was escalated to O-019 and remains historical context, not a prior clearance. No previous O-019 Check exists.

1. **Criteria and rules.** O-019's seven Success Conditions; T-010's four, T-011's five, and T-012's five Done When items; Design sections 1, 2, 5, 6, and 9; FS-01, FS-04..06, DATA-01/03/05, TPL-01/02/04, ARCH-01/03/04, CFG-01/02, TEST-01..05/08/09, and STYLE-01..10 (advisory). The Full gate is `make test-full`.
2. **Changed files and entry points.** Scoped changes are the new `internal/data/task_ids.go`, `internal/data/task_create.go`, and `cmd/create_task.go`; `main.go`; focused source tests in `cmd/init_test.go`, `internal/data/discover_test.go`, `internal/data/project_test.go`, `main_test.go`, and `internal/init`; CI, README, Design, AGENTS, and matching scaffold skills/guidance. Public behavior enters through `savepoint create-task`, `cmd.RunCreateTask`/`ParseCreateTaskArgs`, `data.CreateTaskV2`, `data.AllocateTaskID`, strict V2 loading, and fresh-init/upgrade-assets delivery of the managed guide.
3. **Direct runtime reliance.** `LoadV2Index` and `DecodeTaskV2` validate the complete project and new Task; `replaceV2File` persists the high-water mark; `os.OpenFile(O_EXCL)`, `Sync`, `Rename`, `Close`, and `Remove` implement reservation, creation, and cleanup. CLI output is the final externally visible operation. The native Windows focused test job is the platform boundary. No network service or secret input is involved.
4. **Matrix perimeter.** The table below crosses allocator/creation/CLI/guidance surfaces with valid, missing, malformed, duplicate, occupied, failure, concurrent, retry, high-water boundary, Linux, and Windows cells. Renderer, Unicode display-width, non-finite numeric input, browser/server response, and colour modes are not applicable: the changed behavior does not parse such values or render terminal layouts. Console output success/failure is applicable.
5. **Materiality boundary.** Supported V2 projects, ID-free drafts, documented CLI invocation, and its process/file/output failures are in scope. Manual corruption outside a supported operation, unrelated board behavior, and hypothetical network failures are outside it. A finding must violate a named criterion or guardrail with a credible effect.

## Workflow and Side-Effect Lock

| Order | Operation and side effect | Failure timing / owner / final state and cleanup | Independent oracle |
| --- | --- | --- | --- |
| 1 | Parse flags, resolve V2 project, read and decode ID-free draft | Invalid flags, project, owner, authored ID, or YAML fail before a reservation or new Task | Parser tests; invalid-project/draft fixture tests |
| 2 | Strict-load complete V2 index, acquire exclusive `.savepoint/task-ids.lock`, strict-load again | Invalid index changes no records; held/leftover lock times out and names the path; lock is not removed by a contender | Duplicate-path and lock tests; independent leftover-lock CLI probe |
| 3 | Read active maximum and durable high-water; sync temporary state and rename under lock | Bad high-water refuses unchanged; after successful reservation its number remains retired even if later work fails or the process stops | Bootstrap, malformed-state, deletion and failure tests; independent T-999/T-1000/T-1001 probe |
| 4 | Create Objective `tasks/` directory if needed; decode the injected Task; exclusive-create, write, sync, close new file | Occupied path is preserved; new file is removed on write/validation failure; failed reservation is not reused | Occupied-file, rollback, and preservation tests |
| 5 | Strict-load full index with new Task; release lock | Validation failure removes only the new Task; successful file and high-water remain | Strict-validation rollback and concurrent two-Objective tests |
| 6 | Print `Created <id> at <path>` and return command status | A failed output write returns nonzero after the Task is committed; no rollback or success indication follows | Independent `/dev/full` CLI probe: **I-049** |

External-boundary cells: target and runtime target are the explicit project directory (`ResolveTarget`); lock startup/reuse, refusal, malformed state, success, failure, timeout, retry, cleanup, and partial side effects are covered above. Network response and redirect cells are not applicable. No secrets are handled by this path.

## Coverage Matrix

| Surface / invariant | Cells and independent evidence | Classification |
| --- | --- | --- |
| Allocator identity and boundary | Empty active set, two Objectives with T-005/T-010, higher persisted mark, deleted highest Task, reservation after caller failure, and 999→1000→1001. `TestAllocateTaskID_*`; independent CLI fixture issued T-1000 then T-1001 with `last_issued: 1001`. | Proven |
| Allocator invalid and concurrent state | Duplicate ID across Objectives names both paths; malformed high-water stays unchanged; invalid project writes no lock/state; 16 concurrent callers receive distinct IDs. Named `internal/data` tests passed. | Proven |
| Lock/failure/recovery | Held and leftover locks wait about 500 ms, name the lock path, preserve its bytes, and write no high-water. Independent CLI leftover-lock probe exited 1 and left its two existing Tasks unchanged. Temporary-file replacement and reservation-before-create ordering traced in source. | Proven |
| Draft parser and direct creation | ID-free normal and CRLF draft; preserved body and unknown field; authored ID, mismatched Objective, malformed YAML, missing owner, occupied destination, invalid project, and dangling dependency after file creation. `TestCreateTaskV2*` and `TestCreateExclusiveTaskFileDoesNotOverwriteExistingFile` passed; strict-load rollback retires the ID. | Proven |
| CLI dispatch and output | Normal args/help/errors and success path pass `cmd`/`main` tests. Redirected stdout to `/dev/full` returns exit 1 yet creates `T-001-output-failure.md`; retry can add another Task. | **Issue I-049** |
| Cross-Objective and representation | Concurrent `CreateTaskV2` calls and concurrent CLI processes produce distinct Task IDs, correct Objective links, and a valid complete index. Filename/frontmatter/index agree; original draft bytes outside injected fields remain. | Proven |
| Linux/Windows | `make test-full` passed on Linux with Linux/macOS/Windows builds. Focused native Windows `go test -count=1 ./internal/data -run '^(TestAllocateTaskID_.*\|TestDiscoverV2Records_duplicateTaskID\|TestCreateTaskV2ConcurrentAcrossObjectives)$'` passed (`go1.26.2 windows/amd64`) from a temporary NTFS copy of the scoped package. The matching `windows-latest` CI step is configured. | Proven for scoped platform behavior |
| Planner and shipped guidance | AGENTS permits only the ID-free creation exception, design skill requires it and shows concurrent Objectives, Check/Issue guidance requires resume after other identity writes. Canonical/scaffold pairs compare byte-identical. Fresh-init and real upgrade-assets tests passed; README/Design match code. | Proven |

The CLI output failure is the sole Issue. Remaining cells were completed after it was found; no other criterion remained unverified. The output-failure case is not covered by the existing tests, so green gates do not override it.

## Acceptance and Guardrails

- O-019 Success Conditions 1 and 3–6: Proven. Condition 2: Issue I-049. Condition 7: required commands and Full Check were run, but the integrated result cannot clear while Condition 2 fails.
- T-010 Done When 1–4 and T-012 Done When 1–5: Proven. T-011 Done When 1, 3–5: Proven. T-011 Done When 2: Issue I-049 at command-output failure after persistence.
- FS-01/04..06, DATA-01/03/05, TPL-01/02/04, ARCH-01/03/04, CFG-01/02, TEST-01/03..05/08/09: satisfied on inspected paths. TEST-02's relevant failure-path evidence exposed I-049. The three optional Task-check waivers identify Task, owner actor, reason, and time; they do not supply technical clearance.
- Design section 9 accurately describes lock and high-water behavior, but does not resolve the CLI's ambiguous failure report. No Design rewrite is made in this Check.

## Gates

- `git diff --check`: passed.
- `make build`: passed.
- `make test-full`: passed, including `make test` and Linux, macOS, and Windows build targets; fresh run on 2026-09-24, `go1.26.2 linux/amd64`, exit 0.
- Native Windows focused test command above: passed, `go1.26.2 windows/amd64`, exit 0. Only the scoped package and its test utility were copied to a Windows temporary directory; this was runtime execution, not cross-compilation.
- Independent CLI probes: high-water boundary and leftover lock passed; redirected output failure reproduced I-049.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-049 | Low: output sinks rarely fail | Medium: a nonzero exit can lead a planner to retry and create duplicate work | Medium: directly contradicts the failed-creation promise | Fix in a direct Issue repair, then run a fresh Full Check |

## Observations (non-blocking)

- Windows CI has not executed this uncommitted working tree; the same focused tests were run natively on Windows from a temporary NTFS copy. The configured hosted job will provide later CI evidence.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — `cmd/create_task.go:34-35` has no failing-output regression test; I-049 records the missing branch.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — O-019 spans roughly 950 added tracked lines across CLI, data, tests, and guidance; this is advisory.

## Owner Decision

O-019 is not ready for owner closure. All Tasks remain `done`; remediation belongs under I-049. A later independent Check must verify the repair with a new immutable record. The checker does not mark the Objective or Goal complete.
