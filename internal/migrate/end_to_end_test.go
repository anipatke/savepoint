package migrate_test

// end_to_end_test.go is E51 T009's integration proof. Every other test in this
// package proves one piece; these run both frozen E41 fixtures through the
// real command path into temporary directories and assert the whole-project
// properties that only show up once everything runs together: a clean V2
// load, a fully resolvable reference graph, no invented clearance, an
// accountable destination for every source file, archives that match the
// fixtures' independently recorded hashes, byte-stable converted output, a
// no-op second run, and a write-free preview.
//
// It is an external test package (migrate_test) on purpose: it may only use
// what a real caller can reach, so a passing run is evidence about the
// exported command path rather than about package internals.
//
// TEST-04: every fixture is copied into t.TempDir() first, no test writes to
// internal/data/testdata/migration/, and no test touches a live project.
// TestEndToEnd_frozenFixturesAreNeverMutated asserts that mechanically.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/opencode/savepoint/cmd"
	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/doctor"
	"github.com/opencode/savepoint/internal/migrate"
	"gopkg.in/yaml.v3"
)

// e2eFixtures are the two frozen V1 projects E41 produced for exactly this
// purpose.
var e2eFixtures = []string{"v1-basic", "v1-history"}

const (
	fixtureRoot = "../data/testdata/migration"

	// The injected clock and operation ID. Both are fixed so converted
	// content — which embeds a generated-at date in Issue history and in the
	// v1-to-v2.yml manifest — renders byte-identically on every run, which is
	// what makes the golden comparison below meaningful.
	e2eOperationID = "op-e2e-fixed"
)

var e2eClock = time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC)

// e2ePreparedFixtures holds one converted result per frozen fixture. Unix
// read-only assertions share these protected results; Windows consumers use
// private copies because read-only directory modes do not prevent new files.
// Command, preview, recovery, no-op, and other mutating scenarios keep fresh
// source project copies on every platform.
var e2ePreparedFixtures = map[string]string{}

func TestMain(m *testing.M) {
	preparedRoot, err := os.MkdirTemp("", "savepoint-migration-e2e-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "create shared migration fixture root: %v\n", err)
		os.Exit(1)
	}
	cleanup := func() error {
		if err := makeTreeWritable(preparedRoot); err != nil {
			return fmt.Errorf("make shared migration fixtures removable: %w", err)
		}
		if err := os.RemoveAll(preparedRoot); err != nil {
			return fmt.Errorf("remove shared migration fixtures: %w", err)
		}
		return nil
	}

	for _, fixture := range e2eFixtures {
		root := filepath.Join(preparedRoot, fixture)
		if err := copyFixtureInto(root, fixture); err != nil {
			if cleanupErr := cleanup(); cleanupErr != nil {
				fmt.Fprintf(os.Stderr, "shared fixture cleanup: %v\n", cleanupErr)
			}
			fmt.Fprintf(os.Stderr, "prepare shared migration fixture %s: %v\n", fixture, err)
			os.Exit(1)
		}
		code, out, err := runMigrateCommand(root, "--apply")
		if err != nil || code != 0 || !strings.Contains(out, "migration complete") {
			if cleanupErr := cleanup(); cleanupErr != nil {
				fmt.Fprintf(os.Stderr, "shared fixture cleanup: %v\n", cleanupErr)
			}
			fmt.Fprintf(os.Stderr, "migrate shared fixture %s: code = %d, err = %v\n%s", fixture, code, err, out)
			os.Exit(1)
		}
		if err := makeTreeReadOnly(root); err != nil {
			if cleanupErr := cleanup(); cleanupErr != nil {
				fmt.Fprintf(os.Stderr, "shared fixture cleanup: %v\n", cleanupErr)
			}
			fmt.Fprintf(os.Stderr, "protect shared migration fixture %s: %v\n", fixture, err)
			os.Exit(1)
		}
		e2ePreparedFixtures[fixture] = root
	}

	code := m.Run()
	if err := cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "shared fixture cleanup: %v\n", err)
		if code == 0 {
			code = 1
		}
	}
	os.Exit(code)
}

// --- the real command path ----------------------------------------------

// runMigrate invokes migrate exactly as `savepoint migrate ...` does: the
// real argument parser in cmd, the real dispatch, and the real command body
// in migrate.RunCommand — with only the clock and the operation-ID source
// injected, which is the same seam Plan already requires for determinism.
//
// The one thing above this line is main.go's field-for-field mapping of
// cmd.MigrateOptions onto migrate.CommandOptions, which main_test.go covers
// by running the compiled binary.
func runMigrate(t *testing.T, args ...string) (int, string, error) {
	t.Helper()
	return runMigrateCommand(args...)
}

func runMigrateCommand(args ...string) (int, string, error) {
	var out strings.Builder
	code, err := cmd.RunMigrate(context.Background(), args, &out, func(_ context.Context, opts cmd.MigrateOptions) (int, error) {
		return migrate.RunCommand(migrate.CommandOptions{
			Dir:            opts.Dir,
			Write:          opts.WillWrite(),
			DecisionsFile:  opts.DecisionsFile,
			Recover:        opts.Recover,
			Stdout:         &out,
			Now:            func() time.Time { return e2eClock },
			NewOperationID: func() string { return e2eOperationID },
		})
	})
	return code, out.String(), err
}

// migrateFixture copies fixture into a temp directory and applies the
// migration through the command path, returning the migrated project root.
func migrateFixture(t *testing.T, fixture string) string {
	t.Helper()
	root := copyFixture(t, fixture)
	code, out, err := runMigrate(t, root, "--apply")
	if err != nil || code != 0 {
		t.Fatalf("migrate --apply over %s: code = %d, err = %v\n%s", fixture, code, err, out)
	}
	if !strings.Contains(out, "migration complete") {
		t.Fatalf("migrate --apply over %s did not report completion:\n%s", fixture, out)
	}
	return root
}

func readOnlyMigratedFixture(t *testing.T, fixture string) string {
	t.Helper()
	root, ok := e2ePreparedFixtures[fixture]
	if !ok {
		t.Fatalf("no shared migrated result for fixture %s", fixture)
	}
	if runtime.GOOS != "windows" {
		return root
	}
	copyRoot := filepath.Join(t.TempDir(), fixture)
	if err := copyTree(root, copyRoot); err != nil {
		t.Fatalf("copy prepared migrated fixture %s: %v", fixture, err)
	}
	return copyRoot
}

