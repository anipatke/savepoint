package data

import (
	"errors"
	"testing"
)

func TestDecodeTaskV2_valid(t *testing.T) {
	content := `---
id: T005
title: "Show clear project errors"
objective: O002
planned_by: {role: planner, session: planning-001}
status: in_progress
stage: test
depends_on:
  - task: T003
  - task: T004
    requires: accepted
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.ID != "T005" {
		t.Errorf("ID = %q, want T005", task.ID)
	}
	if task.Title != "Show clear project errors" {
		t.Errorf("Title = %q, want the given title", task.Title)
	}
	if task.Objective != "O002" {
		t.Errorf("Objective = %q, want O002", task.Objective)
	}
	if task.PlannedBy != (Actor{Role: ActorRolePlanner, Session: "planning-001"}) {
		t.Errorf("PlannedBy = %+v, want planner/planning-001", task.PlannedBy)
	}
	if task.Status != ColumnInProgress || task.Stage != StageTest {
		t.Errorf("Status/Stage = %q/%q, want in_progress/test", task.Status, task.Stage)
	}
	if len(task.DependsOn) != 2 {
		t.Fatalf("DependsOn len = %d, want 2", len(task.DependsOn))
	}
	if task.DependsOn[0].Task != "T003" || task.DependsOn[0].Requires != TaskDependencyClear {
		t.Errorf("DependsOn[0] = %+v, want T003/clear (omitted defaults to clear)", task.DependsOn[0])
	}
	if task.DependsOn[1].Task != "T004" || task.DependsOn[1].Requires != TaskDependencyAccepted {
		t.Errorf("DependsOn[1] = %+v, want T004/accepted", task.DependsOn[1])
	}
}

func TestDecodeTaskV2_minimalValidPlanned(t *testing.T) {
	content := `---
id: T001
title: "Bare task"
objective: O001
planned_by: {role: planner, session: planning-001}
status: planned
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Stage != "" {
		t.Errorf("Stage = %q, want empty for planned", task.Stage)
	}
	if len(task.DependsOn) != 0 {
		t.Errorf("DependsOn = %v, want empty", task.DependsOn)
	}
}

