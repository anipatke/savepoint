---
type: audit-findings
audited: 2026-09-14
---

# Audit Findings: E41 Preserve migration source evidence

## Main Findings

### Verdict

**CLEAR** — re-audit after remediation. Both findings from the initial audit are
closed against the same frozen scope, re-verified with evidence independent of
the project's own tests. No finding remains, no criterion is unverified, and no
owner waiver is needed. The remediation touched only
`internal/data/testdata/migration/.gitattributes` and
`internal/data/testdata/migration/v1-history/manifest.yml`; every one of the 30
recorded fixture hashes is byte-for-byte unchanged, and no production Go source
was modified.

**CLEAR TO COMMIT/PUSH.** The epic's deliverable is two frozen V1 projects and
their characterization tests, all of which live in the repository, so nothing
material waits on post-push evidence. Epic closeout (marking the epic audited,
bumping `.savepoint/Design.md` `last_audited`, advancing the router) is still
yours to authorize.

### Closure Map

| Initial finding | Result | Re-verification |
|---|---|---|
| 1. v1-history unprotected from line-ending conversion | **Closed** | `.gitattributes` now carries `v1-history/** -text` alongside `v1-basic/** -text`. Re-ran the original reproduction: committed both fixture trees to a scratch repository, then checked them out with `core.autocrlf=true` and again with `core.eol=crlf`. `v1-history/…/v1.1/…/T001-shared.md` stayed at the manifest-recorded `8526383309a9…` under both settings — previously it became `82efc3ff1737…` — and `F001-awaiting-proof.md` stayed at `deedeb15c595…`. Swept all 32 files inside the two fixture trees with `git check-attr text`: every one reports `text: unset`. |
| 2. Manifest missing proof requirement and duplicate link | **Closed** | The v1-history manifest's F001/F002/F003 entries now carry `disposition` and `proof_needed_summary`, and F003 carries `duplicate_of: F001`. Verified by parsing the manifest with an independent YAML reader (not the test's own helper): all three annotation keys named in the file's header comment are now present in the entries, so the header is accurate too. |

No newly noticed item met the credible-blocker exception, so nothing was
admitted as a new blocking finding. Two small notes appear under Non-Blocking
Observations; neither affects this verdict.

### Materiality Summary

No materiality actions are required. Both previously rated findings are closed;
nothing is deferred, waived, or accepted with residual risk.

### What Is Proven / Not Proven

**Proven.** The two frozen projects exist at the exact planned paths with the
planned shapes: v1-basic carries a done baseline Task and a dependent
in-progress Task still using legacy `phase: implementation`, an unknown
top-level field, an unknown nested `metadata.reviewer.name`, a body comment,
multiline authored plan/criteria text, and CRLF throughout; v1-history carries
the same epic and short Task ID across v1 and v1.1, two release-distinct D001
defects, a fixed-but-unproven finding, an owner-waived finding, a duplicate
finding pointing at F001, one immutable run, an epic audit, and a short custom
Health-Check procedure. Discovery finds the intended release/epic/tasks. The
parser heals legacy `phase: implementation` to `stage: build` while the source
keeps its original key. The short dependency reference resolves release- and
epic-scoped — the v1.1 dependant lands on v1.1's in-progress record even when
v1's done record of the same short ID is in the candidate list, which the test
proves by asserting the resolved status, not merely a non-empty match.
`ParseRawFindingFile` preserves an unrecognized status, `ParseFindingFile` heals
it to `open`, and `DiagnoseFinding` reports it. `LoadAuditRegisterSet` keeps all
three findings distinct. Malformed frontmatter and broken references produce the
existing named parse error and the existing zero-value missing result — no
invented V2 diagnostic. Failure variants mutate only `t.TempDir()` copies, and
frozen bytes are re-verified afterwards. Byte-stability now holds on every
platform for both fixtures, and the reviewable manifest records each finding's
disposition, proof requirement, and duplicate link. No production code changed.

**Not proven / not required.** The frozen `router.md` and `config.yml` records
are byte-inventoried but no test characterizes how the current readers interpret
them; no acceptance criterion asks for that, so it stays an observation for E45
rather than a gap. The manifest's identity annotations are reviewed by a human,
not asserted by code — which the task states deliberately.

### Audit Evidence

- **Scope lock:** unchanged from the initial audit and reused without extension
  — T001 AC1–AC7 and T002 AC1–AC7; guardrails FS-05, DATA-01–05, CFG-02,
  TEST-01–08, TPL-01/03/04, ARCH-04, POL-01/02, STYLE-01–10; the 34 files under
  `internal/data/testdata/migration/` plus the two new `_test.go` files; the 14
  public surfaces from `ListReleases` through `SplitFrontmatterBody`; the axes
  of reference form, release and epic scope, raw-vs-typed-vs-raw-YAML
  representation, line-ending encoding, valid/invalid/unknown status, and
  present/absent optional artifact. Not-applicable cells and their reasons are
  unchanged (no renderer, no width or truncation logic, no writer, scaffold, or
  release path touched).
- **Coverage and workflow result:** every frozen cell reclassified; the two
  cells that previously failed — encoding × v1-history, and manifest-annotation
  coverage × v1-history findings — now pass. No axis, dependency layer, or
  acceptance interpretation was added. Independent oracles: a separate
  Python/`hashlib` sweep recomputed all 30 recorded hashes (0 mismatches), and a
  separate YAML reader confirmed the manifest annotations; neither reuses the
  test suite's own helpers. The only side-effecting operations remain the three
  temp-copy writes, each followed by a read, a parse, and a frozen-byte
  re-verification.