func makeTreeReadOnly(root string) error {
	var paths []string
	if err := filepath.WalkDir(root, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		paths = append(paths, path)
		return nil
	}); err != nil {
		return err
	}
	for i := len(paths) - 1; i >= 0; i-- {
		info, err := os.Stat(paths[i])
		if err != nil {
			return err
		}
		if err := os.Chmod(paths[i], info.Mode().Perm()&^0222); err != nil {
			return fmt.Errorf("chmod %s read-only: %w", paths[i], err)
		}
	}
	return nil
}

func makeTreeWritable(root string) error {
	var paths []string
	if err := filepath.WalkDir(root, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		paths = append(paths, path)
		return nil
	}); err != nil {
		return err
	}
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		mode := info.Mode().Perm() | 0200
		if info.IsDir() {
			mode |= 0300
		}
		if err := os.Chmod(path, mode); err != nil {
			return fmt.Errorf("chmod %s writable: %w", path, err)
		}
	}
	return nil
}

// --- AC1: temp directories, real command path, untouched fixtures --------

func TestEndToEnd_bothFixturesMigrateThroughTheCommandPath(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := migrateFixture(t, fixture)

			version, err := data.ReadSchemaVersion(filepath.Join(root, ".savepoint", "config.yml"))
			if err != nil {
				t.Fatalf("ReadSchemaVersion() error = %v", err)
			}
			if version != data.SchemaVersionV2 {
				t.Fatalf("schema_version = %v, want 2 after a successful apply", version)
			}
		})
	}
}

// TestEndToEnd_frozenFixturesAreNeverMutated is the mechanical half of
// TEST-04: it snapshots both frozen fixture directories, runs preview and
// apply on fresh temporary copies, and asserts the frozen fixtures' bytes and
// modification times are untouched afterwards.
func TestEndToEnd_frozenFixturesAreNeverMutated(t *testing.T) {
	t.Parallel()
	before := map[string]map[string]fileFacts{}
	for _, fixture := range e2eFixtures {
		before[fixture] = snapshot(t, filepath.Join(fixtureRoot, fixture))
	}

	for _, fixture := range e2eFixtures {
		root := copyFixture(t, fixture)
		if _, _, err := runMigrate(t, root); err != nil {
			t.Fatalf("preview over %s: %v", fixture, err)
		}
		migrateFixture(t, fixture)
	}

	for _, fixture := range e2eFixtures {
		assertUnchanged(t, before[fixture], snapshot(t, filepath.Join(fixtureRoot, fixture)),
			"frozen fixture "+fixture)
	}
}

// --- AC2: the migrated project loads clean ------------------------------

func TestEndToEnd_migratedProjectLoadsCleanThroughLoadV2Index(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := readOnlyMigratedFixture(t, fixture)

			index, err := data.LoadV2Index(filepath.Join(root, ".savepoint"))
			if err != nil {
				t.Fatalf("LoadV2Index() after migrating %s reported a diagnostic: %v", fixture, err)
			}
			if len(index.Objectives) == 0 {
				t.Errorf("migrated %s has no Objectives; the fixture's active epic should have converted", fixture)
			}
			if len(index.Tasks) == 0 {
				t.Errorf("migrated %s has no Tasks; the fixture's active work should have converted", fixture)
			}
		})
	}
}

func TestEndToEnd_sharedMigratedResultsAreReadOnly(t *testing.T) {
	if runtime.GOOS == "windows" {
		for _, fixture := range e2eFixtures {
			t.Run(fixture, func(t *testing.T) {
				shared := e2ePreparedFixtures[fixture]
				consumer := readOnlyMigratedFixture(t, fixture)
				probe := filepath.Join(consumer, ".t014-consumer-probe")
				if err := os.WriteFile(probe, []byte("private copy"), 0644); err != nil {
					t.Fatalf("consumer copy is not independently writable: %v", err)
				}
				if err := os.Remove(probe); err != nil {
					t.Fatalf("remove consumer probe: %v", err)
				}
				if _, err := os.Stat(filepath.Join(shared, ".t014-consumer-probe")); !os.IsNotExist(err) {
					t.Fatalf("consumer changed shared result: stat error = %v", err)
				}
			})
		}
		return
	}

	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := e2ePreparedFixtures[fixture]
			info, err := os.Stat(root)
			if err != nil {
				t.Fatalf("stat shared %s root: %v", fixture, err)
			}
			if info.Mode().Perm()&0222 != 0 {
				t.Fatalf("shared %s root is writable: mode = %v", fixture, info.Mode().Perm())
			}
			if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				info, err := entry.Info()
				if err != nil {
					return err
				}
				if info.Mode().Perm()&0222 != 0 {
					return fmt.Errorf("%s remains writable with mode %v", path, info.Mode().Perm())
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			probe := filepath.Join(root, ".t014-consumer-probe")
			if err := os.WriteFile(probe, []byte("must not persist"), 0644); err == nil {
				// A privileged test process can bypass Unix mode bits; remove its
				// probe and rely on the mode checks above in that environment.
				if err := os.Remove(probe); err != nil {
					t.Fatalf("remove privileged write probe: %v", err)
				}
			} else if !os.IsPermission(err) {
				t.Fatalf("probe shared result write: %v", err)
			}
			config := filepath.Join(root, ".savepoint", "config.yml")
			if f, err := os.OpenFile(config, os.O_WRONLY, 0); err == nil {
				_ = f.Close()
			} else if !os.IsPermission(err) {
				t.Fatalf("open shared config for writing: %v", err)
			}
		})
	}
}

