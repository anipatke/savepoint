---
type: audit-findings
audited: 2026-09-17
---

# Audit Findings: E44 Converge follow-up and verify integrated outcomes

## Main Findings

### Verdict

**Audit applied 2026-09-17.** All three findings below and both proposed
documentation edits were applied: `omitempty` added to the Issue `source`,
`resolution`, and `history` frontmatter blocks so an absent optional field is
never written blank; `WriteIssueHistoryV2` now removes the `history` key
instead of writing `history: []` when given no entries, restoring the true
no-op the task file claimed; and `v2-Design.md` section 6 now documents Issue
history in frontmatter (not the body) and lists `history` in the Issue
frontmatter field set, reconciling the divergence T001's Drift Notes flagged.
Regression coverage was added for all three: a byte-level check that a
report-sourced Issue and an accepted resolution never carry a blank `check:`
key, a byte-level check that appending a history entry does not inject blank
keys into an earlier entry, and a no-op test for `WriteIssueHistoryV2` on a
history-free Issue. `.savepoint/Design.md` `last_audited` was advanced to
`v2/E44-issues-objective-checks` with a merged architecture-delta bullet, E44
is marked `status: audited`, and the router has advanced to `epic-design` for
`E45-safe-migration`. `make build && make test` and `go vet ./...` are clean
after the apply.

At first pass (before the apply above), the substance of E44 was already
delivered and correct: Issue records load, link, and close honestly;
Objectives complete only on real integration evidence; dependent work waits at
the Objective level; and doctor names every new problem with usable repair
guidance. Three findings were found, all low-to-medium materiality and all
cheap to fix. One was a planning-document promise the epic made about itself
and never kept; two were file-writing details where the code did something
slightly different from what the task file said it proved. Nothing was a data
loss risk, and no owner-run evidence or waiver was outstanding.

Repository handoff: **READY** — all three findings are closed on the `v2`
branch as of the apply above.

### What Needs Attention

**1. [APPLIED] The epic promised to update the release design and did not.**

E44's own design file states that moving Issue history out of the authored body
and into frontmatter "is reconciled into the release design as part of this
epic, not left as an undocumented divergence," and its Open decisions repeat that
it "must be reconciled there during this epic." That edit never happened.
`v2-Design.md` section 6 still lists Issue frontmatter without `history` and
still names History as a body section. T001's own Drift Notes flagged this
honestly and said it "needs an explicit owner before epic closeout" — no later
task picked it up. This matters because E46 writes the shipped agent guidance
from this design: left alone, the next epic will document a history location the
code does not implement. What should happen next: apply the two-line design edit
proposed below, then treat the divergence as reconciled.
*Evidence:* `.savepoint/releases/v2/v2-Design.md` lines 115 and 121 versus
`internal/data/issue_v2.go` lines 171–187 (`history` decoded from frontmatter);
promise recorded in `E44-Detail.md` "Architectural delta" and "Open decisions".

**2. [APPLIED] Issue records Savepoint writes are not what the task file says they are.**

T004's acceptance criteria say an appended history preserves every earlier entry
"byte-for-byte," and the same task recorded a deliberate convention that an
absent optional field stays absent rather than being written blank. The nested
Issue blocks do not follow it: `source`, `resolution`, and each `history` entry
are re-marshalled from structs with no `omitempty`, so every record Savepoint
creates or patches gains empty placeholder keys, and an earlier history entry is
rewritten with keys it never had. A `report`-sourced Issue is written with
`check: ""` under `source`, and an `accepted` resolution — whose entire contract
is that it must *not* name a proof Check — is written with `check: ""` under
`resolution`. Nothing breaks: the values decode identically, append-only
validation still passes, and the gates are green. But these are user-readable
planning files, and an empty `check:` under a disposition that forbids one is
actively misleading to the next person or agent that hand-edits the record. What
should happen next: add `omitempty` to the five optional nested fields and add a
regression case asserting a created Issue's bytes.
*Evidence:* audit probe — `CreateIssueV2` with `Origin{Kind: report}` emits
`source:\n  kind: report\n  check: ""`; `WriteIssueHistoryV2` appending one entry
rewrites the prior entry with `note: ""` and `check: ""`. Source:
`internal/data/issue_v2.go` lines 148–169.

**3. [APPLIED] An empty history write rewrites a file it should leave alone.**

