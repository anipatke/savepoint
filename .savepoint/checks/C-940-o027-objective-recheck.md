---
id: C-940
scope: {kind: objective, id: O-027}
result: CLEAR
checked_by: {role: checker, session: o027-recheck-20261001}
executed_session: o027-issue-repair-session
checked_at: '2026-09-30T22:21:36Z'
reviewed:
  base_commit: 9fdc11f
  head_commit: 86ab71a67cc697901d4c4b9c3081d37472550458
  files:
    - internal/codehealth/classification.go
    - internal/codehealth/classification_test.go
    - internal/codehealth/config.go
    - internal/codehealth/errors.go
    - internal/codehealth/history.go
    - internal/codehealth/identity.go
    - internal/codehealth/identity_test.go
    - internal/codehealth/model.go
    - internal/codehealth/model_test.go
    - internal/codehealth/primitives.go
    - internal/codehealth/repository.go
    - internal/codehealth/repository_test.go
    - internal/codehealth/snapshot.go
    - internal/codehealth/storage.go
    - internal/codehealth/storage_test.go
    - internal/codehealth/testdata/valid-config-v1.json
    - internal/codehealth/testdata/valid-snapshot-v1.json
  dependencies: [go.mod, go.sum, Makefile]
issues: [I-085, I-086, I-087, I-088, I-089, I-090, I-091]
supersedes: C-939
---

# C-940: O-027 Full Objective Re-check

## Result and independence

**CLEAR.** This re-check of C-939 uses C-939's frozen scope lock, matrix M1–M9 and workflow inventory unchanged. It ran in a fresh checker conversation that neither built the O-027 Tasks nor made the I-085–I-091 repairs; the repairs are uncommitted working-tree changes on top of 86ab71a, identified by the source hashes below. No production code, tests, Task statuses, router selection or Design were changed by this Check. The executed_session value is a collective label: the repair session did not record its own ID in the Issues.

## Closure map

| Issue | C-939 cell | Status | Evidence |
|---|---|---|---|
| I-085 | M1/M7 | Closed | C-939 harness `M7_severity_validation_integration` passes: total 0 with high or critical 1 fails validation (`severity counts total 1, more than the 0 vulnerabilities reported`) and `Assess` returns needs_attention; negative/fractional counts fail with named `ErrIncompatibleUnit`. `hardBlocker` now reads severity before total (classification.go). New `TestVulnerabilitySeverityCountsAreValidated`. |
| I-086 | M1 | Closed | Harness `M1_missing_numbers` passes: omitted and null `value.number` return `ErrMissingValue` while preserving the original ID. Independent probe: null/omitted detail number via `DecodeSnapshot`, omitted/null/empty thresholds via `DecodeConfig` all return `ErrMissingValue`; explicit zero value and zero thresholds round-trip with an unchanged ID. |
| I-087 | M1/M3 | Closed | Harness `M1_trailing_delimiters` passes for all 8 decoder cases. Independent probe: tails `,` `x` `null` `"s"` `:` NUL rejected by both decoders; `Store.LoadConfig` rejects a trailing `]`. Whitespace tails still accepted. |
| I-088 | M2 | Closed | Harness `M2_evidence_tie_permutation` passes. Independent probe: five permutations of a three-way path/line tie plus a no-line entry keep ID and validate. Canonical order is now path, line, note (identity.go). |
| I-089 | M1/M5 | Closed | Harness `M5_malformed_pattern` now returns `ErrInvalidConfig` from `InputScope.Validate`, `Config.Validate` and `ObserveRepository`. Independent probe: exclude `[` rejected; valid `src/f[0-9].go`, `**/*.go`, `src/*`, `src` still see content changes. |
| I-090 | M5 | Closed | Harness `M5_colon_filename` passes (tracked edit → Dirty true, fingerprint changes). Independent probe: untracked `d:x/y.go` included by `d?x` changes fingerprint; `git rm` of a tracked colon file is detected. Unsafe relevant paths now return `ErrInputUnreadable` instead of vanishing; Windows colon rejection kept (`TestSafeRelativeByPlatform`). |
| I-091 | M9/M5 | Closed | Native windows/amd64 run of the codehealth test binary: exit 0, 437 PASS, 10 SKIP, 0 FAIL; `TestObserveUnusualFilenames` PASS. Native per-package run of every test package: all exit 0 (details below). |

No new Issues. Every C-939 criterion previously classified Issue is now Proven by the evidence above; criteria C-939 classified Proven were re-exercised by the harness's `M3_M4_M8_stored_history` (0–12 sweep), `M1_bounds_and_shapes` and `M6_independent_ancestry_sequence` subtests (all pass) and by the full gate. All O-027 success conditions 1–6 are Proven. No materiality actions are required.