func TestEndToEnd_releaseRecordsMappingsAndCutoverGateAgree(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := copyFixture(t, fixture)
			plan := planFor(t, root)
			preview := migrate.FormatPreview(plan)
			releaseTargets := make([]migrate.PlannedTarget, 0)
			sourceHashes := make(map[string]string)
			for _, source := range plan.Sources {
				sourceHashes[source.Path] = source.SHA256
			}
			for _, target := range plan.Targets {
				if target.Kind != migrate.TargetRelease {
					continue
				}
				releaseTargets = append(releaseTargets, target)
				if !strings.Contains(preview, "- release "+target.GlobalID) {
					t.Errorf("preview omitted first-class Release %s:\n%s", target.GlobalID, preview)
				}
			}
			if len(releaseTargets) == 0 {
				t.Fatal("migration plan produced no first-class Release targets")
			}

			if _, out, err := runMigrate(t, root, "--apply"); err != nil {
				t.Fatalf("migrate --apply: %v\n%s", err, out)
			}
			manifest := readWrittenManifest(t, root)
			index := mustIndex(t, root)
			if len(index.Releases) != len(releaseTargets) {
				t.Fatalf("loaded Releases = %d, want %d planned Releases", len(index.Releases), len(releaseTargets))
			}

			for _, target := range releaseTargets {
				var identity *migrate.ManifestIdentity
				for i := range manifest.Identities {
					candidate := &manifest.Identities[i]
					if candidate.Kind == string(migrate.TargetRelease) && candidate.Path == target.Legacy.Path {
						identity = candidate
						break
					}
				}
				if identity == nil || identity.GlobalID != target.GlobalID || identity.TargetPath != target.InstallPath() {
					t.Errorf("Release identity for %s = %+v, want source-qualified live mapping", target.Legacy.Path, identity)
				}

				var archive *migrate.ManifestArchive
				for i := range manifest.Archives {
					candidate := &manifest.Archives[i]
					if candidate.SourcePath == target.Legacy.Path {
						archive = candidate
						break
					}
				}
				if archive == nil {
					t.Errorf("Release source %s has no accountable archive mapping", target.Legacy.Path)
					continue
				}
				if archive.SHA256 != sourceHashes[target.Legacy.Path] {
					t.Errorf("archive hash for %s = %q, want inventory hash %q", target.Legacy.Path, archive.SHA256, sourceHashes[target.Legacy.Path])
				}
				mustExistAt(t, filepath.Join(root, ".savepoint", filepath.FromSlash(target.InstallPath())))
				mustExistAt(t, filepath.Join(root, filepath.FromSlash(archive.ArchivePath)))
			}

			cutover := data.ResolveReleaseCutover(index)
			for releaseID := range index.Releases {
				canonical := data.ResolveReleaseCompletion(index, releaseID)
				if canonical.Allowed {
					continue
				}
				found := false
				for _, blocker := range cutover.Blockers {
					if blocker.ReleaseID == releaseID {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("project cutover gate omitted canonical blocker for Release %s: %+v", releaseID, canonical)
				}
			}

			problems := doctor.RunV2Checks(filepath.Join(root, ".savepoint")).Releases
			for _, blocker := range cutover.Blockers {
				matched := false
				for _, problem := range problems {
					if strings.Contains(problem.Message, "release "+blocker.ReleaseID) && strings.Contains(problem.Message, string(blocker.Gate.Kind)) {
						matched = true
						break
					}
				}
				// Doctor uses stable diagnostic names rather than the raw gate
				// kind; the Release identity must still be present in every
				// canonical refusal.
				if !matched {
					for _, problem := range problems {
						if strings.Contains(problem.Message, "release "+blocker.ReleaseID) {
							matched = true
							break
						}
					}
				}
				if !matched {
					t.Errorf("doctor omitted cutover blocker for Release %s: %+v; problems = %+v", blocker.ReleaseID, blocker.Gate, problems)
				}
			}

			if fixture == "v1-history" {
				historical := index.Releases["R-001"]
				if historical == nil || historical.LegacyCompletion == nil {
					t.Fatalf("R-001 historical Release = %+v, want typed legacy completion", historical)
				}
				if historicalDecision := data.ResolveReleaseCompletion(index, "R-001"); !historicalDecision.AllowedByLegacyCompletion {
					t.Fatalf("R-001 historical decision = %+v, want allowed by archived history", historicalDecision)
				}
				if cutover.Allowed {
					t.Fatalf("v1-history cutover = %+v, want blocked by its active Release", cutover)
				}
			}
		})
	}
}

func TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability(t *testing.T) {
	t.Parallel()
	root := copyRepositoryWorkingTree(t)
	beforePreview := snapshot(t, root)
	code, preview, err := runMigrate(t, root)
	if err != nil || code != 0 {
		t.Fatalf("repository-copy preview: code = %d, err = %v\n%s", code, err, preview)
	}
	assertUnchanged(t, beforePreview, snapshot(t, root), "repository-copy preview")

	plan := planFor(t, root)
	if !plan.Appliable {
		t.Fatalf("repository-copy plan is not appliable: ambiguities = %+v", plan.Ambiguities)
	}
	var plannedReleases []migrate.PlannedTarget
	for _, target := range plan.Targets {
		if target.Kind == migrate.TargetRelease {
			plannedReleases = append(plannedReleases, target)
		}
	}
	if len(plannedReleases) == 0 {
		t.Fatal("repository-copy plan produced no first-class Release")
	}

	code, output, err := runMigrate(t, root, "--apply")
	if err != nil || code != 0 {
		t.Fatalf("repository-copy apply: code = %d, err = %v\n%s", code, err, output)
	}
	index := mustIndex(t, root)
	if len(index.Releases) != len(plannedReleases) {
		t.Fatalf("repository-copy Releases = %d, want %d planned Releases", len(index.Releases), len(plannedReleases))
	}
	manifest := readWrittenManifest(t, root)
	for _, target := range plannedReleases {
		foundLive := false
		for _, identity := range manifest.Identities {
			if identity.Kind == string(migrate.TargetRelease) && identity.Path == target.Legacy.Path && identity.GlobalID == target.GlobalID {
				foundLive = true
				mustExistAt(t, filepath.Join(root, ".savepoint", filepath.FromSlash(identity.TargetPath)))
				break
			}
		}
		if !foundLive {
			t.Errorf("repository-copy Release %s has no manifest/live identity for %s", target.GlobalID, target.Legacy.Path)
		}

		foundArchive := false
		for _, archive := range manifest.Archives {
			if archive.SourcePath == target.Legacy.Path {
				foundArchive = true
				mustExistAt(t, filepath.Join(root, filepath.FromSlash(archive.ArchivePath)))
				break
			}
		}
		if !foundArchive {
			t.Errorf("repository-copy Release %s has no archived source mapping for %s", target.GlobalID, target.Legacy.Path)
		}
	}

	// A second unchanged migration is the real command's no-op path: it must
	// not rewrite the index, router selection, evidence, archive, or mtimes.
	beforeSecond := snapshot(t, root)
	idsBefore := recordIDs(t, root)
	manifestPath := filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml")
	manifestBeforeBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read repository-copy migration manifest: %v", err)
	}
	code, output, err = runMigrate(t, root, "--apply")
	if err != nil || code != 0 {
		t.Fatalf("repository-copy second apply: code = %d, err = %v\n%s", code, err, output)
	}
	assertUnchanged(t, beforeSecond, snapshot(t, root), "repository-copy second apply")
	if got := recordIDs(t, root); !equalStrings(got, idsBefore) {
		t.Fatalf("repository-copy IDs changed on second apply: before %v, after %v", idsBefore, got)
	}
	manifestAfterBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read repository-copy migration manifest after second apply: %v", err)
	}
	if string(manifestAfterBytes) != string(manifestBeforeBytes) {
		t.Fatal("repository-copy migration manifest changed on second apply")
	}
}

// --- AC3: every reference resolves ---------------------------------------

