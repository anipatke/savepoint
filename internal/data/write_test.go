package data

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestWriteObjectiveV2_updatesStatusPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-002
title: "Load V2 work with stable identity"
status: planned
depends_on: [O-001]
release: R-001
owner:
  team: platform
  contact: "team@example.com"
---

# Objective

## Notes

Authored planning notes that must survive the rewrite.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnInProgress

	if err := WriteObjectiveV2(objective); err != nil {
		t.Fatalf("WriteObjectiveV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeObjectiveV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnInProgress {
		t.Errorf("Status = %q, want in_progress", reparsed.Status)
	}
	if len(reparsed.DependsOn) != 1 || reparsed.DependsOn[0] != "O-001" {
		t.Errorf("DependsOn = %v, want [O-001] preserved", reparsed.DependsOn)
	}
	if reparsed.Release != "R-001" {
		t.Errorf("Release = %q, want R-001 preserved", reparsed.Release)
	}
	if !strings.Contains(string(result), "team: platform") {
		t.Error("unknown nested field not preserved")
	}
	if !strings.Contains(string(result), "team@example.com") {
		t.Error("unknown nested field value not preserved")
	}
	if !strings.Contains(string(result), "Authored planning notes that must survive the rewrite.") {
		t.Error("authored body content not preserved")
	}
	if !strings.Contains(string(result), "## Notes") {
		t.Error("authored heading not preserved")
	}
}

func TestWriteObjectiveV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-003
title: "No-op objective"
status: in_progress
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}

	if err := WriteObjectiveV2(objective); err != nil {
		t.Fatalf("WriteObjectiveV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteObjectiveV2_refusesUnsupportedStatusAndLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-004
title: "Guarded objective"
status: planned
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnType("archived")

	err = WriteObjectiveV2(objective)
	if err == nil {
		t.Fatal("WriteObjectiveV2() expected error for unsupported status")
	}
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("WriteObjectiveV2() error = %v, want ErrV2InvalidLifecycle", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %v is not path-qualified with %q", err, path)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("file content changed despite validation refusal")
	}
}

func TestWriteTaskV2_updatesStatusAndStagePreservesDependencies(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-005.md")
	content := `---
id: T-005
title: "Show clear project errors"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
depends_on:
  - task: T-003
  - task: T-004
    requires: accepted
release: R-001
---

# Task

Authored task notes.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnInProgress || reparsed.Stage != StageBuild {
		t.Errorf("Status/Stage = %q/%q, want in_progress/build", reparsed.Status, reparsed.Stage)
	}
	if len(reparsed.DependsOn) != 2 {
		t.Fatalf("DependsOn len = %d, want 2 preserved", len(reparsed.DependsOn))
	}
	if reparsed.DependsOn[1].Task != "T-004" || reparsed.DependsOn[1].Requires != TaskDependencyAccepted {
		t.Errorf("DependsOn[1] = %+v, want T-004/accepted preserved", reparsed.DependsOn[1])
	}
	if !strings.Contains(string(result), "Authored task notes.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteTaskV2_removesStageWhenLeavingInProgress(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-006.md")
	content := `---
id: T-006
title: "Finish up"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: audit
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnDone
	task.Stage = ""

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "stage:") {
		t.Error("stage field should be removed when leaving in_progress")
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnDone || reparsed.Stage != "" {
		t.Errorf("Status/Stage = %q/%q, want done/empty", reparsed.Status, reparsed.Stage)
	}
}

func TestWriteTaskV2_refusesInProgressWithoutStageAndLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-007.md")
	content := `---
id: T-007
title: "Guarded task"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = ""

	err = WriteTaskV2(task)
	if err == nil {
		t.Fatal("WriteTaskV2() expected error for in_progress without stage")
	}
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("WriteTaskV2() error = %v, want ErrV2InvalidLifecycle", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("file content changed despite validation refusal")
	}
}

func TestWriteTaskV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-008.md")
	content := `---
id: T-008
title: "No-op task"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: test
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteTaskV2_preservesCRLFLineEndings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-009.md")
	content := "---\r\nid: T-009\r\ntitle: \"CRLF task\"\r\nobjective: O-002\r\nplanned_by: {role: planner, session: planning-fixture}\r\nstatus: planned\r\n---\r\n\r\n# Task\r\n\r\nAuthored notes.\r\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if strings.Count(string(result), "\n") != strings.Count(string(result), "\r\n") {
		t.Errorf("result did not preserve CRLF line endings throughout: %q", string(result))
	}
	if !strings.Contains(string(result), "Authored notes.") {
		t.Error("authored body content not preserved across CRLF rewrite")
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Status != ColumnInProgress || reparsed.Stage != StageBuild {
		t.Errorf("Status/Stage = %q/%q, want in_progress/build", reparsed.Status, reparsed.Stage)
	}
}

func TestWriteTaskV2_preservesLFLineEndingsByDefault(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-010.md")
	content := "---\nid: T-010\ntitle: \"LF task\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\n---\n\n# Task\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild

	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "\r\n") {
		t.Errorf("LF source should not gain CRLF endings: %q", string(result))
	}
}

