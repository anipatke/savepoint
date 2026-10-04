package data

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const featureConfigBase = `# Savepoint config
verify_strict: false # soft mode
quality_gates:
    lint: null
    block_on_failure: true
theme:
    bg: "#000000"
    accents:
        planned: "#B1A1DF" # purple
# trailing note
schema_version: 2
`

func writeFeatureConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func readFeatureConfig(t *testing.T, path string) *Config {
	t.Helper()
	config, err := NewConfigReader().Read(path)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	return config
}

func TestParallelPlanningDecode(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"absent key", featureConfigBase, false},
		{"empty features", "features:\nschema_version: 2\n", false},
		{"other feature key only", "features:\n  future_thing: true\n", false},
		{"explicit false", "features:\n  parallel_planning: false\n", false},
		{"explicit true", "features:\n  parallel_planning: true\n", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := readFeatureConfig(t, writeFeatureConfig(t, tt.content))
			if got := config.ParallelPlanningEnabled(); got != tt.want {
				t.Errorf("ParallelPlanningEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParallelPlanningDefaultsOffWithoutConfigFile(t *testing.T) {
	config, err := NewConfigReader().Read(filepath.Join(t.TempDir(), "missing.yml"))
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if config.ParallelPlanningEnabled() {
		t.Error("missing config file enabled parallel planning")
	}
}

func TestParallelPlanningDecodeRejectsMalformedValues(t *testing.T) {
	tests := map[string]string{
		"string value":     "features:\n  parallel_planning: maybe\n",
		"mapping value":    "features:\n  parallel_planning:\n    nested: true\n",
		"features scalar":  "features: yes please\n",
		"features list":    "features:\n  - parallel_planning\n",
		"quoted not bool":  "features:\n  parallel_planning: \"sometimes\"\n",
		"sequence boolean": "features:\n  parallel_planning: [true]\n",
		"flow null":        "features: {parallel_planning: null}\n",
		"empty value":      "features:\n  parallel_planning:\n",
		"tilde null":       "features:\n  parallel_planning: ~\n",
		"duplicate key":    "features:\n  parallel_planning: false\n  parallel_planning: true\n",
	}
	for name, content := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewConfigReader().Read(writeFeatureConfig(t, content))
			if !errors.Is(err, ErrMalformedFeaturePreference) {
				t.Fatalf("Read() error = %v, want ErrMalformedFeaturePreference", err)
			}
			if !strings.Contains(err.Error(), "features") {
				t.Errorf("diagnostic %q does not name the key", err)
			}
		})
	}
}