func TestEndToEnd_everyReferenceResolves(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			index := mustIndex(t, readOnlyMigratedFixture(t, fixture))

			for id, task := range index.Tasks {
				if _, ok := index.Objectives[task.Objective]; !ok {
					t.Errorf("task %s names objective %s, which does not exist", id, task.Objective)
				}
				for _, dep := range task.DependsOn {
					if _, ok := index.Tasks[dep.Task]; !ok {
						t.Errorf("task %s depends_on %s, which does not exist", id, dep.Task)
					}
				}
			}

			for id, objective := range index.Objectives {
				for _, dep := range objective.DependsOn {
					if _, ok := index.Objectives[dep]; !ok {
						t.Errorf("objective %s depends_on %s, which does not exist", id, dep)
					}
				}
			}

			for id, issue := range index.Issues {
				for _, taskID := range issue.Tasks {
					if _, ok := index.Tasks[taskID]; !ok {
						t.Errorf("issue %s names task %s, which does not exist", id, taskID)
					}
				}
				for _, checkID := range issue.Checks {
					if _, ok := index.Checks[checkID]; !ok {
						t.Errorf("issue %s names check %s, which does not exist", id, checkID)
					}
				}
				if issue.DuplicateOf != "" {
					if _, ok := index.Issues[issue.DuplicateOf]; !ok {
						t.Errorf("issue %s duplicates %s, which does not exist", id, issue.DuplicateOf)
					}
				}
			}
		})
	}
}

// --- AC4: nothing carries clearance the source did not record ------------

func TestEndToEnd_noConvertedRecordCarriesClearance(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			index := mustIndex(t, readOnlyMigratedFixture(t, fixture))

			// Neither fixture records a single V1 Check, so a migrated
			// project that has one invented it.
			if len(index.Checks) != 0 {
				t.Errorf("migrated %s carries %d Check records; the V1 sources recorded none", fixture, len(index.Checks))
			}

			for id, task := range index.Tasks {
				assertNoEvidence(t, "task "+id, task.Evidence)
				if got := data.ResolveClearance(index, id).State; got != data.ClearanceMissing {
					t.Errorf("task %s clearance = %q, want %q", id, got, data.ClearanceMissing)
				}
			}
			for id, objective := range index.Objectives {
				assertNoEvidence(t, "objective "+id, objective.Evidence)
				if got := data.ResolveClearance(index, id).State; got != data.ClearanceMissing {
					t.Errorf("objective %s clearance = %q, want %q", id, got, data.ClearanceMissing)
				}
			}
		})
	}
}

// TestEndToEnd_convertedActiveTaskIsNotCompletable is the gate-level half of
// AC4: a converted active Task is refused completion under E43's gates
// because migration recorded no clearance for it.
//
// ResolveTaskCompletion only reaches its clearance branch for a Task at
// status in_progress and stage audit, and migration deliberately never puts
// a Task there. The test therefore advances the in-memory record — a probe on
// the loaded struct, never a write and never a migration output — so the
// assertion lands on the clearance decision itself rather than stopping at
// the earlier lifecycle guard.
func TestEndToEnd_convertedActiveTaskIsNotCompletable(t *testing.T) {
	t.Parallel()
	index := mustIndex(t, readOnlyMigratedFixture(t, "v1-basic"))

	taskID := oneActiveTaskID(t, index)
	task := index.Tasks[taskID]
	task.Status = data.ColumnInProgress
	task.Stage = data.StageAudit

	decision := data.ResolveTaskCompletion(index, taskID)
	if decision.Allowed {
		t.Fatalf("converted task %s is completable; migration recorded no clearance for it", taskID)
	}
	if !hasBlocker(decision, data.GateBlockClearanceMissing) {
		t.Fatalf("task %s completion blockers = %+v, want a %s blocker", taskID, decision.Blockers, data.GateBlockClearanceMissing)
	}
}

// --- AC5: every fixture file has an accountable destination --------------

// destination is the fate migration assigned one V1 source file. Every file a
// fixture manifest names must land in exactly one of these; "none" is the
// failure this test exists to report by name.
type destination string

const (
	destConverted destination = "converted to a V2 record"
	destRelocated destination = "relocated or rewritten as a project document"
	destArchived  destination = "archived byte-for-byte under .savepoint/archive/v1/"
	destPreserved destination = "preserved in place, byte-identical"
	destActivated destination = "preserved in place with schema_version activated"
	destNone      destination = ""
)

const configRelPath = ".savepoint/config.yml"

func TestEndToEnd_everyFixtureFileHasAnAccountableDestination(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := readOnlyMigratedFixture(t, fixture)

			// Planning reads the frozen V1 input without writing it. The shared
			// converted result was produced from the same fixed clock and ID.
			plan := planFor(t, filepath.Join(fixtureRoot, fixture, "project"))

			var unaccounted []string
			for _, file := range loadFixtureFiles(t, fixture) {
				if fate := fateOf(t, plan, root, fixture, file.Path); fate == destNone {
					unaccounted = append(unaccounted, file.Path)
				}
			}
			if len(unaccounted) > 0 {
				sort.Strings(unaccounted)
				t.Fatalf("%d file(s) in %s/manifest.yml have no accountable destination:\n  %s",
					len(unaccounted), fixture, strings.Join(unaccounted, "\n  "))
			}
		})
	}
}

// fateOf resolves one source path's destination against the plan and the
// migrated tree.
func fateOf(t *testing.T, plan *migrate.ConversionPlan, root, fixture, path string) destination {
	t.Helper()

	for _, target := range plan.Targets {
		if target.Legacy.Path == path {
			mustExistAt(t, filepath.Join(root, ".savepoint", filepath.FromSlash(target.InstallPath())))
			return destConverted
		}
	}
	for _, doc := range plan.Documents {
		if doc.SourcePath == path {
			mustExistAt(t, filepath.Join(root, ".savepoint", filepath.FromSlash(doc.TargetPath)))
			return destRelocated
		}
	}
	for _, archive := range plan.Archives {
		if archive.SourcePath == path {
			mustExistAt(t, filepath.Join(root, filepath.FromSlash(archive.ArchivePath)))
			return destArchived
		}
	}

	live := filepath.Join(root, filepath.FromSlash(path))
	if _, err := os.Stat(live); err != nil {
		return destNone
	}
	if path == configRelPath {
		version, err := data.ReadSchemaVersion(live)
		if err != nil || version != data.SchemaVersionV2 {
			return destNone
		}
		return destActivated
	}
	if hashFile(t, live) != hashFile(t, filepath.Join(fixtureRoot, fixture, "project", filepath.FromSlash(path))) {
		return destNone
	}
	return destPreserved
}

