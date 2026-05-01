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
