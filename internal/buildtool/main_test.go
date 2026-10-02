package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestVersion_override(t *testing.T) {
	versionOverride = "v1.2.3"
	defer func() { versionOverride = "" }()
	if got := version(); got != "v1.2.3" {
		t.Errorf("version() = %q, want %q", got, "v1.2.3")
	}
}

func TestVersion_env(t *testing.T) {
	versionOverride = ""
	os.Setenv("VERSION", "v2.0.0-env")
	defer os.Unsetenv("VERSION")
	if got := version(); got != "v2.0.0-env" {
		t.Errorf("version() = %q, want %q", got, "v2.0.0-env")
	}
}

func TestVersion_fallback(t *testing.T) {
	versionOverride = ""
	os.Unsetenv("VERSION")
	got := version()
	if got == "" {
		t.Error("version() returned empty string")
	}
}

func TestWriteChecksums(t *testing.T) {
	dir := t.TempDir()

	content := []byte("fake archive content")
	archive := filepath.Join(dir, "savepoint-v1.0.0-linux-amd64.tar.gz")
	if err := os.WriteFile(archive, content, 0o644); err != nil {
		t.Fatal(err)
	}

	dest := filepath.Join(dir, "checksums.txt")
	if err := writeChecksums(dest, []string{archive}); err != nil {
		t.Fatalf("writeChecksums: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}

	h := sha256.Sum256(content)
	wantHash := hex.EncodeToString(h[:])
	wantLine := wantHash + "  savepoint-v1.0.0-linux-amd64.tar.gz"

	lines := strings.Split(strings.TrimSpace(string(got)), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if lines[0] != wantLine {
		t.Errorf("line = %q, want %q", lines[0], wantLine)
	}
}

func TestWriteChecksums_multiple(t *testing.T) {
	dir := t.TempDir()

	names := []string{"a.tar.gz", "b.tar.gz"}
	var paths []string
	for _, name := range names {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, []byte(name), 0o644); err != nil {
			t.Fatal(err)
		}
		paths = append(paths, p)
	}

	dest := filepath.Join(dir, "checksums.txt")
	if err := writeChecksums(dest, paths); err != nil {
		t.Fatalf("writeChecksums: %v", err)
	}

	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(got)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %s", len(lines), got)
	}
	for i, name := range names {
		h := sha256.Sum256([]byte(name))
		want := hex.EncodeToString(h[:]) + "  " + name
		if lines[i] != want {
			t.Errorf("line[%d] = %q, want %q", i, lines[i], want)
		}
	}
}

func TestWriteChecksums_missingFile(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "checksums.txt")
	err := writeChecksums(dest, []string{filepath.Join(dir, "nonexistent.tar.gz")})
	if err == nil {
		t.Error("expected error for missing archive, got nil")
	}
}

func TestExpectedArchiveNamesCoverSixPlatformMatrix(t *testing.T) {
	want := []string{
		"savepoint-v1.0.0-linux-amd64.tar.gz",
		"savepoint-v1.0.0-linux-arm64.tar.gz",
		"savepoint-v1.0.0-darwin-amd64.tar.gz",
		"savepoint-v1.0.0-darwin-arm64.tar.gz",
		"savepoint-v1.0.0-windows-amd64.tar.gz",
		"savepoint-v1.0.0-windows-arm64.tar.gz",
	}
	got := expectedArchiveNames("v1.0.0")
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Fatalf("expected archive names = %v, want %v", got, want)
	}
}

func TestVerifyDistributionStructureAcceptsExactSixArchives(t *testing.T) {
	dir, release := fakeDistribution(t)
	if err := verifyDistributionStructure(dir, release); err != nil {
		t.Fatalf("verifyDistributionStructure: %v", err)
	}
}

