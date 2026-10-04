---
id: C-964
scope: {kind: objective, id: O-037}
result: CLEAR
checked_by: {role: checker, session: recheck-o037-20261004}
executed_session: o037-remediation-4520c93
checked_at: '2026-10-04T02:43:55Z'
health_snapshot: sha256:1dd3b30b344b5b987a944d4e87cb20f82e772cce6ddb261231868b90ec5ee34e
reviewed:
  base_commit: cdfa699
  head_commit: 4520c93a7421dd9cc80bfe6e83cb14445aaac506
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
issues: []
supersedes: C-963
---

# C-964: O-037 Full Objective Recheck

## Closure map

| Prior Issue | Exact frozen cell | Result |
|---|---|---|
| I-131 | M1 null/empty value and duplicate feature keys; M2 malformed | Closed as verified: named refusal through reader and direct node; malformed write leaves exact original bytes |
| I-132 | M2 mixed LF/CRLF, block/flow/no-final-newline | Closed as verified: each original separator survives, off/on round trips preserve unrelated bytes |
| I-133 | M4 visible sequencing wording | Closed as verified: rendered screen says only choice saved now, no suggestions yet, arrive later |
| I-134 | M6 native Windows | Closed as verified: GitHub confirms successful native windows-tests on reviewed 4520c93 |

CLEAR. This checker conversation implemented neither the original work nor its
repairs. The owner supplied independently implemented remediation in 4520c93;
this is the first full recheck. C-963 remains immutable. Its numbered scope lock,
M1–M6 cells, external boundary, workflow inventory, materiality boundary and
recorded nonblocking observations apply without expansion. All scoped code is
committed at 4520c93a7421dd9cc80bfe6e83cb14445aaac506. No implementation, Task
status, Objective status, Design or router was changed in this session.

## Admission ledger

# C-964 admission ledger (C-963 scope retained)
|Recheck item|Prior claim/Issue|Exact frozen cell|Allowed result|
|null/empty and duplicate bool reader/direct node; malformed write unchanged|I-131|M1 null/empty value and duplicate feature keys; M2 malformed|close/remain open/unverified|
|mixed LF/CRLF scalar replacement on/off; insertion/no-final-newline|I-132|M2 mixed LF/CRLF; M2 block/flow/no-final-newline|close/remain open/unverified|
|visible screen explains saved preference now, suggestions later|I-133|M4 visible sequencing wording|close/remain open/unverified|
|native CI run source attribution|I-134|M6 native Windows|close/remain open/unverified|
|decode general supported shapes and source conflicts|unchanged behavior|M1 absent/empty/missing/off/on/types; M2 no-op/stale/missing/read-only|regression/pass|
|full save/restart/close/watch/focus lifecycle journey|unchanged behavior|M3 original complete journey, six close/reopen combinations|regression/pass|
|render width/text classes/state combinations and no IO|unchanged behavior + longer copy|M4 original 40–140 eight classes on/off four states; tiny widths observation|regression/pass; tiny widths observation only|
|upgrade and scaffold defaults|unchanged behavior|M5 fresh/default/legacy/on/off/comments/edited assets/upgrade/repeat|regression/pass|
|waivers/owner report/full gate/health/Design/files|unchanged behavior|M6 policy/gates|pass/unverified/blocking configured health verdict|
|replacement syscall fault source ordering|unchanged supporting limit|M2 temp create/chmod/write/sync/close/rename/cleanup|source trace; original fault-injection limitation retained|


## Frozen matrix results

| Row | Recheck evidence and classification |
|---|---|
| M1 normal/default/missing/types/unknown keys/top-level duplicates | Passed: full original suite, independent Decode and MoreShapes; nonboolean strings/numbers/nonfinite/list/map refused; true/false and absent behave as documented |
| M1 null/empty/duplicate preference and direct decoder | Passed: original independent reproduction now refuses all three; direct-node/write harness proves refusal and byte preservation; regression tests include tilde null |
| M1 receiver reuse / merge | C-963 nonblocking observations retained: direct receiver reuse still retains prior state; normal reader creates fresh Config. Node traversal still does not expand inherited merge-only keys. No new blocking interpretation |
| M2 layout/newlines/comments/Unicode | Passed: original independent Preservation matrix, original + extended exact-byte writer tables, independent EndingRoundTrip crosses both LF/CRLF separators with off/on and Unicode unrelated bytes |
| M2 invalid/layout/refusal/no-op/stale/missing/read-only/symlink | Passed: original full suite and independent ReplacementBoundary/ReadOnly/DirectAndWrite probes; refusal leaves original and cleans temp; unchanged saves preserve mtime |
| M2 replacement inventory | Same atomic helper and order as C-963: create/chmod/write/sync/close/finalCheck/rename/cleanup. Source rereview plus actual final freshness/create failure tests. Original low-level syscall fault-injection limitation retained; no expanded gate |
| M3 all original workflows | Passed: existing save/restart/refresh/stale/retry/focus/quit/repeated-save tests, original complete journey and lifecycle-equality integration; independent six close/reopen-while-pending combinations crossed with enter/space; no write before command/no premature saved report |
| M4 ordinary width/text classes/modes/no IO | Passed: complete original 6,464-case width 40–140 matrix rerun with longer sequencing copy, on/off × four overlay states × eight text classes; original terminal-size tests; independent detached-file renderer and failed readback probe |
| M4 width 1–39 | Original tiny-terminal overflow observation reproduced (10–15 columns adds two); remains nonblocking per C-963, not a remediation demand |
| M4 visible sequencing | Passed: original independent SequencingCopy plus amended production rendered-screen assertion; present delivery is explicitly preference-only and suggestions arrive later |
| M4 environment/output | Original pure renderer/style backend and redirected/non-TTY classifications retain C-963 meaning; full backend suite green, no new IO/cursor behavior |
| M5 scaffold/upgrade/docs | Passed: complete init/upgrade suite rerun through full gate, original fresh/runtime scaffold and preference-off tests, preserved-config/edited-asset fixture table, no lane/worktree data; README/CHANGELOG still agree with current visible screen |
| M6 gates/native/health/waivers/file reality | Passed: fresh local full gate and build, verified native CI/CodeQL, official nonblocking health snapshot, all scoped evidence paths exist and Task waivers intact |
| M6 owner validation | Owner previously reported walkthrough completed. Report retained; acceptance of this Check is still the owner's decision |

