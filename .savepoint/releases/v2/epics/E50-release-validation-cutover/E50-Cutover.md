# E50 maintainer-controlled repository cutover

**Handoff state: NOT CLEAR TO APPLY.** This file is the owner-run runbook for the
live repository. It is not an instruction for an agent to mutate this checkout.
The agent has not run the Savepoint CLI, has not applied the migration, and has
not published, tagged, deployed, or changed a changelog. The primary apply
command appears once below as an owner-only action and is stopped behind the
gates in this document.

## Decision and stop point

The repository is still a V1 project. A read-only preflight found an appliable
plan with no conflicts and no unresolved *blocking* ambiguities. The owner
reconciled the historical release records before this capture: R001, R002,
R003, R004, and R005 are no longer active migration work. R001, R002, R003,
R005 are recorded as completed historical Releases; R004's discarded v1.3
source is preserved at `.savepoint/archive/retired/v1.3/` and is excluded from
the active migration inventory.

The disposable apply now has six Releases (the current v2 release is R006),
one active Objective, and two active Tasks. Canonical
`data.ResolveReleaseCutover` is **Allowed=false** with two named blockers: the
current E50 Objective's T001 is still `in_progress` and T002 is still
`planned`. This is a clear, reproducible Release decision, not an ambiguity to
waive. The candidate also contains 11 Issues; owner acceptance of any material
Issue remains a separate gate below.

Therefore this handoff must stop before live apply. Complete the remaining
current-release Tasks and recheck every material Issue, then repeat the
isolated apply and all preflight evidence. A maintainer may approve the live
operation only when every gate below is green and the exact preview has been
reviewed.

When those gates are green, the one normal owner action is:

```text
savepoint migrate . --apply
```

Do not run that command from this agent session. Do not infer approval from a
clean migration plan: Release readiness, pending-operation state, V2 loading,
exact owner acceptance, and the backup must all be recorded first. If any gate
is red or evidence is stale, stop and return to the dry-run step.

## Independent evidence required before apply

The latest independent E51 re-audit is
`.savepoint/releases/v2/epics/E51-first-class-releases/E51-Audit.md`, dated
2026-09-20. Its verdict is **CLEAR**, with no remaining findings and no waiver;
it explicitly clears E51 for commit/push. T009 is also recorded `status: done`
with its focused tests, repository-copy evidence, and full `make build && make
test` gate. This is an independent audit, not a self-review by the migration
builder. Re-read the audit at the owner checkpoint; a changed audit, source
tree, or design invalidates this handoff.

The current capture was made read-only after the owner-authorized lifecycle
reconciliation and v1.3 retirement, with an injected evidence clock of
`2026-09-20T00:00:00Z` and operation ID `op-t007-reconciliation`. At that
capture:

| Evidence | Value |
| --- | --- |
| source `HEAD` | `a9901f88b8e83ece18e38ecbced625e6e2946f4c` |
| working-tree state at capture | lifecycle reconciliation plus the v1.3 archive move and retirement note; tracked diff SHA-256 `4c55504025dc19b36d392cec86b80df5dcac8eb80464e01925ceb8cd103ddaee` |
| retired v1.3 tree | 40 files preserved at `.savepoint/archive/retired/v1.3/`; pre-move path-qualified hash `88549f35e22f1189f629ddc600c4fcd4d8109a6896014d3d6956272c2630ce65`; post-move content hash `90bba76a269709aad9d7076ea3d2e76551431f3b7986407080eeed90343f8881` |
| migration source inventory | 577 files; canonical digest `c9d408c37cae52aeadff5c01613e4e3a2bbe66b7ccf4854094bc338dc262b42a` |
| deterministic plan | 20 targets (6 Releases, 1 Objective, 2 Tasks, 11 Issues), 2 documents, 544 archives, 1 legacy prerequisite, 0 waived references, 0 conflicts, 141 ambiguities, 0 unresolved blocking IDs |
| fixed-clock preview | 184,840 bytes; full SHA-256 `472f2526bddcb3f013b067d832448f9cbf990b888e95ab92c0b76cef46e717c7` |
| normalized preview | SHA-256 `e7330b0c59767ff2113d801e35bc20a431cab066bf32a6dcaaccea6cbf54c39d`, after removing only the generated `operation:` and `generated:` lines |
| read-only command result | `RunCommand(Write:false)` returned code 0 with no error; no project file changed |