func TestWriteObjectiveV2_writeFailureIsPathQualifiedAndLeavesRecordIntact(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission checks do not apply")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "Objective.md")
	content := `---
id: O-005
title: "Read-only objective"
status: planned
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	objective.Status = ColumnInProgress

	if err := os.Chmod(path, 0444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(path, 0644)

	err = WriteObjectiveV2(objective)
	if err == nil {
		t.Fatal("WriteObjectiveV2() expected error writing to a read-only file")
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("error %v is not path-qualified with %q", err, path)
	}

	if err := os.Chmod(path, 0644); err != nil {
		t.Fatal(err)
	}
	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("existing record was truncated or altered despite the write failure")
	}
}

func TestWriteV2Record_resolvesDiscoveredRelativePathFromProjectRoot(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-first.md", "T-001", "First task", "O-001")

	objectives, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v", err)
	}
	objective := objectives["O-001"]
	task := tasks["T-001"]
	if objective == nil || task == nil {
		t.Fatal("DiscoverV2Records() did not return both V2 records")
	}
	if filepath.IsAbs(objective.Source.Path) || filepath.IsAbs(task.Source.Path) {
		t.Fatalf("source paths must remain project-relative: objective=%q task=%q", objective.Source.Path, task.Source.Path)
	}
	if objective.Source.ProjectRoot != root || task.Source.ProjectRoot != root {
		t.Fatalf("source project roots = %q/%q, want %q", objective.Source.ProjectRoot, task.Source.ProjectRoot, root)
	}

	otherWorkingDir := t.TempDir()
	t.Chdir(otherWorkingDir)

	objective.Status = ColumnInProgress
	if err := WriteObjectiveV2(objective); err != nil {
		t.Fatalf("WriteObjectiveV2() from unrelated cwd error = %v", err)
	}
	task.Status = ColumnInProgress
	task.Stage = StageBuild
	if err := WriteTaskV2(task); err != nil {
		t.Fatalf("WriteTaskV2() from unrelated cwd error = %v", err)
	}

	objectiveBytes, err := os.ReadFile(filepath.Join(root, objective.Source.Path))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(objectiveBytes), "status: in_progress") {
		t.Errorf("objective file at project root was not updated: %s", objectiveBytes)
	}
	taskBytes, err := os.ReadFile(filepath.Join(root, task.Source.Path))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(taskBytes), "status: in_progress") || !strings.Contains(string(taskBytes), "stage: build") {
		t.Errorf("task file at project root was not updated: %s", taskBytes)
	}
}

func TestResolveV2SourcePath_confinesPathsToProjectRoot(t *testing.T) {
	root := t.TempDir()
	inside := filepath.Join(root, "objectives", "O-001", "Objective.md")

	tests := []struct {
		name    string
		path    string
		want    string
		wantErr bool
	}{
		{
			name: "relative inside",
			path: "objectives/O-001/Objective.md",
			want: inside,
		},
		{
			name:    "relative escape",
			path:    filepath.Join("..", "outside.md"),
			wantErr: true,
		},
		{
			name: "absolute inside",
			path: inside,
			want: inside,
		},
		{
			name:    "absolute outside",
			path:    filepath.Join(root, "..", "outside.md"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveV2SourcePath(V2SourceDocument{
				Path:        tt.path,
				ProjectRoot: root,
			})
			if tt.wantErr {
				if !errors.Is(err, ErrV2UnsafePath) {
					t.Fatalf("resolveV2SourcePath() error = %v, want ErrV2UnsafePath", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveV2SourcePath() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("resolveV2SourcePath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveV2SourcePath_allowsAbsolutePathWithoutProjectRoot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.md")

	got, err := resolveV2SourcePath(V2SourceDocument{Path: path})
	if err != nil {
		t.Fatalf("resolveV2SourcePath() error = %v", err)
	}
	if got != path {
		t.Errorf("resolveV2SourcePath() = %q, want %q", got, path)
	}
}

func TestWriteV2Record_refusesStaleLoadedSourceWithoutOverwritingUserEdit(t *testing.T) {
	root := t.TempDir()
	writeV2ObjectiveFixture(t, root, "O-001-first", "O-001", "First objective")
	writeV2TaskFixture(t, root, "O-001-first", "T-001-first.md", "T-001", "First task", "O-001")

	_, tasks, err := DiscoverV2Records(root)
	if err != nil {
		t.Fatalf("DiscoverV2Records() error = %v", err)
	}
	task := tasks["T-001"]
	if task == nil {
		t.Fatal("DiscoverV2Records() did not return T-001")
	}

	path := filepath.Join(root, task.Source.Path)
	userEdit := "---\nid: T-001\ntitle: \"First task\"\nobjective: O-001\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\neditor_note: \"owner edit\"\n---\n\n# First task\n\nOwner edit must survive.\n"
	if err := os.WriteFile(path, []byte(userEdit), 0644); err != nil {
		t.Fatal(err)
	}

	task.Status = ColumnInProgress
	task.Stage = StageBuild
	err = WriteTaskV2(task)
	if !errors.Is(err, ErrV2SourceConflict) {
		t.Fatalf("WriteTaskV2() error = %v, want ErrV2SourceConflict", err)
	}
	if !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("WriteTaskV2() error = %v, want ErrMtimeConflict compatibility marker", err)
	}
	if !strings.Contains(err.Error(), path) {
		t.Errorf("WriteTaskV2() error = %v, want path %q", err, path)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != userEdit {
		t.Fatalf("stale write changed user content:\n got: %s\nwant: %s", got, userEdit)
	}
}

func TestWriteTaskEvidenceV2_setsAllSubBlocksPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-020.md")
	content := `---
id: T-020
title: "No evidence yet"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: build
depends_on:
  - task: T-010
release: R-001
---

# Task

Authored task notes.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Evidence != nil {
		t.Fatalf("Evidence = %+v, want nil before write", task.Evidence)
	}

	assessedAt := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)
	recordedAt := time.Date(2026, 9, 14, 1, 0, 0, 0, time.UTC)
	replanAt := time.Date(2026, 9, 14, 2, 0, 0, 0, time.UTC)
	waiverAt := time.Date(2026, 9, 14, 3, 0, 0, 0, time.UTC)
	task.Evidence = &Evidence{
		CheckWaiver: &CheckWaiver{
			Task:       "T-020",
			Reason:     "Owner completed via board without requesting a Check.",
			Actor:      Actor{Role: ActorRoleOwner, Session: "board-owner"},
			RecordedAt: waiverAt,
		},
		LastCheck: "C-001",
		Freshness: &Freshness{
			State:      FreshnessCurrent,
			Check:      "C-001",
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: assessedAt,
			Basis:      "Reviewed diff against AC.",
		},
		OwnerValidation: &OwnerValidation{Required: true, AcceptedCheck: "C-001", AcceptedBy: Actor{Role: ActorRoleOwner, Session: "owner-1"}},
		Exception: &Exception{
			Requirements: []string{"TEST-08"},
			Reason:       "Owner accepted known risk.",
			Owner:        "owner-1",
			RecordedAt:   recordedAt,
			Check:        "C-001",
		},
		Replan: &Replan{
			Reason:     "Plan needs revisiting.",
			RecordedBy: Actor{Role: ActorRolePlanner, Session: "sess-2"},
			RecordedAt: replanAt,
		},
	}

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Evidence == nil {
		t.Fatal("Evidence = nil, want populated evidence")
	}
	if reparsed.Evidence.LastCheck != "C-001" {
		t.Errorf("LastCheck = %q, want C-001", reparsed.Evidence.LastCheck)
	}
	if reparsed.Evidence.Freshness == nil || reparsed.Evidence.Freshness.State != FreshnessCurrent {
		t.Errorf("Freshness = %+v, want state current", reparsed.Evidence.Freshness)
	}
	if reparsed.Evidence.OwnerValidation == nil || !reparsed.Evidence.OwnerValidation.Required {
		t.Errorf("OwnerValidation = %+v, want required true", reparsed.Evidence.OwnerValidation)
	}
	if reparsed.Evidence.OwnerValidation.AcceptedBy != (Actor{Role: ActorRoleOwner, Session: "owner-1"}) {
		t.Errorf("OwnerValidation.AcceptedBy = %+v, want owner/owner-1", reparsed.Evidence.OwnerValidation.AcceptedBy)
	}
	if reparsed.Evidence.Exception == nil || reparsed.Evidence.Exception.Owner != "owner-1" {
		t.Errorf("Exception = %+v, want owner owner-1", reparsed.Evidence.Exception)
	}
	if reparsed.Evidence.Replan == nil || reparsed.Evidence.Replan.Reason != "Plan needs revisiting." {
		t.Errorf("Replan = %+v, want reason set", reparsed.Evidence.Replan)
	}
	if reparsed.Evidence.CheckWaiver == nil || reparsed.Evidence.CheckWaiver.Task != "T-020" || reparsed.Evidence.CheckWaiver.Actor.Role != ActorRoleOwner {
		t.Errorf("CheckWaiver = %+v, want task T-020 recorded by owner", reparsed.Evidence.CheckWaiver)
	}

	if reparsed.Status != ColumnInProgress || reparsed.Stage != StageBuild {
		t.Errorf("Status/Stage = %q/%q, want in_progress/build preserved", reparsed.Status, reparsed.Stage)
	}
	if len(reparsed.DependsOn) != 1 || reparsed.DependsOn[0].Task != "T-010" {
		t.Errorf("DependsOn = %v, want [T-010] preserved", reparsed.DependsOn)
	}
	if !strings.Contains(string(result), "Authored task notes.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteTaskEvidenceV2_removesReplanKeyRatherThanEmptyValue(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-021.md")
	content := `---
id: T-021
title: "Replan flagged"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: build
last_check: C-001
replan:
  reason: "Plan needs revisiting."
  recorded_by:
    role: planner
    session: sess-2
  recorded_at: '2026-09-14T02:00:00Z'
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	if task.Evidence == nil || task.Evidence.Replan == nil {
		t.Fatalf("Evidence = %+v, want replan present before write", task.Evidence)
	}

	task.Evidence.Replan = nil

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(result), "replan:") {
		t.Error("replan key should be removed, not written empty")
	}

	reparsed, err := DecodeTaskV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeTaskV2() after write error = %v", err)
	}
	if reparsed.Evidence == nil {
		t.Fatal("Evidence = nil, want last_check preserved")
	}
	if reparsed.Evidence.Replan != nil {
		t.Errorf("Replan = %+v, want nil after clearing", reparsed.Evidence.Replan)
	}
	if reparsed.Evidence.LastCheck != "C-001" {
		t.Errorf("LastCheck = %q, want C-001 preserved", reparsed.Evidence.LastCheck)
	}
}

func TestWriteTaskEvidenceV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-022.md")
	content := `---
