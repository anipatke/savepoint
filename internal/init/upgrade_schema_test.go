package init

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/opencode/savepoint/internal/data"
	"github.com/opencode/savepoint/internal/migrate"
	"github.com/opencode/savepoint/internal/testutil"
)

// Tests in this file cover UpgradeProjectAssets' own job: reading a project's
// schema_version and selecting the matching template tree. The upgrade policy
// each tree is walked with (skill conflicts, AGENTS.md merge, audit/policy
// assets, dry run, idempotency) is covered exhaustively in upgrade_test.go
// against upgradeAssetsFromTree, the shared single-tree core both schema
// branches call.

// schemaProject builds a bare Savepoint project directory and, when config is
// non-empty, writes it as the project's .savepoint/config.yml before the
// upgrade runs. An empty config leaves the project exactly as an
// unmigrated V1 project looks: no config.yml at all.
func schemaProject(t *testing.T, config string) string {
	t.Helper()
	dir := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
	if config != "" {
		testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), config)
	}
	return dir
}

// v1v2Templates returns a small V2 template tree distinct from any real V2
// skill name, so a test can tell a walked install apart from one it did not
// expect.
func v1v2Templates() fstest.MapFS {
	return fstest.MapFS{
		"agent-skills/savepoint-idea/SKILL.md": &fstest.MapFile{Data: []byte("# V2 skill")},
	}
}

// dirSnapshot records every file's content under dir, keyed by path relative
// to dir, so a test can prove a refused or no-op upgrade touched nothing.
func dirSnapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	snap := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		snap[rel] = string(content)
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", dir, err)
	}
	return snap
}

func mtimeSnapshot(t *testing.T, dir string) map[string]int64 {
	t.Helper()
	snap := map[string]int64{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		snap[rel] = info.ModTime().UnixNano()
		return nil
	})
	if err != nil {
		t.Fatalf("mtime snapshot %s: %v", dir, err)
	}
	return snap
}

func assertNoChange(t *testing.T, dir string, before map[string]string) {
	t.Helper()
	after := dirSnapshot(t, dir)
	if len(after) != len(before) {
		t.Fatalf("file count changed: before %d, after %d (%v -> %v)", len(before), len(after), before, after)
	}
	for path, want := range before {
		got, ok := after[path]
		if !ok {
			t.Errorf("%s disappeared", path)
			continue
		}
		if got != want {
			t.Errorf("%s content changed", path)
		}
	}
}

// dirDiff reports every relative path whose content differs, or which exists
// in only one of a and b.
func dirDiff(t *testing.T, a, b string) string {
	t.Helper()
	sa := dirSnapshot(t, a)
	sb := dirSnapshot(t, b)
	var diffs []string
	for path, va := range sa {
		vb, ok := sb[path]
		if !ok {
			diffs = append(diffs, fmt.Sprintf("%s only in %s", path, a))
			continue
		}
		if va != vb {
			diffs = append(diffs, fmt.Sprintf("%s differs", path))
		}
	}
	for path := range sb {
		if _, ok := sa[path]; !ok {
			diffs = append(diffs, fmt.Sprintf("%s only in %s", path, b))
		}
	}
	return strings.Join(diffs, "\n")
}

func TestUpgradeProjectAssets_refusesV1WithoutMutation(t *testing.T) {
	dir := schemaProject(t, "")
	v2 := v1v2Templates()
	before := dirSnapshot(t, dir)

	report, err := UpgradeProjectAssets(v2, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if len(report.Actions) != 1 || report.Actions[0].Action != ActionInfo {
		t.Fatalf("V1 refusal report = %+v, want one informational entry", report.Actions)
	}
	if !strings.Contains(report.Actions[0].Note, "savepoint migrate") {
		t.Errorf("V1 refusal note = %q, want migration guidance", report.Actions[0].Note)
	}
	assertNoChange(t, dir, before)
}

func TestUpgradeProjectAssets_selectsV2TreeForSchemaVersion2(t *testing.T) {
	dir := schemaProject(t, "schema_version: 2\n")
	v2 := v1v2Templates()

	report, err := UpgradeProjectAssets(v2, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "agent-skills", "savepoint-idea", "SKILL.md")); err != nil {
		t.Errorf("V2 skill not installed: %v", err)
	}

	for _, e := range report.Actions {
		if e.Action == ActionInfo {
			t.Errorf("V2 project report carries a migrate-route note: %+v", e)
		}
	}
}

