package data

import (
	"os"
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
