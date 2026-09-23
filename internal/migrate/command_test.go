package migrate

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestRunCommand_previewAndRecoveryReportAreReadOnlyOnReadable0555Project(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permissions are not enforced the same way on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory permissions")
	}

	root := copyFixtureProject(t, "v1-basic")
	if err := os.Chmod(root, 0555); err != nil {
		t.Fatalf("chmod project read-only: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(root, 0755) })

	oldProbe := targetWriteabilityProbe
	defer func() { targetWriteabilityProbe = oldProbe }()
	probeCalls := 0
	targetWriteabilityProbe = func(string) error {
		probeCalls++
		return nil
	}

	before := snapshotTree(t, root)
	var preview strings.Builder
	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Stdout:         &preview,
		Now:            fixedClock(time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC)),
		NewOperationID: fixedOperationID("preview"),
	})
	if err != nil || code != 0 {
		t.Fatalf("preview on a readable 0555 project: code = %d, err = %v\n%s", code, err, preview.String())
	}
	if !strings.Contains(preview.String(), "Migration preview") {
		t.Fatalf("preview output = %q, want the preview report", preview.String())
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))

	var recovery strings.Builder
	code, err = RunCommand(CommandOptions{
		Dir:     root,
		Recover: true,
		Stdout:  &recovery,
	})
	if err != nil || code != 0 {
		t.Fatalf("recovery report on a readable 0555 project: code = %d, err = %v\n%s", code, err, recovery.String())
	}
	if !strings.Contains(recovery.String(), "no incomplete migration operation found") {
		t.Fatalf("recovery output = %q, want the no-pending report", recovery.String())
	}
	assertSnapshotsEqual(t, before, snapshotTree(t, root))

	if probeCalls != 0 {
		t.Fatalf("preview/recovery invoked the apply-only writeability probe %d time(s)", probeCalls)
	}
}

func TestRunCommand_applyWriteabilityProbePreservesPreExistingSentinelAndPrefix(t *testing.T) {
	t.Parallel()
	root := copyFixtureProject(t, "v1-basic")
	sentinelPath := filepath.Join(root, ".savepoint-migrate-write-test")
	prefixPath := filepath.Join(root, ".savepoint-migrate-write-test-occupied")
	sentinel := []byte("pre-existing sentinel bytes\x00\n")
	prefix := []byte("pre-existing prefix collision bytes\xff\n")
	if err := os.WriteFile(sentinelPath, sentinel, 0644); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}
	if err := os.WriteFile(prefixPath, prefix, 0644); err != nil {
		t.Fatalf("write prefix collision: %v", err)
	}

	var output strings.Builder
	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Write:          true,
		Stdout:         &output,
		Now:            fixedClock(time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)),
		NewOperationID: fixedOperationID("probe-collision"),
	})
	if err != nil || code != 0 {
		t.Fatalf("apply with pre-existing probe names: code = %d, err = %v\n%s", code, err, output.String())
	}

	for path, want := range map[string][]byte{
		sentinelPath: sentinel,
		prefixPath:   prefix,
	} {
		got, readErr := os.ReadFile(path)
		if readErr != nil {
			t.Fatalf("read preserved probe collision %s: %v", path, readErr)
		}
		if !bytes.Equal(got, want) {
			t.Fatalf("probe collision %s changed: got %q, want %q", path, got, want)
		}
	}

	matches, err := filepath.Glob(filepath.Join(root, ".savepoint-migrate-write-test-*"))
	if err != nil {
		t.Fatalf("glob probe files: %v", err)
	}
	if len(matches) != 1 || matches[0] != prefixPath {
		t.Fatalf("probe files after apply = %v, want only pre-existing %s", matches, prefixPath)
	}
}

