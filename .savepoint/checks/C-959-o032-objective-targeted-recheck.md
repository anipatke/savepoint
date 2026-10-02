---
id: C-959
scope: {kind: objective, id: O-032}
result: NEEDS WORK
checked_by: {role: checker, session: recheck-o032-20261002c-independent}
executed_session: repair-o032-20261002b
checked_at: '2026-10-02T06:50:00Z'
health_snapshot: sha256:0aa2073976fe3e6ef67bfbb7cc3d2288f18067b56e461daaaa78f50bb420397d
reviewed:
  base_commit: d00543d5252fc5557aeba067995427236da2823f
  head_commit: d2791fa1e53a806cb31a9650b4c89d2b24319084
  files:
    - AGENTS.md
    - CHANGELOG.md
    - README.md
    - internal/codehealth/dashboard.go
    - internal/codehealth/discovery_test.go
    - internal/codehealth/runner.go
    - internal/codehealth/runner_safety_test.go
    - internal/codehealth/setup_test.go
    - internal/codehealth/window.go
    - internal/codehealth/window_test.go
  dependencies: []
issues: []
supersedes: C-958
---

# C-959: O-032 Targeted Recheck of I-121 and I-118

## Result and authority

NEEDS WORK. I-121 is technically proven repaired. I-118's original finding is proven, but its replacement repair no longer keeps two things its Proof Needed requires: the documented old-body-damage boundary and the measured cost. It also leaves T-098's accepted criteria, the T-097 decision and Design.md describing a design the code no longer has.

This conversation started after `/clear` and performed no implementation, so it is independent from repair-o032-20261002b. The owner asked for this recheck of I-121 and I-118. Under the convergence method it is the one targeted recheck after C-958's targeted remediation, so the next step is the owner's decision, not another autonomous repair cycle. No new Issue is opened, no Task is retreated, and the router, Design, acceptance criteria and implementation are unchanged. C-957 and C-958 are immutable.

## Frozen scope and admission ledger

C-957's frozen scope lock and C-958's admission ledger are incorporated unchanged. Remediation paths since C-958 are bc21805 (I-121: `tailBuffer.dropped`/`Text`, `cutUserinfo`), 8da136a (I-118: `parseHead`, `readPrefix`, `headBytes` and `snap` removed; `readHead` reads the whole file through `readRecord` and takes `created_at`/`origin` with `json.Unmarshal`) and d2791fa (test-only: `os.Lstat` in two read-only tests). This ledger was written before any probe ran:

| Item | Prior Issue/claim | Exact frozen cell | Allowed result |
|---|---|---|---|
| C-958 harness at 64KiB −1/exact/+12 through real ExecRunner and Collect | I-121, C-958 failure; bc21805 | E secret-safe failure output, 64KiB stderr −1/exact/+1 overflow | Pass or I-121 remains open |
| Every overflow that cuts inside the URI (+1…+len+2), newline and non-whitespace filler, and >2× cap (batch trim) | I-121 proof "normal and oversized stderr" | Same E overflow cell; no new secret class | Same |
| Useful trailing failure text, UTF-8 validity, MaxReasonLen bound | I-121 proof | E control/UTF-8/bounded reason | Same |
| Repeated created_at/origin before and after the normal fields; newest, oldest and window-boundary file | I-118, 8da136a | W valid duplicate headers, order/full-decoder/newest identity | Pass or I-118 remains open |
| compact/reordered/512prefix/wrongtype/trailing/symlink representations | I-118 adjacency | W recorded representation cells | Pass or original in-scope failure |
| Damage inside/outside the window, and no write | I-118 proof "Retain the documented old-body-damage boundary"; T-098 DW4 | W damage and no-write cells | Pass or I-118 remains open |
| single/multi/heavy n=10/100/1,000 benchmarks | I-118 proof "and measured cost"; T-098 DW6 | F benchmark cells | Pass or I-118 remains open |
| I-117 36-cell official/manual count matrix (window.go changed) | I-117 regression | W official count 0–11 × leading/trailing/interleaved | Pass or matching in-scope failure |
| Fresh `make ci`, checksums, current-head native CI, official health after gate | Release gate | R/P | Pass, blocks, or unverified |