T004's acceptance criteria say writing values a record already has is a no-op
with bytes and modification time unchanged. Calling `WriteIssueHistoryV2` with no
entries on an Issue that has no history writes a new `history: []` key and bumps
the file's timestamp. It converges on the second call, so there is no loop, and
no current caller hits it — but it is the one path where the Issue writer changes
a file without changing a value, and it contradicts the same absent-stays-absent
rule as finding 2. What should happen next: remove the key instead of encoding an
empty sequence, and cover it with the existing no-op test pattern.
*Evidence:* audit probe — `WriteIssueHistoryV2(issue, nil)` on a history-free
Issue produced `history: []` and a changed mtime. Source:
`internal/data/write.go`, `WriteIssueHistoryV2`.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation | Status |
|---|---|---|---|---|---|
| 1. Release design never reconciled with the implemented history location | High — the text is already divergent | Medium — E46 generates shipped guidance from this section | Medium | Fix now; two-line documentation edit, block provided below | **Applied** |
| 2. Written Issue records carry empty placeholder keys and rewrite prior history entries | Medium — every Issue create and every managed patch | Low — cosmetic in planning files, misleading on `accepted` resolutions, no data loss | Low | Fix now; combine with finding 3 in one narrow change plus a bytes regression test | **Applied** |
| 3. Empty history write is not a no-op | Low — no current caller passes an empty list | Low — one spurious file rewrite, converges immediately | Low | Combine with finding 2 | **Applied** |

### What Is Proven / Not Proven

**Proven.** Every E44 acceptance criterion except the three above, verified
against current code rather than the task checkboxes. Issue decoding is strict
across identity, title, type, status, source, references, opaque guardrail IDs,
optional absence, resolution shape, and history shape, and it refuses Task
lifecycle vocabulary without healing it. Discovery is confined to
`.savepoint/issues/` through the same traversal, symlink, case-alias,
filename-prefix, and duplicate-ID rules Checks use — now shared through one
generic walker, so the two families cannot drift apart. Link resolution is
bidirectional with the immutable Check as the authoritative side, the reverse
asymmetry is legal, and both link maps are derived from a single walk. Resolution
obligations hold for all three dispositions, including the refusal of a proof
Check on `accepted` and `duplicate`. History is append-only against shortening,
reordering, and editing. Objective completion requires every owned Task done plus
current integration clearance, names every unfinished Task, distinguishes all
five clearance states, refuses acceptance bound to a superseded Check, and
reports an exception as allowed-by-exception rather than as clearance. Objective
dependency readiness is satisfied only by done-plus-current-clearance, reports
cleared-by-exception distinctly, blocks Task start with a dedicated blocker kind
naming both Objectives, and leaves unrelated Objectives untouched. No Issue type,
severity, status, or count changes any Task gate decision. Doctor names all
thirteen new sentinels and both consistency families, never falls through to the
generic repair, derives Issue counts at report time, keeps an open Issue
advisory, and writes nothing.

**Not proven / boundary noted.** T007's criterion covering a `verified` Issue
"whose proof Check is no longer CLEAR" is met by a different mechanism than the
criterion's wording implies: because Checks are immutable, that state only arises
from a hand edit, and the load then fails closed with
`v2-issue-resolution-unusable-proof` rather than being surfaced by
`InspectIssueConsistency`. Verified by probe; the diagnostic is named and
repairable, so this is a wording mismatch, not a coverage gap. Everything else in
scope is proven. No owner waiver is required.

### Audit Evidence

- **Scope lock:** E44's ten acceptance-criterion groups and eleven quality gates;
  the changed source and test files in `internal/data` and `internal/doctor`
  across commits `12da22c..fb7b079` (base `032095a`); the public entry points
  `DecodeIssueV2`, `DiscoverV2Issues`, `LoadV2Index`, the four Issue index
  validators, the three derived listings, `InspectIssueConsistency`,
  `ResolveObjectiveCompletion`, `ResolveObjectiveDependency`,
  `InspectObjectiveConsistency`, `ResolveTaskStart`, `CreateIssueV2`,
  `WriteIssueV2`, `WriteIssueHistoryV2`, `WriteObjectiveEvidenceV2`,
  `CheckProject`, `V2ProblemRepair`, `V2ConsistencyRepair`, and
  `IssuePostureReport`; Guardrails FS/DATA/ARCH/TEST/STYLE as mapped below.
  Out of scope by the epic's own boundary: V1 migration, live consumer rewiring,
  V2 skills, and policy cutover. `.savepoint/audit/` does not exist, so no
  register reconciliation applies.
