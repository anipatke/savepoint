---
id: C-963
scope: {kind: objective, id: O-037}
result: NEEDS WORK
checked_by: {role: checker, session: check-o037-20261004}
executed_session: o037-t108-t109-t110-20261003
checked_at: '2026-10-04T02:31:20Z'
health_snapshot: sha256:5ba3d09d60a2928a2188a7b7fb15c0f69755149d217022532a442a45bef140eb
reviewed:
  base_commit: cdfa699
  head_commit: cdfa699
  files:
    - internal/data/config.go
    - internal/data/feature_preferences.go
    - internal/data/feature_preferences_test.go
    - internal/board/v2/options.go
    - internal/board/v2/options_test.go
    - internal/board/v2/options_integration_test.go
    - internal/board/v2/model.go
    - internal/board/v2/update.go
    - internal/board/v2/load.go
    - internal/board/v2/view.go
    - internal/board/v2/help.go
    - internal/init/upgrade_test.go
    - internal/init/v2_scaffold_test.go
    - templates/project-v2/.savepoint/config.yml
    - README.md
    - CHANGELOG.md
    - .savepoint/Design.md
  dependencies:
    - internal/data/write.go
    - internal/board/v2/watch.go
    - internal/init/upgrade.go
    - .github/workflows/ci.yml
issues: [I-131, I-132, I-133, I-134]
supersedes: null
---

# C-963: O-037 Full Objective Check

NEEDS WORK. Three reproducible defects and one unverified native-platform gate.
This fresh conversation did not implement the reviewed work. Scope is HEAD plus
its current uncommitted O-037 files, not HEAD alone. All owned Tasks are done
with explicit board-owner waivers. No Task status, router, Design, implementation
or owner acceptance was changed by this Check.

## Frozen scope lock

# O-037 initial Full Check scope lock
1. O-037 seven Success Conditions; T-108 criteria 1–5, T-109 1–5, T-110 1–5; named FS/DATA/ARCH/CFG/TPL/TEST rules; full gate, native Windows CI and owner validation. STYLE advisory.
2. Surfaces: ConfigReader.Read, FeaturePreferences.UnmarshalYAML, Config accessors, WriteParallelPlanning; board loadProject, Update, save command/result, View/renderOptions/help/footer, config watcher; init template and upgrade. Changed files are current O-037 diff plus new preference/options tests. README/CHANGELOG/Design sequencing included.
3. Dependencies: existing replaceV2File temp/chmod/write/sync/close/freshness/rename/cleanup, YAML nodes and Bubble Tea command/reload/watch. No lane engine or health provider redesign.
4. Matrix:
|Row|Inputs/states/representations|Environment/boundary/sequences|Oracle|
|M1 Decode|absent/empty/missing, true/false, duplicates, null, strings/numeric/list/map, mixed unknown keys, alias/merge, direct node vs reader|default/off/on/error; mutable node re-decode; LF/CRLF|boolean or named error|
|M2 Writer|block/flow/quoted/tagged/anchored values, empty features, comments/unknown keys, Unicode, LF/CRLF/mixed endings, no final newline; invalid values/layout|off/on/off, no-op, stale before write/final check, missing/directory/read-only/symlink, create/rename failure|exact preserved bytes, requested value, no temp leftovers|
|M3 UI state|sidebar/cards; closed/open/saving/success/refused; enter/space/esc/q/o/ctrl+c; ignored keys, close/reopen pending save, external edit/reload|saved-state before/after command, repeat/retry, restart, failure then recovery, watcher; lifecycle unchanged|disk bytes vs display, command effects, focus|
|M4 render|off/on/saving/success/error/diagnostic; normal, no-color/dumb output, ANSI stripped; ASCII/control/combining/variation/wide/modifier/flag/joined emoji diagnostic|width 1–140, representative heights; detached root rendering; help/footer; supported non-TTY board has no interactive settings|independent ANSI width; no IO; sequencing promise|
|M5 adoption|fresh/legacy/no key/explicit on/off/comments/edited assets; upgrade/repeat|temporary projects, actual scaffold bytes, absent lane data|config byte identity, edited asset identity, no unexpected files|
|M6 policy/gates|Task waivers, owner validation, Design/API ownership, files real, no health/lifecycle changes|fresh full suite/crossbuild/native Windows and official health collection|test results, owner evidence, snapshot verdict|
External boundary matrix: configured/runtime target config.yml agrees; startup/load then command/save then reload; unavailable file/refusal tested; success/error tested; redirects N/A local file; timeout/cancellation N/A synchronous bounded local write (ctrl+c quitting reviewed); malformed config tested; retry/cleanup/partial writes tested; secret output N/A no credential fields printed by new parser.
5. Admission requires supported-path violation of a named criterion/rule touched/promised here, credible consequence. Unrelated old code, lane engine, health internals and general terminal overhaul excluded.