## Original finding closure map

| Issue | Technical assessment | Disposition |
|---|---|---|
| I-121 | **Proven.** 130 of 130 cells pass through the real ExecRunner, Collect and the saved snapshot: 64KiB −1, exact, +12, every overflow from +1 to +len(URI)+2, and 2×+1, 2×+20 and 3×+20 the cap, each with newline and with non-whitespace filler. No credential fragment (`secret` or `alice`) remains in the reason or the immutable snapshot. Every reason ends with the real `failure` text, is valid UTF-8 and stays within MaxReasonLen. C-958's own harness, rerun, passes in all three of its cells. Collect.go:380 is the only consumer of `ToolResult.Stderr`, and it re-sanitizes after redaction, so no later cut can expose a removed credential. | Technically proven. It stays `open` only because Issue capture requires a CLEAR Check for `verified`, and this record is NEEDS WORK. |
| I-118 | **Original finding proven; remaining proof not met.** A valid repeat of `created_at`/`origin` before the normal fields (newest, oldest, boundary) and a same-value repeat after them give the same newest identity, window order and dashboard history as LoadSnapshots. The representation cells match C-957. The I-117 matrix passes 36 of 36. However: (a) **damage boundary:** a truncated or trailing-garbage snapshot outside the window now fails `LoadWindow` and `LoadDashboard` with `is not readable JSON`; the 512-byte head read tolerated both, and T-098 DW4 says "a damaged body outside the window does not fail it"; damage that keeps the JSON valid (renamed field, wrong results value) is still tolerated; (b) **measured cost:** see the cost section below. | Remains open; owner decision required (see Remedy). |

Probe classification: four cells in my first duplicate harness repeated a field after the normal fields with a *different* value. That makes the file invalid for the full decoder (identity mismatch, or a manual snapshot marked permanent), so they test damage, not a valid representation. They are not counted as failures. The window tolerates those files only when the changed time places them outside the window, which is the approved outside-window boundary. The same-value cells above are the valid representation for this case, and they pass.

## Measured cost against T-098 and the I-118 proof

`go test ./internal/codehealth -run '^$' -bench HealthHistoryLoadDashboard -benchmem -count=1` on go1.26.2 linux/amd64 with a Ryzen 7 7800X3D on WSL2 and a warm cache. These are single-run diagnostics:

| LoadDashboard | T-098 evidence (head read) | Now (whole-file read) |
|---|---|---|
| single n=1,000 | 13.2 ms, 4.9 MB | 34.0 ms, 8.1 MB, 29.1k allocs |
| multi n=1,000 | 18 ms | 129 ms, 35.4 MB |
| heavy n=100 / n=1,000 | 22.8 / 34.2 ms | 76.7 / 561.6 ms |
| heavy n=1,000 memory | 23.7 MB | 168.0 MB, 86.9k allocs |

The single n=1,000 load stays inside DW6's 250 ms budget. But each file beyond the window now costs about 148 KB, which is (168.0 − 34.75 MB) / 900 files, close to the 139 KB heavy snapshot size. T-098 evidence stated the per-file constant as about 3.5 KB, independent of snapshot size. T-098's Outcome says "load cost no longer grows with snapshot size", and that is no longer true. DW1 still names a 512-byte head read with a full-decode fallback. T-097's owner-approved decision and Design.md line 30 ("found by reading each file's head for time and origin") still describe the bounded read. README, CHANGELOG and AGENTS.md were updated to match the code, and the CHANGELOG numbers match these measurements.

I-118's history calls this an owner-directed simplification, but there is no recorded owner exception and no amended criterion. A checker may not treat that note as either one, and may not edit the criteria.

## Gate, platform and health evidence

