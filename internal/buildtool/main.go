package main

import (
	"archive/tar"
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type target struct {
	os   string
	arch string
}

var targets = []target{
	{os: "linux", arch: "amd64"},
	{os: "linux", arch: "arm64"},
	{os: "darwin", arch: "amd64"},
	{os: "darwin", arch: "arm64"},
	{os: "windows", arch: "amd64"},
	{os: "windows", arch: "arm64"},
}

var versionOverride string

const npmDistDir = "dist/npm"

type goTestEvent struct {
	Action  string  `json:"Action"`
	Package string  `json:"Package"`
	Test    string  `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string  `json:"Output"`
}

type testTiming struct {
	name    string
	elapsed time.Duration
}

func focusedTestArgs(args []string) ([]string, error) {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return nil, errors.New("usage: make test-focused TEST=<test-name-or-regexp> [PKGS=./package]")
	}
	packages := args[1:]
	if len(packages) == 0 {
		packages = []string{"./..."}
	}
	goArgs := []string{"-json", "-count=1", "-run", args[0]}
	return append(goArgs, packages...), nil
}

func runGoTest(args []string, output io.Writer) error {
	if len(args) == 0 {
		return errors.New("go test arguments are required")
	}
	commandArgs := append([]string{"test"}, args...)
	return runGoTestCommand(exec.Command("go", commandArgs...), output)
}

func runGoTestCommand(cmd *exec.Cmd, output io.Writer) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("open go test output: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start go test: %w", err)
	}

	var packages []testTiming
	var tests []testTiming
	var skipped []string
	var failedPackages []string
	packageOutput := make(map[string]*strings.Builder)
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 64*1024), 10*1024*1024)
	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		var event goTestEvent
		if err := json.Unmarshal(line, &event); err != nil {
			fmt.Fprintln(output, string(line))
			continue
		}
		if event.Output != "" {
			packageName := event.Package
			if packageName == "" {
				packageName = "<go test>"
			}
			if packageOutput[packageName] == nil {
				packageOutput[packageName] = &strings.Builder{}
			}
			packageOutput[packageName].WriteString(event.Output)
		}
		if event.Test == "" && event.Action == "fail" {
			failedPackages = append(failedPackages, event.Package)
		}
		if event.Test == "" && event.Elapsed > 0 && (event.Action == "pass" || event.Action == "fail") {
			packages = append(packages, testTiming{name: event.Package, elapsed: time.Duration(event.Elapsed * float64(time.Second))})
		}
		if event.Test != "" {
			if event.Action == "pass" || event.Action == "fail" {
				tests = append(tests, testTiming{name: event.Package + "." + event.Test, elapsed: time.Duration(event.Elapsed * float64(time.Second))})
			}
			if event.Action == "skip" {
				skipped = append(skipped, event.Package+"."+event.Test)
			}
		}
	}
	scanErr := scanner.Err()
	waitErr := cmd.Wait()
	sort.Strings(failedPackages)
	previousPackage := ""
	for _, packageName := range failedPackages {
		if packageName == previousPackage {
			continue
		}
		previousPackage = packageName
		fmt.Fprintf(output, "\nOutput for failed package %s:\n", packageName)
		if packageOutput[packageName] != nil {
			if _, err := io.WriteString(output, packageOutput[packageName].String()); err != nil {
				return fmt.Errorf("write failed package output: %w", err)
			}
		}
	}
	if stderr.Len() > 0 {
		if _, err := io.Copy(output, &stderr); err != nil {
			return fmt.Errorf("write go test diagnostics: %w", err)
		}
	}
	writeTestTimingSummary(output, packages, tests, skipped)
	if waitErr != nil {
		return waitErr
	}
	if scanErr != nil {
		return fmt.Errorf("read go test output: %w", scanErr)
	}
	return nil
}

func writeTestTimingSummary(output io.Writer, packages, tests []testTiming, skipped []string) {
	fmt.Fprintln(output, "\nGo test timing summary:")
	writeSlowest(output, "packages", packages)
	writeSlowest(output, "tests", tests)
	if len(skipped) > 0 {
		sort.Strings(skipped)
		fmt.Fprintf(output, "Skipped tests (%d): %s\n", len(skipped), strings.Join(skipped, ", "))
	}
}

func writeSlowest(output io.Writer, label string, timings []testTiming) {
	sort.Slice(timings, func(i, j int) bool {
		if timings[i].elapsed == timings[j].elapsed {
			return timings[i].name < timings[j].name
		}
		return timings[i].elapsed > timings[j].elapsed
	})
	limit := len(timings)
	if limit > 10 {
		limit = 10
	}
	fmt.Fprintf(output, "  Slowest %s:\n", label)
	for _, timing := range timings[:limit] {
		fmt.Fprintf(output, "    %-70s %s\n", timing.name, timing.elapsed.Round(time.Millisecond))
	}
}
func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() > 0 {
			os.Exit(exitErr.ExitCode())
		}
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "test":
			return runGoTest(args[1:], os.Stdout)
		case "focused-test":
			focusedArgs, err := focusedTestArgs(args[1:])
			if err != nil {
				return err
			}
			return runGoTest(focusedArgs, os.Stdout)
		}
	}
	flags := flag.NewFlagSet("buildtool", flag.ContinueOnError)
	flags.StringVar(&versionOverride, "version", "", "version to inject into the binary")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: go run ./internal/buildtool [-version vX.Y.Z] <build|clean|build-linux|build-darwin|build-windows|build-npm|build-all|dist|verify-dist|smoke-test|test|focused-test>")
	}

	switch flags.Arg(0) {
	case "build":
		return buildLocal()
	case "clean":
		return clean()
	case "build-linux":
		return buildMatching("linux")
	case "build-darwin":
		return buildMatching("darwin")
	case "build-windows":
		return buildMatching("windows")
	case "build-npm":
		return buildNPM()
	case "build-all":
		return buildAll()
	case "dist":
		return dist()
	case "verify-dist":
		return verifyDistribution("dist", version())
	case "smoke-test":
		return smokeTest()
	default:
		return fmt.Errorf("unknown build target %q", flags.Arg(0))
	}
}

// buildLocal also refreshes the host's npm launcher binary, so `npx savepoint`
// inside this repository cannot run an older build than ./savepoint (I-031).
func buildLocal() error {
	if err := runGoBuild(localExecutable(), runtime.GOOS, runtime.GOARCH); err != nil {
		return err
	}
	for _, target := range targets {
		if target.os == runtime.GOOS && target.arch == runtime.GOARCH {
			output := filepath.Join(npmDistDir, target.os+"-"+target.arch, executableName(target.os))
			return runGoBuild(output, target.os, target.arch)
		}
	}
	return nil
}

func buildNPM() error {
	if err := os.RemoveAll(npmDistDir); err != nil {
		return fmt.Errorf("clean npm dist: %w", err)
	}
	for _, target := range targets {
		output := filepath.Join(npmDistDir, target.os+"-"+target.arch, executableName(target.os))
		if err := runGoBuild(output, target.os, target.arch); err != nil {
			return err
		}
		if target.os == "windows" {
			if err := requireWindowsExecutable(output); err != nil {
				return err
			}
		}
	}
	return nil
}

func clean() error {
	for _, path := range []string{"savepoint", "savepoint.exe", "dist"} {
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("clean %s: %w", path, err)
		}
	}
	return nil
}

func buildMatching(goos string) error {
	for _, target := range targets {
		if target.os != goos {
			continue
		}
		if err := buildTarget(target); err != nil {
			return err
		}
	}
	return nil
}

func buildAll() error {
	for _, target := range targets {
		if err := buildTarget(target); err != nil {
			return err
		}
	}
	return nil
}

func executableName(goos string) string {
	if goos == "windows" {
		return "savepoint.exe"
	}
	return "savepoint"
}

func buildTarget(target target) error {
	output := filepath.Join("dist", target.os+"-"+target.arch, executableName(target.os))
	return runGoBuild(output, target.os, target.arch)
}

func runGoBuild(output, goos, goarch string) error {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil && filepath.Dir(output) != "." {
		return fmt.Errorf("create output dir: %w", err)
	}

	cmd := exec.Command("go", "build", "-ldflags", "-X main.version="+version(), "-o", output, "main.go")
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build %s/%s: %w", goos, goarch, err)
	}
	return nil
}

func requireWindowsExecutable(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open npm executable: %w", err)
	}
	defer f.Close()

	header := make([]byte, 2)
	if _, err := io.ReadFull(f, header); err != nil {
		return fmt.Errorf("read npm executable header: %w", err)
	}
	if string(header) != "MZ" {
		return fmt.Errorf("npm executable %s is not a Windows PE binary", path)
	}
	return nil
}

func dist() error {
	releaseVersion := version()
	if err := buildAll(); err != nil {
		return err
	}
	var archives []string
	for _, target := range targets {
		name := archiveName(releaseVersion, target)
		source := filepath.Join("dist", target.os+"-"+target.arch, executableName(target.os))
		archive := filepath.Join("dist", name)
		if err := writeTarGz(archive, source, executableName(target.os)); err != nil {
			return err
		}
		archives = append(archives, archive)
	}
	if err := writeDistributionChecksums("dist", releaseVersion, archives); err != nil {
		return err
	}
	return verifyDistribution("dist", releaseVersion)
}

func writeChecksums(dest string, archives []string) error {
	var lines strings.Builder
	for _, path := range archives {
		f, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("checksum open %s: %w", path, err)
		}
		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			f.Close()
			return fmt.Errorf("checksum read %s: %w", path, err)
		}
		f.Close()
		lines.WriteString(hex.EncodeToString(h.Sum(nil)))
		lines.WriteString("  ")
		lines.WriteString(filepath.Base(path))
		lines.WriteString("\n")
	}
	if err := os.WriteFile(dest, []byte(lines.String()), 0o644); err != nil {
		return fmt.Errorf("write checksums: %w", err)
	}
	return nil
}

func archiveName(releaseVersion string, target target) string {
	return fmt.Sprintf("savepoint-%s-%s-%s.tar.gz", releaseVersion, target.os, target.arch)
}

func expectedArchiveNames(releaseVersion string) []string {
	names := make([]string, 0, len(targets))
	for _, target := range targets {
		names = append(names, archiveName(releaseVersion, target))
	}
	return names
}

func writeDistributionChecksums(distDir, releaseVersion string, archives []string) error {
	if err := validateArchiveInventory(distDir, releaseVersion); err != nil {
		return err
	}

	expected := make(map[string]struct{}, len(targets))
	for _, name := range expectedArchiveNames(releaseVersion) {
		expected[name] = struct{}{}
	}
	if len(archives) != len(expected) {
		return fmt.Errorf("distribution checksum inventory has %d archives, want %d", len(archives), len(expected))
	}
	seen := make(map[string]struct{}, len(archives))
	for _, path := range archives {
		name := filepath.Base(path)
		if _, ok := expected[name]; !ok {
			return fmt.Errorf("distribution checksum archive %q is not an expected release archive", name)
		}
		if _, ok := seen[name]; ok {
			return fmt.Errorf("distribution checksum archive %q is listed more than once", name)
		}
		seen[name] = struct{}{}
	}
	if len(seen) != len(expected) {
		return fmt.Errorf("distribution checksum inventory is incomplete")
	}

	return writeChecksums(filepath.Join(distDir, "checksums.txt"), archives)
}

func validateArchiveInventory(distDir, releaseVersion string) error {
	expectedNames := expectedArchiveNames(releaseVersion)
	expected := make(map[string]struct{}, len(expectedNames))
	for _, name := range expectedNames {
		expected[name] = struct{}{}
	}

	entries, err := os.ReadDir(distDir)
	if err != nil {
		return fmt.Errorf("read distribution directory: %w", err)
	}
	seen := make(map[string]struct{}, len(expectedNames))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tar.gz") {
			continue
		}
		if _, ok := expected[entry.Name()]; !ok {
			return fmt.Errorf("unexpected distribution archive %q", entry.Name())
		}
		seen[entry.Name()] = struct{}{}
	}
	for _, name := range expectedNames {
		if _, ok := seen[name]; !ok {
			return fmt.Errorf("missing distribution archive %q", name)
		}
	}
	return nil
}

func verifyDistribution(distDir, releaseVersion string) error {
	if err := verifyDistributionStructure(distDir, releaseVersion); err != nil {
		return err
	}
	if err := smokeNativeDistribution(distDir, releaseVersion); err != nil {
		return err
	}
	fmt.Printf("distribution verified: %d archives (%s)\n", len(targets), releaseVersion)
	return nil
}

func verifyDistributionStructure(distDir, releaseVersion string) error {
	if err := validateArchiveInventory(distDir, releaseVersion); err != nil {
		return err
	}
	if err := verifyChecksumManifest(distDir, releaseVersion); err != nil {
		return err
	}
	for _, target := range targets {
		path := filepath.Join(distDir, archiveName(releaseVersion, target))
		if _, _, err := readArchiveBinary(path, target); err != nil {
			return err
		}
	}
	return nil
}

func verifyChecksumManifest(distDir, releaseVersion string) error {
	manifestPath := filepath.Join(distDir, "checksums.txt")
	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("read checksum manifest: %w", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(contents), "\n"), "\n")
	if len(lines) == 1 && strings.TrimSpace(lines[0]) == "" {
		return errors.New("checksum manifest is empty")
	}

	expected := make(map[string]struct{}, len(targets))
	for _, name := range expectedArchiveNames(releaseVersion) {
		expected[name] = struct{}{}
	}
	seen := make(map[string]struct{}, len(expected))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 || len(fields[0]) != sha256.Size*2 {
			return fmt.Errorf("invalid checksum manifest line %q", line)
		}
		if _, err := hex.DecodeString(fields[0]); err != nil {
			return fmt.Errorf("invalid checksum for %q: %w", fields[1], err)
		}
		if _, ok := expected[fields[1]]; !ok {
			return fmt.Errorf("checksum manifest contains unexpected archive %q", fields[1])
		}
		if _, ok := seen[fields[1]]; ok {
			return fmt.Errorf("checksum manifest contains duplicate archive %q", fields[1])
		}
		seen[fields[1]] = struct{}{}

		actual, err := checksumFile(filepath.Join(distDir, fields[1]))
		if err != nil {
			return err
		}
		if !strings.EqualFold(fields[0], actual) {
			return fmt.Errorf("checksum mismatch for %q", fields[1])
		}
	}
	if len(seen) != len(expected) {
		return fmt.Errorf("checksum manifest has %d archives, want %d", len(seen), len(expected))
	}
	return nil
}

func checksumFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("checksum open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("checksum read %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func readArchiveBinary(archivePath string, target target) ([]byte, int64, error) {
	archive, err := os.Open(archivePath)
	if err != nil {
		return nil, 0, fmt.Errorf("open distribution archive %s: %w", archivePath, err)
	}
	defer archive.Close()

	gzipReader, err := gzip.NewReader(archive)
	if err != nil {
		return nil, 0, fmt.Errorf("open distribution gzip %s: %w", archivePath, err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	memberName := executableName(target.os)
	var content []byte
	var mode int64
	memberCount := 0
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("read distribution archive %s: %w", archivePath, err)
		}
		memberCount++
		if memberCount > 1 {
			return nil, 0, fmt.Errorf("distribution archive %s contains more than one member", archivePath)
		}
		if header.Name != memberName {
			return nil, 0, fmt.Errorf("distribution archive %s contains %q, want %q", archivePath, header.Name, memberName)
		}
		if header.Mode&0o111 == 0 {
			return nil, 0, fmt.Errorf("distribution archive %s member is not executable", archivePath)
		}
		content, err = io.ReadAll(tarReader)
		if err != nil {
			return nil, 0, fmt.Errorf("read distribution member %s: %w", archivePath, err)
		}
		mode = header.Mode
	}
	if memberCount != 1 || len(content) == 0 {
		return nil, 0, fmt.Errorf("distribution archive %s does not contain one non-empty executable", archivePath)
	}
	if !hasExpectedBinaryHeader(content, target.os) {
		return nil, 0, fmt.Errorf("distribution archive %s member has an invalid %s executable header", archivePath, target.os)
	}
	return content, mode, nil
}

func hasExpectedBinaryHeader(content []byte, goos string) bool {
	switch goos {
	case "windows":
		return bytes.HasPrefix(content, []byte("MZ"))
	case "linux":
		return bytes.HasPrefix(content, []byte{0x7f, 'E', 'L', 'F'})
	case "darwin":
		if len(content) < 4 {
			return false
		}
		switch string(content[:4]) {
		case string([]byte{0xcf, 0xfa, 0xed, 0xfe}), // 64-bit little-endian
			string([]byte{0xce, 0xfa, 0xed, 0xfe}), // 32-bit little-endian
			string([]byte{0xfe, 0xed, 0xfa, 0xcf}), // 64-bit big-endian
			string([]byte{0xfe, 0xed, 0xfa, 0xce}), // 32-bit big-endian
			string([]byte{0xca, 0xfe, 0xba, 0xbe}), // universal binary
			string([]byte{0xbe, 0xba, 0xfe, 0xca}),
			string([]byte{0xca, 0xfe, 0xba, 0xbf}), // 64-bit universal binary
			string([]byte{0xbf, 0xba, 0xfe, 0xca}):
			return true
		}
	}
	return false
}

func smokeNativeDistribution(distDir, releaseVersion string) error {
	native, ok := nativeTarget()
	if !ok {
		return nil
	}
	archivePath := filepath.Join(distDir, archiveName(releaseVersion, native))
	content, mode, err := readArchiveBinary(archivePath, native)
	if err != nil {
		return err
	}
	tempDir, err := os.MkdirTemp("", "savepoint-dist-smoke-")
	if err != nil {
		return fmt.Errorf("create distribution smoke directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	binaryPath := filepath.Join(tempDir, executableName(native.os))
	fileMode := os.FileMode(mode).Perm()
	if err := os.WriteFile(binaryPath, content, fileMode); err != nil {
		return fmt.Errorf("write distribution smoke binary: %w", err)
	}
	cmd := exec.Command(binaryPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("launch native distribution archive %s: %w", archivePath, err)
	}
	if got := strings.TrimSpace(string(output)); got != releaseVersion {
		return fmt.Errorf("native distribution archive %s reported %q, want %q", archivePath, got, releaseVersion)
	}
	return nil
}

func nativeTarget() (target, bool) {
	for _, target := range targets {
		if target.os == runtime.GOOS && target.arch == runtime.GOARCH {
			return target, true
		}
	}
	return target{}, false
}

func writeTarGz(archivePath, sourcePath, archiveName string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open artifact source: %w", err)
	}
	defer source.Close()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("stat artifact source: %w", err)
	}

	archive, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("create archive: %w", err)
	}
	defer archive.Close()

	gzipWriter := gzip.NewWriter(archive)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return fmt.Errorf("create archive header: %w", err)
	}
	header.Name = archiveName
	// The member is always an executable. Windows reports no execute bits on
	// the built file, so the mode is set rather than copied from the host.
	header.Mode = 0o755
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("write archive header: %w", err)
	}
	if _, err := io.Copy(tarWriter, source); err != nil {
		return fmt.Errorf("write archive content: %w", err)
	}
	return nil
}

func smokeTest() error {
	if err := buildLocal(); err != nil {
		return err
	}
	cmd := exec.Command("."+string(os.PathSeparator)+localExecutable(), "--version")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("smoke test: %w", err)
	}
	fmt.Println("smoke test passed")
	return nil
}

func version() string {
	if versionOverride != "" {
		return versionOverride
	}
	if value := os.Getenv("VERSION"); value != "" {
		return value
	}

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		return strings.TrimSpace(string(output))
	}
	return "v0.0.0"
}

func localExecutable() string {
	if runtime.GOOS == "windows" {
		return "savepoint.exe"
	}
	return "savepoint"
}