- **File reality and drift:** every file named across both tasks' Context Files
  and Context Logs resolves on disk; no phantom files, nothing deleted by the
  remediation. The fixture tree is still 34 files. Both Drift Notes correctly
  say "None": no production module was added, and the AGENTS.md Codebase Map
  entry for `internal/data` already covers this package. `.savepoint/Design.md`
  needs no delta beyond its `last_audited` bump at close.
- **Gates:** `go test ./internal/data -run 'MigrationSource|MigrationHistory'
  -count=1 -v` — pass, 7/7 named tests and 5/5 subtests, uncached. `make build`
  — pass. `make test` — pass, all packages. `go vet ./...` — clean.
  `git diff --check` — clean. `gofmt -l` on both new test files — clean.
  `git status` shows no modified production Go source.

### Guardrails Verification

- **Rule IDs checked:** FS-05, DATA-01, DATA-02, DATA-03, DATA-04, DATA-05,
  CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-06, TEST-07, TEST-08,
  TPL-01, TPL-03, TPL-04, ARCH-04, POL-01, POL-02, STYLE-01–10. No blockers
  open; no waiver requested or required.
- **Health check mode:** Full. The project has no `.savepoint/Health-Check.md`,
  so that procedure step is skipped per policy and its absence is not a finding;
  guardrails were applied directly.
- **Evidence:** CFG-02 is now satisfied — the platform-sensitive fixture
  encoding is explicit in `.gitattributes` and proven with `git check-attr` plus
  the two conversion-mode checkouts described in the closure map. TEST-08
  satisfied on Linux by the full gate, and the same gate would now hold on a
  Windows clone because the fixture bytes no longer shift. FS-05 satisfied:
  fixture paths are built with `filepath.Join` and compared via
  `filepath.ToSlash`. TEST-04 satisfied: all mutation goes to `t.TempDir()` and
  nothing reads the developer's real project. TEST-01/02/06 satisfied: each
  criterion maps to a named test with happy and failure paths and exact test
  names recorded. DATA-01–05 not at risk — no writer ran and no production
  loader changed. TPL-01 re-verified independently this pass: all nine canonical
  `agent-skills/*/SKILL.md` files remain byte-identical to their
  `templates/project/` copies. TPL-04 not applicable — these are testdata, not
  scaffold assets. ARCH-04 satisfied: fixtures and tests sit inside the existing
  `internal/data` purpose. POL-01/02 satisfied throughout.
- **File reality evidence:** see Audit Evidence above.
- **Waivers or unresolved findings:** none.

### Non-Blocking Observations

Carried forward unchanged; none is a finding, and none requires work inside this
epic.

- The frozen `router.md` and `config.yml` in both fixtures are inventoried but
  never parsed. A converter will have to read them, so characterizing the
  current router/config interpretation is worth a line in E45's plan.
- `internal/data/migration_source_test.go:24` refers to
  `TestMigrationSourceHistory*`; the tests T002 actually added are named
  `TestMigrationHistory*`.
- The inventory walk covers each fixture's `project/` subtree only, matching
  what the manifests record. A stray file dropped directly in `v1-basic/` or
  `v1-history/` would not be reported.
- `internal/data/testdata/migration/README.md` and `.gitattributes` are the only
  two files in the fixture directory without the `-text` attribute. Neither is
  hash-recorded in any manifest and neither is frozen source, so no test depends
  on their bytes.
- `gofmt -l internal/data` flags `router.go` and `task_test.go`. Both are
  pre-existing in `HEAD` (a missing trailing newline in `router.go`), untouched
  by this epic and outside its scope. Worth a one-line cleanup in unrelated work.

### Applied Outcome

**1. CRLF protection applied.** `internal/data/testdata/migration/.gitattributes`
now covers both `v1-basic/**` and `v1-history/**` with `-text`, and its comment
records that every new fixture directory needs a rule. Every file in both trees
reports `text: unset`, so Git line-ending conversion cannot rewrite their frozen
bytes under `core.autocrlf` or `core.eol`.

**2. Manifest review annotations applied.** The v1-history manifest records
F001's fixed-but-unverified disposition and its proof requirement, F002's owner
waiver and absence of a proof requirement, and F003's duplicate disposition,
`duplicate_of: F001`, and proof relationship. The annotations are additive and
change no recorded fixture hash.

## Code Style Review

Unchanged from the initial audit — the remediation touched no Go code. All ten
rules are Guideline severity; the three unchecked boxes are observations and
were never a blocker or a `NEEDS WORK` cause.

- [x] STYLE-01 **One job per file** — split files when responsibilities mix. Fixture-specific tests are split per fixture; shared manifest helpers live in one place.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases. Happy path, malformed frontmatter, malformed run, broken reference, and unknown status are all covered.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [ ] STYLE-05 **Build only what is needed** — no speculative abstractions. `migrationManifestFile` declares seven annotation fields (`ScopedID`, `RawStatus`, `RawPhase`, `ExpectedStageAfterParse`, `DependsOn`, `Encoding`, `ExpectedClassification`) that no assertion reads; `yaml.Unmarshal` ignores unknown keys, so the struct does not need them to parse — as the newly added `disposition`, `proof_needed_summary`, and `duplicate_of` keys demonstrate.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data. Every read and parse failure is reported with the path.
- [ ] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config. The manifest-vs-disk inventory walk is duplicated almost line-for-line between `TestMigrationSourceBasicInventory` and `assertHistoryInventoryMatchesManifest`; one `assertInventoryMatchesManifest(t, fixture)` helper would serve both.
- [ ] STYLE-08 **Comments explain why** — not what the code already says. Comment quality is otherwise high; one stale reference (see observations) names a test that does not exist.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic. Provenance and expectations live in `manifest.yml` and `README.md`.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes. Purely additive; no production code touched.
