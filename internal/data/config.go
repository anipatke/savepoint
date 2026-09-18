package data

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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
	AgentLauncher AgentLauncher `yaml:"agent_launcher"`
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
	Theme:         defaultTheme,
	AgentLauncher: AgentLauncher{Terminal: TerminalConfig{Mode: TerminalModeAuto}},
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
	config.AgentLauncher = fillLauncherDefaults(config.AgentLauncher)

	if err := config.AgentLauncher.Validate(); err != nil {
		return nil, fmt.Errorf("invalid launcher config: %w", err)
	}

	return &config, nil
}

// SchemaVersion identifies the project record schema a .savepoint project
// declares via config.yml. It is independent of package, release, and
// .upgrade-manifest.yml versions.
type SchemaVersion int

const (
	// SchemaVersionV1 is selected when config.yml has no explicit
	// schema_version. It signals transitional V1 discovery behavior, not a
	// declared version.
	SchemaVersionV1 SchemaVersion = 0
	SchemaVersionV2 SchemaVersion = 2
)

type schemaVersionDoc struct {
	SchemaVersion yaml.Node `yaml:"schema_version"`
}

// ReadSchemaVersion reads only the schema_version field from the config.yml
// at path, independent of theme/quality_gates/launcher parsing. Absence of
// the file or the field selects SchemaVersionV1. A present but non-integer
// value returns ErrMalformedSchemaVersion; a present integer other than 2
// returns ErrUnsupportedSchemaVersion. Both errors identify path and the
// supplied value.
func ReadSchemaVersion(path string) (SchemaVersion, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return SchemaVersionV1, nil
	}
	if err != nil {
		return SchemaVersionV1, fmt.Errorf("failed to read config: %w", err)
	}

	var doc schemaVersionDoc
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return SchemaVersionV1, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	if doc.SchemaVersion.Kind == 0 {
		return SchemaVersionV1, nil
	}

	var version int
	if err := doc.SchemaVersion.Decode(&version); err != nil {
		return SchemaVersionV1, fmt.Errorf("%w: %s: schema_version %q", ErrMalformedSchemaVersion, path, doc.SchemaVersion.Value)
	}

	if SchemaVersion(version) != SchemaVersionV2 {
		return SchemaVersionV1, fmt.Errorf("%w: %s: schema_version %d", ErrUnsupportedSchemaVersion, path, version)
	}

	return SchemaVersionV2, nil
}

// ActivateSchemaV2 sets schema_version to 2 in content — the raw bytes of a
// project's config.yml — preserving every other key, comment, and
// structural detail through a YAML node patch rather than a decode/re-encode
// of the typed Config, matching DATA-01. content may be nil or empty, in
// which case the result is a fresh document containing only
// schema_version: 2, for a V1 project that never had a config.yml at all.
// ActivateSchemaV2 is a pure transform: it never touches the filesystem,
// and it is migrate's apply step that performs the actual replace, as the
// operation's last write.
func ActivateSchemaV2(content []byte) ([]byte, error) {
	var doc yaml.Node
	if len(strings.TrimSpace(string(content))) > 0 {
		if err := yaml.Unmarshal(content, &doc); err != nil {
			return nil, fmt.Errorf("parse config.yml: %w", err)
		}
	}
	if doc.Kind == 0 {
		doc = yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Kind: yaml.MappingNode}}}
	}
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, fmt.Errorf("config.yml: unexpected yaml structure")
	}

	setSchemaVersionNode(doc.Content[0], SchemaVersionV2)

	out, err := yaml.Marshal(&doc)
	if err != nil {
		return nil, fmt.Errorf("marshal config.yml: %w", err)
	}
	return out, nil
}

// setSchemaVersionNode sets mapping's schema_version key to version as an
// integer-tagged scalar, in place, adding the key if absent. It is the only
// place ActivateSchemaV2 mutates the parsed document, so every other
// mapping entry — theme, quality_gates, and any unknown key — passes
// through the marshal untouched.
func setSchemaVersionNode(mapping *yaml.Node, version SchemaVersion) {
	valStr := strconv.Itoa(int(version))
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		if mapping.Content[i].Value == "schema_version" {
			mapping.Content[i+1] = &yaml.Node{Kind: yaml.ScalarNode, Value: valStr, Tag: "!!int"}
			return
		}
	}
	mapping.Content = append(mapping.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "schema_version", Tag: "!!str"},
		&yaml.Node{Kind: yaml.ScalarNode, Value: valStr, Tag: "!!int"},
	)
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
