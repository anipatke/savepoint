package data

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// ReleaseID names the stable identity of a V2 Release. The filesystem slug is
// presentation only; this value is the record's identity and is never
// derived from the slug. It remains a string alias so transitional migration
// readers can consume the pre-E51 Objective projection without conversion.
type ReleaseID = string

var releaseIDPatternV2 = regexp.MustCompile(`^R[0-9]{3,}$`)
var legacyCompletionHashPatternV2 = regexp.MustCompile(`^[0-9a-fA-F]{64}$`)

// LegacyCompletionReference records why a migrated historical Release may
// retain status: done without pretending that archived evidence is a new V2
// CLEAR Check. Migration owns the referenced bytes and their archive.
type LegacyCompletionReference struct {
	SourcePath  string
	ArchivePath string
	SHA256      string
}

// LegacyCompletion is a short alias for callers that name the frontmatter
// block rather than its reference role.
type LegacyCompletion = LegacyCompletionReference

// ReleaseV2 is a strict V2 Release record. Release membership is deliberately
// absent: it is derived from ObjectiveV2.Release when the project index is
// built, so there is only one authoritative membership edge.
type ReleaseV2 struct {
	ID                ReleaseID
	Title             string
	Status            ColumnType
	Outcome           string
	Why               string
	SuccessConditions string
	Boundaries        string
	Evidence          *Evidence
	LegacyCompletion  *LegacyCompletionReference
	Source            V2SourceDocument
}

type releaseV2Frontmatter struct {
	ID                    string     `yaml:"id"`
	Title                 string     `yaml:"title"`
	Status                ColumnType `yaml:"status"`
	evidenceV2Frontmatter `yaml:",inline"`
	LegacyCompletion      *legacyCompletionV2Frontmatter `yaml:"legacy_completion"`
}

type legacyCompletionV2Frontmatter struct {
	SourcePath  string `yaml:"source_path"`
	ArchivePath string `yaml:"archive_path"`
	SHA256      string `yaml:"sha256"`
	// These aliases let migration adapters consume plain-language manifests;
	// managed writes always emit the canonical names above.
	Reference  string `yaml:"reference"`
	Archive    string `yaml:"archive"`
	SourceHash string `yaml:"source_hash"`
}

// DecodeReleaseV2 strictly decodes a Release's identity and lifecycle, while
// retaining the authored body in Source. The four required body headings are
// projected into fields for consumers, but Source remains the preservation
// boundary for later managed writes.
func DecodeReleaseV2(path, content string) (*ReleaseV2, error) {
	doc, err := ParseV2Document(path, content)
	if err != nil {
		return nil, err
	}

	var fields releaseV2Frontmatter
	if err := doc.Frontmatter.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrV2Malformed, path, err)
	}

	if !releaseIDPatternV2.MatchString(fields.ID) {
		return nil, fmt.Errorf("%w: %s: release id %q must match R plus at least three digits", ErrV2InvalidID, path, fields.ID)
	}
	if strings.TrimSpace(fields.Title) == "" {
		return nil, fmt.Errorf("%w: %s: release %s missing required field title", ErrV2MissingField, path, fields.ID)
	}
	if !IsCanonicalTaskStatus(fields.Status) {
		return nil, fmt.Errorf("%w: %s: release %s status %q; use planned, in_progress, or done", ErrV2InvalidLifecycle, path, fields.ID, fields.Status)
	}

	sections, err := decodeReleaseBody(path, fields.ID, doc.Body)
	if err != nil {
		return nil, err
	}

	evidence, err := decodeEvidenceV2(path, "release", fields.ID, fields.evidenceV2Frontmatter)
	if err != nil {
		return nil, err
	}
	legacyCompletion, err := decodeLegacyCompletionV2(path, fields.ID, fields.LegacyCompletion)
	if err != nil {
		return nil, err
	}
	if legacyCompletion != nil && fields.Status != ColumnDone {
		return nil, fmt.Errorf("%w: %s: release %s legacy_completion is valid only with status done", ErrV2ReleaseLegacyMalformed, path, fields.ID)
	}

	return &ReleaseV2{
		ID:                fields.ID,
		Title:             fields.Title,
		Status:            fields.Status,
		Outcome:           sections["Outcome"],
		Why:               sections["Why"],
		SuccessConditions: sections["Success Conditions"],
		Boundaries:        sections["Boundaries"],
		Evidence:          evidence,
		LegacyCompletion:  legacyCompletion,
		Source:            doc,
	}, nil
}