func TestWriteParallelPlanningPreservesAuthoredBytes(t *testing.T) {
	tests := []struct {
		name    string
		content string
		enable  bool
		want    string
	}{
		{
			name:    "adds block at end of file",
			content: featureConfigBase,
			enable:  true,
			want:    featureConfigBase + "features:\n    parallel_planning: true\n",
		},
		{
			name:    "adds block after file without trailing newline",
			content: "theme:\n  bg: \"#000\"",
			enable:  true,
			want:    "theme:\n  bg: \"#000\"\nfeatures:\n  parallel_planning: true\n",
		},
		{
			name:    "replaces value keeping comment",
			content: "a: 1\nfeatures:\n  parallel_planning: false # owner note\nz: 2\n",
			enable:  true,
			want:    "a: 1\nfeatures:\n  parallel_planning: true # owner note\nz: 2\n",
		},
		{
			name:    "turns off keeping indentation",
			content: "features:\n    parallel_planning: True\n",
			enable:  false,
			want:    "features:\n    parallel_planning: false\n",
		},
		{
			name:    "adds key to existing block before next section's comments",
			content: "features:\n  future_thing: true\n\n# theme section\ntheme:\n  bg: x\n",
			enable:  true,
			want:    "features:\n  future_thing: true\n  parallel_planning: true\n\n# theme section\ntheme:\n  bg: x\n",
		},
		{
			name:    "fills empty features key",
			content: "features:\ntheme:\n  bg: x\n",
			enable:  true,
			want:    "features:\n  parallel_planning: true\ntheme:\n  bg: x\n",
		},
		{
			name:    "flow mapping value",
			content: "features: {parallel_planning: false}\n",
			enable:  true,
			want:    "features: {parallel_planning: true}\n",
		},
		{
			name:    "keeps CRLF line endings",
			content: "a: 1\r\nfeatures:\r\n  parallel_planning: false\r\n",
			enable:  true,
			want:    "a: 1\r\nfeatures:\r\n  parallel_planning: true\r\n",
		},
		{
			name:    "keeps mixed endings when replacing a value",
			content: "schema_version: 2\r\n# keep LF\nfeatures:\r\n  parallel_planning: false\r\n",
			enable:  true,
			want:    "schema_version: 2\r\n# keep LF\nfeatures:\r\n  parallel_planning: true\r\n",
		},
		{
			name:    "keeps mixed endings with LF value line",
			content: "a: 1\r\nfeatures:\n  parallel_planning: true\n# keep CRLF\r\n",
			enable:  false,
			want:    "a: 1\r\nfeatures:\n  parallel_planning: false\n# keep CRLF\r\n",
		},
		{
			name:    "appends block to CRLF file without trailing newline",
			content: "a: 1\r\nb: 2",
			enable:  true,
			want:    "a: 1\r\nb: 2\r\nfeatures:\r\n    parallel_planning: true\r\n",
		},
		{
			name:    "inserted key copies the preceding line ending",
			content: "# top\nfeatures:\r\n  future_thing: true\r\nz: 2\n",
			enable:  true,
			want:    "# top\nfeatures:\r\n  future_thing: true\r\n  parallel_planning: true\r\nz: 2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeFeatureConfig(t, tt.content)
			config := readFeatureConfig(t, path)
			if err := WriteParallelPlanning(path, tt.enable, config.FeatureSource()); err != nil {
				t.Fatalf("WriteParallelPlanning() error = %v", err)
			}
			got, _ := os.ReadFile(path)
			if string(got) != tt.want {
				t.Errorf("file = %q, want %q", got, tt.want)
			}
			if reloaded := readFeatureConfig(t, path); reloaded.ParallelPlanningEnabled() != tt.enable {
				t.Errorf("reload = %v, want %v", reloaded.ParallelPlanningEnabled(), tt.enable)
			}
			assertNoTempFiles(t, filepath.Dir(path))
		})
	}
}

func TestWriteParallelPlanningToggleRoundTripKeepsUnrelatedConfig(t *testing.T) {
	path := writeFeatureConfig(t, featureConfigBase)
	for _, enable := range []bool{true, false, true} {
		config := readFeatureConfig(t, path)
		if err := WriteParallelPlanning(path, enable, config.FeatureSource()); err != nil {
			t.Fatalf("WriteParallelPlanning(%v) error = %v", enable, err)
		}
	}
	got, _ := os.ReadFile(path)
	if !strings.HasPrefix(string(got), featureConfigBase) {
		t.Errorf("unrelated config changed:\n%s", got)
	}
	config := readFeatureConfig(t, path)
	if !config.ParallelPlanningEnabled() || config.QualityGates.BlockOnFailure != true {
		t.Errorf("config after toggles = %+v", config)
	}
}

func TestWriteParallelPlanningUnchangedValueIsNoOp(t *testing.T) {
	for _, content := range []string{featureConfigBase, "features:\n  parallel_planning: true\n"} {
		path := writeFeatureConfig(t, content)
		config := readFeatureConfig(t, path)
		old := time.Now().Add(-time.Hour)
		if err := os.Chtimes(path, old, old); err != nil {
			t.Fatal(err)
		}
		before, _ := os.Stat(path)
		enabled := config.ParallelPlanningEnabled()
		if err := WriteParallelPlanning(path, enabled, config.FeatureSource()); err != nil {
			t.Fatalf("WriteParallelPlanning() error = %v", err)
		}
		info, _ := os.Stat(path)
		if !info.ModTime().Equal(before.ModTime()) {
			t.Errorf("no-op rewrote %q", content)
		}
		got, _ := os.ReadFile(path)
		if string(got) != content {
			t.Errorf("no-op changed file to %q", got)
		}
	}
}