- **Coverage and workflow result:** every mandatory cell classified. Axes
  exercised were public surfaces, input shape (normal, empty, missing, blank,
  duplicate, wrong vocabulary, malformed YAML, no frontmatter), state (all three
  Issue statuses, all three dispositions, all six history kinds, all five
  clearance states, done/not-done Objectives), boundaries (`I1`/`I12` rejected,
  `I123`/`I1234` accepted, `I999` → `I1000` allocation, filename length),
  sequences (the end-to-end `TestE44_EpicScenario`: dependent Objective blocked,
  Objective closed under checker authority, Issue verified, later Check
  staling both), and representations (decode, index, write round trip, doctor
  render). The write workflow was walked operation by operation —
  next-ID allocation, marshal, decode-validate, temp file, chmod, write, sync,
  close, hard-link, deferred cleanup — including a genuine link failure
  (filename too long) that left no file and no temp behind. Independent probes
  beyond the suite: gate neutrality rebuilt from scratch with five Issues of
  every type versus none (identical decisions on start, advance, completion);
  dependency readiness against a hand-edited done Objective with an unfinished
  Task (satisfied by design, and named by `InspectObjectiveConsistency`); doctor
  re-derivation with no writes across five files; the load-closed path for a
  no-longer-CLEAR proof Check. Findings 2 and 3 came from those probes, not from
  the existing tests.
- **File reality and drift:** every file named in all seven task context logs
  exists on disk; no phantom or discarded file. Two new files
  (`internal/data/issue_v2.go`, `internal/data/objective_gate_v2.go`) match the
  epic's declared component table. `gofmt -l` flags five files, none touched by
  E44 — matching the task logs' claims exactly. T001's Drift Note is real and
  becomes finding 1; T006's "no drift" is accurate. One unrelated change rode
  along in the epic's commits (`.claude/settings.local.json` permission entries).
- **Gates:** `make build` pass; `make test` pass across all nine packages;
  `go vet ./...` clean; `git diff --check` clean; focused
  `go test ./internal/data/ ./internal/doctor/` pass; 184 E44 `internal/data`
  test runs and 24 `internal/doctor` test runs enumerated and green.

### Guardrails Verification

- **Rule IDs checked:** FS-01, FS-04, FS-05, FS-06; DATA-01, DATA-02, DATA-03,
  DATA-04; TPL-01, TPL-02; ARCH-03, ARCH-04; TEST-01 through TEST-08;
  STYLE-01 through STYLE-10.
- **Health check mode:** Full. `.savepoint/Health-Check.md` is absent from this
  project, so no project-specific evidence procedure applies; per AGENTS.md its
  absence is not a finding.
- **Evidence:** FS-01/FS-04 — writes are patch-in-place with validate-before-
  replace and staleness re-checks; the one no-op exception is finding 3. FS-05 —
  every path built with `filepath.Join`; a title containing `../` slugged to
  `a-escape` with no traversal. FS-06 — absent `issues/` loads empty; an
  unwritable target and a failed link both refuse cleanly with no partial file.
  DATA-01 — unknown top-level frontmatter keys and authored bodies survive
  managed patches (verified by probe); managed nested blocks are wholly owned by
  the writer, consistent with the E43 evidence pattern. DATA-02 — Issue status is
  a separate vocabulary and Task lifecycle values are rejected, not healed, with
  a `stage` on an Issue refused outright. DATA-03/DATA-04 — every new failure is
  a named sentinel with a doctor name and a manual repair line; nothing is
  self-healed. DATA-05 is untouched: E44 changes no completion authority.
  TPL-01 — canonical and scaffolded skills remain byte-identical (E44 touched no
  skill). TPL-02 was finding 1, now applied. ARCH-03 — decisions re-resolve from the index at
  call time; no cached judgement. ARCH-04 — no new package; the AGENTS.md
  Codebase Map's `internal/data` entry still predates V2 records generally, which
  is not E44's regression (see observations). TEST-01 through TEST-08 — every
  changed behavior has named test-case evidence recorded in the task files and
  verified to exist and pass; happy and failure paths both covered; all tests use
  `t.TempDir()`; `make build && make test` passes.