Workflow inventory:
|Order|Operation/effect|Failure/final state/cleanup|Oracle|
|1|load/index/config decode|board diagnostic; no writes|expected decoded preference/error|
|2|open/toggle; set Saving; return explicit command|no disk change until command; repeat ignored|file snapshot|
|3|lstat/read/hash/decode/splice/redecode|fatal refusal, original bytes unchanged|byte identity|
|4|temp create/chmod/write/sync/close|fatal; original intact; cleanup temp|existing writer tests and failure fixture|
|5|final freshness check/rename|refusal or complete replacement; cleanup|byte identity/temp inventory|
|6|readback/result/display/reload|only confirm matching saved value; failure notice; close restores origin|disk/display/focus|
|7|watch/reload/restart|external choice reflected, selection/lifecycle preserved|fresh instance and file snapshot|
|8|init/upgrade/adoption|scaffold off; owner config never rewritten|actual template/upgrade fixture|


Initial-pass corrections: quoted string 'true' is a non-boolean and correctly
refused, rather than a valid boolean as the first scratch harness assumed.
Repeated direct UnmarshalYAML on an already-populated receiver is not used by
the supported fresh ConfigReader path; retained as a nonblocking observation.
Widths 1–39 were completed separately after the first width-1 probe stopped the
render loop; widths 40–140 were fully rerun. Tiny-terminal observations below
do not redefine ordinary terminal support.

## Completed coverage matrix

| Row/cells | Classification and evidence |
|---|---|
| M1 absent/empty/missing, off/on, wrong type/list/map/numeric/non-finite, mixed unknown keys, top-level duplicate | Passed: existing decode/malformed tests plus independent MoreShapes. Quoted booleans refused correctly. Alias scalar refused by explicit Kind check; inherited merge keys ignored (see observation). |
| M1 null/empty value and duplicate feature keys | Issue I-131: independent Decode accepts all three without diagnostic. |
| M1 direct mutable-node reuse | Observation: preserves previous true when re-decoding an unrelated mapping; fresh reader never reuses it. |
| M2 block/flow, comments, absent/empty mapping, homogeneous LF/CRLF, no final newline, Unicode | Passed: exact-byte existing writer table and independent Unicode flow fixture. Tagged/anchored scalar writes safely refused and original preserved. Unsupported flow insertion safely refused. |
| M2 mixed LF/CRLF | Issue I-132: unrelated comment newline changed on successful save. |
| M2 on/off/repeat, stale load, source missing/appeared, malformed, directory permissions/read-only/symlink | Passed: named feature writer tests; independent read-only board probe (works even as root); no-op mtime and exact-byte checks. Missing directory and final freshness failure independently tested at replacement boundary. |
| M2 temp create/chmod/write/sync/close/rename/cleanup | Source-traced same-directory durable replacement; create/freshness failures exercised, cleanup asserted; short-write/sync/close OS fault injection unverified. These are supporting limitations, not evidence of a defect. |
| M3 complete journey and on/off lifecycle equality | Passed: OptionsJourneyOnTemporaryProject and OptionsLeaveLifecycleAndBoardIdenticalWhenOn, full suite; no lane data. |
| M3 pending save then close/reopen, enter/space crossed with esc/q/o, later second toggle | Passed: independent UISequence six combinations; reducer performs no write before command and no premature Saved display. |
| M3 repeated saving, conflict/retry, external watch/reload, focus/quit/readback failure | Passed: options named tests, watcher/reload suite, independent NoRenderingIO and ReadOnly probes. Invalid config remains diagnostic and refuses a write. |
| M4 widths 40–140 × eight text classes × on/off × normal/saving/failure/success | Passed: deterministic independent RenderMatrix with ANSI-aware width oracle, 6,464 combinations. Detached config render unchanged (no IO). Existing sizing cases include 40×14. |
| M4 widths 1–39 | Observation: options adds two overflow columns at 10–15; broad tiny-board truncation not an admitted material blocker. Other sizes retain baseline. |
| M4 normal/no-color/dumb/redirected | Render uses existing styles/overlay and no cursor commands; pure render/ANSI-stripped probes passed. Non-TTY board has no interactive settings by design; full board backend suite passed. No new renderer IO. |
| M4 visible sequencing wording | Issue I-133: independent rendered screen says suggestions are shown already. |
| M5 fresh/default/legacy/on/off/comments/edited assets/upgrade/repeat | Passed: scaffold loads through existing runtime tests, new scaffold/default and upgrade preservation tables, existing upgrade idempotence suite. No advice/worktree files created. README/CHANGELOG correctly say later delivery. |
| M6 full gate and Code Health | Passed: fresh full suite/crossbuild and official nonblocking snapshot. |
| M6 native Windows | Unverified / I-134: no current native CI evidence linked to dirty implementation. |
| M6 owner walkthrough | Owner replied completed during this Check. Recorded as report of walkthrough completion only; no accepted_check or acceptance on owner's behalf. |