func decodeLegacyCompletionV2(path, id string, raw *legacyCompletionV2Frontmatter) (*LegacyCompletionReference, error) {
	if raw == nil {
		return nil, nil
	}

	sourcePath := firstNonBlankV2(raw.SourcePath, raw.Reference)
	if strings.TrimSpace(sourcePath) == "" {
		return nil, fmt.Errorf("%w: %s: release %s legacy_completion requires source_path", ErrV2ReleaseLegacyMalformed, path, id)
	}
	archivePath := firstNonBlankV2(raw.ArchivePath, raw.Archive)
	if strings.TrimSpace(archivePath) == "" {
		return nil, fmt.Errorf("%w: %s: release %s legacy_completion requires archive_path", ErrV2ReleaseLegacyMalformed, path, id)
	}
	if !isSafeRelativeV2Path(archivePath) {
		return nil, fmt.Errorf("%w: %s: release %s legacy_completion archive_path %q is not a confined relative path", ErrV2ReleaseLegacyMalformed, path, id, archivePath)
	}

	hash := firstNonBlankV2(raw.SHA256, raw.SourceHash)
	if !legacyCompletionHashPatternV2.MatchString(hash) {
		return nil, fmt.Errorf("%w: %s: release %s legacy_completion sha256 must be 64 hexadecimal characters", ErrV2ReleaseLegacyMalformed, path, id)
	}

	return &LegacyCompletionReference{SourcePath: sourcePath, ArchivePath: archivePath, SHA256: hash}, nil
}

func firstNonBlankV2(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func isSafeRelativeV2Path(path string) bool {
	clean := filepath.Clean(path)
	if clean == "." || filepath.IsAbs(path) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return false
	}
	return path != ""
}

func decodeReleaseBody(path, id, body string) (map[string]string, error) {
	const requiredCount = 4
	required := []string{"Outcome", "Why", "Success Conditions", "Boundaries"}
	sections := make(map[string]string, requiredCount)
	var current string
	var content []string
	fenceTicks := 0

	flush := func() {
		if current != "" {
			sections[current] = strings.TrimSpace(strings.Join(content, "\n"))
		}
		content = nil
	}

	for _, line := range strings.Split(normalizeLineEndings(body), "\n") {
		trimmed := strings.TrimSpace(line)
		if ticks := leadingBackticksV2(trimmed); ticks >= 3 {
			if fenceTicks == 0 {
				fenceTicks = ticks
			} else if ticks >= fenceTicks {
				fenceTicks = 0
			}
			if current != "" {
				content = append(content, line)
			}
			continue
		}
		if fenceTicks > 0 {
			if current != "" {
				content = append(content, line)
			}
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			flush()
			current = strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
			if _, wanted := sections[current]; wanted {
				return nil, fmt.Errorf("%w: %s: release %s body section %q appears more than once", ErrV2ReleaseMissingSection, path, id, current)
			}
			continue
		}
		if current != "" {
			content = append(content, line)
		}
	}
	flush()

	for _, name := range required {
		if _, ok := sections[name]; !ok {
			return nil, fmt.Errorf("%w: %s: release %s body requires heading ## %s", ErrV2ReleaseMissingSection, path, id, name)
		}
	}
	return sections, nil
}

func leadingBackticksV2(line string) int {
	count := 0
	for count < len(line) && line[count] == '`' {
		count++
	}
	return count
}
