---
type: audit-findings
audited: 2026-09-19
---

# Audit Findings: E48 Recover current work and the next action

## Main Findings

### Verdict

**NEEDS WORK — NOT READY TO COMMIT/PUSH.** Three findings remain, all in repository code and all fixable now. The epic's hardest promises hold up under independent testing: `resume` writes nothing on any path I could reach, its evidence wording is honest about what was recorded versus what was verified, a router line naming a record that no longer exists never has a similar record substituted for it, and doctor no longer calls a project unsound because it has open Issues. What does not hold is the claim the epic leads with — that there is *one* derived answer. Two project states that a user will genuinely be in produce the wrong next action, and one of them produces a *different* answer for the same project depending on what the router happens to name. No owner-run evidence or waiver is outstanding; the gates pass today.

### What Needs Attention

**1. An Objective with no Tasks yet is told to record its integration Check.**

Why it matters: this is the normal state of every project between "we defined an Objective" and "we broke it into Tasks" — exactly the moment someone reopens a project and runs `resume`. Instead of "plan Tasks under this Objective", they are told to record the integration Check for work that does not exist yet. The rung's own contract says it means "every Task the Objective owns is done"; with zero Tasks that is vacuously true, so the Objective falls into it. The same Objective, when the router does *not* name it, correctly renders "Plan Tasks under Objective O001" — so the product already knows the right answer and gives the wrong one.

What should happen next: make the Objective-integration rung require at least one owned Task, so an Objective that has not been broken down falls through to the planning rung. Add a matrix row for it. The change is in `## Proposed Changes`; I verified it produces "Plan Tasks under Objective O001." and leaves the real integration case untouched, with the full suite passing.

Evidence: `internal/data/next.go:394` (`resolveObjectiveIntegrationRung`). Project with one Objective, no Tasks, router `state: design / objective: O001 / task: none` → `savepoint resume` prints `Next action: Record the Objective O001 integration Check.`

**2. An outstanding integration Check disappears when the router selects nothing.**

Why it matters: with every Task of an Objective done and its integration Check missing, a router line that names no selection makes `resume` print "Plan the next Objective — for a project with nothing underway yet." That statement is false, and the work it hides is the evidence step the whole V2 model is built around. Adding the Objective back to the router line flips the same project to "Record the Objective O001 integration Check." The epic states the router is a hint that never decides anything; here it decides which rung the project lands on.

What should happen next: evaluate the integration rung project-wide — over Objectives that are not yet done and own at least one Task — before the "next ready" search, matching the documented rung order. The proposed change carries the necessary `status != done` guard so an already-closed Objective with thin evidence does not outrank all new work forever; that guard is a judgement call worth confirming. Verified: the unselected project now reports the integration Check, the design-phase and ready-Task projects are unaffected, and the full suite passes.

Evidence: `internal/data/next.go:282` (`resolveLadder` evaluates rung seven only for the in-view Objective). Same project, router `objective: none` → "Plan the next Objective…"; router `objective: O001` → "Record the Objective O001 integration Check."

**3. `resume` reports its directory failures under the name `migrate`.**

Why it matters: `savepoint resume /typo` prints `migrate: target directory does not exist: /typo`. The body is actionable, the attribution is wrong, and in this epic the confusion is specific rather than cosmetic — a pending migration is a real state `resume` reports on, so a user who mistypes a path is told something that reads like a migration problem. Every other failure path in the command is correctly prefixed `resume:`; this one returns the shared resolver's error unwrapped.

What should happen next: name the two target diagnostics for the command the user ran. The existing tests assert only the substrings, so they keep passing.