## Admission ledger

Each probe maps to a frozen C-939 cell: I-085 → M1 detail severity consistency / M7 contradictory total-severity; I-086 → M1 missing/null; I-087 → M1 trailing JSON delimiters, M3 loader path; I-088 → M2 evidence same path/line permutation; I-089 → M1/M5 malformed glob; I-090 → M5 unusual POSIX colon filename; I-091 → M9 native Windows evidence. No probe outside the lock was treated as blocking.

## Commands and platform evidence

Run 2026-09-30 UTC (2026-10-01 local), Go 1.26.2 linux/amd64 and Go 1.26.2 windows/amd64:

- `go test -overlay <scratch>/overlay.json -count=1 -run TestO027Independent -v ./internal/codehealth` with C-939's embedded harness unchanged: **PASS**, every subtest.
- `go test -overlay <scratch>/overlay2.json -count=1 -run TestRecheckProbes -v ./internal/codehealth` with the harness below: **PASS**. (First run failed on two probe-authoring mistakes: a summary/result count mismatch, and a literal `d:x` include pattern, which config validation rejects by design; both probes were corrected before this result.)
- `make test-full`: **PASS**, exit 0, fresh; host tests plus linux/darwin/windows builds.
- `go test -race -count=1 ./internal/codehealth`: **PASS**. `go vet ./internal/codehealth`: **PASS**. `git diff --check`: **PASS**.
- Native Windows: each package's test binary cross-compiled with `GOOS=windows GOARCH=amd64 go test -c`, run from Windows PowerShell with `-test.v -test.count=1` and cwd at the package directory. Results: root 129 pass/2 skip; cmd 54; board 14; board/v2 328; buildtool 25; codehealth 437 pass/10 skip; doctor 95; init 292/2 skip; migrate 259/6 skip; resume 71; styles 8; data 899 pass/8 skip. All exit 0, no FAIL.
- `internal/data` initially failed once when run with cwd on the `\\wsl.localhost` UNC share: `testing.Chdir: chdir .` could not restore a UNC cwd in cleanup. Rerun from a temporary local copy at `C:\Users\Public\savepoint-wincheck` (deleted afterwards): exit 0, 899 pass. This is a runner-location artifact, not a code defect; data is outside O-027 changes.
- `go test ./...` with native Windows Go directly on the UNC share could not start (`RLock go.mod: Incorrect function`); the per-package binaries above replace it.
- Windows Defender quarantined the first test binary written to `C:\Users\Public`; the owner added a folder exclusion for the checker scratch directory, and the binary was rebuilt there.
- Native Windows codehealth skips: symlink privilege unavailable (TestObserveSymlinkHashesTargetWithoutFollowing, four TestSymlinkEscapesAreRefused subtests), directory-permission scenarios (four tests/subtests), and `TestObserveFingerprintsPOSIXColonFilenames` (colon illegal on Windows). These match C-939's disclosed skip classes and are not claimed as passing Windows safety tests.

Scratch harnesses, binaries and logs live in the checker's session scratchpad and are disposable; the harness below is the durable reproduction artifact.

## Guardrails and reconciliation

DATA-03 / CFG-01 now pass for I-085/I-086/I-087/I-089; CFG-02 / CFG-03 pass for I-090/I-091 with native Windows evidence; TEST-01/02 gaps are covered by new named tests and independent probes; TEST-08 full gate passes. go.mod, go.sum and Makefile are unchanged since 9fdc11f (DEP-01). C-939's Design.md reconciliation note (no Code Health section in Design) is unchanged and remains a planning follow-up, not a blocker.

## Observations (non-blocking)

- Scope patterns that name a literal colon (for example `d:x`) are refused by path-form validation with a named error; a wildcard such as `d?x` matches such a directory on POSIX. No silent loss results.
- `Assess` on unvalidated input with a negative severity count (total 2, high -1) returns watch; validation rejects that record first, consistent with C-939's M7 note that direct models are expected to pass validation before use.
- On Windows, a checkout containing a POSIX colon path that is relevant to scope now makes `ObserveRepository` return `ErrInputUnreadable` rather than silently omitting it. This is the intended explicit failure.
- C-939's unexecuted OS fault-timing limitations (M3/M4) and the T-054 User Check wording observation carry forward unchanged.
- T-053 and T-054 declare `owner_validation.required`; this Check records no owner acceptance. The owner decides Objective acceptance.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Reviewed source fingerprints

