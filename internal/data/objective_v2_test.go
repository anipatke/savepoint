package data

import (
	"errors"
	"fmt"
	"testing"
)

func TestDecodeObjectiveV2_valid(t *testing.T) {
	content := `---
id: O-002
title: "Load V2 work with stable identity"
status: in_progress
depends_on: [O-001]
release: R-001
---

# Objective`

	objective, err := DecodeObjectiveV2("O-002-identity/Objective.md", content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if objective.ID != "O-002" {
		t.Errorf("ID = %q, want O-002", objective.ID)
	}
	if objective.Title != "Load V2 work with stable identity" {
		t.Errorf("Title = %q, want the given title", objective.Title)
	}
	if objective.Status != ColumnInProgress {
		t.Errorf("Status = %q, want in_progress", objective.Status)
	}
	if len(objective.DependsOn) != 1 || objective.DependsOn[0] != "O-001" {
		t.Errorf("DependsOn = %v, want [O-001]", objective.DependsOn)
	}
	if objective.Release != "R-001" {
		t.Errorf("Release = %q, want R-001", objective.Release)
	}
}

func TestDecodeObjectiveV2_minimalValid(t *testing.T) {
	content := `---
id: O-010
title: "Bare objective"
status: planned
---

# Objective`

	objective, err := DecodeObjectiveV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if len(objective.DependsOn) != 0 {
		t.Errorf("DependsOn = %v, want empty", objective.DependsOn)
	}
	if objective.Release != "" {
		t.Errorf("Release = %q, want empty", objective.Release)
	}
}

func TestDecodeObjectiveV2_malformedID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"missing digits", "O"},
		{"too few digits", "O01"},
		{"unhyphenated identity", "O001"},
		{"hyphenated but too few digits", "O-01"},
		{"wrong prefix letter", "T-002"},
		{"lowercase prefix", "o002"},
		{"trailing garbage", "O-002x"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: \"" + tt.id + "\"\ntitle: \"Objective\"\nstatus: planned\n---\n\n# Objective"
			_, err := DecodeObjectiveV2("test.md", content)
			if !errors.Is(err, ErrV2InvalidID) {
				t.Fatalf("DecodeObjectiveV2() error = %v, want ErrV2InvalidID", err)
			}
		})
	}
}

func TestDecodeObjectiveV2_missingTitle(t *testing.T) {
	content := `---
id: O-002
status: planned
---

# Objective`

	_, err := DecodeObjectiveV2("test.md", content)
	if !errors.Is(err, ErrV2MissingField) {
		t.Fatalf("DecodeObjectiveV2() error = %v, want ErrV2MissingField", err)
	}
}

func TestDecodeObjectiveV2_whitespaceOnlyTitle(t *testing.T) {
	for _, title := range []string{"   ", "\t\t", "\n\t"} {
		t.Run(fmt.Sprintf("title-%q", title), func(t *testing.T) {
			content := "---\nid: O-002\ntitle: \u0022" + title + "\u0022\nstatus: planned\n---\n\n# Objective"
			_, err := DecodeObjectiveV2("test.md", content)
			if !errors.Is(err, ErrV2MissingField) {
				t.Fatalf("DecodeObjectiveV2() error = %v, want ErrV2MissingField", err)
			}
		})
	}
}

func TestDecodeObjectiveV2_unknownStatusNotHealed(t *testing.T) {
	tests := []struct {
		name   string
		status string
	}{
		{"garbage value", "review"},
		{"missing status", ""},
		{"legacy todo alias", "todo"},
		{"legacy complete alias", "complete"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: O-002\ntitle: \"Objective\"\nstatus: \"" + tt.status + "\"\n---\n\n# Objective"
			_, err := DecodeObjectiveV2("test.md", content)
			if !errors.Is(err, ErrV2InvalidLifecycle) {
				t.Fatalf("DecodeObjectiveV2() error = %v, want ErrV2InvalidLifecycle (no healing to planned)", err)
			}
		})
	}
}