// --- AC6: archives hash equal to the fixture manifest --------------------

func TestEndToEnd_archivedFilesMatchTheFixtureManifestHashes(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := readOnlyMigratedFixture(t, fixture)
			manifest := readWrittenManifest(t, root)

			// The manifest is read as the expectation. Nothing here
			// recomputes a fixture hash and compares it to itself.
			expected := map[string]string{}
			for _, file := range loadFixtureFiles(t, fixture) {
				expected[file.Path] = file.SHA256
			}

			checked := 0
			for _, archive := range manifest.Archives {
				want, recorded := expected[archive.SourcePath]
				if !recorded {
					continue // not a manifest-listed source (e.g. a test-added file)
				}
				got := hashFile(t, filepath.Join(root, filepath.FromSlash(archive.ArchivePath)))
				if got != want {
					t.Errorf("archive %s sha256 = %s, want %s recorded in %s/manifest.yml for %s",
						archive.ArchivePath, got, want, fixture, archive.SourcePath)
				}
				checked++
			}
			if checked == 0 {
				t.Fatalf("no archived file in %s was covered by the fixture manifest", fixture)
			}
		})
	}
}

// --- AC7: the recurring T001-shared stays two distinct records -----------

const sharedTaskID = "E01-example/T001-shared"

// TestEndToEnd_recurringSharedTaskIsRecordedPerRelease covers the half of AC7
// the frozen fixture can prove directly. v1-history declares
// E01-example/T001-shared in both release v1 and release v1.1. v1's copy is
// done under a done epic, so per the migration contract's "completed work is archived, not
// converted" rule it receives no V2 identity — the two recurrences therefore
// appear as one converted identity and one archive entry, each recording the
// release it came from, and never collapse onto a single record.
func TestEndToEnd_recurringSharedTaskIsRecordedPerRelease(t *testing.T) {
	t.Parallel()
	root := readOnlyMigratedFixture(t, "v1-history")
	manifest := readWrittenManifest(t, root)

	releases := map[string]string{}
	for _, identity := range manifest.Identities {
		if identity.OriginalID == sharedTaskID {
			releases[identity.Release] = "converted as " + identity.GlobalID
		}
	}
	for _, archive := range manifest.Archives {
		if archive.OriginalID == sharedTaskID {
			releases[archive.Release] = "archived at " + archive.ArchivePath
		}
	}

	for _, release := range []string{"v1", "v1.1"} {
		if _, ok := releases[release]; !ok {
			t.Errorf("v1-to-v2.yml records no destination for %s from release %s; recorded: %+v",
				sharedTaskID, release, releases)
		}
	}
	if len(releases) != 2 {
		t.Fatalf("the two recurrences of %s resolved to %d recorded destination(s), want 2: %+v",
			sharedTaskID, len(releases), releases)
	}
}

// TestEndToEnd_recurringSharedTaskAllocatesDistinctGlobalIDs covers the other
// half of AC7: when both recurrences are active, the same short ID across two
// releases allocates two distinct global Task IDs rather than colliding. The
// frozen v1-history fixture cannot show this directly — its v1 copy is
// settled history — so the test reopens v1's epic and task on the temporary
// copy, which is a statement about identity allocation, not about the
// fixture.
func TestEndToEnd_recurringSharedTaskAllocatesDistinctGlobalIDs(t *testing.T) {
	t.Parallel()
	root := copyFixture(t, "v1-history")
	reopen(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "E01-Detail.md"),
		"status: done", "status: in_progress")
	reopen(t, filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T001-shared.md"),
		"status: done", "status: planned")

	if _, out, err := runMigrate(t, root, "--apply"); err != nil {
		t.Fatalf("migrate --apply: %v\n%s", err, out)
	}

	manifest := readWrittenManifest(t, root)
	byRelease := map[string]string{}
	for _, identity := range manifest.Identities {
		if identity.OriginalID != sharedTaskID || identity.Kind != "task" {
			continue
		}
		if existing, ok := byRelease[identity.Release]; ok {
			t.Fatalf("release %s allocated %s twice: %s and %s", identity.Release, sharedTaskID, existing, identity.GlobalID)
		}
		byRelease[identity.Release] = identity.GlobalID
	}

	if len(byRelease) != 2 {
		t.Fatalf("%s converted to %d Tasks, want one per release: %+v", sharedTaskID, len(byRelease), byRelease)
	}
	if byRelease["v1"] == byRelease["v1.1"] {
		t.Fatalf("both releases' %s allocated the same global ID %s", sharedTaskID, byRelease["v1"])
	}

	index := mustIndex(t, root)
	for release, id := range byRelease {
		if _, ok := index.Tasks[id]; !ok {
			t.Errorf("release %s allocated task %s, which is not in the loaded V2 index", release, id)
		}
	}
}

// --- AC8: legacy facts stay resolvable through v1-to-v2.yml --------------

func TestEndToEnd_legacyFactsAreResolvableThroughTheManifest(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := readOnlyMigratedFixture(t, fixture)
			manifest := readWrittenManifest(t, root)

			archived := map[string]migrate.ManifestArchive{}
			for _, archive := range manifest.Archives {
				archived[archive.SourcePath] = archive
				mustExistAt(t, filepath.Join(root, filepath.FromSlash(archive.ArchivePath)))
			}
			converted := map[string]bool{}
			for _, identity := range manifest.Identities {
				converted[identity.Path] = true
			}

			// The general rule behind AC8: a V1 source that received no V2
			// record and is no longer at its original path can only stay
			// resolvable through v1-to-v2.yml. Nothing may simply disappear.
			for _, file := range loadFixtureFiles(t, fixture) {
				if converted[file.Path] {
					continue
				}
				if _, stillLive := os.Stat(filepath.Join(root, filepath.FromSlash(file.Path))); stillLive == nil {
					continue // preserved in place; resolvable by inspection
				}
				if _, ok := archived[file.Path]; !ok {
					t.Errorf("%s has no V2 record, is gone from its V1 path, and v1-to-v2.yml resolves it to nothing", file.Path)
				}
			}

			// The specific legacy facts the migration contract names: completed Tasks, and
			// verified or waived findings. Each is settled history that gets no
			// V2 identity, so each must resolve to an archive entry.
			for _, file := range loadFixtureFiles(t, fixture) {
				settled := (file.Role == "task" && file.RawStatus == "done") ||
					(file.Role == "finding" && (file.RawStatus == "verified" || file.RawStatus == "waived"))
				if !settled {
					continue
				}
				if converted[file.Path] {
					t.Errorf("%s is settled V1 history (%s %s) but was converted to a V2 record", file.Path, file.Role, file.RawStatus)
					continue
				}
				if _, ok := archived[file.Path]; !ok {
					t.Errorf("%s is settled V1 history (%s %s) but v1-to-v2.yml resolves it to nothing", file.Path, file.Role, file.RawStatus)
				}
			}
		})
	}
}