func TestDistributionInventoryRejectsMissingExtraRenamedAndChanged(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(t *testing.T, dir, release string)
	}{
		{name: "missing", mutate: func(t *testing.T, dir, release string) {
			name := expectedArchiveNames(release)[0]
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "extra", mutate: func(t *testing.T, dir, release string) {
			if err := os.WriteFile(filepath.Join(dir, "renamed.tar.gz"), []byte("extra"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "unsupported", mutate: func(t *testing.T, dir, release string) {
			if err := os.WriteFile(filepath.Join(dir, "savepoint-"+release+"-freebsd-amd64.tar.gz"), []byte("extra"), 0o644); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "renamed", mutate: func(t *testing.T, dir, release string) {
			name := expectedArchiveNames(release)[0]
			if err := os.Rename(filepath.Join(dir, name), filepath.Join(dir, "renamed.tar.gz")); err != nil {
				t.Fatal(err)
			}
		}},
		{name: "changed", mutate: func(t *testing.T, dir, release string) {
			name := expectedArchiveNames(release)[0]
			path := filepath.Join(dir, name)
			f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.Write([]byte("changed")); err != nil {
				f.Close()
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}
		}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir, release := fakeDistribution(t)
			tc.mutate(t, dir, release)
			if err := verifyDistributionStructure(dir, release); err == nil {
				t.Fatal("verifyDistributionStructure succeeded for invalid inventory")
			}
		})
	}
}

func TestTargets_includesWindows(t *testing.T) {
	var gotAMD64, gotARM64 bool
	for _, tgt := range targets {
		if tgt.os != "windows" {
			continue
		}
		switch tgt.arch {
		case "amd64":
			gotAMD64 = true
		case "arm64":
			gotARM64 = true
		}
	}
	if !gotAMD64 {
		t.Error("targets missing windows/amd64")
	}
	if !gotARM64 {
		t.Error("targets missing windows/arm64")
	}
}

func TestTargets_preservesLinuxDarwin(t *testing.T) {
	want := map[string]bool{
		"linux/amd64":  false,
		"linux/arm64":  false,
		"darwin/amd64": false,
		"darwin/arm64": false,
	}
	for _, tgt := range targets {
		key := tgt.os + "/" + tgt.arch
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for key, found := range want {
		if !found {
			t.Errorf("targets missing %s", key)
		}
	}
}

func TestExecutableName(t *testing.T) {
	if got := executableName("windows"); got != "savepoint.exe" {
		t.Errorf("executableName(windows) = %q, want savepoint.exe", got)
	}
	if got := executableName("linux"); got != "savepoint" {
		t.Errorf("executableName(linux) = %q, want savepoint", got)
	}
	if got := executableName("darwin"); got != "savepoint" {
		t.Errorf("executableName(darwin) = %q, want savepoint", got)
	}
}

func TestWriteTarGzPreservesWindowsExecutableName(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "savepoint.exe")
	if err := os.WriteFile(source, []byte("binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	archive := filepath.Join(dir, "savepoint-windows-amd64.tar.gz")
	if err := writeTarGz(archive, source, executableName("windows")); err != nil {
		t.Fatalf("writeTarGz: %v", err)
	}

	f, err := os.Open(archive)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	header, err := tr.Next()
	if err != nil {
		t.Fatal(err)
	}
	if header.Name != "savepoint.exe" {
		t.Fatalf("archive member = %q, want savepoint.exe", header.Name)
	}

	content, err := io.ReadAll(tr)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "binary" {
		t.Fatalf("archive content = %q, want binary", content)
	}
}

func TestRequireWindowsExecutableAcceptsMZHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "savepoint.exe")
	if err := os.WriteFile(path, []byte("MZfake-pe"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := requireWindowsExecutable(path); err != nil {
		t.Fatalf("requireWindowsExecutable: %v", err)
	}
}

func TestRequireWindowsExecutableRejectsELFHeader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "savepoint.exe")
	if err := os.WriteFile(path, []byte{0x7f, 'E', 'L', 'F'}, 0o755); err != nil {
		t.Fatal(err)
	}

	err := requireWindowsExecutable(path)
	if err == nil {
		t.Fatal("expected ELF header to be rejected")
	}
	if !strings.Contains(err.Error(), "not a Windows PE binary") {
		t.Fatalf("error = %q, want Windows PE binary message", err)
	}
}

func TestLocalExecutable(t *testing.T) {
	got := localExecutable()
	if runtime.GOOS == "windows" {
		if got != "savepoint.exe" {
			t.Errorf("localExecutable() = %q, want %q", got, "savepoint.exe")
		}
	} else {
		if got != "savepoint" {
			t.Errorf("localExecutable() = %q, want %q", got, "savepoint")
		}
	}
}

func TestFocusedTestArgsRequiresPattern(t *testing.T) {
	if _, err := focusedTestArgs(nil); err == nil {
		t.Fatal("expected an empty focused-test pattern to fail")
	}
	if _, err := focusedTestArgs([]string{"   "}); err == nil {
		t.Fatal("expected a whitespace focused-test pattern to fail")
	}
}

func TestFocusedTestArgsAddsUncachedTimingAndPackage(t *testing.T) {
	got, err := focusedTestArgs([]string{"TestResume", "./internal/resume"})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-json", "-count=1", "-run", "TestResume", "./internal/resume"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("focusedTestArgs() = %#v, want %#v", got, want)
	}
}

func TestRunGoTestCommandReportsTimingAndPropagatesFailure(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=^TestGoTestCommandHelperProcess$")
	cmd.Env = append(os.Environ(), "SAVEPOINT_BUILDTOOL_TEST_HELPER=1")
	var output bytes.Buffer
	err := runGoTestCommand(cmd, &output)
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("runGoTestCommand() error = %v, want child exit error", err)
	}
	if got := exitErr.ExitCode(); got != 7 {
		t.Fatalf("child exit code = %d, want 7", got)
	}
	for _, want := range []string{"internal/migrate.TestSlowFixture", "internal/migrate.TestBroken", "injected diagnostic", "Go test timing summary", "Slowest packages", "Slowest tests"} {
		if !strings.Contains(output.String(), want) {
			t.Errorf("output does not contain %q:\n%s", want, output.String())
		}
	}
}

func TestGoTestCommandHelperProcess(t *testing.T) {
	if os.Getenv("SAVEPOINT_BUILDTOOL_TEST_HELPER") != "1" {
		return
	}
	fmt.Println(`{"Action":"pass","Package":"internal/migrate","Test":"TestSlowFixture","Elapsed":0.35}`)
	fmt.Println(`{"Action":"fail","Package":"internal/migrate","Test":"TestBroken","Elapsed":0.5}`)
	fmt.Println(`{"Action":"fail","Package":"internal/migrate","Elapsed":0.8}`)
	fmt.Fprintln(os.Stderr, "injected diagnostic")
	os.Exit(7)
}
func fakeDistribution(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	sourceDir := filepath.Join(dir, "sources")
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		t.Fatal(err)
	}
	release := "v1.0.0"
	archives := make([]string, 0, len(targets))
	for _, target := range targets {
		content := []byte{0x7f, 'E', 'L', 'F', 'f', 'a', 'k', 'e'}
		if target.os == "windows" {
			content = []byte{'M', 'Z', 'f', 'a', 'k', 'e'}
		} else if target.os == "darwin" {
			content = []byte{0xcf, 0xfa, 0xed, 0xfe, 'f', 'a', 'k', 'e'}
		}
		source := filepath.Join(sourceDir, target.os+"-"+target.arch)
		if err := os.WriteFile(source, content, 0o755); err != nil {
			t.Fatal(err)
		}
		archive := filepath.Join(dir, archiveName(release, target))
		if err := writeTarGz(archive, source, executableName(target.os)); err != nil {
			t.Fatal(err)
		}
		archives = append(archives, archive)
	}
	if err := writeDistributionChecksums(dir, release, archives); err != nil {
		t.Fatal(err)
	}
	return dir, release
}