func TestDecodeObjectiveV2_invalidDependency(t *testing.T) {
	content := `---
id: O-002
title: "Objective"
status: planned
depends_on: [O1]
---

# Objective`

	_, err := DecodeObjectiveV2("test.md", content)
	if !errors.Is(err, ErrV2InvalidDependency) {
		t.Fatalf("DecodeObjectiveV2() error = %v, want ErrV2InvalidDependency", err)
	}
}

func TestDecodeObjectiveV2_malformedYAML(t *testing.T) {
	content := `---
id: [broken
---

# Objective`

	_, err := DecodeObjectiveV2("test.md", content)
	if err == nil {
		t.Fatal("DecodeObjectiveV2() expected error for malformed YAML")
	}
}

func TestDecodeObjectiveV2_noFrontmatter(t *testing.T) {
	_, err := DecodeObjectiveV2("test.md", "# No frontmatter here")
	if !errors.Is(err, ErrNoFrontmatter) {
		t.Fatalf("DecodeObjectiveV2() error = %v, want ErrNoFrontmatter", err)
	}
}

// TestDecodeObjectiveV2_evidenceValid proves the shared evidence block
// decodes on Objective records through the same decoder as Task records.
func TestDecodeObjectiveV2_evidenceValid(t *testing.T) {
	content := `---
id: O-002
title: "Objective"
status: planned
exception:
  requirements: [TEST-08]
  reason: "owner accepted the tradeoff"
  owner: ani
  recorded_at: "2026-09-15T00:00:00Z"
  check: C-003
---

# Objective`

	objective, err := DecodeObjectiveV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if objective.Evidence == nil || objective.Evidence.Exception == nil {
		t.Fatal("DecodeObjectiveV2() Evidence.Exception = nil, want decoded exception")
	}
	if objective.Evidence.Exception.Check != "C-003" {
		t.Errorf("Evidence.Exception.Check = %q, want C-003", objective.Evidence.Exception.Check)
	}
}

// TestDecodeObjectiveV2_noEvidenceIsNil proves an Objective carrying none of
// the evidence fields decodes with a nil Evidence.
func TestDecodeObjectiveV2_noEvidenceIsNil(t *testing.T) {
	content := `---
id: O-002
title: "Objective"
status: planned
---

# Objective`

	objective, err := DecodeObjectiveV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if objective.Evidence != nil {
		t.Errorf("Evidence = %+v, want nil", objective.Evidence)
	}
}

func TestDecodeObjectiveV2_releaseIsTypedReference(t *testing.T) {
	base := `---
id: O-002
title: "Objective"
status: planned
release: %s
---

# Objective`

	withRelease, err := DecodeObjectiveV2("test.md", fmt.Sprintf(base, "R-001"))
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	withoutRelease, err := DecodeObjectiveV2("test.md", `---
id: O-002
title: "Objective"
status: planned
---

# Objective`)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}

	// Release is an optional typed reference and does not replace Objective
	// identity or ownership.
	if withRelease.ID != withoutRelease.ID || withRelease.Title != withoutRelease.Title {
		t.Fatalf("release value changed identity fields: %+v vs %+v", withRelease, withoutRelease)
	}
	if withRelease.Release != "R-001" || withoutRelease.Release != "" {
		t.Fatalf("Release fields = %q / %q, want R-001 / empty", withRelease.Release, withoutRelease.Release)
	}
}

func TestDecodeObjectiveV2_preservesTransitionalPackagingText(t *testing.T) {
	content := `---
id: O-002
title: "Objective"
status: planned
release: v2
---

# Objective`

	objective, err := DecodeObjectiveV2("objectives/O-002-objective/Objective.md", content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v, want transitional compatibility", err)
	}
	if objective.Release != "v2" {
		t.Errorf("Release = %q, want legacy v2 label retained for migration compatibility", objective.Release)
	}
}
