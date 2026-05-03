package init

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateTarget_missing(t *testing.T) {
	err := ValidateTarget(filepath.Join(t.TempDir(), "nonexistent"), false)
	if err == nil {
		t.Fatal("ValidateTarget() expected error for missing directory")
	}
	if !errors.Is(err, ErrTargetMissing) {
		t.Fatalf("ValidateTarget() error type = %v, want ErrTargetMissing", err)
	}
}

func TestValidateTarget_notADirectory(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "file")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(filePath, false)
	if err == nil {
		t.Fatal("ValidateTarget() expected error for non-directory")
	}
	if !errors.Is(err, ErrNotADirectory) {
		t.Fatalf("ValidateTarget() error type = %v, want ErrNotADirectory", err)
	}
}

func TestValidateTarget_empty(t *testing.T) {
	dir := t.TempDir()

	err := ValidateTarget(dir, false)
	if err != nil {
		t.Fatalf("ValidateTarget() error = %v, want nil for empty directory", err)
	}
}

func TestValidateTarget_hasCompatibleFiles(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"package.json", ".git", "README.md"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
	}

	err := ValidateTarget(dir, false)
	if err != nil {
		t.Fatalf("ValidateTarget() error = %v, want nil for compatible files", err)
	}
}

func TestValidateTarget_existingSavepoint(t *testing.T) {
	dir := t.TempDir()
	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.Mkdir(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, false)
	if err == nil {
		t.Fatal("ValidateTarget() expected error for existing .savepoint")
	}
	if !errors.Is(err, ErrAlreadyInit) {
		t.Fatalf("ValidateTarget() error type = %v, want ErrAlreadyInit", err)
	}
}

func TestValidateTarget_existingSavepointWithForce(t *testing.T) {
	dir := t.TempDir()
	savepointDir := filepath.Join(dir, ".savepoint")
	if err := os.Mkdir(savepointDir, 0755); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, true)
	if err != nil {
		t.Fatalf("ValidateTarget() with --force error = %v, want nil", err)
	}
}

func TestValidateTarget_conflictingFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, false)
	if err == nil {
		t.Fatal("ValidateTarget() expected error for conflicting AGENTS.md")
	}
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("ValidateTarget() error type = %v, want ErrConflict", err)
	}
}

func TestValidateTarget_conflictingFileWithForce(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("existing"), 0644); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, true)
	if err != nil {
		t.Fatalf("ValidateTarget() with --force error = %v, want nil", err)
	}
}

func TestValidateTarget_conflictingAgentSkillsDirectory(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "agent-skills"), 0755); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, false)
	if err == nil {
		t.Fatal("ValidateTarget() expected error for conflicting agent-skills directory")
	}
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("ValidateTarget() error type = %v, want ErrConflict", err)
	}
}

func TestValidateTarget_conflictingAgentSkillsDirectoryWithForce(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "agent-skills"), 0755); err != nil {
		t.Fatal(err)
	}

	err := ValidateTarget(dir, true)
	if err != nil {
		t.Fatalf("ValidateTarget() with --force error = %v, want nil", err)
	}
}

func TestValidateTarget_emptyStringResolvesToDot(t *testing.T) {
	err := ValidateTarget("", false)
	if err != nil {
		t.Fatalf("ValidateTarget(\"\") error = %v, want nil (resolves to cwd)", err)
	}
}