- **File reality evidence:** as recorded under Audit Evidence; all declared
  source and test files were reviewed from current file reality, not from the
  task logs.
- **Waivers or unresolved findings:** no waivers. Findings 1–3 above, all applied on 2026-09-17.

### Non-Blocking Observations

These are outside the findings bar and change nothing about the verdict:

- `internal/doctor/checks.go`'s doc comment on `v2ConsistencyProblems` still
  describes only `InspectTaskConsistency` and says "Each Problem's File names the
  Task's source record," although the function now also reports Objective and
  Issue consistency. Stale comment only; the code is correct.
- On a V2 project whose Issue records fail the load, the report's Issue Posture
  section prints "(not a V2 project)". The real error is reported directly above
  it by the Project Check, so nothing is hidden, but the wording is misleading.
- `createV2RecordFile` still names its temporary file `.savepoint-v2-check-*`
  even when creating an Issue. Invisible in practice — the temp file is always
  removed — but the name no longer matches the generalized function.
- `WriteIssueV2` validates only that its output decodes; it cannot check
  index-level obligations, so an API caller can write a record that makes the
  next load fail closed. This is the documented fail-closed contract and matches
  how a hand edit behaves, and there is no live consumer yet — worth a note when
  E46 wires the write path into skills.
- The AGENTS.md Codebase Map's `internal/data` row does not mention the V2 record
  families at all. That gap predates E44 (E42 and E43 introduced them), so it is
  not this epic's regression, but E44 is the third epic to widen the package
  without updating the row.
- Two `.claude/settings.local.json` permission entries were committed inside the
  E44 range. Harmless local tooling config, unrelated to the epic.

## Code Style Review

STYLE rules are Guideline severity and advisory; no unchecked box below is a
blocker or a cause of the NEEDS WORK verdict.

- [ ] STYLE-01 **One job per file** — split files when responsibilities mix.
  `issue_v2.go` carries decoding, index-level obligations, consistency
  inspection, and derived listings in one 661-line file, and `write.go` now holds
  the Task, Check, Objective, and Issue writers. Both are consistent with the
  existing V2 layout, and E43's audit made the same observation about `write.go`.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
  Among the most thorough coverage in the repository: every vocabulary value,
  every clearance state, every disposition, both link directions, and an
  end-to-end scenario.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions. T003
  explicitly declined to add helper functions with no obligation to enforce.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and
  external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state,
  or config. Check and Issue discovery, the evidence patch builder, the
  create-only file sequence, the actor/timestamp decoders, and the status/type
  vocabularies were each unified rather than copied.
- [ ] STYLE-08 **Comments explain why** — not what the code already says.
  Comments are unusually strong throughout. Finding 2's `WriteIssueHistoryV2`
  "byte-for-byte" mismatch is now resolved by the applied `omitempty` fix. The
  `v2ConsistencyProblems` doc comment in `internal/doctor/checks.go` (still
  describing only `InspectTaskConsistency`, though the function now also
  reports Objective and Issue consistency) remains a non-blocking observation,
  not applied here.
- [ ] STYLE-09 **Content lives in data** — keep copy/config out of logic. The new
  diagnostic names and repair strings are string literals inside `switch`
  statements, matching the existing pattern E43's audit already flagged.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving
  changes. The diff is large but confined to the epic's declared files, and each
  refactor of existing code was behaviour-preserving.

## Proposed Changes

### Target File
.savepoint/releases/v2/v2-Design.md

### Replace
```md
Issue frontmatter: `id`, `title`, `type: defect|drift|guardrail|verification|other`, `status: open|in_progress|resolved`, `source`, `tasks`, `checks`, `guardrail_ids`, optional severity, `resolution`, and `duplicate_of`. Body: Summary, Evidence, Proof Needed, History. Type is descriptive, never sufficient by itself to block a Task.
```

### With
```md
Issue frontmatter: `id`, `title`, `type: defect|drift|guardrail|verification|other`, `status: open|in_progress|resolved`, `source`, `tasks`, `checks`, `guardrail_ids`, optional severity, `resolution`, `duplicate_of`, and `history`. Body: Summary, Evidence, Proof Needed. Type is descriptive, never sufficient by itself to block a Task.
```