// TestEndToEnd_legacyPrerequisiteIsTyped covers the specific legacy fact the migration contract
// calls out: an active Task that depended on completed, archived work keeps
// that dependency as a typed legacy prerequisite rather than a fabricated
// depends_on or a fabricated Check.
func TestEndToEnd_legacyPrerequisiteIsTyped(t *testing.T) {
	t.Parallel()
	root := readOnlyMigratedFixture(t, "v1-basic")
	manifest := readWrittenManifest(t, root)

	if len(manifest.LegacyPrerequisites) == 0 {
		t.Fatal("v1-basic's T002-follow-up depended on the completed T001-original, " +
			"but v1-to-v2.yml records no legacy prerequisite")
	}

	index := mustIndex(t, root)
	for _, prereq := range manifest.LegacyPrerequisites {
		if _, ok := index.Tasks[prereq.Task]; !ok {
			t.Errorf("legacy prerequisite names task %s, which is not in the V2 index", prereq.Task)
		}
		mustExistAt(t, filepath.Join(root, filepath.FromSlash(prereq.ArchivePath)))
		if strings.TrimSpace(prereq.Evidence) == "" {
			t.Errorf("legacy prerequisite for %s records no evidence", prereq.Task)
		}
	}

	// The archived prerequisite must not have leaked into the V2 dependency
	// graph, which is the whole reason it is typed separately.
	for id, task := range index.Tasks {
		for _, dep := range task.DependsOn {
			if _, ok := index.Tasks[dep.Task]; !ok {
				t.Errorf("task %s depends_on %s, which is archived rather than converted", id, dep.Task)
			}
		}
	}
}

// --- AC9: golden converted output ---------------------------------------

// regenerateGoldenEnv deliberately regenerates the golden files.
//
// A golden mismatch is a report that converted output changed. It is never
// fixed by regenerating in place, because that would erase the very evidence
// the golden exists to produce. The procedure is:
//
//  1. Read the diff the failing test prints and decide whether the change is
//     intended. An unintended change is a bug in conversion, not in the
//     golden file.
//  2. Only for an intended change, run:
//     SAVEPOINT_REGENERATE_MIGRATION_GOLDEN=1 go test ./internal/migrate/ -run TestEndToEnd_goldenConvertedOutput
//  3. Review the regenerated file as part of the change, and record the
//     behavior change in the task that caused it.
//
// The regeneration path always fails the test after writing, so a run with
// the variable set can never be mistaken for a passing verification.
const regenerateGoldenEnv = "SAVEPOINT_REGENERATE_MIGRATION_GOLDEN"

// goldenFile captures every file migration produced for one fixture, keyed by
// its project-relative path. Content is stored as a YAML string so CRLF and
// trailing whitespace survive the round trip exactly.
type goldenFile struct {
	Fixture string        `yaml:"fixture"`
	Note    string        `yaml:"note"`
	Files   []goldenEntry `yaml:"files"`
}

type goldenEntry struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

const goldenNote = "Converted output of this fixture under the fixed clock and operation id in " +
	"end_to_end_test.go. Regenerate only through the documented procedure at regenerateGoldenEnv; " +
	"a mismatch is a behavior change to review, not a file to refresh."

func TestEndToEnd_goldenConvertedOutput(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := readOnlyMigratedFixture(t, fixture)
			plan := planFor(t, filepath.Join(fixtureRoot, fixture, "project"))

			got := goldenFile{Fixture: fixture, Note: goldenNote, Files: convertedOutput(t, plan, root)}
			path := filepath.Join("testdata", "golden", fixture+".yml")

			if os.Getenv(regenerateGoldenEnv) != "" {
				writeGolden(t, path, got)
				t.Fatalf("%s was set: %s was regenerated. Review the diff and rerun without the variable.",
					regenerateGoldenEnv, path)
			}

			want := readGolden(t, path)
			assertGoldenEqual(t, path, want, got)
		})
	}
}

// TestEndToEnd_goldenIsReproducible is the byte-for-byte rerun half of AC9:
// a second independent migration of the same fixture, under the same injected
// clock and operation ID, produces identical converted content.
func TestEndToEnd_goldenIsReproducible(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			first := convertedOutputOf(t, fixture)
			second := convertedOutputOf(t, fixture)
			assertGoldenEqual(t, "rerun of "+fixture, goldenFile{Files: first}, goldenFile{Files: second})
		})
	}
}

// convertedOutput collects everything migration produced: each converted V2
// record, each relocated or rewritten project document, and the v1-to-v2.yml
// manifest that records the identity map. Archived copies are excluded — they
// are byte copies of the frozen sources, already proven by
// TestEndToEnd_archivedFilesMatchTheFixtureManifestHashes.
func convertedOutput(t *testing.T, plan *migrate.ConversionPlan, root string) []goldenEntry {
	t.Helper()

	paths := []string{".savepoint/migrations/v1-to-v2.yml"}
	for _, target := range plan.Targets {
		paths = append(paths, ".savepoint/"+target.InstallPath())
	}
	for _, doc := range plan.Documents {
		paths = append(paths, ".savepoint/"+doc.TargetPath)
	}
	sort.Strings(paths)

	entries := make([]goldenEntry, 0, len(paths))
	for _, path := range paths {
		content, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("read converted output %s: %v", path, err)
		}
		entries = append(entries, goldenEntry{Path: path, Content: string(content)})
	}
	return entries
}

func convertedOutputOf(t *testing.T, fixture string) []goldenEntry {
	t.Helper()
	root := copyFixture(t, fixture)
	plan := planFor(t, root)
	if _, out, err := runMigrate(t, root, "--apply"); err != nil {
		t.Fatalf("migrate --apply: %v\n%s", err, out)
	}
	return convertedOutput(t, plan, root)
}

func assertGoldenEqual(t *testing.T, label string, want, got goldenFile) {
	t.Helper()

	wantByPath := map[string]string{}
	for _, entry := range want.Files {
		wantByPath[entry.Path] = entry.Content
	}
	gotByPath := map[string]string{}
	for _, entry := range got.Files {
		gotByPath[entry.Path] = entry.Content
	}

	for path, wantContent := range wantByPath {
		gotContent, ok := gotByPath[path]
		if !ok {
			t.Errorf("%s: migration no longer produces %s — a behavior change to review, not a golden to refresh", label, path)
			continue
		}
		if gotContent != wantContent {
			t.Errorf("%s: converted content of %s changed — a behavior change to review, not a golden to refresh.\n--- recorded ---\n%s\n--- produced ---\n%s",
				label, path, wantContent, gotContent)
		}
	}
	for path := range gotByPath {
		if _, ok := wantByPath[path]; !ok {
			t.Errorf("%s: migration now produces %s, which is not recorded — a behavior change to review, not a golden to refresh", label, path)
		}
	}
}