```text
adaffc3de752d6ac17471a32ed3937cef73ee788ea29c8b0fad17766c019197e  internal/codehealth/classification.go
2862c8660cc716968d0a534ba327e22aeae548029696a3ec9e709c1c9f52564c  internal/codehealth/classification_test.go
99f347b1afd95e680645ff85ab76f66809e0e1d2bd3eafa68ada94390a10f69e  internal/codehealth/config.go
b9422b1b3dfbe1aec51c38c4a4c6d933dafc8f221009f9143de96ac89bb88880  internal/codehealth/errors.go
3ce4cbb64e8ea53390aa8768e6278ebd1e65ec28394c7e23172891c74855d431  internal/codehealth/history.go
e8c95dcd9b1c07efb359f8bc62616fe6ad53e6f47f7c4772bf4bd4d8051e4a9d  internal/codehealth/identity.go
ecd3573050ef1d7b782c9cc90bc5330ab8eb9fc7e1bfa2d5144c0394136eacba  internal/codehealth/identity_test.go
e53c24da5e762e232e8f6bb3af96dcc53dea2fe763e3976c8cd7b0bdb1f71c90  internal/codehealth/model.go
12cf22a610f346a3971a6ec6d6ab248ee8a0c2167ec15a74b581b7d73d312392  internal/codehealth/model_test.go
4a2e5d0daea608aaa232073f687dc350e4ae9fe640a2ee8b3846dad0ce90d89b  internal/codehealth/primitives.go
983d7c5d097b536fccfe824b1a95957a11cad98cd645f5e68f93be37d001d219  internal/codehealth/repository.go
447e4900e2341e9addba7f21fd6aa8d2a756b6c179460eb94e83261b1ee30946  internal/codehealth/repository_test.go
7145de2b8c97ec21756834830fde68be30a6b4226bb8ed129352a692b8efff69  internal/codehealth/snapshot.go
64633bcbf20f9ec6f6f47a8e9113da13e7ef1a5c0b195e2fa3c28e4ec0373b06  internal/codehealth/storage.go
180086661f39715904b1971ff23a94de37deeb1bdeb357afaa9eed88da818122  internal/codehealth/storage_test.go
83f5d16ff2641d34e5a3087fcae6e86773ee3acf0f7d1a761d02e625a2a498b1  internal/codehealth/testdata/valid-config-v1.json
cc0d7df9cacbd5c0a93630785b65a25f56e8205610a660425a35730e9ab3ea06  internal/codehealth/testdata/valid-snapshot-v1.json
```

## Re-check probe harness

Used with a Go overlay mapping `internal/codehealth/recheck_probe_test.go` to this file, alongside C-939's harness; not added to source.