### Target File
.savepoint/releases/v2/v2-Design.md

### Replace
```md
- History records observations, repair attempts, rechecks, deferrals, and owner decisions. Deferral stays open with reason, not another lifecycle state. Check records replace immutable audit run history; issue listings/counts are derived without a separate register.
```

### With
```md
- History is an append-only frontmatter list of `{at, actor, kind, note, check}` entries recording observations, repair attempts, rechecks, deferrals, reopenings, and owner decisions. It lives in frontmatter, not the body, because it must be machine-appended, ordered, dated, and provably append-only: a managed write that would shorten, reorder, or edit a recorded entry is refused. Deferral stays open with reason, not another lifecycle state. Check records replace immutable audit run history; issue listings/counts are derived without a separate register.
```

### Target File
internal/data/issue_v2.go

### Replace
```go
type issueOriginFrontmatter struct {
	Kind  string                   `yaml:"kind"`
	Check string                   `yaml:"check"`
	Actor evidenceActorFrontmatter `yaml:"actor"`
	At    string                   `yaml:"at"`
}

type issueResolutionFrontmatter struct {
	Disposition string                   `yaml:"disposition"`
	Check       string                   `yaml:"check"`
	Actor       evidenceActorFrontmatter `yaml:"actor"`
	At          string                   `yaml:"at"`
	Reason      string                   `yaml:"reason"`
}

type issueHistoryFrontmatter struct {
	At    string                   `yaml:"at"`
	Actor evidenceActorFrontmatter `yaml:"actor"`
	Kind  string                   `yaml:"kind"`
	Note  string                   `yaml:"note"`
	Check string                   `yaml:"check"`
}
```

### With
```go
// The optional fields carry omitempty so a written record never records a
// blank value for a field the record does not declare, matching the rule that
// an absent optional field stays absent. An accepted resolution in particular
// must not appear to name a proof check.
type issueOriginFrontmatter struct {
	Kind  string                   `yaml:"kind"`
	Check string                   `yaml:"check,omitempty"`
	Actor evidenceActorFrontmatter `yaml:"actor"`
	At    string                   `yaml:"at"`
}

type issueResolutionFrontmatter struct {
	Disposition string                   `yaml:"disposition"`
	Check       string                   `yaml:"check,omitempty"`
	Actor       evidenceActorFrontmatter `yaml:"actor"`
	At          string                   `yaml:"at"`
	Reason      string                   `yaml:"reason,omitempty"`
}

type issueHistoryFrontmatter struct {
	At    string                   `yaml:"at"`
	Actor evidenceActorFrontmatter `yaml:"actor"`
	Kind  string                   `yaml:"kind"`
	Note  string                   `yaml:"note,omitempty"`
	Check string                   `yaml:"check,omitempty"`
}
```

### Target File
internal/data/write.go

### Replace
```go
	node, err := encodeV2Node(issueHistoryV2Frontmatter(entries))
	if err != nil {
		return fmt.Errorf("encode issue history: %w", err)
	}

	return writeV2Record(&issue.Source, []v2FieldPatch{{Key: "history", Node: node}}, func(content string) error {
		_, err := DecodeIssueV2(issue.Source.Path, content)
		return err
	})
```

### With
```go
	// An empty history removes the key rather than recording an empty
	// sequence, so writing no entries to a record that has none stays a true
	// no-op instead of rewriting the file with `history: []`.
	patch := v2FieldPatch{Key: "history", Remove: true}
	if len(entries) > 0 {
		node, err := encodeV2Node(issueHistoryV2Frontmatter(entries))
		if err != nil {
			return fmt.Errorf("encode issue history: %w", err)
		}
		patch = v2FieldPatch{Key: "history", Node: node}
	}

	return writeV2Record(&issue.Source, []v2FieldPatch{patch}, func(content string) error {
		_, err := DecodeIssueV2(issue.Source.Path, content)
		return err
	})
```

Applying the two Go blocks should be accompanied by regression cases in
`internal/data/write_test.go`: one asserting the exact frontmatter bytes of a
`report`-sourced `CreateIssueV2` record and of an `accepted` resolution patch
(no `check:` key present), one asserting a prior history entry's rendered YAML
is unchanged across an append, and one asserting `WriteIssueHistoryV2` with no
entries leaves a history-free record's bytes and mtime untouched.