Evidence: `main.go:181` (`runResume` returns `migrate.ResolveTarget`'s error unchanged); sentinels at `internal/migrate/command.go:29-30`.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| 1. Objective with no Tasks told to record its integration Check | High — every project between design and task breakdown | Medium — wrong instruction, nothing written or lost | High | Fix now; one guard plus a matrix row |
| 2. Integration Check invisible when the router selects nothing | Medium — any project whose selection was cleared after its last Task closed | Medium — outstanding evidence hidden behind a false "nothing underway" | Medium | Fix now, combined with finding 1; confirm the `status != done` guard |
| 3. `resume` failures attributed to `migrate` | Medium — a mistyped path is ordinary | Low — message text only, correct exit code and no write | Low | Fix now; small, contained |

### What Is Proven / Not Proven

**Proven.** The no-write guarantee, independently: a shell-level byte-and-mtime snapshot around `resume` on a freshly scaffolded project, plus repeated invocations, showed no file added, changed, or touched — matching the in-repo snapshot helper, which does compare content hashes *and* mtimes *and* file counts. Deterministic, ANSI-free, non-TTY output across repeated runs. Strict router decoding: unknown state, missing state, malformed `O###`/`T###`, a Task selected with no Objective, an unknown key, and a wrong YAML type each produce a distinct named diagnostic and a nonzero exit, with `none` correctly read as "not selected". Selection honesty: a router naming a non-existent Task, and a Task whose own record names a different Objective, are both reported as themselves with the available next action alongside — no near-match substitution anywhere. Evidence wording: `missing`, `needs_work`, `stale`, `unknown` and exception-allowed completion each render in their own words, with the Check ID, assessor, date and recorded basis when a freshness record exists; no rendering claims resume verified anything. Rung behaviour verified end to end through the real CLI for pending migration, replan, Task dependency, Objective dependency, execute, check-needed, owner-validation, integration, ready and plan. V1 projects get the schema message and a nonzero exit. Doctor: a V2 project whose only finding is an open Issue reports ALL CLEAN with exit 0, a done Task with no Check reports under Missing Evidence with exit 1, and the four categories print in fixed order with empty ones omitted.

**Not proven.** Two acceptance claims are not fully met by the code: the ladder "returns exactly one next action for any project state" (findings 1 and 2 return the wrong one) and "a V2 project resolves the same selection and next action regardless of consumer, with the projection carrying no presentation state" — the projection carries no presentation state, but finding 2 shows the answer moving with the router hint. Everything else in the epic's quality-gate list is proven.

### Audit Evidence

- **Scope lock:** E48 T001–T007 acceptance criteria; guardrails FS-03, FS-06, DATA-02, DATA-03, ARCH-01, ARCH-03, ARCH-04, CFG-01, TEST-01/02/04/06/08; changed files `internal/data/next.go`, `router_v2.go`, `internal/resume/{resume,evidence}.go`, `cmd/resume.go`, `main.go`, `internal/doctor/{report,checks}.go`, `AGENTS.md`, and their tests. Public entry points: `ReadStateV2`, `ResolveSelection`, `ResolveNext`, `resume.Render`, `cmd.ParseResumeArgs`/`RunResume`, `runResume`, `DiagnosticReport.HealthFindings`/`HasProblems`/`Format`. Relied-on behaviour held constant: the E43/E44 gate resolvers, `LoadV2Index`, `migrate.ResolveTarget`/`PendingOperation`. Boundary: a project on disk loaded through `LoadProject` is the supported path.
- **Coverage and workflow result:** all nine ladder rungs plus five clearance states, exception, replan, both dependency kinds and both selection diagnostics exercised against real on-disk projects through the built binary, not the in-repo fixtures. Argument handling (`--help`, unknown flag, two directories, default directory), missing directory, non-Savepoint directory, absent router, V1 project and pending migration all exercised. Findings 1 and 2 came out of the "selected versus unselected Objective" cell, which no existing test covers. Not applicable, with reason: the Unicode/text-class axis (nothing in `internal/resume` measures width, truncates, or aligns — output is unpadded sentences); external-boundary axis (no server, subprocess, or network in scope — confirmed by the package's own import test); write/rollback axis (the command performs no write to roll back).
- **File reality and drift:** every file named in the seven context logs exists; the throwaway scratch test T006 records is declared discarded. No `## Drift Notes` were required and none of the Codebase Map rows the epic promised are missing. Doctor's V1-only structural checks are now gated off for V2 projects — a real behaviour change beyond "route existing findings into categories", documented in T006's log with its root cause but not in `E48-Detail.md`; it is necessary for the advisory-backlog criterion and does not change V1 behaviour.
- **Gates:** `make build && make test` pass (11 packages), `go vet ./...` clean, `gofmt -l` clean over every scoped file, `git diff --check` clean. Both proposed code changes were applied to a scratch copy of the repository and the full suite passed with them — which also confirms no current test asserts the behaviour the findings describe.

### Guardrails Verification

- **Rule IDs checked:** FS-03, FS-06, DATA-02, DATA-03, ARCH-01, ARCH-03, ARCH-04, CFG-01, TEST-01, TEST-02, TEST-04, TEST-06, TEST-08.
- **Health check mode:** Full. `.savepoint/Health-Check.md` is absent from this project, so the evidence-mode step is skipped per AGENTS.md; absence is not a finding. `.savepoint/Guardrails.md` is present and applied.
- **Evidence:** FS-03 and TEST-04 proven by the independent snapshot described above and by the working tree being unchanged after the suite ran. DATA-02: no new lifecycle vocabulary — `NextKind` and `HealthCategory` are projection and report vocabularies owned by `internal/data` and `internal/doctor`. DATA-03 proven by the router-decoding probes. ARCH-01 holds and is enforced by a package-wide import test over `cmd/`. ARCH-03: the directory resolves identically with and without an argument. ARCH-04: the Codebase Map rows for `main.go`, `cmd/`, `internal/data` and `internal/resume` are present and accurate. CFG-01 and FS-06 are met except for finding 3's misattributed message. TEST-06 is satisfied — the tasks name exact test files and cases where they rely on existing coverage.
- **File reality evidence:** see Audit Evidence above.
- **Waivers or unresolved findings:** no waiver requested; findings 1–3 are open.

### Non-Blocking Observations

- `ResolveNext` panics on a zero or partly filled input (nil index or nil router). Not reachable through `resume` — `LoadProject` always yields a non-nil index and `ReadStateV2` a non-nil router — but E49's board is the named second consumer of this exact call, so a nil guard is worth adding when it adopts the projection.
- The `unknown` clearance sentence says "no freshness assessment has ever been recorded for it" even in the one case where a freshness assessment *does* exist: a CLEAR Check whose current assessment lacks independent checker provenance. I confirmed both shapes are rejected by the V2 decoders at load, so it is unreachable from a project on disk today — a latent wording bug, not a live one.
- On a V2 project, doctor prints "✓ no problems" for the Structure, Dependency, Audit State, Orphan and Defect checks, which no longer run at all for V2. Claiming a clean result for a check that did not execute is the same class of over-claim this epic removed from resume.
- Each advisory Issue prints with a ✗ mark under Pending Semantic Review in a report that concludes ALL CLEAN. The glyph reads as failure for the one category T006 exists to make read as normal.
- T007's "no project matches zero rungs or two" property is structurally guaranteed — `ResolveNext` returns a single value — so what the matrix really proves is that ten fixtures reach the nine expected rungs. Both findings above live in the untested space between "reaches a rung" and "reaches the right rung"; a matrix row for each closes it.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix. `internal/resume` splits layout from phrasing exactly as designed; `next.go` carries both selection resolution and the ladder, which the epic's own component table specifies.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [ ] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases. The two rung branches behind findings 1 and 2 have no test; the whole suite passes with both fixed.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments. Typed rung kinds, blocker kinds, diagnostic kinds and health categories throughout.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data. Finding 3 and the nil-input observation are the two rough edges; neither swallows an error.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config. No gate rule is restated: every readiness claim is a resolver value carried forward, and `HasProblems` now derives from one finding pass.
- [x] STYLE-08 **Comments explain why** — not what the code already says. Consistently strong across the epic.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic. Every rendered phrase has one home in `evidence.go`.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes.

`STYLE` rules are Guideline severity; the unchecked box above is an observation and not a cause of this verdict.

## Proposed Changes

### Target File
internal/data/next.go

### Replace
```go
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return Next{}, false
	}

	decision := ResolveObjectiveCompletion(index, objectiveID)
```

### With
```go
	objective, ok := index.Objectives[objectiveID]
	if !ok {
		return Next{}, false
	}

	// An Objective that owns no Task has not been broken down yet, so
	// "every owned Task is done" is vacuously true and this rung would
	// ask for an integration Check over work that does not exist. The
	// real next action is planning, which rung eight answers.
	if len(index.ObjectiveTasks[objectiveID]) == 0 {
		return Next{}, false
	}

	decision := ResolveObjectiveCompletion(index, objectiveID)
```

### Target File
internal/data/next.go

### Replace
```go
	if next, ok := resolveReadyRung(index); ok {
		return next
	}

	return Next{Kind: NextPlanObjective}
}
```

### With
```go
	if next, ok := resolveProjectWideIntegrationRung(index); ok {
		return next
	}

	if next, ok := resolveReadyRung(index); ok {
		return next
	}

	return Next{Kind: NextPlanObjective}
}

// resolveProjectWideIntegrationRung answers rung seven for an Objective the
// router does not name: the lowest-ID Objective that is not yet done, owns at
// least one Task, and whose own integration clearance is what remains unmet.
// Without it, an outstanding integration Check is reported only while the
// router happens to select its Objective, so the same project state would
// answer differently depending on a hint that decides nothing. An Objective
// already recorded done is skipped: thin evidence behind a closed Objective is
// doctor's missing-evidence report, not the project's next action.
func resolveProjectWideIntegrationRung(index *V2Index) (Next, bool) {
	for _, id := range slices.Sorted(maps.Keys(index.Objectives)) {
		objective := index.Objectives[id]
		if objective.Status == ColumnDone || len(index.ObjectiveTasks[id]) == 0 {
			continue
		}
		if next, ok := resolveObjectiveIntegrationRung(index, id); ok {
			return next, true
		}
	}
	return Next{}, false
}
```

### Target File
main.go

### Replace
```go
import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
```

### With
```go
import (
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
```

### Target File
main.go

### Replace
```go
	root, err := migrate.ResolveTarget(dir)
	if err != nil {
		return 1, err
	}
```

### With
```go
	root, err := migrate.ResolveTarget(dir)
	if err != nil {
		return 1, resumeTargetError(dir, err)
	}
```

### Target File
main.go

### Replace
```go
// schemaVersionLabel names the schema_version a V1-vs-V2 branch read, for the
```

### With
```go
// resumeTargetError names migrate.ResolveTarget's two target diagnostics for
// the command the user actually ran. Resolution is shared read-only behavior
// (ARCH-03), but its sentinels spell "migrate:", so without this a failed
// `savepoint resume /typo` reports a migrate error for a resume invocation.
func resumeTargetError(dir string, err error) error {
	switch {
	case errors.Is(err, migrate.ErrTargetMissing):
		return fmt.Errorf("resume: target directory does not exist: %s", dir)
	case errors.Is(err, migrate.ErrTargetNotSavepoint):
		return fmt.Errorf("resume: target directory is not a Savepoint project: %s has no .savepoint directory", dir)
	default:
		return fmt.Errorf("resume: %w", err)
	}
}

// schemaVersionLabel names the schema_version a V1-vs-V2 branch read, for the
```

Alongside the two `internal/data/next.go` changes, add one matrix row each in
`main_resume_matrix_test.go`: an Objective selected with no Tasks yet, expecting
`NextReady` and "Plan Tasks under Objective O001.", and an Objective whose Tasks
are all done with the router selecting nothing, expecting
`NextObjectiveIntegration`. Both fail before the changes and pass after them.
