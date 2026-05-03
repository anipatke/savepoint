package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"path/filepath"
	"runtime"
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
}

var versionOverride string

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("buildtool", flag.ContinueOnError)
	flags.StringVar(&versionOverride, "version", "", "version to inject into the binary")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("usage: go run ./internal/buildtool [-version vX.Y.Z] <build|clean|build-linux|build-darwin|build-all|dist|smoke-test>")
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
	case "build-all":
		return buildAll()
	case "dist":
		return dist()
	case "smoke-test":
		return smokeTest()
	default:
		return fmt.Errorf("unknown build target %q", flags.Arg(0))
	}
}

func buildLocal() error {
	return runGoBuild(localExecutable(), runtime.GOOS, runtime.GOARCH)
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

func buildTarget(target target) error {
	output := filepath.Join("dist", target.os+"-"+target.arch, "savepoint")
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

func dist() error {
	if err := buildAll(); err != nil {
		return err
	}
	for _, target := range targets {
		name := "savepoint-" + version() + "-" + target.os + "-" + target.arch + ".tar.gz"
		source := filepath.Join("dist", target.os+"-"+target.arch, "savepoint")
		archive := filepath.Join("dist", name)
		if err := writeTarGz(archive, source, "savepoint"); err != nil {
			return err
		}
	}
	return nil
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