func TestRunGoTestWithReportsWritesReportsOnPass(t *testing.T) {
	dir := t.TempDir()
	writeFixtureModule(t, dir, "package fixture\n\nfunc F() int { return 1 }\n", "package fixture\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) { _ = F() }\n")
	var output bytes.Buffer
	t.Chdir(dir)
	if err := runGoTestWithReports(dir, []string{"-json", "-count=1", "./..."}, &output); err != nil {
		t.Fatalf("runGoTestWithReports() error = %v\n%s", err, output.String())
	}
	assertReports(t, dir, true)
	if !strings.Contains(output.String(), "Go test timing summary") {
		t.Errorf("summary missing:\n%s", output.String())
	}
}

func TestRunGoTestWithReportsWritesReportsOnFailure(t *testing.T) {
	dir := t.TempDir()
	writeFixtureModule(t, dir, "package fixture\n\nfunc F() int { return 1 }\n", "package fixture\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) { _ = F(); t.Fatal(\"boom\") }\n")
	var output bytes.Buffer
	t.Chdir(dir)
	err := runGoTestWithReports(dir, []string{"-json", "-count=1", "./..."}, &output)
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("error = %v, want child exit error", err)
	}
	assertReports(t, dir, true)
	data, _ := os.ReadFile(filepath.Join(dir, goTestReportName))
	if !strings.Contains(string(data), `"Action":"fail"`) {
		t.Errorf("report lacks failure event:\n%s", data)
	}
	if !strings.Contains(output.String(), "Go test timing summary") {
		t.Errorf("summary missing:\n%s", output.String())
	}
}

