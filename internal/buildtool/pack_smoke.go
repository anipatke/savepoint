package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const smokeProjectManifestJSON = `{"name":"savepoint-pack-smoke","version":"0.0.0","private":true}` + "\n"

func packSmoke() error {
	if err := buildNPM(); err != nil {
		return err
	}

	host, err := hostTarget()
	if err != nil {
		return err
	}

	workDir, err := os.MkdirTemp("", "savepoint-pack-smoke-*")
	if err != nil {
		return fmt.Errorf("create work dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	tarballDir := filepath.Join(workDir, "tarballs")
	if err := os.MkdirAll(tarballDir, 0o755); err != nil {
		return fmt.Errorf("create tarball dir: %w", err)
	}

	rootTarball, err := npmPack(".", tarballDir)
	if err != nil {
		return fmt.Errorf("pack root package: %w", err)
	}
	platformDir := filepath.Join(npmDistDir, host.os+"-"+host.arch)
	platformTarball, err := npmPack(platformDir, tarballDir)
	if err != nil {
		return fmt.Errorf("pack %s/%s package: %w", host.os, host.arch, err)
	}

	installDir := filepath.Join(workDir, "install")
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return fmt.Errorf("create install dir: %w", err)
	}
	if err := os.WriteFile(filepath.Join(installDir, "package.json"), []byte(smokeProjectManifestJSON), 0o644); err != nil {
		return fmt.Errorf("write smoke project manifest: %w", err)
	}

	if err := runNpm(installDir, "install", "--no-audit", "--no-fund", "--silent", platformTarball, rootTarball); err != nil {
		return fmt.Errorf("npm install tarballs: %w", err)
	}

	bin := installedBinaryPath(installDir, host.os)
	if err := runCmd(installDir, bin, "--version"); err != nil {
		return fmt.Errorf("savepoint --version smoke: %w", err)
	}

	fixtureDir := filepath.Join(workDir, "fixture")
	if err := os.MkdirAll(fixtureDir, 0o755); err != nil {
		return fmt.Errorf("create fixture dir: %w", err)
	}
	if err := runCmd(installDir, bin, "init", fixtureDir); err != nil {
		return fmt.Errorf("savepoint init fixture: %w", err)
	}
	if err := runCmd(installDir, bin, "upgrade-assets", fixtureDir, "--dry-run"); err != nil {
		return fmt.Errorf("savepoint upgrade-assets --dry-run smoke: %w", err)
	}

	fmt.Println("pack-smoke passed")
	return nil
}

func hostTarget() (target, error) {
	want := target{os: runtime.GOOS, arch: runtime.GOARCH}
	for _, t := range targets {
		if t == want {
			return t, nil
		}
	}
	return target{}, fmt.Errorf("host %s/%s not in supported targets", want.os, want.arch)
}

func installedBinaryPath(installDir, goos string) string {
	if goos == "windows" {
		return filepath.Join(installDir, "node_modules", ".bin", "savepoint.cmd")
	}
	return filepath.Join(installDir, "node_modules", ".bin", "savepoint")
}

func npmPack(packageDir, destDir string) (string, error) {
	absDest, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("resolve destination: %w", err)
	}
	cmd := exec.Command(npmExecutable(), "pack", "--json", "--pack-destination", absDest)
	cmd.Dir = packageDir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("npm pack %s: %w", packageDir, err)
	}
	filename, err := parseNpmPackFilename(stdout.Bytes())
	if err != nil {
		return "", fmt.Errorf("parse npm pack output for %s: %w", packageDir, err)
	}
	return filepath.Join(absDest, filename), nil
}

func parseNpmPackFilename(output []byte) (string, error) {
	trimmed := bytes.TrimSpace(output)
	if len(trimmed) == 0 {
		return "", fmt.Errorf("empty npm pack output")
	}
	var entries []struct {
		Filename string `json:"filename"`
	}
	if err := json.Unmarshal(trimmed, &entries); err != nil {
		return "", fmt.Errorf("unmarshal npm pack json: %w", err)
	}
	if len(entries) == 0 || entries[0].Filename == "" {
		return "", fmt.Errorf("npm pack json missing filename")
	}
	return entries[0].Filename, nil
}

func runNpm(dir string, args ...string) error {
	return runCmd(dir, npmExecutable(), args...)
}

func npmExecutable() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}

func runCmd(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
