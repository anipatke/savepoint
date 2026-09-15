package data

import (
	"fmt"
	"regexp"
	"strings"
)

// objectiveIDPattern anchors V2 global Objective identity: O plus at least
// three digits. Shared with task_v2.go for validating Objective ownership
// and dependency references.
var objectiveIDPattern = regexp.MustCompile(`^O[0-9]{3,}$`)

// ObjectiveV2 is a strict V2 Objective record decoded by DecodeObjectiveV2.
// Unlike V1 Task parsing, nothing here is healed into a valid-looking
// default; every missing or malformed field returns a named diagnostic.
type ObjectiveV2 struct {
	ID        string
	Title     string
	Status    ColumnType
	DependsOn []string // O### references
	Release   string   // optional filtering/packaging metadata only
	Evidence  *Evidence
	Source    V2SourceDocument
}

type objectiveV2Frontmatter struct {
	ID                    string     `yaml:"id"`
	Title                 string     `yaml:"title"`
	Status                ColumnType `yaml:"status"`
	DependsOn             []string   `yaml:"depends_on"`
	Release               string     `yaml:"release"`
	evidenceV2Frontmatter `yaml:",inline"`
}

// DecodeObjectiveV2 strictly decodes a V2 Objective record from content.
// It requires a valid global O### ID, a non-empty title, a canonical
// lifecycle status, and well-formed Objective dependency IDs. Unknown or
// missing values are rejected rather than defaulted into a completion-
// capable state; that healing belongs to V1 only.
func DecodeObjectiveV2(path, content string) (*ObjectiveV2, error) {
	doc, err := ParseV2Document(path, content)
	if err != nil {
		return nil, err
	}

	var fields objectiveV2Frontmatter
	if err := doc.Frontmatter.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrV2Malformed, path, err)
	}

	if !objectiveIDPattern.MatchString(fields.ID) {
		return nil, fmt.Errorf("%w: %s: objective id %q must match O plus at least three digits", ErrV2InvalidID, path, fields.ID)
	}

	if strings.TrimSpace(fields.Title) == "" {
		return nil, fmt.Errorf("%w: %s: objective %s missing required field title", ErrV2MissingField, path, fields.ID)
	}

	if !IsCanonicalTaskStatus(fields.Status) {
		return nil, fmt.Errorf("%w: %s: objective %s status %q; use planned, in_progress, or done", ErrV2InvalidLifecycle, path, fields.ID, fields.Status)
	}

	dependsOn := make([]string, 0, len(fields.DependsOn))
	for _, ref := range fields.DependsOn {
		if !objectiveIDPattern.MatchString(ref) {
			return nil, fmt.Errorf("%w: %s: objective %s depends_on %q must match O plus at least three digits", ErrV2InvalidDependency, path, fields.ID, ref)
		}
		dependsOn = append(dependsOn, ref)
	}

	evidence, err := decodeEvidenceV2(path, "objective", fields.ID, fields.evidenceV2Frontmatter)
	if err != nil {
		return nil, err
	}

	return &ObjectiveV2{
		ID:        fields.ID,
		Title:     fields.Title,
		Status:    fields.Status,
		DependsOn: dependsOn,
		Release:   fields.Release,
		Evidence:  evidence,
		Source:    doc,
	}, nil
}
