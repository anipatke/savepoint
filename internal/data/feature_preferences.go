package data

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	featuresKey         = "features"
	parallelPlanningKey = "parallel_planning"
)

// ErrMalformedFeaturePreference reports a features value config.yml carries
// that this package cannot interpret. It identifies the key and the supplied
// value so the owner can correct the file by hand.
var ErrMalformedFeaturePreference = errors.New("malformed feature preference")

// FeaturePreferences holds the project's optional-feature choices from the
// features mapping in config.yml. Every field defaults to off when absent.
type FeaturePreferences struct {
	// ParallelPlanning opts in to optional lane advice. It never gates
	// lifecycle, Code Health, or completion.
	ParallelPlanning bool
}

// UnmarshalYAML reads features.parallel_planning as a boolean and ignores
// other keys, so a newer project's choices pass through an older reader.
func (f *FeaturePreferences) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode && node.Tag == "!!null" {
		return nil
	}
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("%w: %s must be a mapping (line %d)", ErrMalformedFeaturePreference, featuresKey, node.Line)
	}
	seen := false
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != parallelPlanningKey {
			continue
		}
		if seen {
			return fmt.Errorf("%w: %s.%s appears more than once (line %d)", ErrMalformedFeaturePreference, featuresKey, parallelPlanningKey, node.Content[i].Line)
		}
		seen = true
		value := node.Content[i+1]
		var enabled bool
		if value.Kind != yaml.ScalarNode || value.Tag == "!!null" || value.Decode(&enabled) != nil {
			return fmt.Errorf("%w: %s.%s must be true or false (line %d)", ErrMalformedFeaturePreference, featuresKey, parallelPlanningKey, value.Line)
		}
		f.ParallelPlanning = enabled
	}
	return nil
}

// FeatureSource identifies the config.yml bytes a Config was decoded from. A
// write names it so a file edited since the load is refused rather than
// overwritten.
type FeatureSource struct {
	hash [sha256.Size]byte
}

func newFeatureSource(content []byte) FeatureSource {
	return FeatureSource{hash: sha256.Sum256(content)}
}

// FeatureSource returns the identity of the config.yml bytes c was read from.
// A Config read from a missing file carries the zero value, which no existing
// file matches, so a write against it is refused as a conflict.
func (c *Config) FeatureSource() FeatureSource {
	return c.source
}

// ParallelPlanningEnabled reports the saved preference, false when absent.
func (c *Config) ParallelPlanningEnabled() bool {
	return c.Features.ParallelPlanning
}

// WriteParallelPlanning saves features.parallel_planning in the config.yml at
// path, touching no other byte of the file: comments, key order, indentation,
// and unrelated keys survive. expected must be the source of the Config the
// caller last loaded; a file that differs is refused with ErrMtimeConflict
// (wrapped in ErrV2SourceConflict) before anything is written. Saving the value
// the file already resolves to, including false for an absent key, is a no-op
// that leaves the file untouched. The edited text must decode back to the
// requested value before it replaces the file, and the file is replaced
// through a same-directory temporary file, so a failed write leaves the
// original intact.
func WriteParallelPlanning(path string, enabled bool, expected FeatureSource) error {
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("%w: %s is a symlink", ErrV2UnsafePath, path)
	}
	if info.IsDir() {
		return fmt.Errorf("write %s: target is a directory", path)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if newFeatureSource(raw) != expected {
		return fmt.Errorf("%w: %w: %s", ErrV2SourceConflict, ErrMtimeConflict, path)
	}

	current, err := parseConfig(raw)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if current.Features.ParallelPlanning == enabled {
		return nil
	}

	updated, err := spliceParallelPlanning(string(raw), enabled)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	check, err := parseConfig([]byte(updated))
	if err != nil {
		return fmt.Errorf("%s: edited config no longer decodes: %w", path, err)
	}
	if check.Features.ParallelPlanning != enabled {
		return fmt.Errorf("%w: %s: edit did not produce %s.%s: %t", ErrMalformedFeaturePreference, path, featuresKey, parallelPlanningKey, enabled)
	}

	if info.Mode().Perm()&0222 == 0 {
		return fmt.Errorf("write %s: %w", path, os.ErrPermission)
	}
	return replaceV2File(path, []byte(updated), info.Mode(), func() error {
		latest, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		if newFeatureSource(latest) != expected {
			return fmt.Errorf("%w: %w: %s", ErrV2SourceConflict, ErrMtimeConflict, path)
		}
		return nil
	})
}