Original external-boundary and workflow rows 1–8 retain their classification:
load/decode before explicit command; narrow validated write before atomic
replacement; readback before confirmation; reload after outcome; watcher picks
up config edits; restart reads saved state; upgrade never rewrites owner config.
No operations were added to health/lifecycle/lane behavior. Success and failure
paths keep their semantic outcome; secondary cleanup cannot hide primary errors.

## Acceptance and Design reconciliation

O-037 Success Conditions 1–7 and T-108/T-109/T-110 Done When 1–5 are Proven
against the frozen coverage. C-963's failed boolean/preservation/sequencing
criteria now pass; required native Windows evidence is current. Full review
includes all three owner-completed Tasks with recorded Task-check waivers.
Their waived optional Checks do not substitute for this integration Check.

Design's typed preference, data-owned narrow writer, IO through Bubble Tea
commands, load-only rendering, conflict refusal, preservation and upgrade
boundaries now match the reviewed behavior. No speculative lane advice or health
option is claimed. No Design edit was needed or made by the checker. Guardrails
named in C-963 are satisfied for the admitted scope; style remains advisory.

## Adversarial pass

Replayed original bypass/failure findings: custom decoder no longer silently
accepts null or duplicate keys; malformed reader/direct node/write representations
all refuse. Mixed separators remain exact in both toggle directions. Untouched
comments/Unicode survive. Stale/read-only saves and readback diagnostics never
confirm unsaved state. Close/reopen while a save is pending and later retry
preserve saved-state and focus behavior. Longer copy remains within original
normal terminal matrix. Original alias/tag safe-refusal and mutable-node/merge
observations do not become new requirements. No new dependency layer or
adversarial axis was introduced.

## Evidence and gates

- Fresh `make build && make test-full`, 2026-10-04 UTC, Go 1.26.2 linux/amd64:
  exit 0; full host suite and linux/darwin/windows builds. No executor evidence
  reused. `git diff --check`: exit 0. Code/tests/fixtures/dependencies/gates are
  unchanged from 4520c93; only new Check/Issue resolution and health metadata
  are written after the fresh gate.
- `go test -overlay /tmp/o037-overlay.json ./internal/data ./internal/board/v2
  -run 'TestO037(Independent(Decode|Preservation|NodeReuse|UISequence|RenderMatrix|SequencingCopy|NoRenderingIO|ReadOnly|ReplacementBoundary|MoreShapes)|Recheck)'
  -count=1`: exit 0; data 0.055s, board 7.945s. Original scratch harness retained
  with explicit frozen-cell direct-node/malformed-write and mixed-endings
  round-trip cases. No repository source/test changes. Initial sandbox cache
  refusal resolved with approved escalation.
- Original TinyWidth probe rerun separately, output retained in
  /tmp/o037-recheck-observations.txt; unchanged C-963 observation only.
- GitHub API independently verified [CI run 37171504754](https://github.com/anipatke/savepoint/actions/runs/37171504754):
  headSha 4520c93a7421dd9cc80bfe6e83cb14445aaac506, conclusion success;
  ci job completed 2026-10-04T02:35:56Z, native windows-tests completed
  2026-10-04T02:37:34Z. All job steps succeeded, including native full suite.
  Draft PR #22's head matches this exact SHA. CodeQL actions/go/JavaScript/
  Python Analyze and CodeQL check all report SUCCESS (run 37171502918).
- Official `savepoint health check O-037`, after fresh full gate with approved
  network access: frontmatter snapshot created; all five configured instances
  collected, Code Health does not block clearance. No unhealthy measurement
  inferred from an unavailable tool, and no optional observations promoted to
  Issues.

## Materiality and observations

All four prior Issues are proven resolved; no materiality actions required.
C-963's tiny-terminal width, receiver reuse, merge-only YAML interpretation,
and un-injected short-write/sync/close limitations remain nonblocking with the
same boundaries. No new Issue or blocking perimeter was added. Owner walkthrough
completion report retained; it is not acceptance naming C-964.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [ ] STYLE-09 **Content lives in data** — user-facing copy remains in options.go constants/functions, matching existing board patterns; unchanged advisory observation.
- [x] STYLE-10 **Small diffs**

## Owner handoff

Technical clearance is current. All owned Tasks remain done. The owner may now
accept/close O-037; no Objective closure or accepted_check was recorded on the
owner's behalf. Temporary draft PR #22 may be closed under the owner's explicit
request; it was used only to run CI, not to merge implementation.
