package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	machO64Magic    uint32 = 0xFEEDFACF
	machOCPUX86_64  uint32 = 0x01000007
	machOCPUARM64   uint32 = 0x0100000C
	elfMachineAMD64 uint16 = 0x3E
	elfMachineARM64 uint16 = 0xB7
)

var elfMachineByArch = map[string]uint16{
	"amd64": elfMachineAMD64,
	"arm64": elfMachineARM64,
}

var machOCPUByArch = map[string]uint32{
	"amd64": machOCPUX86_64,
	"arm64": machOCPUARM64,
}

func validateBinaryFormat(t target, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s/%s artifact %s: %w", t.os, t.arch, path, err)
	}
	defer f.Close()

	switch t.os {
	case "windows":
		return validatePEHeader(f, t, path)
	case "linux":
		return validateELFHeader(f, t, path)
	case "darwin":
		return validateMachOHeader(f, t, path)
	default:
		return fmt.Errorf("no binary validator for %s/%s artifact %s", t.os, t.arch, path)
	}
}

func validatePEHeader(r io.Reader, t target, path string) error {
	header := make([]byte, 2)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("read %s/%s header %s: %w", t.os, t.arch, path, err)
	}
	if string(header) != "MZ" {
		return fmt.Errorf("%s/%s artifact %s is not a Windows PE binary", t.os, t.arch, path)
	}
	return nil
}

func validateELFHeader(r io.Reader, t target, path string) error {
	header := make([]byte, 20)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("read %s/%s header %s: %w", t.os, t.arch, path, err)
	}
	if header[0] != 0x7F || header[1] != 'E' || header[2] != 'L' || header[3] != 'F' {
		return fmt.Errorf("%s/%s artifact %s is not a Linux ELF binary", t.os, t.arch, path)
	}
	wantMachine, ok := elfMachineByArch[t.arch]
	if !ok {
		return fmt.Errorf("no ELF machine mapping for %s/%s artifact %s", t.os, t.arch, path)
	}
	machine := binary.LittleEndian.Uint16(header[18:20])
	if machine != wantMachine {
		return fmt.Errorf("%s/%s artifact %s ELF machine %#x does not match expected %#x", t.os, t.arch, path, machine, wantMachine)
	}
	return nil
}

func validateMachOHeader(r io.Reader, t target, path string) error {
	header := make([]byte, 8)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("read %s/%s header %s: %w", t.os, t.arch, path, err)
	}
	magic := binary.LittleEndian.Uint32(header[0:4])
	if magic != machO64Magic {
		return fmt.Errorf("%s/%s artifact %s is not a 64-bit Mach-O binary", t.os, t.arch, path)
	}
	wantCPU, ok := machOCPUByArch[t.arch]
	if !ok {
		return fmt.Errorf("no Mach-O CPU mapping for %s/%s artifact %s", t.os, t.arch, path)
	}
	cpu := binary.LittleEndian.Uint32(header[4:8])
	if cpu != wantCPU {
		return fmt.Errorf("%s/%s artifact %s Mach-O CPU %#x does not match expected %#x", t.os, t.arch, path, cpu, wantCPU)
	}
	return nil
}