func TestRunCommand_recoverApplyAcrossInvocationsUsesRecordedOperation(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	firstAt := time.Date(2026, 9, 19, 10, 11, 12, 0, time.UTC)
	interruptFirstCommand(t, root, firstAt, "op-recorded")

	pending, err := PendingOperation(root)
	if err != nil || pending == nil {
		t.Fatalf("PendingOperation() = %+v, %v, want an interrupted operation", pending, err)
	}
	if pending.Journal.Plan == nil {
		t.Fatal("pending journal has no original plan; recovery would need a fresh invocation's volatile state")
	}
	if pending.Journal.Plan.OperationID != "op-recorded" || !pending.Journal.Plan.GeneratedAt.Equal(firstAt) {
		t.Fatalf("recorded plan = operation %q at %v, want op-recorded at %v", pending.Journal.Plan.OperationID, pending.Journal.Plan.GeneratedAt, firstAt)
	}

	var output strings.Builder
	laterAt := firstAt.Add(24 * time.Hour)
	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Write:          true,
		Recover:        true,
		Stdout:         &output,
		Now:            fixedClock(laterAt),
		NewOperationID: fixedOperationID("op-fresh-must-not-be-used"),
	})
	if err != nil || code != 0 {
		t.Fatalf("separate Recover+Write invocation: code = %d, err = %v\n%s", code, err, output.String())
	}
	if !strings.Contains(output.String(), "resumed and completed migration operation op-recorded") {
		t.Fatalf("recovery output = %q, want the recorded operation ID", output.String())
	}

	manifestBytes, err := os.ReadFile(filepath.Join(root, ".savepoint", "migrations", "v1-to-v2.yml"))
	if err != nil {
		t.Fatalf("read migration manifest: %v", err)
	}
	manifest, err := UnmarshalManifest(manifestBytes)
	if err != nil {
		t.Fatalf("parse migration manifest: %v", err)
	}
	if manifest.OperationID != "op-recorded" || !manifest.GeneratedAt.Equal(firstAt) {
		t.Fatalf("manifest volatile fields = operation %q at %v, want op-recorded at %v", manifest.OperationID, manifest.GeneratedAt, firstAt)
	}
	if report, reportErr := PendingOperation(root); reportErr != nil || report != nil {
		t.Fatalf("PendingOperation() after recovery = %+v, %v, want none", report, reportErr)
	}
}

func TestRunCommand_recoverApplyAcrossInvocationsRejectsEditedSource(t *testing.T) {
	root := copyFixtureProject(t, "v1-basic")
	firstAt := time.Date(2026, 9, 19, 11, 0, 0, 0, time.UTC)
	interruptFirstCommand(t, root, firstAt, "op-edited-source")

	source := filepath.Join(root, ".savepoint", "releases", "v1", "epics", "E01-example", "tasks", "T002-follow-up.md")
	original, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read source before edit: %v", err)
	}
	edited := append(append([]byte{}, original...), []byte("\nuser edit after interruption\n")...)
	if err := os.WriteFile(source, edited, 0644); err != nil {
		t.Fatalf("edit source after interruption: %v", err)
	}

	var output strings.Builder
	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Write:          true,
		Recover:        true,
		Stdout:         &output,
		Now:            fixedClock(firstAt.Add(48 * time.Hour)),
		NewOperationID: fixedOperationID("op-must-not-be-used"),
	})
	if code == 0 || !errors.Is(err, ErrSourceConflict) {
		t.Fatalf("recovery after a source edit: code = %d, err = %v, want ErrSourceConflict", code, err)
	}
	got, err := os.ReadFile(source)
	if err != nil {
		t.Fatalf("read edited source after refusal: %v", err)
	}
	if string(got) != string(edited) {
		t.Fatal("source edit was not preserved after recovery refused it")
	}
}

func interruptFirstCommand(t *testing.T, root string, at time.Time, operationID string) {
	t.Helper()
	oldHook := afterPublishWriteHook
	calls := 0
	stop := errors.New("test interruption after first published output")
	afterPublishWriteHook = func(string) error {
		calls++
		if calls == 1 {
			return stop
		}
		return nil
	}
	defer func() { afterPublishWriteHook = oldHook }()

	var output strings.Builder
	code, err := RunCommand(CommandOptions{
		Dir:            root,
		Write:          true,
		Stdout:         &output,
		Now:            fixedClock(at),
		NewOperationID: fixedOperationID(operationID),
	})
	if code == 0 || !errors.Is(err, stop) {
		t.Fatalf("interrupted first invocation: code = %d, err = %v, want the injected interruption", code, err)
	}
}