func TestUpgradeProjectAssets_v1ProjectCarriesMigrateRouteNote(t *testing.T) {
	dir := schemaProject(t, "")
	v2 := v1v2Templates()

	report, err := UpgradeProjectAssets(v2, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	var infoEntries []UpgradeEntry
	for _, e := range report.Actions {
		if e.Action == ActionInfo {
			infoEntries = append(infoEntries, e)
		}
	}
	if len(infoEntries) != 1 {
		t.Fatalf("info entries = %d, want exactly 1: %+v", len(infoEntries), infoEntries)
	}
	if !strings.Contains(infoEntries[0].Note, "savepoint migrate") {
		t.Errorf("info note = %q, want it to name savepoint migrate", infoEntries[0].Note)
	}
}

func TestUpgradeProjectAssets_malformedSchemaVersionRefusesCleanly(t *testing.T) {
	dir := schemaProject(t, "schema_version: [not, a, scalar]\n")
	v2 := v1v2Templates()
	before := dirSnapshot(t, dir)

	_, err := UpgradeProjectAssets(v2, dir, false, false)
	if !errors.Is(err, data.ErrMalformedSchemaVersion) {
		t.Fatalf("err = %v, want ErrMalformedSchemaVersion", err)
	}
	configPath := filepath.Join(dir, ".savepoint", "config.yml")
	if !strings.Contains(err.Error(), configPath) {
		t.Errorf("err = %q, want it to name the config path %q", err.Error(), configPath)
	}

	assertNoChange(t, dir, before)
	if _, err := os.Stat(filepath.Join(dir, ".savepoint", ".upgrade-manifest.yml")); !os.IsNotExist(err) {
		t.Errorf("refused upgrade created a manifest, stat err = %v", err)
	}
}

func TestUpgradeProjectAssets_unsupportedSchemaVersionRefusesCleanly(t *testing.T) {
	dir := schemaProject(t, "schema_version: 99\n")
	v2 := v1v2Templates()
	before := dirSnapshot(t, dir)

	_, err := UpgradeProjectAssets(v2, dir, false, false)
	if !errors.Is(err, data.ErrUnsupportedSchemaVersion) {
		t.Fatalf("err = %v, want ErrUnsupportedSchemaVersion", err)
	}
	configPath := filepath.Join(dir, ".savepoint", "config.yml")
	if !strings.Contains(err.Error(), configPath) {
		t.Errorf("err = %q, want it to name the config path %q", err.Error(), configPath)
	}

	assertNoChange(t, dir, before)
	if _, err := os.Stat(filepath.Join(dir, ".savepoint", ".upgrade-manifest.yml")); !os.IsNotExist(err) {
		t.Errorf("refused upgrade created a manifest, stat err = %v", err)
	}
}

func TestUpgradeProjectAssets_dryRunRefusesUnreadableSchemaVersionTheSameWay(t *testing.T) {
	for _, config := range []string{"schema_version: [oops]\n", "schema_version: 99\n"} {
		t.Run(config, func(t *testing.T) {
			dir := schemaProject(t, config)
			v2 := v1v2Templates()
			before := dirSnapshot(t, dir)

			_, err := UpgradeProjectAssets(v2, dir, true, false)
			if err == nil {
				t.Fatal("dry-run error = nil, want a refusal")
			}
			assertNoChange(t, dir, before)
		})
	}
}

func TestUpgradeProjectAssets_dryRunReachesSameTreeSelectionAndWritesNothing(t *testing.T) {
	cases := []struct{ name, config string }{
		{name: "v1 implicit", config: ""},
		{name: "v2 explicit", config: "schema_version: 2\n"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dryDir := schemaProject(t, c.config)
			v2 := v1v2Templates()
			before := dirSnapshot(t, dryDir)

			dryReport, err := UpgradeProjectAssets(v2, dryDir, true, false)
			if err != nil {
				t.Fatalf("dry-run error = %v", err)
			}
			assertNoChange(t, dryDir, before)

			realDir := schemaProject(t, c.config)
			realReport, err := UpgradeProjectAssets(v2, realDir, false, false)
			if err != nil {
				t.Fatalf("real run error = %v", err)
			}

			if len(dryReport.Actions) != len(realReport.Actions) {
				t.Fatalf("dry-run actions = %d, real run actions = %d", len(dryReport.Actions), len(realReport.Actions))
			}
			for i := range dryReport.Actions {
				if dryReport.Actions[i].Path != realReport.Actions[i].Path || dryReport.Actions[i].Action != realReport.Actions[i].Action {
					t.Errorf("action[%d] dry=%+v real=%+v", i, dryReport.Actions[i], realReport.Actions[i])
				}
			}
		})
	}
}