- Fresh `make ci` on 2026-10-02 (go1.26.2 linux/amd64): exit 0. It ran `go run ./internal/buildtool test -reports -json -count=1 ./...` (uncached), six target builds, dist, the npm wrapper and the package dry-run. `sha256sum -c checksums.txt` passes all six archives. The temporary harness was removed before the gate, so the gate ran on the reviewed tree.
- `git diff --check d00543d..HEAD`: code is clean. The only hits are trailing spaces inside C-958's pasted benchmark log.
- Native evidence: run https://github.com/anipatke/savepoint/actions/runs/36974396208 on exact head d2791fa1e53a806cb31a9650b4c89d2b24319084 has windows-tests (110734974216) and ci (110734974304) both successful. That is one green run of the `os.Lstat` test change, which is test-only and does not prove the NTFS-lag cause. I read the run metadata only and did not dispatch anything.
- After the gate, the official `./savepoint health check O-032` exited 0 and created the snapshot cited above. Code Health does not block clearance. The optional OSV result again reports two unknown-severity groups that need review; this is advisory and unchanged from C-958.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-118 (a) damage boundary | Low: saves are atomic, so a truncated or garbage-tailed old file needs outside damage | Low: the dashboard shows an error instead of hiding the file, and the collector already fails on the same history | Low | Owner decision together with (b) |
| I-118 (b) cost and records | Medium: any project with large snapshots and long history | Medium: 0.5 s and 168 MB per dashboard load at the heavy end, paid two or three times per board open (T-097 finding); it also contradicts accepted criteria and Design | Medium | Owner chooses: amend T-097, T-098 and Design through `savepoint-design` to the new design (covering DW1, DW4 and DW6), record an explicit owner exception naming those criteria, or restore a bounded read under I-118 |

## Design reconciliation and observations

- Design.md line 30 still describes the head read and "not read" wording. The code, README and AGENTS.md say "not checked". This drift is part of I-118 (b) and is the planner's to reconcile, not the checker's.
- `TestLoadDashboardReadsNoSnapshotBodyOutsideTheWindow` now proves that bodies outside the window are not validated, but it no longer proves they are not read. Its name overstates what it checks. Non-blocking.
- Untracked `.coverage.out.*` and `.go-test.json.*` files were present at the repo root before this run and were left alone.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — no permanent test covers a JSON-invalid body outside the window, which is the branch whose behaviour changed (`window.go:readHead`).
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — the head read now uses the decoder's own JSON rule.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

Style is advisory and did not affect the verdict.

## Owner handoff

The convergence limit has been reached. The owner decides on I-118, choosing one of the three options in the Materiality table, and may then request a final Check. I-121 needs no further repair; it will close as `verified` in the next CLEAR Check, or the owner may resolve it as `accepted` from the board. I-115, I-116, I-117, I-119 and I-120 keep C-958's proven assessments, and their own code is unchanged since then (window.go changed for I-118, and the I-117 matrix was rerun and passes).

## Reproducible evidence appendix

The scratch harness was deleted from the tree after the run.

### o032_c959_probe_test.go