All cells classified; no early stop after findings. Finite external and workflow
inventory above reviewed through file reality, source ordering, named tests and
independent probes. Persistence never touches lane/lifecycle/health records.
Failure before replacement preserves the original; primary IO errors remain
visible with secondary cleanup appended. No network boundary in settings.

## Acceptance coverage and Design

| Scope | Classification |
|---|---|
| O-037 SC1 | Issue I-131 (malformed present values silently default / ambiguous duplicate boolean) |
| O-037 SC2 | Proven: documented keyboard action, one option, ignorable explanation |
| O-037 SC3 | Issues I-131/I-132; ordinary conflict/atomic-save/failure behavior proven |
| O-037 SC4 | Proven: load/watch/reload/selection/focus unchanged in supported journeys |
| O-037 SC5 | Proven: scaffold off and upgrade configuration byte preservation |
| O-037 SC6 | Issue I-133; independent gate remains unverified I-134 |
| O-037 SC7 | Proven: no health option, provider run, deletion or waiver added |
| T-108 DW1/2/3 | Issues I-131/I-132; remaining normal/safe-refusal cases proven |
| T-108 DW4 | Proven: narrow API and one preference, no generic registry/gates |
| T-108 DW5 | Host/full evidence and waiver proven; native evidence unverified I-134 |
| T-109 DW1/2/3/4 | Proven except O-037 sequencing requirement I-133; save/restart/close/health isolation proven |
| T-109 DW5 | Fresh full and waiver proven; native evidence unverified I-134 |
| T-110 DW1/2/3/4 | Proven adoption/integration/docs, except inherited preservation I-132 and visible-copy I-133 |
| T-110 DW5 | Host/full/waiver proven; owner reports walkthrough complete; native evidence unverified I-134 |

Design ownership matches implementation: typed data reader/writer, command IO,
load-only preferences, unchanged watcher boundary and preserved upgrade config.
Its preservation/non-boolean claims are not yet fully true (I-131/I-132).
No Design repair performed. Existing authored and scoped evidence files exist;
untouched referenced files verified. Untracked root report temporaries and
AGENTS.md.new predate the review and are outside this Objective's changes.

## Adversarial pass