func TestWriteParallelPlanningRefusesStaleSource(t *testing.T) {
	path := writeFeatureConfig(t, featureConfigBase)
	config := readFeatureConfig(t, path)
	edited := featureConfigBase + "extra: external edit\n"
	if err := os.WriteFile(path, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	err := WriteParallelPlanning(path, true, config.FeatureSource())
	if !errors.Is(err, ErrV2SourceConflict) || !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("error = %v, want source conflict", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != edited {
		t.Errorf("stale write changed file: %q", got)
	}
}

func TestWriteParallelPlanningRefusesAMissingConfigSource(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "config.yml")
	config, _ := NewConfigReader().Read(missing)
	if err := WriteParallelPlanning(missing, true, config.FeatureSource()); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error = %v, want not exist", err)
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Error("write created config.yml")
	}
	// A file that appears after a missing-file load is a conflict, not a write.
	path := writeFeatureConfig(t, featureConfigBase)
	if err := WriteParallelPlanning(path, true, config.FeatureSource()); !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("error = %v, want conflict", err)
	}
}

func TestWriteParallelPlanningRefusesMalformedConfig(t *testing.T) {
	for name, content := range map[string]string{
		"invalid yaml":    "a: [unclosed\n",
		"bad value":       "features:\n  parallel_planning: maybe\n",
		"flow no key":     "features: {other: 1}\n",
		"features scalar": "features: 3\n",
	} {
		t.Run(name, func(t *testing.T) {
			path := writeFeatureConfig(t, content)
			err := WriteParallelPlanning(path, true, newFeatureSource([]byte(content)))
			if err == nil {
				t.Fatal("expected refusal")
			}
			got, _ := os.ReadFile(path)
			if string(got) != content {
				t.Errorf("refused write changed file: %q", got)
			}
		})
	}
}

func TestWriteParallelPlanningFailedReplaceLeavesOriginal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permission bits do not block file creation on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("root ignores directory permissions")
	}
	path := writeFeatureConfig(t, featureConfigBase)
	config := readFeatureConfig(t, path)
	dir := filepath.Dir(path)
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0o755)

	if err := WriteParallelPlanning(path, true, config.FeatureSource()); err == nil {
		t.Fatal("expected write failure in read-only directory")
	}
	got, _ := os.ReadFile(path)
	if string(got) != featureConfigBase {
		t.Errorf("failed write changed file: %q", got)
	}
	if readFeatureConfig(t, path).ParallelPlanningEnabled() {
		t.Error("failed write reads back as saved")
	}
}

func TestWriteParallelPlanningRefusesReadOnlyFileAndSymlink(t *testing.T) {
	path := writeFeatureConfig(t, featureConfigBase)
	config := readFeatureConfig(t, path)

	if runtime.GOOS != "windows" {
		link := filepath.Join(filepath.Dir(path), "link.yml")
		if err := os.Symlink(path, link); err != nil {
			t.Fatal(err)
		}
		if err := WriteParallelPlanning(link, true, config.FeatureSource()); !errors.Is(err, ErrV2UnsafePath) {
			t.Errorf("symlink error = %v, want ErrV2UnsafePath", err)
		}
	}

	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := WriteParallelPlanning(path, true, config.FeatureSource()); !errors.Is(err, os.ErrPermission) {
		t.Errorf("read-only error = %v, want permission", err)
	}
	got, _ := os.ReadFile(path)
	if string(got) != featureConfigBase {
		t.Errorf("refused write changed file: %q", got)
	}
}

func assertNoTempFiles(t *testing.T, dir string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".savepoint-v2-write-") {
			t.Errorf("left temporary file %s", entry.Name())
		}
	}
}