id: T-022
title: "Fully evidenced task"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: in_progress
stage: build
last_check: C-001
freshness:
  state: current
  check: C-001
  assessed_by:
    role: checker
    session: sess-1
  assessed_at: '2026-09-14T00:00:00Z'
  basis: "Reviewed diff against AC."
owner_validation:
  required: true
  accepted_check: C-001
  accepted_by:
    role: owner
    session: owner-1
exception:
  requirements:
    - TEST-08
  reason: "Owner accepted known risk."
  owner: owner-1
  recorded_at: '2026-09-14T01:00:00Z'
  check: C-001
replan:
  reason: "Plan needs revisiting."
  recorded_by:
    role: planner
    session: sess-2
  recorded_at: '2026-09-14T02:00:00Z'
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteTaskEvidenceV2_refusesStaleSourceWithoutOverwritingUserEdit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-023.md")
	content := `---
id: T-023
title: "Guarded evidence write"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}

	userEdit := "---\nid: T-023\ntitle: \"Guarded evidence write\"\nobjective: O-002\nplanned_by: {role: planner, session: planning-fixture}\nstatus: planned\neditor_note: \"owner edit\"\n---\n\n# Task\n\nOwner edit must survive.\n"
	if err := os.WriteFile(path, []byte(userEdit), 0644); err != nil {
		t.Fatal(err)
	}

	task.Evidence = &Evidence{LastCheck: "C-001"}

	err = WriteTaskEvidenceV2(task)
	if !errors.Is(err, ErrV2SourceConflict) {
		t.Fatalf("WriteTaskEvidenceV2() error = %v, want ErrV2SourceConflict", err)
	}
	if !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("WriteTaskEvidenceV2() error = %v, want ErrMtimeConflict compatibility marker", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != userEdit {
		t.Fatalf("stale write changed user content:\n got: %s\nwant: %s", got, userEdit)
	}
}

func TestWriteTaskEvidenceV2_preservesCRLFLineEndings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-024.md")
	content := "---\r\nid: T-024\r\ntitle: \"CRLF evidence task\"\r\nobjective: O-002\r\nplanned_by: {role: planner, session: planning-fixture}\r\nstatus: planned\r\n---\r\n\r\n# Task\r\n\r\nAuthored notes.\r\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Evidence = &Evidence{LastCheck: "C-001"}

	if err := WriteTaskEvidenceV2(task); err != nil {
		t.Fatalf("WriteTaskEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(result), "\n") != strings.Count(string(result), "\r\n") {
		t.Errorf("result did not preserve CRLF line endings throughout: %q", string(result))
	}
	if !strings.Contains(string(result), "Authored notes.") {
		t.Error("authored body content not preserved across CRLF rewrite")
	}
}

func TestWriteTaskEvidenceV2_rejectsMalformedEvidenceLeavesFileUntouched(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "T-025.md")
	content := `---
