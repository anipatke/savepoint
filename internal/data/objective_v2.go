package data

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// ObjectivePriority is an Objective's planning priority within its Goal.
// Priority changes ordering only; it does not affect lifecycle or gates.
type ObjectivePriority string

const (
	ObjectivePriorityCritical ObjectivePriority = "critical"
	ObjectivePriorityHigh     ObjectivePriority = "high"
	ObjectivePriorityMedium   ObjectivePriority = "medium"
	ObjectivePriorityLow      ObjectivePriority = "low"
)

func isCanonicalObjectivePriority(priority ObjectivePriority) bool {
	switch priority {
	case ObjectivePriorityCritical, ObjectivePriorityHigh, ObjectivePriorityMedium, ObjectivePriorityLow:
		return true
	default:
		return false
	}
}

func objectivePriorityOrder(priority ObjectivePriority) int {
	switch priority {
	case ObjectivePriorityCritical:
		return 0
	case ObjectivePriorityHigh:
		return 1
	case ObjectivePriorityMedium, "":
		return 2
	case ObjectivePriorityLow:
		return 3
	default:
		return 4
	}
}

// ObjectiveV2 is a strict V2 Objective record decoded by DecodeObjectiveV2.
// Required identity and lifecycle fields remain strict; optional planning
// priority and rank use their documented defaults when omitted.
type ObjectiveV2 struct {
	ID        string
	Title     string
	Status    ColumnType
	Priority  ObjectivePriority
	Rank      int       // zero means unranked
	DependsOn []string  // O-### references
	Release   ReleaseID // optional R-### or G-### Goal reference
	Evidence  *Evidence
	Source    V2SourceDocument
}

type objectiveV2Frontmatter struct {
	ID                    string     `yaml:"id"`
	Title                 string     `yaml:"title"`
	Status                ColumnType `yaml:"status"`
	Priority              yaml.Node  `yaml:"priority"`
	Rank                  yaml.Node  `yaml:"rank"`
	DependsOn             []string   `yaml:"depends_on"`
	Release               string     `yaml:"release"`
	evidenceV2Frontmatter `yaml:",inline"`
}

// DecodeObjectiveV2 strictly decodes a V2 Objective record from content.
// It requires a valid global O-### ID, a non-empty title, a canonical
// lifecycle status, and well-formed Objective dependency IDs. Missing
// priority defaults to medium and missing rank to unranked; explicitly
// malformed ordering metadata is rejected.
func DecodeObjectiveV2(path, content string) (*ObjectiveV2, error) {
	doc, err := ParseV2Document(path, content)
	if err != nil {
		return nil, err
	}

	var fields objectiveV2Frontmatter
	if err := doc.Frontmatter.Decode(&fields); err != nil {
		return nil, fmt.Errorf("%w: %s: %v", ErrV2Malformed, path, err)
	}

	if !matchesV2Identity(fields.ID, 'O') {
		return nil, fmt.Errorf("%w: %s: objective id %q must match O- plus at least three digits", ErrV2InvalidID, path, fields.ID)
	}

	if strings.TrimSpace(fields.Title) == "" {
		return nil, fmt.Errorf("%w: %s: objective %s missing required field title", ErrV2MissingField, path, fields.ID)
	}

	if !IsCanonicalTaskStatus(fields.Status) {
		return nil, fmt.Errorf("%w: %s: objective %s status %q; use planned, in_progress, or done", ErrV2InvalidLifecycle, path, fields.ID, fields.Status)
	}
	priority, err := decodeObjectivePriority(path, fields.ID, fields.Priority)
	if err != nil {
		return nil, err
	}
	rank, err := decodeObjectiveRank(path, fields.ID, fields.Rank)
	if err != nil {
		return nil, err
	}
	dependsOn := make([]string, 0, len(fields.DependsOn))
	for _, ref := range fields.DependsOn {
		if !matchesV2Identity(ref, 'O') {
			return nil, fmt.Errorf("%w: %s: objective %s depends_on %q must match O- plus at least three digits", ErrV2InvalidDependency, path, fields.ID, ref)
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
		Priority:  priority,
		Rank:      rank,
		DependsOn: dependsOn,
		Release:   ReleaseID(fields.Release),
		Evidence:  evidence,
		Source:    doc,
	}, nil
}

func decodeObjectivePriority(path, objectiveID string, node yaml.Node) (ObjectivePriority, error) {
	if node.Kind == 0 {
		return ObjectivePriorityMedium, nil
	}
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		return "", fmt.Errorf("%w: %s: objective %s priority must be one of critical, high, medium, or low", ErrV2Malformed, path, objectiveID)
	}
	priority := ObjectivePriority(node.Value)
	if !isCanonicalObjectivePriority(priority) {
		return "", fmt.Errorf("%w: %s: objective %s priority %q; use critical, high, medium, or low", ErrV2Malformed, path, objectiveID, node.Value)
	}
	return priority, nil
}

func decodeObjectiveRank(path, objectiveID string, node yaml.Node) (int, error) {
	if node.Kind == 0 {
		return 0, nil
	}
	if node.Kind != yaml.ScalarNode || node.Tag != "!!int" {
		return 0, fmt.Errorf("%w: %s: objective %s rank must be a positive integer", ErrV2Malformed, path, objectiveID)
	}
	var rank int
	if err := node.Decode(&rank); err != nil || rank <= 0 {
		return 0, fmt.Errorf("%w: %s: objective %s rank must be a positive integer", ErrV2Malformed, path, objectiveID)
	}
	return rank, nil
}