The inventory digest is the sorted hash of each source path and exact-byte
SHA-256 (`path NUL sha256 newline`) over the migration inventory: `AGENTS.md`,
`.savepoint/` excluding `.migration/` and `archive/`, and `agent-skills/`.
The preview digest is only comparable after normalizing the random operation
ID and timestamp. The full fixed-clock value above is a reproducibility anchor;
an owner-run CLI preview will otherwise have different first two lines.

The same preflight was run against a disposable repository copy and the copy
was removed after comparison. It returned code 0, the same 577-file inventory,
the same plan counts, the same inventory digest, and the same normalized
preview digest. Live and copy preview output was byte-equal; no live migration
was used as a test fixture. The intentional differences from the earlier
capture are the owner-authorized historical lifecycle edits and removal of the
discarded v1.3 tree from active migration input.

The owner must recapture after this file and the task handoff are committed;
neither this uncommitted evidence nor the prior captures is permission to apply
a stale plan. Any other source edit, changed revision, changed decision file,
or changed operation state requires the same recapture.

## Gates that all must pass

1. **Independent audit.** E51-Audit remains CLEAR, and the E50 validation
   trials and three named agent scenarios remain present and internally
   consistent.
2. **Clean source boundary.** Record the exact commit, `git status --short`,
   inventory digest, and normalized preview digest. A user edit, an unknown
   file, or an uncommitted change not listed in the evidence stops the run.
3. **No operation in progress.** There is no `.savepoint/.migration/<opID>/`
   pending operation, no multiple-operation diagnostic, and no stale staged or
   backed-up path. Never start a fresh apply over an incomplete operation.
4. **Plan is appliable.** The dry-run has no conflicts and no unresolved
   blocking ambiguity. Advisory ambiguities must be reviewed and recorded;
   they are not silently converted into decisions.
5. **V2 candidate loads.** An isolated apply of the exact plan must load a
   structurally valid V2 index. Invalid IDs, unsafe paths, dangling ownership,
   dependency cycles, malformed records, or a failed recovery simulation stop
   the handoff.
6. **Release decision is canonical and allowed.** Use only
   `data.ResolveReleaseCutover`, which delegates to
   `ResolveReleaseCompletion`. Do not add a checklist or a parallel E50 rule.
   Every declared Release must be allowed; stale/missing evidence, a missing
   owner acceptance, an unexcepted material Issue, an incomplete objective, or
   a Release with no valid objective blocks the handoff.
7. **Owner acceptance.** The maintainer has reviewed the complete preview,
   backup hashes, expected write map, recovery plan, rollback limits, and the
   post-apply verification list, and has explicitly approved the one command.
8. **Scope.** No publish, tag, deployment, release announcement, or changelog
   edit is part of migration. Those are separate owner decisions.

The current candidate fails gate 6 on the two current E50 Task states. The
failure is intentionally recorded below rather than hidden behind a generic
“not ready” statement.

## Recorded Release decision from the isolated apply

The disposable apply used the fixed operation ID `op-t007-reconciliation` and
a copy of the source tree. It completed into schema 2 with 6 Releases, 1
Objective, 2 Tasks, 0 Checks, and 11 Issues. `data.ResolveReleaseCutover`
returned `Allowed=false` and 2 blockers:

| Release | Decision | Canonical blockers |
| --- | --- | --- |
| R001 | allowed (historical) | v1 is `done`; its audited epics have no live Objective members |
| R002 | allowed (historical) | v1.1 is `done`; E17 is audited with all six Tasks done |
| R003 | allowed (historical) | v1.2 is `done`; audited epics have no live Objective members |
| R004 | not in active migration | v1.3 was explicitly retired; the 40-file source tree is preserved under `.savepoint/archive/retired/v1.3/` |
| R005 | allowed (historical) | v1.4 is `done`; audited epics have no live Objective members |
| R006 | not allowed (2) | current v2/E50 Objective O001: T001 is `in_progress`; T002 is `planned` |

