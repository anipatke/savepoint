package data

import (
	"fmt"
	"regexp"
	"strings"
)

// ReleaseID names the stable identity of a V2 Release. The filesystem slug is
// presentation only; this value is the record's identity and is never
// derived from the slug. It remains a string alias so transitional migration
// readers can consume the pre-E51 Objective projection without conversion.
type ReleaseID = string

var releaseIDPatternV2 = regexp.MustCompile(`^R[0-9]{3,}$`)

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
	Source            V2SourceDocument
}

type releaseV2Frontmatter struct {
	ID     string     `yaml:"id"`
	Title  string     `yaml:"title"`
	Status ColumnType `yaml:"status"`
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

	return &ReleaseV2{
		ID:                fields.ID,
		Title:             fields.Title,
		Status:            fields.Status,
		Outcome:           sections["Outcome"],
		Why:               sections["Why"],
		SuccessConditions: sections["Success Conditions"],
		Boundaries:        sections["Boundaries"],
		Source:            doc,
	}, nil
}

func decodeReleaseBody(path, id, body string) (map[string]string, error) {
	const requiredCount = 4
	required := []string{"Outcome", "Why", "Success Conditions", "Boundaries"}
	sections := make(map[string]string, requiredCount)
	var current string
	var content []string

	flush := func() {
		if current != "" {
			sections[current] = strings.TrimSpace(strings.Join(content, "\n"))
		}
		content = nil
	}

	for _, line := range strings.Split(normalizeLineEndings(body), "\n") {
		trimmed := strings.TrimSpace(line)
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