Custom UnmarshalYAML bypasses ordinary mapping duplicate/null diagnostics:
I-131. Whole-text newline normalization violates narrow-write preservation:
I-132. Readback diagnostics never confirm success; stale and read-only writes
leave saved off state and exact original bytes. Closing/reopening during a
pending save succeeds without losing final preference; ordinary repeated
presses while saving are ignored. Config edits reflected on reload, filesystem
mutation after validation is guarded again before rename. Nonboolean strings,
numbers, maps/lists and nonfinite values refuse safely. UI copy was assessed
against current implementation independently from correct README wording:
I-133. No privilege, authentication, billing, secret or multi-tenant boundary
exists here. Remaining native-platform proof is not inferred from cross-build.

## Commands and evidence

- Fresh `make test-full`, Go 1.26.2 linux/amd64, 2026-10-03/04 UTC: exit 0;
  full host suite and linux/darwin/windows cross-builds. `make build`: exit 0.
  No prior executor gate reused. `git diff --check`: exit 0.
- `go test -overlay /tmp/o037-overlay.json ./internal/data -run
  TestO037Independent -v -count=1`: independent null/duplicate and mixed-ending
  reproductions fail as recorded; numeric/list, mixed unknown, Unicode,
  tagged/anchored safe refusal pass. Source stored in /tmp/o037-data_test.go;
  exact reproductions are preserved in Issues, so scratch paths are not needed
  to reproduce them.
- Board independent UISequence and NoRenderingIO pass. Initial copy probe
  fails I-133; initial width-1 probe and completed tiny-width pass produce
  the nonblocking observations. Separate complete supported RenderMatrix,
  UISequence and NoRenderingIO command exit 0 (3.501s), logged in
  /tmp/o037-board-matrix.txt. ReadOnly passes. ReplacementBoundary and corrected
  MoreShapes pass; one combined rerun's board package was blocked by Go cache
  sandbox access, while prior board probes completed with approved escalation.
- Scratch tests use Go overlay; no repository test/implementation edits.
- First official health collection saved
  sha256:8b64148f400bff81fa051b6ec6a82713e572958db0302f72d4830de3ce7623c3,
  optional OSV failed because sandbox DNS was blocked. After approved network
  rerun, the frontmatter snapshot collected all five configured instances:
  Code Health does not block clearance. The failed first immutable snapshot
  is retained; collection failure was not labeled unhealthy code.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-131 | Medium (manual YAML edits) | Low (wrong preference/hidden malformed input) | Medium | Narrow decoder validation repair |
| I-132 | Low (mixed endings) | Low (unrelated authored bytes change) | Low | Narrow splice preservation repair |
| I-133 | High (every settings visit) | Low (misleading feature availability) | Medium | Correct visible wording |
| I-134 | High (no attributable CI evidence supplied) | Medium (platform clearance unknown) | Medium | Supply native CI proof after repairs |

## Nonblocking observations and limitations

- Tiny terminals 10–15 columns overflow two more columns than baseline. Normal
  supported-size matrix is sound; no unrelated terminal overhaul requested.
- Reusing FeaturePreferences.UnmarshalYAML receiver can retain true on absent
  keys; normal ConfigReader creates a fresh Config each time.
- YAML merge-only inherited feature mappings are not expanded by custom node
  traversal. Explicit documented scalar booleans work; alias/tag layouts are
  conservatively refused without rewriting. No generalized YAML editor promised.
- No new short-write/sync/close fault injection; atomicity ordering verified
  in source, and real create/refusal/freshness failure paths exercised.
- Owner walkthrough reported completed; this is not owner acceptance naming a
  Check. Only the owner can close/accept the Objective after clearance.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — null/duplicates and mixed-ending cases absent from existing preference tests; I-131/I-132.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [ ] STYLE-09 **Content lives in data** — user-facing explanation/notice copy remains in options.go constants/functions, following existing board patterns.
- [x] STYLE-10 **Small diffs**

## Handoff

Repair I-131–I-133 directly and supply I-134's platform evidence, then run an
independent recheck against this frozen matrix. Keep all completed Tasks done.
Do not manufacture acceptance or retire the mandatory Objective Check.