// spliceParallelPlanning returns text with features.parallel_planning set to
// enabled by editing only the lines that carry it: the existing value in
// place, a new key under an existing features block, or a new features block
// at the end of the file. Layouts it cannot edit without restyling the file
// are refused rather than rewritten.
func spliceParallelPlanning(text string, enabled bool) (string, error) {
	// Lines keep their own trailing "\r", so each original separator survives
	// and a new line copies the separator of the line it follows.
	fallbackCR := ""
	if strings.Contains(text, "\r\n") {
		fallbackCR = "\r"
	}
	lines := strings.Split(text, "\n")
	value := fmt.Sprintf("%t", enabled)

	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(text), &doc); err != nil {
		return "", fmt.Errorf("failed to parse config YAML: %w", err)
	}
	if !isV2MappingDocument(&doc) {
		if strings.TrimSpace(text) != "" {
			return "", fmt.Errorf("%w: config.yml is not a mapping", ErrV2Malformed)
		}
		return strings.Join(appendFeaturesBlock(nil, 4, value, fallbackCR), "\n"), nil
	}
	root := doc.Content[0]
	indent := v2FrontmatterIndent(&doc)

	featuresIdx := -1
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == featuresKey {
			featuresIdx = i
		}
	}
	if featuresIdx < 0 {
		return strings.Join(appendFeaturesBlock(lines, indent, value, fallbackCR), "\n"), nil
	}

	key, features := root.Content[featuresIdx], root.Content[featuresIdx+1]
	if features.Kind == yaml.ScalarNode && features.Tag == "!!null" && features.Value == "" {
		insertAt := key.Line
		return strings.Join(insertLine(lines, insertAt, strings.Repeat(" ", indent)+parallelPlanningKey+": "+value, fallbackCR), "\n"), nil
	}
	if features.Kind != yaml.MappingNode || features.Style&yaml.FlowStyle != 0 && !hasKey(features, parallelPlanningKey) {
		return "", fmt.Errorf("%w: %s layout cannot be edited in place; set %s.%s by hand", ErrMalformedFeaturePreference, featuresKey, featuresKey, parallelPlanningKey)
	}

	for i := 0; i+1 < len(features.Content); i += 2 {
		if features.Content[i].Value != parallelPlanningKey {
			continue
		}
		scalar := features.Content[i+1]
		span := len(scalar.Value)
		if scalar.Style&(yaml.DoubleQuotedStyle|yaml.SingleQuotedStyle) != 0 {
			span += 2
		}
		line := []rune(lines[scalar.Line-1])
		if scalar.Column-1+span > len(line) {
			return "", fmt.Errorf("%w: %s.%s value cannot be edited in place", ErrMalformedFeaturePreference, featuresKey, parallelPlanningKey)
		}
		lines[scalar.Line-1] = string(line[:scalar.Column-1]) + value + string(line[scalar.Column-1+span:])
		return strings.Join(lines, "\n"), nil
	}

	if len(features.Content) == 0 {
		return "", fmt.Errorf("%w: %s layout cannot be edited in place; set %s.%s by hand", ErrMalformedFeaturePreference, featuresKey, featuresKey, parallelPlanningKey)
	}
	end := len(lines)
	if featuresIdx+2 < len(root.Content) {
		end = root.Content[featuresIdx+2].Line - 1
	}
	last := end
	for last > key.Line && isBlankOrComment(lines[last-1]) {
		last--
	}
	childIndent := features.Content[0].Column - 1
	return strings.Join(insertLine(lines, last, strings.Repeat(" ", childIndent)+parallelPlanningKey+": "+value, fallbackCR), "\n"), nil
}

func hasKey(mapping *yaml.Node, key string) bool {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return true
		}
	}
	return false
}

func isBlankOrComment(line string) bool {
	trimmed := strings.TrimSpace(line)
	return trimmed == "" || strings.HasPrefix(trimmed, "#")
}

// lineCR returns "\r" when the 1-based line at ends with CRLF. The final
// element of a file with no trailing newline has no separator to copy, so it
// falls back to the file's own convention.
func lineCR(lines []string, at int, fallback string) string {
	if at < 1 || at >= len(lines) {
		return fallback
	}
	if strings.HasSuffix(lines[at-1], "\r") {
		return "\r"
	}
	return ""
}

// insertLine returns lines with line placed after the 1-based line number at,
// ending like the line it follows.
func insertLine(lines []string, at int, line, fallbackCR string) []string {
	out := append([]string{}, lines[:at]...)
	out = append(out, line+lineCR(lines, at, fallbackCR))
	return append(out, lines[at:]...)
}

// appendFeaturesBlock adds a features block after the last non-blank line,
// keeping the file's trailing newline.
func appendFeaturesBlock(lines []string, indent int, value, fallbackCR string) []string {
	trailing := len(lines) > 0 && lines[len(lines)-1] == ""
	body := lines
	if trailing {
		body = lines[:len(lines)-1]
	}
	cr := fallbackCR
	if len(body) > 0 && trailing {
		cr = lineCR(lines, len(body), fallbackCR)
	}
	out := append([]string{}, body...)
	if len(out) > 0 && !trailing {
		out[len(out)-1] += cr
	}
	out = append(out, featuresKey+":"+cr, strings.Repeat(" ", indent)+parallelPlanningKey+": "+value+cr, "")
	return out
}
