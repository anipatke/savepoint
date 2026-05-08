package data

import (
	"testing"
)

func TestRouterReaderReadState(t *testing.T) {
	r := NewRouterReader()
	content := "## Current state\n\n```yaml\nstate: building\nrelease: v1\nepic: E01-go-setup\ntask: E01-go-setup/T002-entrypoint\nnext_action: \"Start T002-entrypoint\"\n```\n"

	state, err := r.ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}

	if state.State != "building" {
		t.Errorf("State = %v, want building", state.State)
	}
	if state.Epic != "E01-go-setup" {
		t.Errorf("Epic = %v, want E01-go-setup", state.Epic)
	}
	if state.Task != "E01-go-setup/T002-entrypoint" {
		t.Errorf("Task = %v, want E01-go-setup/T002-entrypoint", state.Task)
	}
}

func TestRouterReaderCapitalizedHeading(t *testing.T) {
	r := NewRouterReader()
	content := "## Current State\n\n```yaml\nstate: task-building\nrelease: v1\nepic: E01\ntask: E01/T001\nnext_action: \"Do thing\"\n```\n"

	state, err := r.ReadState(content)
	if err != nil {
		t.Fatalf("ReadState() error = %v", err)
	}
	if state.State != "task-building" {
		t.Errorf("State = %v, want task-building", state.State)
	}
}

func TestRouterReaderMissing(t *testing.T) {
	r := NewRouterReader()
	content := "# No state block here"

	_, err := r.ReadState(content)
	if err == nil {
		t.Error("ReadState() expected error for missing state block")
	}
}