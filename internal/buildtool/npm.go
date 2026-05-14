package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var npmOSMap = map[string]string{
	"linux":   "linux",
	"darwin":  "darwin",
	"windows": "win32",
}

var npmCPUMap = map[string]string{
	"amd64": "x64",
	"arm64": "arm64",
}

type npmManifest struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	License     string   `json:"license"`
	Repository  string   `json:"repository,omitempty"`
	OS          []string `json:"os"`
	CPU         []string `json:"cpu"`
	Files       []string `json:"files"`
}

func platformPackageName(t target) string {
	return "savepoint-" + t.os + "-" + t.arch
}

func buildNPMManifest(t target, version string) (npmManifest, error) {
	npmOS, ok := npmOSMap[t.os]
	if !ok {
		return npmManifest{}, fmt.Errorf("npm os mapping missing for %q", t.os)
	}
	npmCPU, ok := npmCPUMap[t.arch]
	if !ok {
		return npmManifest{}, fmt.Errorf("npm cpu mapping missing for %q", t.arch)
	}
	return npmManifest{
		Name:        platformPackageName(t),
		Version:     version,
		Description: fmt.Sprintf("Savepoint native binary for %s/%s", t.os, t.arch),
		License:     "MIT",
		Repository:  "https://github.com/anipatke/savepoint",
		OS:          []string{npmOS},
		CPU:         []string{npmCPU},
		Files:       []string{"bin"},
	}, nil
}

func writeNPMManifest(dir string, manifest npmManifest) error {
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal npm manifest: %w", err)
	}
	data = append(data, '\n')
	path := filepath.Join(dir, "package.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write npm manifest: %w", err)
	}
	return nil
}

func readRootPackageVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read package.json: %w", err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", fmt.Errorf("parse package.json: %w", err)
	}
	if pkg.Version == "" {
		return "", errors.New("package.json missing version")
	}
	return pkg.Version, nil
}