id: T-025
title: "Guarded malformed evidence"
objective: O-002
planned_by: {role: planner, session: planning-fixture}
status: planned
---

# Task`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	task, err := DecodeTaskV2(path, content)
	if err != nil {
		t.Fatalf("DecodeTaskV2() error = %v", err)
	}
	task.Evidence = &Evidence{
		Freshness: &Freshness{
			State:      FreshnessState("bogus"),
			Check:      "C-001",
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC),
			Basis:      "Reviewed diff against AC.",
		},
	}

	err = WriteTaskEvidenceV2(task)
	if err == nil {
		t.Fatal("WriteTaskEvidenceV2() expected error for malformed freshness state")
	}
	if !errors.Is(err, ErrV2EvidenceMalformed) {
		t.Fatalf("WriteTaskEvidenceV2() error = %v, want ErrV2EvidenceMalformed", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != content {
		t.Error("file content changed despite validation refusal")
	}
}

func TestWriteObjectiveEvidenceV2_setsSubBlocksPreservesUnknownFieldsAndBody(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "O-020.md")
	content := `---
id: O-020
title: "No evidence yet"
status: in_progress
depends_on: [O-001]
release: R-001
---

# Objective

Authored objective notes.`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}
	if objective.Evidence != nil {
		t.Fatalf("Evidence = %+v, want nil before write", objective.Evidence)
	}

	assessedAt := time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC)
	objective.Evidence = &Evidence{
		LastCheck: "C-001",
		Freshness: &Freshness{
			State:      FreshnessCurrent,
			Check:      "C-001",
			AssessedBy: Actor{Role: ActorRoleChecker, Session: "sess-1"},
			AssessedAt: assessedAt,
			Basis:      "Objective integration Check is current.",
		},
	}

	if err := WriteObjectiveEvidenceV2(objective); err != nil {
		t.Fatalf("WriteObjectiveEvidenceV2() error = %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	reparsed, err := DecodeObjectiveV2(path, string(result))
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() after write error = %v", err)
	}
	if reparsed.Evidence == nil || reparsed.Evidence.LastCheck != "C-001" {
		t.Errorf("Evidence = %+v, want LastCheck C-001", reparsed.Evidence)
	}
	if reparsed.Evidence.Freshness == nil || reparsed.Evidence.Freshness.State != FreshnessCurrent {
		t.Errorf("Freshness = %+v, want state current", reparsed.Evidence.Freshness)
	}
	if len(reparsed.DependsOn) != 1 || reparsed.DependsOn[0] != "O-001" {
		t.Errorf("DependsOn = %v, want [O-001] preserved", reparsed.DependsOn)
	}
	if reparsed.Release != "R-001" {
		t.Errorf("Release = %q, want R-001 preserved", reparsed.Release)
	}
	if !strings.Contains(string(result), "Authored objective notes.") {
		t.Error("authored body content not preserved")
	}
}

func TestWriteObjectiveEvidenceV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "O-021.md")
	content := `---
id: O-021
title: "No-op objective evidence"
status: in_progress
last_check: C-001
freshness:
  state: current
  check: C-001
  assessed_by:
    role: checker
    session: sess-1
  assessed_at: "2026-09-14T00:00:00Z"
  basis: "Objective integration Check is current."
---

# Objective`

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

	objective, err := DecodeObjectiveV2(path, content)
	if err != nil {
		t.Fatalf("DecodeObjectiveV2() error = %v", err)
	}

	if err := WriteObjectiveEvidenceV2(objective); err != nil {
		t.Fatalf("WriteObjectiveEvidenceV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	afterBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", before.ModTime(), after.ModTime())
	}
	if string(beforeBytes) != string(afterBytes) {
		t.Error("file bytes changed on no-op write")
	}
}

// routerV2FixtureContent includes an old router's long quoted next_action to
// prove that selection writes preserve the retired key byte-for-byte. The
// reader tolerates it for compatibility, while the writer changes selection
// keys only and never creates the retired field.
func routerV2FixtureContent() string {
	return `# Agent State Machine

This file routes the agent. The active skill is the canonical workflow source.

## Read order

1. This file
2. The matching skill for ` + "`state`" + `

## Current state

` + "```" + `yaml
state: design
objective: none
task: none
next_action: "Turn a rough idea — a single sentence is a valid starting point, no prepared requirements document needed — into .savepoint/Idea.md through a short back-and-forth with the owner."
` + "```" + `

## State Meanings

- ` + "`idea`" + `: intent and boundary are being defined.

` + "```" + `yaml
project_note: a second fenced block the writer must not touch
` + "```" + `
`
}

