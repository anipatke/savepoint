package main

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFixture(t *testing.T, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, content, 0o755); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

func elfHeader(machine uint16) []byte {
	h := make([]byte, 20)
	h[0], h[1], h[2], h[3] = 0x7F, 'E', 'L', 'F'
	binary.LittleEndian.PutUint16(h[18:20], machine)
	return h
}

func machOHeader(cpu uint32) []byte {
	h := make([]byte, 8)
	binary.LittleEndian.PutUint32(h[0:4], machO64Magic)
	binary.LittleEndian.PutUint32(h[4:8], cpu)
	return h
}

func TestValidateBinaryFormat_acceptsValid(t *testing.T) {
	cases := []struct {
		name    string
		target  target
		fixture string
		content []byte
	}{
		{"windows-amd64 PE", target{os: "windows", arch: "amd64"}, "savepoint.exe", append([]byte("MZ"), make([]byte, 30)...)},
		{"windows-arm64 PE", target{os: "windows", arch: "arm64"}, "savepoint.exe", append([]byte("MZ"), make([]byte, 30)...)},
		{"linux-amd64 ELF", target{os: "linux", arch: "amd64"}, "savepoint", elfHeader(elfMachineAMD64)},
		{"linux-arm64 ELF", target{os: "linux", arch: "arm64"}, "savepoint", elfHeader(elfMachineARM64)},
		{"darwin-amd64 Mach-O", target{os: "darwin", arch: "amd64"}, "savepoint", machOHeader(machOCPUX86_64)},
		{"darwin-arm64 Mach-O", target{os: "darwin", arch: "arm64"}, "savepoint", machOHeader(machOCPUARM64)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, tc.fixture, tc.content)
			if err := validateBinaryFormat(tc.target, path); err != nil {
				t.Fatalf("validateBinaryFormat: %v", err)
			}
		})
	}
}

func TestValidateBinaryFormat_rejectsMismatch(t *testing.T) {
	cases := []struct {
		name    string
		target  target
		content []byte
		wantMsg string
	}{
		{"windows ELF rejected", target{os: "windows", arch: "amd64"}, elfHeader(elfMachineAMD64), "not a Windows PE binary"},
		{"windows Mach-O rejected", target{os: "windows", arch: "arm64"}, machOHeader(machOCPUARM64), "not a Windows PE binary"},
		{"linux PE rejected", target{os: "linux", arch: "amd64"}, append([]byte("MZ"), make([]byte, 30)...), "not a Linux ELF binary"},
		{"linux wrong arch rejected", target{os: "linux", arch: "amd64"}, elfHeader(elfMachineARM64), "ELF machine"},
		{"darwin PE rejected", target{os: "darwin", arch: "amd64"}, append([]byte("MZ"), make([]byte, 30)...), "not a 64-bit Mach-O binary"},
		{"darwin wrong arch rejected", target{os: "darwin", arch: "amd64"}, machOHeader(machOCPUARM64), "Mach-O CPU"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "artifact.bin", tc.content)
			err := validateBinaryFormat(tc.target, path)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
			if !strings.Contains(err.Error(), tc.wantMsg) {
				t.Fatalf("error = %q, want substring %q", err, tc.wantMsg)
			}
			if !strings.Contains(err.Error(), path) {
				t.Fatalf("error = %q, want artifact path %q", err, path)
			}
			if !strings.Contains(err.Error(), tc.target.os) {
				t.Fatalf("error = %q, want target os %q", err, tc.target.os)
			}
		})
	}
}

func TestValidateBinaryFormat_unknownOS(t *testing.T) {
	path := writeFixture(t, "artifact.bin", []byte("anything"))
	err := validateBinaryFormat(target{os: "plan9", arch: "amd64"}, path)
	if err == nil {
		t.Fatal("expected error for unknown os")
	}
	if !strings.Contains(err.Error(), "no binary validator") {
		t.Fatalf("error = %q, want no validator message", err)
	}
}

func TestValidateBinaryFormat_missingFile(t *testing.T) {
	err := validateBinaryFormat(target{os: "linux", arch: "amd64"}, filepath.Join(t.TempDir(), "missing"))
	if err == nil {
		t.Fatal("expected error for missing artifact")
	}
}

func TestValidateBinaryFormat_truncatedHeader(t *testing.T) {
	cases := []struct {
		name   string
		target target
		bytes  []byte
	}{
		{"windows truncated", target{os: "windows", arch: "amd64"}, []byte{'M'}},
		{"linux truncated", target{os: "linux", arch: "amd64"}, []byte{0x7F, 'E', 'L'}},
		{"darwin truncated", target{os: "darwin", arch: "arm64"}, []byte{0xCF, 0xFA}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeFixture(t, "artifact.bin", tc.bytes)
			err := validateBinaryFormat(tc.target, path)
			if err == nil {
				t.Fatal("expected error for truncated header")
			}
		})
	}
}