These are the exact resolver outputs, not a second E50 policy. The 11 Issues
must be reviewed by the owner after the candidate is loaded; this handoff does
not claim that any Issue is immaterial, waived, accepted, or resolved. Because
the current decision is clear but false, the correct action is to stop, not to
invent a waiver or apply to discover the same blockers in the live checkout.

## Owner procedure

### 1. Freeze and make the rollback snapshot

Perform this outside the agent session, from the repository root, after the
runbook and task handoff have been committed. Use the parent directory so the
backup cannot enter the migration inventory:

```text
REV=$(git rev-parse HEAD)
BACKUP_DIR="../savepoint-cutover-backups/${REV}"
mkdir -p "$BACKUP_DIR"
git status --short
git bundle create "$BACKUP_DIR/repository.bundle" --all
git diff --binary > "$BACKUP_DIR/worktree.patch"
git diff --staged --binary > "$BACKUP_DIR/index.patch"
tar --exclude=.git -cf "$BACKUP_DIR/worktree.tar" .
sha256sum "$BACKUP_DIR"/* > "$BACKUP_DIR/SHA256SUMS"
```

The expected location is `../savepoint-cutover-backups/<clean-HEAD-sha>/`,
outside the project root. Keep the bundle, worktree snapshot, patches, and
hash file until the owner has accepted the post-apply verification. If the
tree is not clean, stop and explain every intentional difference before
continuing; do not overwrite an existing backup directory.

Also confirm that no pending operation exists and that no owner is editing the
project. A pending operation is recoverable state, not disposable temporary
data.

### 2. Run the fresh, write-free preview

The owner now runs the actual CLI (the agent deliberately did not):

```text
savepoint migrate . --dry-run > "$BACKUP_DIR/migration-preview.txt"
printf '%s\n' "$?"
```

The exit code must be 0. `--dry-run` must not create a manifest, archive,
router/config change, migration directory, or other project file. Compare the
complete output, not only its summary. Normalize only the two generated lines:
remove the `operation:` line and the `generated:` line, then compare the
remaining SHA-256 with the recaptured repository-copy result. Every difference
must be listed with its reason (for example, this runbook's intentional source
addition); an unexplained difference is a hard stop. Preserve the unmodified
preview and the normalized comparison in the backup directory.

The reconciled result is 20 targets, 2 documents, 544 archives, 1 legacy
prerequisite, 0 waivers, 0 conflicts, and 141 advisory/resolved ambiguities.
Do not treat these hashes as valid after the owner recaptures the committed
tree.

### 3. Re-run the isolated candidate and Release gate

Before touching the live checkout, apply the exact reviewed plan to a fresh
disposable copy. Confirm that the copy loads as V2, that its operation reaches
the schema activation boundary, that the archive/manifest hashes are exact,
and that a second apply reports “already at schema_version: 2; nothing to
migrate.” Resolve every structural diagnostic and every recovery simulation
failure in the copy first.

Then record the canonical Release decision and the Issue review. The current
two-blocker decision above is a stop. Only an explicit `Allowed=true` result,
fresh evidence for every Release, and owner acceptance of all material Issues
can move to the live stop point. No task status or Release field may be changed
just to make this gate pass.

### 4. Stop for explicit approval

At this point attach the clean revision, backup checksum file, complete dry-run,
copy comparison, V2 load result, Release decision, Issue decisions, and the
owner's acceptance to the change record. Stop. The owner—not the agent—decides
whether to execute the primary apply command shown at the top of this file.

### 5. Apply once, with the journal intact

After explicit approval, the owner runs the primary command once. `migrate`
creates a unique `.savepoint/.migration/<opID>/` journal, backs up each path
before replacement/removal, stages and verifies every output, installs the
verified files, verifies them again, removes an original only after its archive
copy is verified, and activates `schema_version: 2` last. The schema activation
is the commit point. Do not interrupt or edit the project while this operation
is running.

Expected writes are the exact paths in the reviewed plan:

- `.savepoint/config.yml` changes to `schema_version: 2` only at the final
  activation boundary;