func TestDecodeTaskV2_plannedByValidation(t *testing.T) {
	tests := []struct {
		name        string
		plannedBy   string
		wantErr     error
		wantMessage string
	}{
		{name: "missing", plannedBy: "", wantErr: ErrV2MissingField},
		{name: "malformed shape", plannedBy: "planner-001", wantErr: ErrV2Malformed},
		{name: "missing role", plannedBy: "{session: planning-001}", wantErr: ErrV2MissingField},
		{name: "missing session", plannedBy: "{role: planner}", wantErr: ErrV2MissingField},
		{name: "blank session", plannedBy: "{role: planner, session: '   '}", wantErr: ErrV2MissingField},
		{name: "wrong role", plannedBy: "{role: executor, session: execution-001}", wantErr: ErrV2TaskMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plannedBy := ""
			if tt.plannedBy != "" {
				plannedBy = "planned_by: " + tt.plannedBy + "\n"
			}
			content := "---\nid: T005\ntitle: \"Task\"\nobjective: O002\n" + plannedBy + "status: planned\n---\n\n# Task"
			_, err := DecodeTaskV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeTaskV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeTaskV2_malformedID(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{"missing digits", "T"},
		{"too few digits", "T01"},
		{"wrong prefix letter", "O001"},
		{"lowercase prefix", "t001"},
		{"empty", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: \"" + tt.id + "\"\ntitle: \"Task\"\nobjective: O001\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\n---\n\n# Task"
			_, err := DecodeTaskV2("test.md", content)
			if !errors.Is(err, ErrV2InvalidID) {
				t.Fatalf("DecodeTaskV2() error = %v, want ErrV2InvalidID", err)
			}
		})
	}
}

// TestDecodeTaskV2_missingTitleIsNotBackfilledFromObjective proves V2 never
// repeats V1 parser.go's firstNonEmpty(Title, Objective) fallback: a missing
// title fails even though a valid O### objective owner is present.
func TestDecodeTaskV2_missingTitleIsNotBackfilledFromObjective(t *testing.T) {
	content := `---
id: T005
objective: O002
status: planned
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if !errors.Is(err, ErrV2MissingField) {
		t.Fatalf("DecodeTaskV2() error = %v, want ErrV2MissingField", err)
	}
	if task != nil {
		t.Fatalf("DecodeTaskV2() task = %+v, want nil on missing title", task)
	}
}

func TestDecodeTaskV2_whitespaceOnlyTitle(t *testing.T) {
	for _, title := range []string{"   ", "\t\t", "\n\t"} {
		t.Run("whitespace", func(t *testing.T) {
			content := "---\nid: T005\ntitle: \u0022" + title + "\u0022\nobjective: O002\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\n---\n\n# Task"
			_, err := DecodeTaskV2("test.md", content)
			if !errors.Is(err, ErrV2MissingField) {
				t.Fatalf("DecodeTaskV2() error = %v, want ErrV2MissingField", err)
			}
		})
	}
}

func TestDecodeTaskV2_missingOrMalformedObjective(t *testing.T) {
	tests := []struct {
		name      string
		objective string
	}{
		{"missing", ""},
		{"malformed id", "O1"},
		{"wrong prefix", "T002"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: T005\ntitle: \"Task\"\nobjective: \"" + tt.objective + "\"\nplanned_by: {role: planner, session: planning-001}\nstatus: planned\n---\n\n# Task"
			_, err := DecodeTaskV2("test.md", content)
			if !errors.Is(err, ErrV2InvalidOwnership) {
				t.Fatalf("DecodeTaskV2() error = %v, want ErrV2InvalidOwnership", err)
			}
		})
	}
}

// TestDecodeTaskV2_rejectsMultipleObjectiveOwners proves the schema enforces
// "exactly one objective: O### owner" structurally: a sequence where a
// scalar owner is expected fails to decode rather than picking one.
func TestDecodeTaskV2_rejectsMultipleObjectiveOwners(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: [O001, O002]
status: planned
---

# Task`

	_, err := DecodeTaskV2("test.md", content)
	if !errors.Is(err, ErrV2Malformed) {
		t.Fatalf("DecodeTaskV2() error = %v, want ErrV2Malformed for a non-scalar objective", err)
	}
}

func TestDecodeTaskV2_lifecycleNotHealed(t *testing.T) {
	tests := []struct {
		name    string
		fields  string
		wantErr error
	}{
		{
			name:    "missing status",
			fields:  "",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "legacy todo status is rejected, not healed to planned",
			fields:  "status: todo\n",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "legacy complete status is rejected, not healed to done",
			fields:  "status: complete\n",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "in_progress missing stage is rejected, not defaulted to build",
			fields:  "status: in_progress\n",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "in_progress with unknown stage is rejected, not healed to build",
			fields:  "status: in_progress\nstage: review\n",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "legacy implementation stage is rejected, not healed to build",
			fields:  "status: in_progress\nstage: implementation\n",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			name:    "stale stage outside in_progress is rejected, not silently cleared",
			fields:  "status: done\nstage: build\n",
			wantErr: ErrV2InvalidLifecycle,
		},
		{
			// V1's phase→stage compatibility alias must not leak into V2:
			// phase is not a recognized taskV2Frontmatter field, so it is
			// silently ignored by decode and stage is still required.
			name:    "legacy phase field is not treated as a stage alias",
			fields:  "status: in_progress\nphase: test\n",
			wantErr: ErrV2InvalidLifecycle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\nid: T005\ntitle: \"Task\"\nobjective: O002\nplanned_by: {role: planner, session: planning-001}\n" + tt.fields + "---\n\n# Task"
			_, err := DecodeTaskV2("test.md", content)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("DecodeTaskV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeTaskV2_dependencyRequiresDefaultsOnlyWhenOmitted(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
depends_on:
  - task: T003
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.DependsOn[0].Requires != TaskDependencyClear {
		t.Fatalf("Requires = %q, want clear default for omitted requirement", task.DependsOn[0].Requires)
	}
}

func TestDecodeTaskV2_dependencyRequiresRejectsUnknownValues(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
depends_on:
  - task: T003
    requires: maybe
---

# Task`

	_, err := DecodeTaskV2("test.md", content)
	if !errors.Is(err, ErrV2InvalidDependency) {
		t.Fatalf("DecodeTaskV2() error = %v, want ErrV2InvalidDependency", err)
	}
}

func TestDecodeTaskV2_dependencyMalformedTaskID(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
depends_on:
  - task: T1
---

# Task`

	_, err := DecodeTaskV2("test.md", content)
	if !errors.Is(err, ErrV2InvalidDependency) {
		t.Fatalf("DecodeTaskV2() error = %v, want ErrV2InvalidDependency", err)
	}
}

func TestDecodeTaskV2_malformedYAML(t *testing.T) {
	content := `---
id: [broken
---

# Task`

	_, err := DecodeTaskV2("test.md", content)
	if err == nil {
		t.Fatal("DecodeTaskV2() expected error for malformed YAML")
	}
}

func TestDecodeTaskV2_noFrontmatter(t *testing.T) {
	_, err := DecodeTaskV2("test.md", "# No frontmatter here")
	if !errors.Is(err, ErrNoFrontmatter) {
		t.Fatalf("DecodeTaskV2() error = %v, want ErrNoFrontmatter", err)
	}
}

func TestDecodeTaskV2_ownershipComesOnlyFromObjective(t *testing.T) {
	// A task has no Release field of its own. Its only ownership edge remains
	// the explicit Objective reference.
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Objective != "O002" {
		t.Errorf("Objective = %q, want O002", task.Objective)
	}
}

// TestDecodeTaskV2_typeIsolatedFromV1Task proves TaskV2 is a distinct type
// from Task (task.go), not an extension of the active V1 Task contract.
func TestDecodeTaskV2_typeIsolatedFromV1Task(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
---

# Task`

	taskV2, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	var _ *TaskV2 = taskV2
	var _ *Task = &Task{} // V1 Task remains its own, unrelated type.
}

// TestDecodeTaskV2_evidenceValid proves the shared evidence block decodes
// from a real Task document alongside its existing fields.
func TestDecodeTaskV2_evidenceValid(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
last_check: C001
freshness:
  state: current
  check: C001
  assessed_by: {role: checker, session: sess-1}
  assessed_at: "2026-09-15T00:00:00Z"
  basis: reviewed the diff
owner_validation:
  required: true
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Evidence == nil {
		t.Fatal("DecodeTaskV2() Evidence = nil, want decoded evidence")
	}
	if task.Evidence.LastCheck != "C001" {
		t.Errorf("Evidence.LastCheck = %q, want C001", task.Evidence.LastCheck)
	}
	if task.Evidence.Freshness == nil || task.Evidence.Freshness.State != FreshnessCurrent {
		t.Errorf("Evidence.Freshness = %+v, want state current", task.Evidence.Freshness)
	}
	if task.Evidence.OwnerValidation == nil || !task.Evidence.OwnerValidation.Required {
		t.Errorf("Evidence.OwnerValidation = %+v, want Required true", task.Evidence.OwnerValidation)
	}
}

// TestDecodeTaskV2_noEvidenceIsNil proves a Task carrying none of the
// evidence fields decodes with a nil Evidence rather than a defaulted one.
func TestDecodeTaskV2_noEvidenceIsNil(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
---

# Task`

	task, err := DecodeTaskV2("test.md", content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Evidence != nil {
		t.Errorf("Evidence = %+v, want nil", task.Evidence)
	}
}

// TestDecodeTaskV2_malformedEvidencePropagatesDiagnostic proves a malformed
// evidence sub-block fails DecodeTaskV2 through the shared decoder, not just
// decodeEvidenceV2 in isolation.
func TestDecodeTaskV2_malformedEvidencePropagatesDiagnostic(t *testing.T) {
	content := `---
id: T005
title: "Task"
objective: O002
planned_by: {role: planner, session: planning-001}
status: planned
freshness:
  state: expired
  check: C001
  assessed_by: {role: checker, session: sess-1}
  assessed_at: "2026-09-15T00:00:00Z"
  basis: reviewed the diff
---

# Task`

	_, err := DecodeTaskV2("test.md", content)
	if !errors.Is(err, ErrV2EvidenceMalformed) {
		t.Fatalf("DecodeTaskV2() error = %v, want ErrV2EvidenceMalformed", err)
	}
}

// TestParseTaskFile_v1BehaviorUnaffectedByV2Types is a regression check that
// V1 Task parsing, including its legacy status healing, is unchanged now
// that strict V2 decoders exist in the same package.
func TestParseTaskFile_v1BehaviorUnaffectedByV2Types(t *testing.T) {
	p := NewParser()
	content := `---
id: E06/T001
status: todo
objective: "Style the board"
---

# Task`

	task, err := p.ParseTaskFile("test.md", content)
	if err != nil {
		t.Fatalf("ParseTaskFile() error = %v", err)
	}
	if task.Column != ColumnPlanned {
		t.Fatalf("Task.Column = %v, want %v (legacy todo alias still heals in V1)", task.Column, ColumnPlanned)
	}
	if task.Title != "Style the board" {
		t.Fatalf("Task.Title = %q, want V1's objective-as-title fallback to still apply", task.Title)
	}
}