func writeRouterV2Fixture(t *testing.T, content string) (root, path string, mtime time.Time) {
	t.Helper()
	root = t.TempDir()
	path = filepath.Join(root, "router.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return root, path, info.ModTime()
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(raw)
}

func TestWriteRouterStateV2_setsSelectionAndPreservesEveryOtherByte(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	want := strings.Replace(content, "objective: none", "objective: O-001", 1)
	want = strings.Replace(want, "task: none", "task: T-005", 1)

	if got := readFileString(t, path); got != want {
		t.Errorf("file content changed beyond the selection keys.\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteRouterStateV2_setsIssueAndPreservesEveryOtherByte(t *testing.T) {
	content := routerV2FixtureContent()
	content = strings.Replace(content, "task: none", "task: none\nissue: none", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005", Issue: "I-042"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	want := strings.Replace(content, "objective: none", "objective: O-001", 1)
	want = strings.Replace(want, "task: none", "task: T-005", 1)
	want = strings.Replace(want, "issue: none", "issue: I-042", 1)
	if got := readFileString(t, path); got != want {
		t.Errorf("file content changed beyond the selection keys.\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriteRouterStateV2_addsIssueAlone(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\n```\n"
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Issue: "I-042"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}
	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Issue != "I-042" || state.Objective != "" || state.Task != "" {
		t.Fatalf("selection = issue %q objective %q task %q, want I-042 alone", state.Issue, state.Objective, state.Task)
	}
}

func TestRouterSelectionAfterClosureV2_advancesToLowestUnfinishedEvenWhenBlocked(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{
			"O-001": {ID: "O-001", Status: ColumnInProgress},
		},
		Tasks: map[string]*TaskV2{
			"T-001": {ID: "T-001", Objective: "O-001", Status: ColumnDone},
			"T-002": {ID: "T-002", Objective: "O-001", Status: ColumnPlanned,
				DependsOn: []TaskDependencyV2{{Task: "T-004", Requires: TaskDependencyClear}}},
			"T-004": {ID: "T-004", Objective: "O-001", Status: ColumnInProgress, Stage: StageBuild},
			"T-900": {ID: "T-900", Objective: "O-001", Status: ColumnDone},
		},
		ObjectiveTasks: map[string][]string{"O-001": {"T-001", "T-002", "T-004", "T-900"}},
	}
	current := RouterSelectionV2{Release: "R-006", Objective: "O-001", Task: "T-900", Issue: "I-042"}
	closed := RouterSelectionV2{Objective: "O-001", Task: "T-900"}

	got, changed := RouterSelectionAfterClosureV2(index, current, closed)
	if !changed {
		t.Fatal("RouterSelectionAfterClosureV2() changed = false, want selected Task advanced")
	}
	want := RouterSelectionV2{Release: "R-006", Objective: "O-001", Task: "T-002", Issue: "I-042"}
	if got != want {
		t.Fatalf("RouterSelectionAfterClosureV2() = %+v, want %+v", got, want)
	}
	if decision := ResolveTaskStart(index, "T-002"); decision.Allowed || len(decision.Blockers) == 0 {
		t.Fatalf("ResolveTaskStart(T-002) = %+v, want blocked next Task", decision)
	}
}

func TestRouterSelectionAfterClosureV2_clearsMatchingSelectionsAndLeavesOthersAlone(t *testing.T) {
	index := &V2Index{
		Objectives: map[string]*ObjectiveV2{"O-001": {ID: "O-001", Status: ColumnDone}},
		Tasks: map[string]*TaskV2{
			"T-001": {ID: "T-001", Objective: "O-001", Status: ColumnDone},
		},
		ObjectiveTasks: map[string][]string{"O-001": {"T-001"}},
	}
	current := RouterSelectionV2{Release: "R-006", Objective: "O-001", Task: "T-001", Issue: "I-042"}

	got, changed := RouterSelectionAfterClosureV2(index, current, RouterSelectionV2{Objective: "O-001", Task: "T-001"})
	if !changed {
		t.Fatal("RouterSelectionAfterClosureV2() changed = false, want completed Task selection cleared")
	}
	if want := (RouterSelectionV2{Release: "R-006", Objective: "O-001", Issue: "I-042"}); got != want {
		t.Fatalf("all-done Task selection = %+v, want %+v", got, want)
	}

	got, changed = RouterSelectionAfterClosureV2(index, current, RouterSelectionV2{Objective: "O-002", Task: "T-001"})
	if changed || got != current {
		t.Fatalf("non-matching Task closure = %+v, changed %t; want original selection and no change", got, changed)
	}

	got, changed = RouterSelectionAfterClosureV2(index, current, RouterSelectionV2{Objective: "O-001"})
	if !changed {
		t.Fatal("RouterSelectionAfterClosureV2() changed = false, want matching Objective selection cleared")
	}
	if want := (RouterSelectionV2{Release: "R-006", Issue: "I-042"}); got != want {
		t.Fatalf("Objective selection = %+v, want %+v", got, want)
	}

	got, changed = RouterSelectionAfterClosureV2(index, current, RouterSelectionV2{Objective: "O-002"})
	if changed || got != current {
		t.Fatalf("non-matching Objective closure = %+v, changed %t; want original selection and no change", got, changed)
	}
}

func TestWriteRouterStateV2_setsReleaseContextWithoutChangingRouterProse(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, mtime := writeRouterV2Fixture(t, content)
	before, err := NewRouterReader().ReadStateV2(content)
	if err != nil {
		t.Fatalf("ReadStateV2() before release write error = %v", err)
	}

	if err := WriteRouterStateV2(root, RouterSelectionV2{Release: "R-001", Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() after release write error = %v", err)
	}
	if state.Release != "R-001" || state.Objective != "O-001" || state.Task != "T-005" {
		t.Fatalf("selections = release %q objective %q task %q, want R-001/O-001/T-005", state.Release, state.Objective, state.Task)
	}
	if state.State != before.State || state.HasRetiredNextAction != before.HasRetiredNextAction {
		t.Fatalf("router lifecycle/retired-key presence changed: before state %q key-present %t, after state %q key-present %t", before.State, before.HasRetiredNextAction, state.State, state.HasRetiredNextAction)
	}
	got := readFileString(t, path)
	if !strings.Contains(got, "project_note: a second fenced block the writer must not touch") || !strings.Contains(got, "This file routes the agent.") {
		t.Error("release selection write did not preserve surrounding document bytes")
	}
}

func TestWriteRouterStateV2_repeatingReleaseSelectionLeavesBytesAndMtimeUnchanged(t *testing.T) {
	content := strings.NewReplacer(
		"objective: none", "objective: O-001",
		"task: none", "task: T-005",
	).Replace(routerV2FixtureContent())
	content = strings.Replace(content, "state: design", "state: task\nrelease: R-001", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)

	time.Sleep(10 * time.Millisecond)
	if err := WriteRouterStateV2(root, RouterSelectionV2{Release: "R-001", Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}
	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(mtime) || readFileString(t, path) != content {
		t.Fatal("repeating unchanged Release selection changed bytes or mtime")
	}
}

func TestWriteRouterStateV2_clearingSelectionWritesTheNoneSentinel(t *testing.T) {
	content := strings.NewReplacer(
		"objective: none", "objective: O-001",
		"task: none", "task: T-005",
	).Replace(routerV2FixtureContent())
	content = strings.Replace(content, "task: T-005", "task: T-005\nissue: I-042", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	got := readFileString(t, path)
	if !strings.Contains(got, "objective: none") || !strings.Contains(got, "task: none") || !strings.Contains(got, "issue: none") {
		t.Errorf("cleared selection not written as the none sentinel:\n%s", got)
	}

	state, err := NewRouterReader().ReadStateV2(got)
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "" || state.Task != "" {
		t.Errorf("cleared selection read back as objective %q task %q", state.Objective, state.Task)
	}
}

func TestWriteRouterStateV2_doesNotAddSentinelKeysTheDocumentLacks(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: idea\n```\n"
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	if got := readFileString(t, path); got != content {
		t.Errorf("clearing an absent selection changed the file:\n%s", got)
	}
}

func TestWriteRouterStateV2_addsSelectionKeysToRecordARealSelection(t *testing.T) {
	content := "## Current state\n\n```yaml\nstate: task\n```\n"
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-007", Task: "T-042"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "O-007" || state.Task != "T-042" {
		t.Errorf("selection read back as objective %q task %q, want O-007/T-042", state.Objective, state.Task)
	}
	if state.State != RouterPhaseTask {
		t.Errorf("state changed to %q", state.State)
	}
}

func TestWriteRouterStateV2_refusesMalformedSelectionAndLeavesFileUntouched(t *testing.T) {
	cases := []struct {
		name      string
		selection RouterSelectionV2
		wantErr   error
	}{
		{"objective wrong family", RouterSelectionV2{Objective: "T-001"}, ErrV2InvalidID},
		{"release wrong family", RouterSelectionV2{Release: "O-001"}, ErrV2InvalidID},
		{"release too few digits", RouterSelectionV2{Release: "R01"}, ErrV2InvalidID},
		{"objective too few digits", RouterSelectionV2{Objective: "O1"}, ErrV2InvalidID},
		{"objective with trailing slug", RouterSelectionV2{Objective: "O-001-recovery"}, ErrV2InvalidID},
		{"two objectives", RouterSelectionV2{Objective: "O-001 O-002"}, ErrV2InvalidID},
		{"task wrong family", RouterSelectionV2{Objective: "O-001", Task: "C-001"}, ErrV2InvalidID},
		{"task too few digits", RouterSelectionV2{Objective: "O-001", Task: "T4"}, ErrV2InvalidID},
		{"issue wrong family", RouterSelectionV2{Issue: "C-001"}, ErrV2InvalidID},
		{"issue too few digits", RouterSelectionV2{Issue: "I4"}, ErrV2InvalidID},
		{"issue with trailing slug", RouterSelectionV2{Issue: "I-001-recovery"}, ErrV2InvalidID},
		{"task without objective", RouterSelectionV2{Task: "T-005"}, ErrV2InvalidOwnership},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			content := routerV2FixtureContent()
			root, path, mtime := writeRouterV2Fixture(t, content)

			err := WriteRouterStateV2(root, tc.selection, mtime)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("WriteRouterStateV2() error = %v, want %v", err, tc.wantErr)
			}
			if got := readFileString(t, path); got != content {
				t.Error("file changed despite a refused selection")
			}
		})
	}
}

func TestWriteRouterStateV2_noOpLeavesBytesAndMtimeUnchanged(t *testing.T) {
	content := strings.NewReplacer(
		"objective: none", "objective: O-001",
		"task: none", "task: T-005",
	).Replace(routerV2FixtureContent())
	root, path, mtime := writeRouterV2Fixture(t, content)

	time.Sleep(10 * time.Millisecond)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	after, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if !after.ModTime().Equal(mtime) {
		t.Errorf("ModTime changed on no-op write: before %v, after %v", mtime, after.ModTime())
	}
	if got := readFileString(t, path); got != content {
		t.Error("file bytes changed on no-op write")
	}
}

func TestWriteRouterStateV2_refusesStaleMtimeAndLeavesFileUntouched(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, _ := writeRouterV2Fixture(t, content)

	err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, time.Now().Add(-time.Hour))
	if !errors.Is(err, ErrMtimeConflict) {
		t.Fatalf("WriteRouterStateV2() error = %v, want ErrMtimeConflict", err)
	}
	if got := readFileString(t, path); got != content {
		t.Error("file changed despite an mtime conflict")
	}
}

func TestWriteRouterStateV2_refusesUnreadableDocuments(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantErr string
	}{
		{"no anchor", "# Agent State Machine\n\nNo current state here.\n", "no Current state block found"},
		{"no yaml fence", "## Current state\n\nstate: idea\n", "no yaml code block found"},
		{"unterminated fence", "## Current state\n\n```yaml\nstate: idea\n", "no closing code block found"},
		{"anchor is not a mapping", "## Current state\n\n```yaml\n- state: idea\n```\n", "not a mapping"},
		{"anchor is not yaml", "## Current state\n\n```yaml\nstate: [idea\n```\n", "router state"},
		// The anchor's key set is closed on the way in, so a key the reader
		// would reject is refused on the way out too rather than being
		// written into a document that would then fail to load.
		{"anchor carries an unknown key", "## Current state\n\n```yaml\nstate: idea\nproject_note: extra\n```\n", "project_note"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, path, mtime := writeRouterV2Fixture(t, tc.content)

			err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, mtime)
			if err == nil {
				t.Fatal("WriteRouterStateV2() expected an error")
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("WriteRouterStateV2() error = %v, want it to mention %q", err, tc.wantErr)
			}
			if got := readFileString(t, path); got != tc.content {
				t.Error("file changed despite a refused write")
			}
		})
	}
}

func TestWriteRouterStateV2_refusesMissingFile(t *testing.T) {
	root := t.TempDir()

	err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, time.Now())
	if err == nil {
		t.Fatal("WriteRouterStateV2() expected an error for a missing router.md")
	}
	if !strings.Contains(err.Error(), "router.md") {
		t.Errorf("WriteRouterStateV2() error = %v, want it to name router.md", err)
	}
}

func TestWriteRouterStateV2_refusesARouterItCouldNotReadBack(t *testing.T) {
	content := strings.Replace(routerV2FixtureContent(), "state: design", "state: audit-pending", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)

	err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, mtime)
	if !errors.Is(err, ErrV2InvalidLifecycle) {
		t.Fatalf("WriteRouterStateV2() error = %v, want ErrV2InvalidLifecycle", err)
	}
	if got := readFileString(t, path); got != content {
		t.Error("file changed despite content that could not be read back")
	}
}

func TestWriteRouterStateV2_roundTripsThroughReadStateV2(t *testing.T) {
	content := routerV2FixtureContent()
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-012", Task: "T-034"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	state, err := NewRouterReader().ReadStateV2(readFileString(t, path))
	if err != nil {
		t.Fatalf("ReadStateV2() error = %v", err)
	}
	if state.Objective != "O-012" || state.Task != "T-034" {
		t.Errorf("round trip gave objective %q task %q, want O-012/T-034", state.Objective, state.Task)
	}
	if state.State != RouterPhaseDesign {
		t.Errorf("state changed to %q, want design", state.State)
	}
	if !state.HasRetiredNextAction {
		t.Error("HasRetiredNextAction = false after selection write, want true")
	}
}

func TestWriteRouterStateV2_preservesRetiredNextActionLineByteForByte(t *testing.T) {
	content := routerV2FixtureContent()
	var nextActionLine string
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "next_action:") {
			nextActionLine = line
			break
		}
	}
	if nextActionLine == "" {
		t.Fatal("router fixture has no next_action line")
	}

	root, path, mtime := writeRouterV2Fixture(t, content)
	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}
	if got := readFileString(t, path); strings.Count(got, nextActionLine) != 1 {
		t.Errorf("selection write did not preserve the retired line exactly once:\n%s", got)
	}
}