func TestRunGoTestFocusedWritesNoReports(t *testing.T) {
	dir := t.TempDir()
	writeFixtureModule(t, dir, "package fixture\n\nfunc F() int { return 1 }\n", "package fixture\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) { _ = F() }\n")
	t.Chdir(dir)
	var output bytes.Buffer
	args, err := focusedTestArgs([]string{"TestF"})
	if err != nil {
		t.Fatal(err)
	}
	if err := runGoTest(args, &output); err != nil {
		t.Fatalf("runGoTest() error = %v\n%s", err, output.String())
	}
	assertReports(t, dir, false)
}

func TestRunGoTestWithReportsLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	writeFixtureModule(t, dir, "package fixture\n", "package fixture\n\nimport \"testing\"\n\nfunc TestF(t *testing.T) {}\n")
	t.Chdir(dir)
	var output bytes.Buffer
	_ = runGoTestWithReports(dir, []string{"-json", "-count=1", "./..."}, &output)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			t.Errorf("temp file left behind: %s", entry.Name())
		}
	}
}

func writeFixtureModule(t *testing.T, dir, source, test string) {
	t.Helper()
	for name, content := range map[string]string{
		"go.mod":          "module fixture\n\ngo 1.21\n",
		"fixture.go":      source,
		"fixture_test.go": test,
	} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func assertReports(t *testing.T, dir string, want bool) {
	t.Helper()
	for _, name := range []string{goTestReportName, goCoverReportName} {
		info, err := os.Stat(filepath.Join(dir, name))
		if want && (err != nil || info.Size() == 0) {
			t.Errorf("%s missing or empty: %v", name, err)
		}
		if !want && err == nil {
			t.Errorf("%s written, want none", name)
		}
	}
}

func TestRunGoTestWithReportsKeepsCompleteReportsWhenInterrupted(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses a shell script as a fake go command")
	}
	dir := t.TempDir()
	bin := t.TempDir()
	script := "#!/bin/sh\necho '{\"Action\":\"start\",\"Package\":\"fixture\"}'\nkill -TERM $$\n"
	if err := os.WriteFile(filepath.Join(bin, "go"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	previous := map[string]string{goTestReportName: "complete test stream\n", goCoverReportName: "mode: set\n"}
	for name, content := range previous {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var output bytes.Buffer
	if err := runGoTestWithReports(dir, []string{"-json", "./..."}, &output); err == nil {
		t.Fatal("an interrupted run must return an error")
	}
	for name, want := range previous {
		got, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || string(got) != want {
			t.Errorf("%s = %q (%v), want the previous complete report %q", name, got, err, want)
		}
	}
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			t.Errorf("temp file left behind: %s", entry.Name())
		}
	}
}

// TestMain lets the test binary stand in for the go command: copied onto PATH
// as go, it prints a start event and exits nonzero without a package result,
// like a Windows child ended by TerminateProcess.
func TestMain(m *testing.M) {
	if os.Getenv("BUILDTOOL_FAKE_GO_INTERRUPTED") == "1" {
		fmt.Println(`{"Action":"start","Package":"fixture"}`)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestRunGoTestWithReportsKeepsCompleteReportsWhenChildExitsEarly(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		t.Fatal(err)
	}
	bin := t.TempDir()
	name := "go"
	if runtime.GOOS == "windows" {
		name = "go.exe"
	}
	if err := os.WriteFile(filepath.Join(bin, name), data, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("BUILDTOOL_FAKE_GO_INTERRUPTED", "1")
	dir := t.TempDir()
	for _, report := range []string{goTestReportName, goCoverReportName} {
		if err := os.WriteFile(filepath.Join(dir, report), []byte("previous\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := runGoTestWithReports(dir, []string{"-json", "./..."}, io.Discard); err == nil {
		t.Fatal("an early-exiting run must return an error")
	}
	for _, report := range []string{goTestReportName, goCoverReportName} {
		if got, _ := os.ReadFile(filepath.Join(dir, report)); string(got) != "previous\n" {
			t.Errorf("%s = %q, want the previous report", report, got)
		}
	}
}

func TestGoTestStreamComplete(t *testing.T) {
	cases := []struct {
		name, stream  string
		emptyOK, want bool
	}{
		{"finished package", `{"Action":"start","Package":"a"}` + "\n" + `{"Action":"fail","Package":"a"}` + "\n", false, true},
		{"start only", `{"Action":"start","Package":"a"}` + "\n", false, false},
		{"one package unfinished", `{"Action":"start","Package":"a"}` + "\n" + `{"Action":"start","Package":"b"}` + "\n" + `{"Action":"pass","Package":"a"}` + "\n", true, false},
		{"test result is not a package result", `{"Action":"start","Package":"a"}` + "\n" + `{"Action":"pass","Package":"a","Test":"T"}` + "\n", false, false},
		{"empty after failure", "", false, false},
		{"empty after success", "", true, true},
	}
	for _, c := range cases {
		path := filepath.Join(t.TempDir(), "stream")
		if err := os.WriteFile(path, []byte(c.stream), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := goTestStreamComplete(path, c.emptyOK); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}

func streamHelperCommand(mode string) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=^TestGoTestStreamHelperProcess$")
	cmd.Env = append(os.Environ(), "SAVEPOINT_BUILDTOOL_STREAM_HELPER="+mode)
	return cmd
}

func TestGoTestStreamHelperProcess(t *testing.T) {
	switch os.Getenv("SAVEPOINT_BUILDTOOL_STREAM_HELPER") {
	case "failed-package":
		fmt.Println(`not json at all`)
		fmt.Println(`{"Action":"output","Package":"fixture/b","Output":"b detail\n"}`)
		fmt.Println(`{"Action":"output","Package":"fixture/a","Output":"a detail\n"}`)
		fmt.Println(`{"Action":"fail","Package":"fixture/b","Elapsed":0.1}`)
		fmt.Println(`{"Action":"fail","Package":"fixture/a","Elapsed":0.1}`)
		os.Exit(1)
	case "oversized":
		fmt.Println(`{"Action":"start","Package":"fixture"}`)
		fmt.Println(strings.Repeat("x", 11*1024*1024))
		os.Exit(0)
	case "ok":
		fmt.Println(`{"Action":"pass","Package":"fixture","Elapsed":0.1}`)
		os.Exit(0)
	}
}

func TestRunGoTestStreamPrintsMalformedLinesAndFailedPackageOutputInOrder(t *testing.T) {
	var output bytes.Buffer
	err := runGoTestStream(streamHelperCommand("failed-package"), &output, nil)
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		t.Fatalf("runGoTestStream() error = %v, want child exit 1", err)
	}
	got := output.String()
	malformed := strings.Index(got, "not json at all")
	a := strings.Index(got, "Output for failed package fixture/a:\na detail")
	b := strings.Index(got, "Output for failed package fixture/b:\nb detail")
	if malformed < 0 || a < 0 || b < 0 || !(malformed < a && a < b) {
		t.Fatalf("output order wrong (malformed=%d a=%d b=%d):\n%s", malformed, a, b, got)
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("disk full") }

func TestRunGoTestStreamReportsReportWriteFailureAfterChildSucceeds(t *testing.T) {
	var output bytes.Buffer
	err := runGoTestStream(streamHelperCommand("ok"), &output, failingWriter{})
	if err == nil || !strings.Contains(err.Error(), "write go test report: disk full") {
		t.Fatalf("runGoTestStream() error = %v, want report write failure", err)
	}
	if !strings.Contains(output.String(), "Go test timing summary") {
		t.Errorf("timing summary missing after write failure:\n%s", output.String())
	}
}

func TestRunGoTestStreamReportsOversizedEventLine(t *testing.T) {
	var output bytes.Buffer
	err := runGoTestStream(streamHelperCommand("oversized"), &output, nil)
	if err == nil || !strings.Contains(err.Error(), "read go test output") {
		t.Fatalf("runGoTestStream() error = %v, want read failure", err)
	}
}
