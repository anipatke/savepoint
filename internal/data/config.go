package data

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Theme struct {
	BG       string            `yaml:"bg"`
	Surface  string            `yaml:"surface"`
	Surface2 string            `yaml:"surface_2"`
	Border   string            `yaml:"border"`
	Text     string            `yaml:"text"`
	Accents  map[string]string `yaml:"accents"`
}

type QualityGates struct {
	Lint           *string `yaml:"lint"`
	Typecheck      *string `yaml:"typecheck"`
	Test           *string `yaml:"test"`
	BlockOnFailure bool    `yaml:"block_on_failure"`
	Timeout        string  `yaml:"gate_timeout"`
}

type Config struct {
	Theme         Theme         `yaml:"theme"`
	QualityGates  QualityGates  `yaml:"quality_gates"`
}

var defaultTheme = Theme{
	BG:       "#1a1b26",
	Surface:  "#24283b",
	Surface2: "#414868",
	Border:   "#565f89",
	Text:     "#c0caf5",
	Accents: map[string]string{
		"planned":     "#7aa2f7",
		"in_progress": "#bb9af7",
		"done":        "#9ece6a",
		"blocked":     "#f7768e",
		"epic":        "#2ac3de",
	},
}

var defaultConfig = Config{
	Theme: defaultTheme,
}

type ConfigReader struct{}

func NewConfigReader() *ConfigReader {
	return &ConfigReader{}
}

func (r *ConfigReader) Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &defaultConfig, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	config.Theme = fillThemeDefaults(config.Theme)

	return &config, nil
}

func fillThemeDefaults(theme Theme) Theme {
	if theme.BG == "" {
		theme.BG = defaultTheme.BG
	}
	if theme.Surface == "" {
		theme.Surface = defaultTheme.Surface
	}
	if theme.Surface2 == "" {
		theme.Surface2 = defaultTheme.Surface2
	}
	if theme.Border == "" {
		theme.Border = defaultTheme.Border
	}
	if theme.Text == "" {
		theme.Text = defaultTheme.Text
	}
	if theme.Accents == nil {
		theme.Accents = make(map[string]string)
	}
	for k, v := range defaultTheme.Accents {
		if _, ok := theme.Accents[k]; !ok {
			theme.Accents[k] = v
		}
	}
	return theme
}