func TestWriteRouterStateV2_doesNotAddRetiredNextAction(t *testing.T) {
	content := routerV2FixtureContent()
	var nextActionLine string
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "next_action:") {
			nextActionLine = line
			break
		}
	}
	content = strings.Replace(content, nextActionLine+"\n", "", 1)
	root, path, mtime := writeRouterV2Fixture(t, content)
	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001", Task: "T-005"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}
	if got := readFileString(t, path); strings.Contains(got, "next_action:") {
		t.Errorf("selection write added retired next_action:\n%s", got)
	}
}

func TestWriteRouterStateV2_preservesCRLFLineEndings(t *testing.T) {
	content := strings.ReplaceAll(routerV2FixtureContent(), "\n", "\r\n")
	root, path, mtime := writeRouterV2Fixture(t, content)

	if err := WriteRouterStateV2(root, RouterSelectionV2{Objective: "O-001"}, mtime); err != nil {
		t.Fatalf("WriteRouterStateV2() error = %v", err)
	}

	got := readFileString(t, path)
	if strings.Contains(strings.ReplaceAll(got, "\r\n", ""), "\n") {
		t.Error("CRLF router gained bare LF line endings")
	}
	want := strings.Replace(content, "objective: none", "objective: O-001", 1)
	if got != want {
		t.Errorf("CRLF router changed beyond the selection key:\n%q", got)
	}
}

