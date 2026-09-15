package data

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigReaderDefault(t *testing.T) {
	r := NewConfigReader()
	config, err := r.Read("nonexistent.yml")
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if config.Theme.BG != defaultTheme.BG {
		t.Errorf("Theme.BG = %v, want %v", config.Theme.BG, defaultTheme.BG)
	}
}

func TestConfigReaderRead(t *testing.T) {
	content := `theme:
  bg: "#000000"
  surface: "#111111"
  text: "#ffffff"
  accents:
    planned: "#222222"
`
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	r := NewConfigReader()
	config, err := r.Read(tmpfile.Name())
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}

	if config.Theme.BG != "#000000" {
		t.Errorf("Theme.BG = %v, want #000000", config.Theme.BG)
	}
	if config.Theme.Surface2 != defaultTheme.Surface2 {
		t.Errorf("Theme.Surface2 = %v, want default %v", config.Theme.Surface2, defaultTheme.Surface2)
	}
	if config.Theme.Accents["planned"] != "#222222" {
		t.Errorf("Theme.Accents[planned] = %v, want #222222", config.Theme.Accents["planned"])
	}
}

func TestFillThemeDefaults_PartialAccents(t *testing.T) {
	theme := Theme{
		BG:      "#000000",
		Accents: map[string]string{"planned": "#ff0000"},
	}
	result := fillThemeDefaults(theme)
	if result.Accents["planned"] != "#ff0000" {
		t.Errorf("Accents[planned] = %v, want #ff0000 (user value preserved)", result.Accents["planned"])
	}
	if result.Accents["in_progress"] != defaultTheme.Accents["in_progress"] {
		t.Errorf("Accents[in_progress] = %v, want default %v", result.Accents["in_progress"], defaultTheme.Accents["in_progress"])
	}
	if result.Accents["done"] != defaultTheme.Accents["done"] {
		t.Errorf("Accents[done] = %v, want default %v", result.Accents["done"], defaultTheme.Accents["done"])
	}
	if result.Accents["blocked"] != defaultTheme.Accents["blocked"] {
		t.Errorf("Accents[blocked] = %v, want default %v", result.Accents["blocked"], defaultTheme.Accents["blocked"])
	}
	if result.Accents["epic"] != defaultTheme.Accents["epic"] {
		t.Errorf("Accents[epic] = %v, want default %v", result.Accents["epic"], defaultTheme.Accents["epic"])
	}
}

func TestFillThemeDefaults_NilAccents(t *testing.T) {
	theme := Theme{
		BG:      "#000000",
		Accents: nil,
	}
	result := fillThemeDefaults(theme)
	for k, v := range defaultTheme.Accents {
		if result.Accents[k] != v {
			t.Errorf("Accents[%s] = %v, want default %v", k, result.Accents[k], v)
		}
	}
}

func TestFillThemeDefaults_EmptyAccents(t *testing.T) {
	theme := Theme{
		BG:      "#000000",
		Accents: map[string]string{},
	}
	result := fillThemeDefaults(theme)
	for k, v := range defaultTheme.Accents {
		if result.Accents[k] != v {
			t.Errorf("Accents[%s] = %v, want default %v", k, result.Accents[k], v)
		}
	}
}

func TestReadSchemaVersion(t *testing.T) {
	cases := []struct {
		name    string
		content string
		absent  bool
		want    SchemaVersion
		wantErr error
	}{
		{
			name:   "absent file selects V1",
			absent: true,
			want:   SchemaVersionV1,
		},
		{
			name:    "absent field selects V1",
			content: "theme:\n  bg: \"#000000\"\n",
			want:    SchemaVersionV1,
		},
		{
			name:    "explicit version 2 selects V2",
			content: "schema_version: 2\n",
			want:    SchemaVersionV2,
		},
		{
			name:    "malformed non-integer version fails named",
			content: "schema_version: not-a-number\n",
			want:    SchemaVersionV1,
			wantErr: ErrMalformedSchemaVersion,
		},
		{
			name:    "unsupported explicit version fails named",
			content: "schema_version: 3\n",
			want:    SchemaVersionV1,
			wantErr: ErrUnsupportedSchemaVersion,
		},
		{
			name:    "version 1 is an unsupported explicit version",
			content: "schema_version: 1\n",
			want:    SchemaVersionV1,
			wantErr: ErrUnsupportedSchemaVersion,
		},
		{
			name:    "unrelated version-shaped fields do not select V2",
			content: "theme:\n  bg: \"#000000\"\nagent_launcher:\n  terminal:\n    mode: auto\n# package version: 1.3.1, release: v2, upgrade-manifest schema: 2\n",
			want:    SchemaVersionV1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var path string
			if tc.absent {
				path = filepath.Join(t.TempDir(), "config.yml")
			} else {
				dir := t.TempDir()
				path = filepath.Join(dir, "config.yml")
				if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			got, err := ReadSchemaVersion(path)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ReadSchemaVersion() error = %v, want wrapping %v", err, tc.wantErr)
				}
				if !strings.Contains(err.Error(), path) {
					t.Errorf("ReadSchemaVersion() error = %v, want it to identify path %v", err, path)
				}
			} else if err != nil {
				t.Fatalf("ReadSchemaVersion() unexpected error = %v", err)
			}
			if got != tc.want {
				t.Errorf("ReadSchemaVersion() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestConfigReaderMalformedYAML(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "config-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte("theme: [broken")); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	r := NewConfigReader()
	_, err = r.Read(tmpfile.Name())
	if err == nil {
		t.Fatal("Read() expected malformed YAML error")
	}
}