func readGolden(t *testing.T, path string) goldenFile {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v\nIf this is a new fixture, create it with %s=1 (see regenerateGoldenEnv).", path, err, regenerateGoldenEnv)
	}
	var golden goldenFile
	if err := yaml.Unmarshal(raw, &golden); err != nil {
		t.Fatalf("parse golden %s: %v", path, err)
	}
	return golden
}

func writeGolden(t *testing.T, path string, golden goldenFile) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("create %s: %v", filepath.Dir(path), err)
	}
	raw, err := yaml.Marshal(golden)
	if err != nil {
		t.Fatalf("marshal golden %s: %v", path, err)
	}
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatalf("write golden %s: %v", path, err)
	}
}

// --- AC10: a second full run is a true no-op -----------------------------

func TestEndToEnd_secondFullRunChangesNothing(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := migrateFixture(t, fixture)

			before := snapshot(t, root)
			idsBefore := recordIDs(t, root)

			code, out, err := runMigrate(t, root, "--apply")
			if err != nil || code != 0 {
				t.Fatalf("second migrate --apply: code = %d, err = %v\n%s", code, err, out)
			}
			if !strings.Contains(out, "already") {
				t.Errorf("second run output does not report the project as already migrated:\n%s", out)
			}

			assertUnchanged(t, before, snapshot(t, root), "second full run over migrated "+fixture)
			if got := recordIDs(t, root); !equalStrings(idsBefore, got) {
				t.Errorf("record IDs changed on the second run:\n before %v\n after  %v", idsBefore, got)
			}
		})
	}
}

// --- AC11: preview writes nothing ----------------------------------------

func TestEndToEnd_previewWritesNothing(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			for _, args := range [][]string{{}, {"--dry-run"}, {"--apply", "--dry-run"}} {
				label := "default"
				if len(args) > 0 {
					label = strings.Join(args, " ")
				}
				t.Run(label, func(t *testing.T) {
					root := copyFixture(t, fixture)
					before := snapshot(t, root)

					code, out, err := runMigrate(t, append([]string{root}, args...)...)
					if err != nil || code != 0 {
						t.Fatalf("migrate %s: code = %d, err = %v\n%s", label, code, err, out)
					}
					if !strings.Contains(out, "Migration preview") {
						t.Errorf("migrate %s did not print a preview:\n%s", label, out)
					}

					assertUnchanged(t, before, snapshot(t, root), "preview ("+label+") over "+fixture)
				})
			}
		})
	}
}

// --- AC12: absent optional artifacts are not findings --------------------

// TestEndToEnd_absentOptionalArtifactsAreNotFindings asserts each fixture's
// manifest-declared absent_by_design paths produce no ambiguity, no conflict,
// and no mention in the preview report. v1-basic has no audit/ register and
// no Health-Check.md; v1-history has no AGENTS.md. Absence is expected V1
// shape, not a missing input.
func TestEndToEnd_absentOptionalArtifactsAreNotFindings(t *testing.T) {
	t.Parallel()
	for _, fixture := range e2eFixtures {
		t.Run(fixture, func(t *testing.T) {
			root := copyFixture(t, fixture)

			code, out, err := runMigrate(t, root)
			if err != nil || code != 0 {
				t.Fatalf("preview over %s: code = %d, err = %v\n%s", fixture, code, err, out)
			}

			plan := planFor(t, root)
			if len(plan.Conflicts) != 0 {
				t.Errorf("%s previewed with %d conflict(s): %+v", fixture, len(plan.Conflicts), plan.Conflicts)
			}

			for _, absent := range loadFixtureManifestFile(t, fixture).AbsentByDesign {
				path := fixtureRelPathOf(absent.Path)
				for _, ambiguity := range plan.Ambiguities {
					if strings.Contains(ambiguity.Path, path) || strings.Contains(ambiguity.Detail, path) {
						t.Errorf("%s: absent-by-design %s produced ambiguity %s: %s",
							fixture, path, ambiguity.ID, ambiguity.Detail)
					}
				}
				if strings.Contains(out, path) {
					t.Errorf("%s: preview report mentions absent-by-design %s:\n%s", fixture, path, out)
				}
			}

			// The absences never block: the project still migrates.
			migrateFixture(t, fixture)
		})
	}
}

// --- shared helpers ------------------------------------------------------

func planFor(t *testing.T, root string) *migrate.ConversionPlan {
	t.Helper()
	plan, err := migrate.Plan(root, nil,
		func() time.Time { return e2eClock },
		func() string { return e2eOperationID })
	if err != nil {
		t.Fatalf("Plan(%s) error = %v", root, err)
	}
	return plan
}

func mustIndex(t *testing.T, root string) *data.V2Index {
	t.Helper()
	index, err := data.LoadV2Index(filepath.Join(root, ".savepoint"))
	if err != nil {
		t.Fatalf("LoadV2Index(%s) error = %v", root, err)
	}
	return index
}

func readWrittenManifest(t *testing.T, root string) *migrate.ManifestV1ToV2 {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))
	if err != nil {
		t.Fatalf("read v1-to-v2.yml: %v", err)
	}
	manifest, err := migrate.UnmarshalManifest(raw)
	if err != nil {
		t.Fatalf("parse v1-to-v2.yml: %v", err)
	}
	return manifest
}

func assertNoEvidence(t *testing.T, label string, evidence *data.Evidence) {
	t.Helper()
	if evidence == nil {
		return
	}
	if evidence.LastCheck != "" {
		t.Errorf("%s carries last_check %s; no V1 source recorded one", label, evidence.LastCheck)
	}
	if evidence.Freshness != nil {
		t.Errorf("%s carries a freshness assessment; no V1 source recorded one", label)
	}
	if evidence.OwnerValidation != nil {
		t.Errorf("%s carries owner_validation; no V1 source recorded one", label)
	}
	if evidence.Exception != nil {
		t.Errorf("%s carries an exception; no V1 source recorded one", label)
	}
}

