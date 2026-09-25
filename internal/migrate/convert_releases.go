package migrate

// convert_releases.go renders one V1 release PRD into the first-class V2
// Release record planned by plan.go. The authored promise stays visible in a
// dedicated fenced legacy-source section, while the exact source bytes are
// independently preserved by the ArchiveEntry paired with the target.

import (
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
	"gopkg.in/yaml.v3"
)

const continuationGoalTitle = "Continued work after migration"

type releaseOutputFrontmatter struct {
	ID               string                  `yaml:"id"`
	Title            string                  `yaml:"title"`
	Status           string                  `yaml:"status"`
	LegacyCompletion *legacyCompletionOutput `yaml:"legacy_completion,omitempty"`
	LegacyFields     map[string]any          `yaml:"legacy_fields,omitempty"`
}

type legacyCompletionOutput struct {
	SourcePath  string `yaml:"source_path"`
	ArchivePath string `yaml:"archive_path"`
	SHA256      string `yaml:"sha256"`
}

// ConvertRelease renders and validates a V2 Release without writing it. The
// source PRD's body is included verbatim inside a fence so authored headings
// remain readable and cannot be mistaken for the V2 record's required
// sections. A settled V1 status receives only a typed archive reference; no
// V2 Check, freshness, or owner acceptance is invented.
func ConvertRelease(root string, plan *ConversionPlan, target PlannedTarget) (string, error) {
	if target.Kind != TargetRelease {
		return "", fmt.Errorf("convert release: target %s is not a release target", target.GlobalID)
	}
	if target.ReleaseStatus == "" {
		return "", fmt.Errorf("convert release %s: missing planned release status", target.GlobalID)
	}
	if target.Generated {
		return convertGeneratedGoal(target)
	}

	raw, err := readSourceFile(root, target.Legacy.Path)
	if err != nil {
		return "", err
	}
	crlf := strings.Contains(raw, "\r\n")
	fm, body, err := data.SplitFrontmatterBody(raw)
	if err != nil {
		return "", fmt.Errorf("convert release %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	var source releaseRawFrontmatter
	if err := yaml.Unmarshal([]byte(fm), &source); err != nil {
		return "", fmt.Errorf("convert release %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}
	legacyFields, err := unknownFields(fm, []string{"version", "name", "title", "status"})
	if err != nil {
		return "", fmt.Errorf("convert release %s: %s: %w", target.GlobalID, target.Legacy.Path, err)
	}

	title := firstNonEmptyString(source.Title, source.Name)
	if title == "" {
		if h1, ok := extractH1(body); ok {
			title = h1
		}
	}
	if title == "" {
		title = target.Legacy.Release
	}

	var legacy *legacyCompletionOutput
	if target.ReleaseStatus == string(data.ColumnDone) {
		legacy = &legacyCompletionOutput{
			SourcePath:  target.Legacy.Path,
			ArchivePath: archivePathFor(target.Legacy.Path),
			SHA256:      hashBytes([]byte(raw)),
		}
	}

	out := releaseOutputFrontmatter{
		ID:               target.GlobalID,
		Title:            title,
		Status:           target.ReleaseStatus,
		LegacyCompletion: legacy,
		LegacyFields:     legacyFields,
	}
	marshalled, err := yaml.Marshal(&out)
	if err != nil {
		return "", fmt.Errorf("convert release %s: marshal yaml: %w", target.GlobalID, err)
	}

	fence := releaseSourceFence(body)
	content := "---\n" + strings.TrimSpace(string(marshalled)) + "\n---\n"
	content += "## Outcome\n\n"
	content += fmt.Sprintf("The V1 release promise is preserved in the Legacy Source section below. Source: `%s`.\n\n", target.Legacy.Path)
	content += "## Why\n\n"
	content += "This Release carries the V1 delivery boundary forward with a stable V2 identity.\n\n"
	content += "## Success Conditions\n\n"
	content += "- Converted Objectives sourced from this V1 release reference this Release identity.\n"
	content += "- Historical completion, when present, remains typed archive evidence rather than a V2 Check.\n\n"
	content += "## Boundaries\n\n"
	content += "Release membership is derived from Objective records; Tasks remain owned by Objectives.\n\n"
	content += "## Legacy Source (verbatim)\n\n"
	content += fmt.Sprintf("Source path: `%s`\n\nSHA-256: `%s`\n\nOriginal frontmatter:\n\n", target.Legacy.Path, hashBytes([]byte(raw)))
	content += "```yaml\n" + fm + "\n```\n\n"
	content += fence + "markdown\n" + body
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += fence + "\n"

	if crlf {
		content = strings.ReplaceAll(content, "\n", "\r\n")
	}
	if _, err := data.DecodeReleaseV2(target.TargetPath, content); err != nil {
		return "", fmt.Errorf("convert release %s: rendered content did not decode: %w", target.GlobalID, err)
	}
	return content, nil
}

func convertGeneratedGoal(target PlannedTarget) (string, error) {
	title := target.GeneratedTitle
	if title == "" {
		title = continuationGoalTitle
	}
	frontmatter, err := yaml.Marshal(&releaseOutputFrontmatter{
		ID:     target.GlobalID,
		Title:  title,
		Status: target.ReleaseStatus,
	})
	if err != nil {
		return "", fmt.Errorf("convert generated Goal %s: marshal yaml: %w", target.GlobalID, err)
	}
	content := "---\n" + strings.TrimSpace(string(frontmatter)) + "\n---\n"
	content += "## Outcome\n\nTODO: Define the outcome this Goal will deliver after migration.\n\n"
	content += "## Why\n\nTODO: Record why this continued work matters.\n\n"
	content += "## Success Conditions\n\n- TODO: Define observable conditions for completion.\n\n"
	content += "## Boundaries\n\n- TODO: Define scope limits and exclusions.\n"
	if _, err := data.DecodeReleaseV2(target.TargetPath, content); err != nil {
		return "", fmt.Errorf("convert generated Goal %s: rendered content did not decode: %w", target.GlobalID, err)
	}
	return content, nil
}

func releaseSourceFence(body string) string {
	max := 3
	for _, line := range strings.Split(body, "\n") {
		ticks := 0
		for ticks < len(line) && line[ticks] == '`' {
			ticks++
		}
		if ticks > max {
			max = ticks
		}
	}
	return strings.Repeat("`", max+1)
}