```go
package codehealth

// Temporary C-959 checker harness. Deleted after the run; reproduced in the
// Check record appendix.

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

const c959Secret = "o032-secret-for-probe"
const c959URI = "https://alice:" + c959Secret + "@proxy.example.invalid"

// E cell: the C-958 harness (-1/exact/+12) plus every overflow that cuts inside
// the URI, and overflows past 2x the cap (batch trim), through ExecRunner and
// Collect with the persisted snapshot.
func TestC959CredentialCapMatrix(t *testing.T) {
	var sizes []int
	sizes = append(sizes, maxStderrBytes-1, maxStderrBytes, maxStderrBytes+12)
	for over := 1; over <= len(c959URI)+2; over++ {
		sizes = append(sizes, maxStderrBytes+over)
	}
	sizes = append(sizes, 2*maxStderrBytes+1, 2*maxStderrBytes+20, 3*maxStderrBytes+20)
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, size := range sizes {
		for _, filler := range []string{"\n", "x"} {
			t.Run(fmt.Sprintf("%d/%q", size-maxStderrBytes, filler), func(t *testing.T) {
				t.Setenv("C959_SIZE", strconv.Itoa(size))
				t.Setenv("C959_FILLER", filler)
				cfg := cfgOf(lizardInstance("local", exe, "a/**"))
				cfg.Capabilities[0].Args = []string{"-test.run=^TestC959CredentialTool$"}
				root := project(t)
				c := collect(t, root, cfg, Readers{ProviderLizardCSV: okReader(1, UnitCCN)}, ExecRunner{})
				r := collectedFor(t, c, CapabilityComplexity, "local").Result
				if r.Outcome != OutcomeFailed {
					t.Errorf("outcome %s", r.Outcome)
				}
				if strings.Contains(r.Reason, "secret") || strings.Contains(r.Reason, "alice") {
					t.Errorf("credential fragment persisted in reason %q", r.Reason)
				}
				if !strings.HasSuffix(r.Reason, "failure") || !utf8.ValidString(r.Reason) || len(r.Reason) > MaxReasonLen {
					t.Errorf("useful/valid/bounded tail lost: %q", r.Reason)
				}
				snaps, err := NewStore(root).LoadSnapshots()
				if err != nil {
					t.Fatal(err)
				}
				data, _ := json.Marshal(snaps)
				if strings.Contains(string(data), c959Secret) {
					t.Error("immutable snapshot contains credential")
				}
			})
		}
	}
}

func TestC959CredentialTool(t *testing.T) {
	sz := os.Getenv("C959_SIZE")
	if sz == "" {
		return
	}
	size, _ := strconv.Atoi(sz)
	filler := os.Getenv("C959_FILLER")
	// The URI first, so an overflow cuts into it from the front.
	pad := size - len(c959URI) - len(" failure")
	os.Stderr.WriteString(c959URI + strings.Repeat(filler, pad) + " failure")
	os.Exit(127)
}

// W cell: repeated created_at/origin before and after the normal fields.
func TestC959DuplicateHeaderMatchesFullDecoder(t *testing.T) {
	for _, key := range []string{"created_at", "origin"} {
		for _, where := range []string{"before", "after"} {
			for _, target := range []string{"newest", "oldest", "boundary"} {
				t.Run(key+"/"+where+"/"+target, func(t *testing.T) {
					store, root, _ := windowHistory(t, 40)
					all, err := store.LoadSnapshots()
					if err != nil {
						t.Fatal(err)
					}
					s := all[len(all)-1]
					switch target {
					case "oldest":
						s = all[0]
					case "boundary":
						win, _, _ := store.LoadWindow()
						s = win[0]
					}
					path := snapshotFileOf(t, root, s)
					data, _ := os.ReadFile(path)
					field := `"created_at":"2000-01-01T00:00:00Z"`
					if key == "origin" {
						field = `"origin":"manual"`
					}
					var altered string
					if where == "before" {
						altered = strings.Replace(string(data), "{", "{"+field+",", 1)
					} else {
						i := strings.LastIndex(string(data), "}")
						altered = string(data[:i]) + "," + field + "}"
					}
					if err := os.WriteFile(path, []byte(altered), 0o644); err != nil {
						t.Fatal(err)
					}
					full, ferr := store.LoadSnapshots()
					win, _, werr := store.LoadWindow()
					if (ferr == nil) != (werr == nil) {
						t.Fatalf("full err %v window err %v", ferr, werr)
					}
					if ferr != nil {
						return
					}
					if full[len(full)-1].ID != win[len(win)-1].ID {
						t.Errorf("newest differs")
					}
					tail := full[len(full)-len(win):]
					if !reflect.DeepEqual(snapshotIDs(tail), snapshotIDs(win)) {
						t.Errorf("window order differs from full order")
					}
					got, err := LoadDashboard(root)
					if err != nil {
						t.Fatal(err)
					}
					fd := fullDashboard(t, root)
					if got.SnapshotID != fd.SnapshotID || !reflect.DeepEqual(got.History, fd.History) {
						t.Errorf("dashboard newest/history differs")
					}
				})
			}
		}
	}
}

// W representation cells, original C-957 set.
func TestC959RepresentationMatrix(t *testing.T) {
	for _, shape := range []string{"compact", "reordered", "512prefix", "wrongtype", "trailing", "symlink"} {
		t.Run(shape, func(t *testing.T) {
			store, root, _ := windowHistory(t, 8)
			all, _ := store.LoadSnapshots()
			path := snapshotFileOf(t, root, all[len(all)-1])
			data, _ := os.ReadFile(path)
			valid := true
			switch shape {
			case "compact":
				var v any
				json.Unmarshal(data, &v)
				data, _ = json.Marshal(v)
			case "reordered":
				var v map[string]json.RawMessage
				json.Unmarshal(data, &v)
				data, _ = json.Marshal(v)
			case "512prefix":
				data = append([]byte(strings.Repeat(" ", 512)), data...)
			case "wrongtype":
				data = []byte(strings.Replace(string(data), `"origin": "official"`, `"origin": 1`, 1))
				valid = false
			case "trailing":
				data = append(data, ']')
				valid = false
			case "symlink":
				target := path + ".outside"
				os.WriteFile(target, data, 0o644)
				os.Remove(path)
				if err := os.Symlink(target, path); err != nil {
					t.Skip(err)
				}
				valid = false
			}
			if shape != "symlink" {
				os.WriteFile(path, data, 0o644)
			}
			_, _, err := store.LoadWindow()
			_, ferr := store.LoadSnapshots()
			if (err == nil) != valid || (ferr == nil) != valid {
				t.Errorf("valid=%v window err=%v full err=%v", valid, err, ferr)
			}
		})
	}
}

// W damage cells against T-098 Done When 4: damage outside the window does not
// fail the load; inside it does; nothing is written.
func TestC959DamageOutsideWindow(t *testing.T) {
	for _, kind := range []string{"renamed-field", "truncated", "bad-results-value", "trailing-garbage"} {
		for _, inside := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/inside=%v", kind, inside), func(t *testing.T) {
				store, root, _ := windowHistory(t, 40)
				all, _ := store.LoadSnapshots()
				s := all[0]
				if inside {
					s = all[len(all)-1]
				}
				path := snapshotFileOf(t, root, s)
				data, _ := os.ReadFile(path)
				switch kind {
				case "renamed-field":
					data = []byte(strings.Replace(string(data), `"results"`, `"resultz"`, 1))
				case "truncated":
					data = data[:len(data)*2/3]
				case "bad-results-value":
					data = []byte(strings.Replace(string(data), `"results": [`, `"results": [7,`, 1))
				case "trailing-garbage":
					data = append(data, []byte("garbage")...)
				}
				os.WriteFile(path, data, 0o644)
				before := dirState(t, root)
				_, _, err := store.LoadWindow()
				if inside && err == nil {
					t.Error("damage inside the window loaded")
				}
				if !inside && err != nil {
					t.Errorf("damage outside the window failed the load: %v", err)
				}
				if !reflect.DeepEqual(before, dirState(t, root)) {
					t.Error("load wrote")
				}
			})
		}
	}
}

func dirState(t *testing.T, root string) map[string]string {
	t.Helper()
	m := map[string]string{}
	entries, _ := os.ReadDir(root + "/.savepoint/health/snapshots")
	for _, e := range entries {
		b, _ := os.ReadFile(root + "/.savepoint/health/snapshots/" + e.Name())
		m[e.Name()] = string(b)
	}
	return m
}

// W count cells (I-117 regression): official 0–11 x leading/trailing/interleaved.
func TestC959WindowCountMatrix(t *testing.T) {
	for officials := 0; officials <= 11; officials++ {
		for _, shape := range []string{"leading", "trailing", "interleaved"} {
			t.Run(fmt.Sprintf("%d/%s", officials, shape), func(t *testing.T) {
				store, root := newProject(t)
				p := historyProfiles[0]
				cfg := historyConfig(p)
				store.SaveConfig(cfg)
				n := 0
				save := func(origin Origin) {
					s := historySnapshot(cfg, p, n)
					n++
					s.Origin = origin
					s.Retention = RetentionPrunable
					if origin == OriginOfficial {
						s.Retention = RetentionPermanent
					}
					s.ID = s.ComputeID()
					mustSave(t, store, s)
				}
				if shape == "leading" {
					save(OriginManual)
					save(OriginManual)
				}
				for i := 0; i < officials; i++ {
					save(OriginOfficial)
					if shape == "interleaved" {
						save(OriginManual)
					}
				}
				if shape == "trailing" {
					save(OriginManual)
					save(OriginManual)
				}
				if n == 0 {
					save(OriginManual)
				}
				all, _ := store.LoadSnapshots()
				win, cut, err := store.LoadWindow()
				if err != nil {
					t.Fatal(err)
				}
				if officials <= 10 {
					if cut || !reflect.DeepEqual(snapshotIDs(win), snapshotIDs(all)) {
						t.Errorf("inside-window promise broken")
					}
					got, _ := LoadDashboard(root)
					if !reflect.DeepEqual(got, fullDashboard(t, root)) {
						t.Errorf("dashboard differs")
					}
				} else if len(win) < 10 {
					t.Errorf("window lost officials")
				}
			})
		}
	}
}

// W cell, valid representation: the real created_at/origin repeated after the
// normal fields with the same value (the decoder still accepts the file).
func TestC959DuplicateSameValueAfter(t *testing.T) {
	for _, key := range []string{"created_at", "origin"} {
		for _, target := range []string{"newest", "oldest", "boundary"} {
			t.Run(key+"/"+target, func(t *testing.T) {
				store, root, _ := windowHistory(t, 40)
				all, _ := store.LoadSnapshots()
				s := all[len(all)-1]
				switch target {
				case "oldest":
					s = all[0]
				case "boundary":
					win, _, _ := store.LoadWindow()
					s = win[0]
				}
				path := snapshotFileOf(t, root, s)
				data, _ := os.ReadFile(path)
				v, _ := json.Marshal(s.CreatedAt)
				if key == "origin" {
					v, _ = json.Marshal(s.Origin)
				}
				i := strings.LastIndex(string(data), "}")
				os.WriteFile(path, []byte(string(data[:i])+`,"`+key+`":`+string(v)+"}"), 0o644)
				full, err := store.LoadSnapshots()
				if err != nil {
					t.Fatalf("decoder rejected same-value repeat: %v", err)
				}
				win, _, err := store.LoadWindow()
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(snapshotIDs(full[len(full)-len(win):]), snapshotIDs(win)) {
					t.Error("order differs")
				}
				got, _ := LoadDashboard(root)
				fd := fullDashboard(t, root)
				if got.SnapshotID != fd.SnapshotID || !reflect.DeepEqual(got.History, fd.History) {
					t.Error("dashboard differs")
				}
			})
		}
	}
}
```

### c959-probe summary

```text
CredentialCapMatrix passing cells: 130
--- PASS: TestC959CredentialCapMatrix (3.35s)
--- PASS: TestC959CredentialTool (0.00s)
--- FAIL: TestC959DuplicateHeaderMatchesFullDecoder (0.95s)
--- PASS: TestC959RepresentationMatrix (0.11s)
--- FAIL: TestC959DamageOutsideWindow (0.60s)
--- PASS: TestC959WindowCountMatrix (0.82s)
--- PASS: TestC959DuplicateSameValueAfter (0.50s)
FAIL
FAIL	github.com/opencode/savepoint/internal/codehealth	6.328s
FAIL
    o032_c959_probe_test.go:117: full err 363552ff89fc60ae03d7bc1b93c8c07606a5e662e34aa4647088a19a46a13456.json: identity does not match content: id: got sha256:363552ff89fc60ae03d7bc1b93c8c07606a5e662e34aa4647088a19a46a13456, content hashes to sha256:d02839b48d765b72045eba118a2de52451dad24420140d02d0b72e22c1e54c4a window err <nil>
=== RUN   TestC959DuplicateHeaderMatchesFullDecoder/created_at/after/oldest
    o032_c959_probe_test.go:117: full err b8b57466f51c31260eaccf170886d5558b637454b065f26adef2e08882448f65.json: identity does not match content: id: got sha256:b8b57466f51c31260eaccf170886d5558b637454b065f26adef2e08882448f65, content hashes to sha256:61733a57b25cbf1fc4805a17eecac74d1f80bb7bbe42adf79a1d0eddffaa168c window err <nil>
=== RUN   TestC959DuplicateHeaderMatchesFullDecoder/created_at/after/boundary
    o032_c959_probe_test.go:117: full err b8e6924e95f30f9423caf2711266b2306bf6104906a2f9e2221dc0cbe7022432.json: identity does not match content: id: got sha256:b8e6924e95f30f9423caf2711266b2306bf6104906a2f9e2221dc0cbe7022432, content hashes to sha256:1d498ebb9a347ed015d2479657619b75fb0b4c1660fb70b78d9040016165e70f window err <nil>
=== RUN   TestC959DuplicateHeaderMatchesFullDecoder/origin/before/newest
=== RUN   TestC959DuplicateHeaderMatchesFullDecoder/origin/before/oldest
    o032_c959_probe_test.go:117: full err b8b57466f51c31260eaccf170886d5558b637454b065f26adef2e08882448f65.json: invalid origin or retention: retention: manual snapshots must be prunable, got "permanent" window err <nil>
=== RUN   TestC959DuplicateHeaderMatchesFullDecoder/origin/after/boundary
--- FAIL: TestC959DuplicateHeaderMatchesFullDecoder (0.95s)
    o032_c959_probe_test.go:221: damage outside the window failed the load: malformed record: b8b57466f51c31260eaccf170886d5558b637454b065f26adef2e08882448f65.json: is not readable JSON: unexpected end of JSON input
=== RUN   TestC959DamageOutsideWindow/truncated/inside=true
=== RUN   TestC959DamageOutsideWindow/bad-results-value/inside=false
    o032_c959_probe_test.go:221: damage outside the window failed the load: malformed record: b8b57466f51c31260eaccf170886d5558b637454b065f26adef2e08882448f65.json: is not readable JSON: invalid character 'g' after top-level value
=== RUN   TestC959DamageOutsideWindow/trailing-garbage/inside=true
--- FAIL: TestC959DamageOutsideWindow (0.60s)
```

### c959-bench.log

```text
BenchmarkHealthHistoryLoadDashboard/single/n=10-16         	     997	   1208591 ns/op	      5266 bytes/snapshot	     52655 fixture-bytes	  676453 B/op	    4425 allocs/op
BenchmarkHealthHistoryLoadDashboard/single/n=100-16        	     265	   4533668 ns/op	      5257 bytes/snapshot	    525682 fixture-bytes	 1633424 B/op	    8379 allocs/op
BenchmarkHealthHistoryLoadDashboard/single/n=1000-16       	      32	  33966327 ns/op	      5256 bytes/snapshot	   5256194 fixture-bytes	 8085665 B/op	   29111 allocs/op
BenchmarkHealthHistoryLoadDashboard/multi/n=10-16          	     225	   5321638 ns/op	     27929 bytes/snapshot	    279287 fixture-bytes	 3519240 B/op	   22023 allocs/op
BenchmarkHealthHistoryLoadDashboard/multi/n=100-16         	      62	  19707814 ns/op	     27893 bytes/snapshot	   2789275 fixture-bytes	 7935377 B/op	   34357 allocs/op
BenchmarkHealthHistoryLoadDashboard/multi/n=1000-16        	       8	 129405341 ns/op	     27890 bytes/snapshot	  27890063 fixture-bytes	35387060 B/op	   55111 allocs/op
BenchmarkHealthHistoryLoadDashboard/heavy/n=10-16          	      60	  20017627 ns/op	    139429 bytes/snapshot	   1394287 fixture-bytes	14826095 B/op	   43265 allocs/op
BenchmarkHealthHistoryLoadDashboard/heavy/n=100-16         	      15	  76664466 ns/op	    139393 bytes/snapshot	  13939275 fixture-bytes	34751704 B/op	   66150 allocs/op
BenchmarkHealthHistoryLoadDashboard/heavy/n=1000-16        	       2	 561573816 ns/op	    139390 bytes/snapshot	 139390063 fixture-bytes	168018836 B/op	   86900 allocs/op
ok  	github.com/opencode/savepoint/internal/codehealth	20.317s
```

### health check

```text
Objective: O-032
Snapshot: sha256:0aa2073976fe3e6ef67bfbb7cc3d2288f18067b56e461daaaa78f50bb420397d (created)
Code Health does not block clearance.
- complexity via lizard-csv [optional]: reported only; no finding: no blocking finding.
- coverage via go-cover-profile [optional]: reported only; no finding: no blocking finding.
- dependency_vulnerabilities via osv-scanner-json [optional]: reported only; unhealthy measurement: 2 of unknown severity need review.
- duplication via jscpd-json [optional]: reported only; no finding: no blocking finding.
- tests via go-test-json [optional]: reported only; no finding: no blocking finding.
```