```go
package codehealth

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestRecheckProbes(t *testing.T) {
	// I-087: more tail classes through both decoders and store loaders.
	t.Run("I087_tails", func(t *testing.T) {
		for _, tail := range []string{",", "x", "null", " \"s\"", " :", "\x00"} {
			for name, fx := range map[string]string{"config": "valid-config-v1.json", "snap": "valid-snapshot-v1.json"} {
				raw := append(readFixture(t, fx), tail...)
				var err error
				if name == "config" { _, err = DecodeConfig(raw) } else { _, err = DecodeSnapshot(raw) }
				if !errors.Is(err, ErrMalformedRecord) { t.Errorf("%s tail %q: %v", name, tail, err) }
			}
		}
		st, root := newProject(t)
		p := filepath.Join(root, ".savepoint", "health", "config.json")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil { t.Fatal(err) }
		writeFile(t, p, string(readFixture(t, "valid-config-v1.json"))+"]")
		if _, err := st.LoadConfig(); !errors.Is(err, ErrMalformedRecord) { t.Errorf("LoadConfig with ]: %v", err) }
	})
	// I-086: detail and threshold omission through whole-record decoders; explicit zero kept.
	t.Run("I086_whole_records", func(t *testing.T) {
		s := validSnapshot(t)
		r := vulnResult(0, Detail{Key: "high", Number: 0}, Detail{Key: "critical", Number: 0})
		a := Assess(OriginOfficial, r, nil, nil)
		s.Results = []CapabilityResult{r}
		s.Summary = Summary{Overall: a.Classification, Capabilities: []CapabilitySummary{a.Summary()}}
		s = reseal(s)
		if err := s.Validate(); err != nil { t.Fatal(err) }
		raw, _ := json.Marshal(s)
		got, err := DecodeSnapshot(raw)
		if err != nil || got.Results[0].Value.Number != 0 || got.ID != s.ID { t.Fatalf("explicit zero round trip: %v", err) }
		for _, mut := range []string{`"number":0}`} {
			// drop the number of the last detail only
			i := strings.LastIndex(string(raw), mut)
			bad := string(raw[:i]) + `"number":null}` + string(raw[i+len(mut):])
			if _, err := DecodeSnapshot([]byte(bad)); !errors.Is(err, ErrMissingValue) { t.Errorf("null detail: %v", err) }
			bad = string(raw[:i-1]) + "}" + string(raw[i+len(mut):])
			if _, err := DecodeSnapshot([]byte(bad)); !errors.Is(err, ErrMissingValue) { t.Errorf("omitted detail (%s): %v", bad[len(bad)-200:], err) }
		}
		cfg := string(readFixture(t, "valid-config-v1.json"))
		for _, bad := range []string{
			strings.Replace(cfg, `{"good": 80, "watch": 60}`, `{"good": 80}`, 1),
			strings.Replace(cfg, `{"good": 80, "watch": 60}`, `{"good": null, "watch": 60}`, 1),
			strings.Replace(cfg, `{"good": 80, "watch": 60}`, `{}`, 1),
		} {
			if bad == cfg { t.Fatal("mutation missed") }
			if _, err := DecodeConfig([]byte(bad)); !errors.Is(err, ErrMissingValue) { t.Errorf("threshold: %v", err) }
		}
		if _, err := DecodeConfig([]byte(strings.Replace(cfg, `{"good": 0, "watch": 2}`, `{"good": 0, "watch": 0}`, 1))); err != nil { t.Errorf("explicit zero thresholds: %v", err) }
	})
	// I-088: three-way ties, all permutations, exact duplicates.
	t.Run("I088_permutations", func(t *testing.T) {
		ev := []EvidenceRef{{Path: "r.json", Line: 1, Note: "b"}, {Path: "r.json", Line: 1, Note: "a"}, {Path: "r.json", Line: 1, Note: "c"}, {Path: "r.json", Note: "z"}}
		s := validSnapshot(t); s.Results[0].Evidence = ev; s = reseal(s)
		if err := s.Validate(); err != nil { t.Fatal(err) }
		perms := [][]int{{0,1,2,3},{3,2,1,0},{1,0,3,2},{2,3,0,1},{3,0,2,1}}
		for _, p := range perms {
			c := s; c.Results = slices.Clone(s.Results)
			var e []EvidenceRef
			for _, i := range p { e = append(e, ev[i]) }
			c.Results[0].Evidence = e
			if c.ComputeID() != s.ID || c.Validate() != nil { t.Errorf("perm %v changed identity", p) }
		}
	})
	// I-089: exclusion + config pattern propagation into ObserveRepository and valid patterns still match.
	t.Run("I089_valid_patterns_match", func(t *testing.T) {
		root := newRepo(t)
		write(t, root, "src/f1.go", "x")
		for _, p := range []string{"src/f[0-9].go", "**/*.go", "src/*", "src"} {
			o := observe(t, root, InputScope{Include: []string{p}})
			write(t, root, "src/f1.go", "y"+p)
			o2 := observe(t, root, InputScope{Include: []string{p}})
			if o.Identity.InputFingerprint == o2.Identity.InputFingerprint { t.Errorf("valid pattern %q did not see change", p) }
		}
		if _, err := ObserveRepository(context.Background(), GitRunner{}, root, InputScope{Exclude: []string{"["}}); !errors.Is(err, ErrInvalidConfig) { t.Errorf("exclude [: %v", err) }
	})
	// I-090: colon dir, untracked colon file in scope, staged deletion.
	t.Run("I090_colon_variants", func(t *testing.T) {
		root := newRepo(t)
		write(t, root, "d:x/y.go", "one")
		o1 := observe(t, root, InputScope{Include: []string{"d?x"}})
		write(t, root, "d:x/y.go", "two")
		o2 := observe(t, root, InputScope{Include: []string{"d?x"}})
		if o1.Identity.InputFingerprint == o2.Identity.InputFingerprint { t.Error("untracked colon dir change ignored") }
		git(t, root, "add", "-A"); git(t, root, "commit", "-qm", "c")
		clean := observe(t, root, InputScope{})
		git(t, root, "rm", "-q", "d:x/y.go")
		gone := observe(t, root, InputScope{})
		if !gone.Identity.Dirty || gone.Identity.InputFingerprint == clean.Identity.InputFingerprint { t.Error("deleted colon file unnoticed") }
	})
}
```