func oneActiveTaskID(t *testing.T, index *data.V2Index) string {
	t.Helper()
	ids := make([]string, 0, len(index.Tasks))
	for id, task := range index.Tasks {
		if task.Status != data.ColumnDone {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		t.Fatal("migrated project has no active Task")
	}
	sort.Strings(ids)
	return ids[0]
}

func hasBlocker(decision data.GateDecision, kind data.GateBlockKind) bool {
	for _, blocker := range decision.Blockers {
		if blocker.Kind == kind {
			return true
		}
	}
	return false
}

// recordIDs is every global identity the migrated project carries, sorted, so
// an idempotence test can prove no ID churned as well as no byte.
func recordIDs(t *testing.T, root string) []string {
	t.Helper()
	index := mustIndex(t, root)
	var ids []string
	for id := range index.Releases {
		ids = append(ids, "R:"+id)
	}
	for id := range index.Objectives {
		ids = append(ids, "O:"+id)
	}
	for id := range index.Tasks {
		ids = append(ids, "T:"+id)
	}
	for id := range index.Issues {
		ids = append(ids, "I:"+id)
	}
	for id := range index.Checks {
		ids = append(ids, "C:"+id)
	}
	sort.Strings(ids)
	return ids
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func reopen(t *testing.T, path, old, new string) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(raw), old) {
		t.Fatalf("%s does not contain %q", path, old)
	}
	updated := strings.Replace(string(raw), old, new, 1)
	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func mustExistAt(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s to exist: %v", path, err)
	}
}

// --- fixture manifests ---------------------------------------------------

type fixtureManifestDoc struct {
	Fixture        string                `yaml:"fixture"`
	Files          []fixtureFile         `yaml:"files"`
	AbsentByDesign []fixtureAbsentRecord `yaml:"absent_by_design"`
}

type fixtureFile struct {
	Path      string `yaml:"path"`
	SHA256    string `yaml:"sha256"`
	Role      string `yaml:"role"`
	RawStatus string `yaml:"raw_status"`
}

type fixtureAbsentRecord struct {
	Path   string `yaml:"path"`
	Reason string `yaml:"reason"`
}

func loadFixtureManifestFile(t *testing.T, fixture string) fixtureManifestDoc {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(fixtureRoot, fixture, "manifest.yml"))
	if err != nil {
		t.Fatalf("read %s manifest: %v", fixture, err)
	}
	var doc fixtureManifestDoc
	if err := yaml.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parse %s manifest: %v", fixture, err)
	}
	return doc
}

// loadFixtureFiles returns the manifest's file records with project-relative
// paths, the form Inventory, Classify, and the plan all use.
func loadFixtureFiles(t *testing.T, fixture string) []fixtureFile {
	t.Helper()
	doc := loadFixtureManifestFile(t, fixture)
	if len(doc.Files) == 0 {
		t.Fatalf("%s/manifest.yml lists no files", fixture)
	}
	files := make([]fixtureFile, len(doc.Files))
	for i, file := range doc.Files {
		file.Path = fixtureRelPathOf(file.Path)
		files[i] = file
	}
	return files
}

// fixtureRelPathOf turns a manifest path ("project/AGENTS.md") into the
// project-relative form ("AGENTS.md").
func fixtureRelPathOf(manifestPath string) string {
	return strings.TrimPrefix(filepath.ToSlash(manifestPath), "project/")
}

// --- filesystem helpers --------------------------------------------------

type fileFacts struct {
	sha256  string
	modTime time.Time
}

func snapshot(t *testing.T, root string) map[string]fileFacts {
	t.Helper()
	facts := map[string]fileFacts{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		facts[filepath.ToSlash(rel)] = fileFacts{sha256: hashFile(t, path), modTime: info.ModTime()}
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return facts
}

func assertUnchanged(t *testing.T, before, after map[string]fileFacts, label string) {
	t.Helper()
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Errorf("%s: %s was removed", label, path)
			continue
		}
		if got.sha256 != want.sha256 {
			t.Errorf("%s: %s content changed", label, path)
		}
		if !got.modTime.Equal(want.modTime) {
			t.Errorf("%s: %s modification time changed: before %v, after %v", label, path, want.modTime, got.modTime)
		}
	}
	for path := range after {
		if _, ok := before[path]; !ok {
			t.Errorf("%s: %s was created", label, path)
		}
	}
}

func hashFile(t *testing.T, path string) string {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()

	sum := sha256.New()
	if _, err := io.Copy(sum, f); err != nil {
		t.Fatalf("hash %s: %v", path, err)
	}
	return hex.EncodeToString(sum.Sum(nil))
}

// copyFixture copies fixture's project/ tree into a fresh temporary
// directory for any test that needs a mutable project. The frozen fixtures
// under internal/data/testdata/migration/ are never touched.
func copyFixture(t *testing.T, fixture string) string {
	t.Helper()
	dst := t.TempDir()
	if err := copyFixtureInto(dst, fixture); err != nil {
		t.Fatalf("copy fixture %s: %v", fixture, err)
	}
	return dst
}

func copyFixtureInto(dst, fixture string) error {
	src := filepath.Join(fixtureRoot, fixture, "project")
	return filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0644)
	})
}

// copyRepositoryWorkingTree copies the repository's working tree, excluding
// only .git metadata. Migration is exercised against this temporary copy so
// no test can ever activate schema_version or archive sources in the live
// repository.
func copyRepositoryWorkingTree(t *testing.T) string {
	t.Helper()
	src, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve repository root: %v", err)
	}
	dst := t.TempDir()
	err = filepath.WalkDir(src, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if rel == ".git" {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), content, 0644)
	})
	if err != nil {
		t.Fatalf("copy repository working tree: %v", err)
	}
	// The live repository is V2 after E50. Reconstruct the pre-migration
	// project boundary from the byte-preserved archive so this integration test
	// continues to exercise a real repository-scale V1 conversion without
	// treating the migrated checkout as a mutable fixture.
	if err := restoreArchivedV1Boundary(t, src, dst); err != nil {
		t.Fatalf("restore archived V1 boundary: %v", err)
	}
	return dst
}

func restoreArchivedV1Boundary(t *testing.T, repositoryRoot, destinationRoot string) error {
	t.Helper()
	archivedRoot := filepath.Join(repositoryRoot, ".savepoint", "archive", "v1")
	for _, pair := range []struct {
		source string
		target string
	}{
		{source: filepath.Join(archivedRoot, ".savepoint"), target: filepath.Join(destinationRoot, ".savepoint")},
		{source: filepath.Join(archivedRoot, "agent-skills"), target: filepath.Join(destinationRoot, "agent-skills")},
	} {
		if err := os.RemoveAll(pair.target); err != nil {
			return err
		}
		if err := copyTree(pair.source, pair.target); err != nil {
			return err
		}
	}
	return nil
}

func copyTree(sourceRoot, destinationRoot string) error {
	return filepath.WalkDir(sourceRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(sourceRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destinationRoot, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, 0644)
	})
}