func TestWriteObjectiveGroupOrderV2_preservesContentAndNoOpFiles(t *testing.T) {
	root := writeObjectiveOrderProjectFixture(t)
	firstPath := writeObjectiveOrderRecordFixture(t, root, "O-001", ObjectivePriorityCritical, 2,
		"owner:\n  team: platform\n", "# O-001\n\nKeep this authored note.\n")
	secondPath := writeObjectiveOrderRecordFixture(t, root, "O-002", ObjectivePriorityCritical, 1, "", "# O-002\n")
	thirdPath := writeObjectiveOrderRecordFixture(t, root, "O-003", ObjectivePriorityCritical, 3, "", "# O-003\n")
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	thirdBeforeInfo, err := os.Stat(thirdPath)
	if err != nil {
		t.Fatal(err)
	}
	thirdBeforeBytes, err := os.ReadFile(thirdPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := WriteObjectiveGroupOrderV2(index, "R-001", ObjectivePriorityCritical, []string{"O-001", "O-002", "O-003"}); err != nil {
		t.Fatalf("WriteObjectiveGroupOrderV2() error = %v", err)
	}
	updated, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() after reorder error = %v", err)
	}
	if updated.Objectives["O-001"].Rank != 1 || updated.Objectives["O-002"].Rank != 2 || updated.Objectives["O-003"].Rank != 3 {
		t.Errorf("updated ranks = (%d, %d, %d), want (1, 2, 3)", updated.Objectives["O-001"].Rank, updated.Objectives["O-002"].Rank, updated.Objectives["O-003"].Rank)
	}
	firstAfter, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"team: platform", "Keep this authored note."} {
		if !strings.Contains(string(firstAfter), want) {
			t.Errorf("Objective rewrite lost %q:\n%s", want, firstAfter)
		}
	}
	thirdAfterInfo, err := os.Stat(thirdPath)
	if err != nil {
		t.Fatal(err)
	}
	thirdAfterBytes, err := os.ReadFile(thirdPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(thirdAfterBytes) != string(thirdBeforeBytes) || !thirdAfterInfo.ModTime().Equal(thirdBeforeInfo.ModTime()) {
		t.Error("renumber changed a record whose priority and rank already matched")
	}

	firstInfo, err := os.Stat(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondInfo, err := os.Stat(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	firstBytes := string(firstAfter)
	secondBytes, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	if err := WriteObjectiveGroupOrderV2(updated, "R-001", ObjectivePriorityCritical, []string{"O-001", "O-002", "O-003"}); err != nil {
		t.Fatalf("repeating WriteObjectiveGroupOrderV2() error = %v", err)
	}
	firstInfoAfter, err := os.Stat(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondInfoAfter, err := os.Stat(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	firstBytesAfter, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondBytesAfter, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(firstBytesAfter) != firstBytes || string(secondBytesAfter) != string(secondBytes) ||
		!firstInfo.ModTime().Equal(firstInfoAfter.ModTime()) || !secondInfo.ModTime().Equal(secondInfoAfter.ModTime()) {
		t.Error("repeating the same Objective order changed file bytes or modification times")
	}
}

func TestWriteObjectiveGroupOrderV2_staleMemberWritesNothing(t *testing.T) {
	root := writeObjectiveOrderProjectFixture(t)
	firstPath := writeObjectiveOrderRecordFixture(t, root, "O-001", ObjectivePriorityCritical, 2, "", "# O-001\n")
	secondPath := writeObjectiveOrderRecordFixture(t, root, "O-002", ObjectivePriorityCritical, 1, "", "# O-002\n")
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	firstBefore, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondRaw, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	stale := strings.Replace(string(secondRaw), "rank: 1\n", "rank: 1\nexternal_note: edited after load\n", 1)
	if stale == string(secondRaw) {
		t.Fatal("second Objective has no rank to extend")
	}
	if err := os.WriteFile(secondPath, []byte(stale), 0644); err != nil {
		t.Fatal(err)
	}

	err = WriteObjectiveGroupOrderV2(index, "R-001", ObjectivePriorityCritical, []string{"O-001", "O-002"})
	if !errors.Is(err, ErrV2SourceConflict) {
		t.Fatalf("WriteObjectiveGroupOrderV2() error = %v, want ErrV2SourceConflict", err)
	}
	firstAfter, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(firstAfter) != string(firstBefore) {
		t.Error("stale later member caused an earlier Objective to be written")
	}
}

func TestWriteObjectiveGroupOrderV2_partialFailureRemainsLoadableAndDiagnosed(t *testing.T) {
	root := writeObjectiveOrderProjectFixture(t)
	firstPath := writeObjectiveOrderRecordFixture(t, root, "O-001", ObjectivePriorityCritical, 3, "", "# O-001\n")
	secondPath := writeObjectiveOrderRecordFixture(t, root, "O-002", ObjectivePriorityCritical, 1, "", "# O-002\n")
	index, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() error = %v", err)
	}
	secondBefore, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(secondPath, 0444); err != nil {
		t.Fatal(err)
	}

	err = WriteObjectiveGroupOrderV2(index, "R-001", ObjectivePriorityCritical, []string{"O-001", "O-002"})
	if !errors.Is(err, os.ErrPermission) {
		t.Fatalf("WriteObjectiveGroupOrderV2() error = %v, want permission failure on the second file", err)
	}
	firstAfter, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondAfter, err := os.ReadFile(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(firstAfter), "rank: 1") {
		t.Errorf("first file rank = %q, want successful first replacement", firstAfter)
	}
	if string(secondAfter) != string(secondBefore) {
		t.Error("second file changed despite its replacement failure")
	}

	reloaded, err := LoadV2Index(root)
	if err != nil {
		t.Fatalf("LoadV2Index() after partial write error = %v, want loadable state", err)
	}
	if got, want := OrderedObjectiveIDsForGoal(reloaded, "R-001"), []string{"O-001", "O-002"}; !reflect.DeepEqual(got, want) {
		t.Errorf("partial order = %v, want deterministic order %v", got, want)
	}
	if len(reloaded.DuplicateObjectiveRanks) != 1 || !reflect.DeepEqual(reloaded.DuplicateObjectiveRanks[0].ObjectiveIDs, []string{"O-001", "O-002"}) {
		t.Errorf("partial duplicate-rank facts = %+v, want a diagnostic for O-001 and O-002", reloaded.DuplicateObjectiveRanks)
	}
}

func writeObjectiveOrderProjectFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	releaseDir := filepath.Join(root, "releases", "R-001-goal")
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", releaseDir, err)
	}
	release := "---\nid: R-001\ntitle: \"Ordering Goal\"\nstatus: planned\n---\n\n" +
		"## Outcome\n\nOrder the Goal's Objectives.\n\n" +
		"## Why\n\nPlanning order is useful.\n\n" +
		"## Success Conditions\n\n- Objectives have a deterministic order.\n\n" +
		"## Boundaries\n\nMembership remains derived from Objective records.\n"
	if err := os.WriteFile(filepath.Join(releaseDir, "Release.md"), []byte(release), 0644); err != nil {
		t.Fatalf("WriteFile(Release.md) error = %v", err)
	}
	return root
}

func writeObjectiveOrderRecordFixture(t *testing.T, root, id string, priority ObjectivePriority, rank int, extraFrontmatter, body string) string {
	t.Helper()
	dir := filepath.Join(root, "objectives", id+"-ordering")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", dir, err)
	}
	content := fmt.Sprintf("---\nid: %s\ntitle: \"%s ordering\"\nstatus: planned\nrelease: R-001\npriority: %s\nrank: %d\n", id, id, priority, rank)
	content += extraFrontmatter + "---\n\n" + body
	path := filepath.Join(dir, "Objective.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", path, err)
	}
	return path
}