func TestUpgradeProjectAssets_neverWritesSchemaVersion(t *testing.T) {
	cases := []struct{ name, config string }{
		{name: "v1", config: ""},
		{name: "v2", config: "schema_version: 2\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := schemaProject(t, c.config)
			v2 := v1v2Templates()
			configPath := filepath.Join(dir, ".savepoint", "config.yml")

			var before []byte
			if c.config != "" {
				var err error
				before, err = os.ReadFile(configPath)
				if err != nil {
					t.Fatal(err)
				}
			}

			if _, err := UpgradeProjectAssets(v2, dir, false, false); err != nil {
				t.Fatalf("UpgradeProjectAssets() error = %v", err)
			}

			after, err := os.ReadFile(configPath)
			if c.config == "" {
				if !os.IsNotExist(err) {
					t.Errorf("upgrade wrote config.yml to a V1 project that had none: err = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("config.yml disappeared: %v", err)
			}
			if string(after) != string(before) {
				t.Errorf("config.yml changed: before %q, after %q", before, after)
			}
		})
	}
}

func TestUpgradeProjectAssets_secondRunIsNoOp(t *testing.T) {
	cases := []struct{ name, config string }{
		{name: "v1", config: ""},
		{name: "v2", config: "schema_version: 2\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := schemaProject(t, c.config)
			v2 := v1v2Templates()

			if _, err := UpgradeProjectAssets(v2, dir, false, false); err != nil {
				t.Fatalf("first run error = %v", err)
			}

			before := dirSnapshot(t, dir)
			beforeTimes := mtimeSnapshot(t, dir)

			report, err := UpgradeProjectAssets(v2, dir, false, false)
			if err != nil {
				t.Fatalf("second run error = %v", err)
			}

			assertNoChange(t, dir, before)
			afterTimes := mtimeSnapshot(t, dir)
			for path, mt := range beforeTimes {
				if afterTimes[path] != mt {
					t.Errorf("%s mtime changed on a no-op rerun", path)
				}
			}

			for _, e := range report.Actions {
				if e.Action != ActionUnchanged && e.Action != ActionSkipped && e.Action != ActionInfo {
					t.Errorf("second run action %v for %q, want unchanged/skipped/info", e.Action, e.Path)
				}
			}
		})
	}
}

// TestUpgradeProjectAssets_refusesPendingMigrationOnBothTrees proves the
// pending-migration guard fires before either the V1 refusal or the V2 asset
// walk, so an incomplete operation always owns the project.
func TestUpgradeProjectAssets_refusesPendingMigrationOnBothTrees(t *testing.T) {
	cases := []struct{ name, config string }{
		{name: "v1", config: ""},
		{name: "v2", config: "schema_version: 2\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := schemaProject(t, c.config)
			if _, err := migrate.CreateOperation(dir, "op-1", nil, []migrate.JournalEntry{{Path: "objectives/O001.md", Action: migrate.ActionCreate}}, time.Now()); err != nil {
				t.Fatalf("CreateOperation() error = %v", err)
			}
			v2 := v1v2Templates()

			_, err := UpgradeProjectAssets(v2, dir, false, false)
			if err == nil {
				t.Fatal("UpgradeProjectAssets() error = nil, want refusal while migration operation is incomplete")
			}
			if !strings.Contains(err.Error(), "op-1") {
				t.Errorf("error = %q, want it to name the operation op-1", err.Error())
			}

			if _, err := os.Stat(filepath.Join(dir, "agent-skills", "savepoint-idea", "SKILL.md")); !os.IsNotExist(err) {
				t.Errorf("refused upgrade wrote agent-skills/savepoint-idea/SKILL.md, stat err = %v", err)
			}
		})
	}
}

func TestUpgradeReport_formatIncludesInfoCount(t *testing.T) {
	r := &UpgradeReport{
		Actions: []UpgradeEntry{
			{Action: ActionInfo, Note: noteMigrateRoute},
		},
	}
	output := r.Format()
	if !strings.Contains(output, "Info: 1") {
		t.Errorf("missing info count: %q", output)
	}
	if !strings.Contains(output, noteMigrateRoute) {
		t.Errorf("missing info note: %q", output)
	}
}