- the V2 router, `.savepoint/Idea.md`, objective/task/check/issue/release
  records, and any planned documents are installed at their identity paths;
- `.savepoint/archive/v1/` receives the 544 byte-preserved legacy sources in
  the reconciled plan; the retired v1.3 tree is already preserved separately
  under `.savepoint/archive/retired/v1.3/` and is not migration input;
- `.savepoint/migrations/v1-to-v2.yml` records the operation ID, source hashes,
  identity mapping, archives, prerequisites, waived references, and decisions;
- the operation journal's `backup/`, `staging/`, and `operation.yml` remain
  available while the publish is in progress and are removed only according to
  the successful-operation cleanup contract.

No file outside the reviewed plan is an expected write. Authored bodies and
source bytes are preserved in the archive/manifest mapping; do not hand-edit
generated records during the operation.

## Recovery paths

If the process is interrupted, leave `.savepoint/.migration/<opID>/` in place.
Do not remove its journal, staging, or backup files and do not start a fresh
plan. The owner may first request a read-only report:

```text
savepoint migrate . --recover
```

The report names the operation, the first unverified path (if any), and the
recoverable backup/staged copies. After reviewing that report and re-confirming
the source/backup boundary, the owner may resume the recorded operation with:

```text
savepoint migrate . --recover --apply
```

Recovery must refuse a changed source, changed installed output, multiple
pending operations, or a missing persisted plan. Preserve the refusal and
return to the backup/owner decision; never bypass it with a manual copy or a
new operation ID. If schema activation already completed, treat the project as
V2 and verify it; do not replay the V1 plan.

## Post-apply verification

The owner records the operation ID and performs all of the following before
calling the cutover complete:

```text
make build && make test
savepoint doctor .
savepoint board .
savepoint resume .
```

The expected checks are:

- `config.yml` declares schema 2 and the V2 index loads without a structural
  diagnostic;
- the router is on the V2 lifecycle and ordinary board/doctor/resume paths do
  not parse V1 records;
- the manifest exists and every source/archive pair has the recorded exact
  SHA-256; the 544 current archived files are present and no unplanned source
  vanished;
- there is no incomplete `.savepoint/.migration/` operation after successful
  cleanup; if one remains, use the recovery path and do not claim completion;
- Release output is the canonical `ResolveReleaseCutover` result, with every
  declared Release allowed and every material Issue explicitly accepted or
  resolved; a generic “migration complete” message is not Release evidence;
- a second migration preview reports that the project is already at schema 2
  and performs no write; no legacy reader is reachable from ordinary runtime;
- the final worktree diff, manifest, operation report, doctor output, board /
  resume output, and test output are attached to the handoff.

If any check fails, preserve the complete evidence and stop. Do not mark a task
done, close E50, or publish a release while a recovery, structural, Release,
Issue, freshness, or owner-acceptance gate remains unresolved.

## Rollback limits and excluded actions

There is no in-command rollback operation. The operation journal's backup is a
recovery aid for an interrupted publish and is not a durable post-success
rollback copy. Before apply, the durable rollback material is the named bundle,
worktree snapshot, patches, and checksums in
`../savepoint-cutover-backups/<clean-HEAD-sha>/`. If a completed migration must
be abandoned, the owner must stop all writers and restore a complete snapshot
under their normal repository-recovery procedure; never reverse only
`config.yml`, delete archives, or hand-delete a migration journal. Re-run the
full preview and gates after any restoration.

Migration does not publish artifacts, create tags, deploy software, update a
changelog, announce a release, or alter external services. Those actions are
out of scope and require separate explicit owner decisions.

## Handoff conclusion

The runbook, independent E51 audit reference, source/plan evidence, disposable
apply, and owner recovery boundary are prepared. The current live handoff is
**blocked** by the explicit two-blocker Release decision and the outstanding
owner review of candidate Issues. The maintainer must complete the current E50
Tasks, resolve those conditions, recapture the dry-run after committing the
reconciliation, review the backup, and then make the separate approval decision
at the stop point. Until that happens, the repository remains V1 and no live
migration command is authorized by this task.
