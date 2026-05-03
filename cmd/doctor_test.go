package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunDoctorHelp(t *testing.T) {
	var stdout bytes.Buffer
	called := false

	code, err := RunDoctor(context.Background(), []string{"--help"}, &stdout, func(DoctorOptions) (int, error) {
		called = true
		return 0, nil
	})

	if err != nil {
		t.Fatalf("RunDoctor() error = %v", err)
	}
	if called {
		t.Fatal("RunDoctor() called runner for help")
	}
	if code != 0 {
		t.Fatalf("RunDoctor() code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "doctor [--epic <epic>]") {
		t.Fatalf("help output = %q", stdout.String())
	}
}

func TestRunDoctorNoArgs(t *testing.T) {
	got := runDoctorOptions(t, nil)

	if got.Epic != "" {
		t.Fatalf("Epic = %q, want empty", got.Epic)
	}
}

func TestRunDoctorEpic(t *testing.T) {
	got := runDoctorOptions(t, []string{"--epic", "E03"})

	if got.Epic != "E03" {
		t.Fatalf("Epic = %q, want E03", got.Epic)
	}
}

func TestRunDoctorEpicMissingValue(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunDoctor(context.Background(), []string{"--epic"}, &stdout, func(DoctorOptions) (int, error) {
		return 0, nil
	})

	if err == nil {
		t.Fatal("RunDoctor() error = nil, want missing value error")
	}
	if !strings.Contains(err.Error(), "--epic requires a value") {
		t.Fatalf("error = %q", err.Error())
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestRunDoctorRejectsUnknownFlag(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunDoctor(context.Background(), []string{"--bogus"}, &stdout, func(DoctorOptions) (int, error) {
		return 0, nil
	})

	if err == nil {
		t.Fatal("RunDoctor() error = nil, want unknown flag error")
	}
	if !strings.Contains(err.Error(), "unknown doctor flag") {
		t.Fatalf("error = %q", err.Error())
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestRunDoctorRejectsPositionalArgs(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunDoctor(context.Background(), []string{"extra"}, &stdout, func(DoctorOptions) (int, error) {
		return 0, nil
	})

	if err == nil {
		t.Fatal("RunDoctor() error = nil, want positional arg error")
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func TestRunDoctorExitCode1(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunDoctor(context.Background(), nil, &stdout, func(DoctorOptions) (int, error) {
		return 1, nil
	})

	if err != nil {
		t.Fatalf("RunDoctor() error = %v", err)
	}
	if code != 1 {
		t.Fatalf("code = %d, want 1", code)
	}
}

func TestRunDoctorExitCode2(t *testing.T) {
	var stdout bytes.Buffer

	code, err := RunDoctor(context.Background(), nil, &stdout, func(DoctorOptions) (int, error) {
		return 2, nil
	})

	if err != nil {
		t.Fatalf("RunDoctor() error = %v", err)
	}
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
}

func runDoctorOptions(t *testing.T, args []string) DoctorOptions {
	t.Helper()

	var stdout bytes.Buffer
	var got DoctorOptions
	code, err := RunDoctor(context.Background(), args, &stdout, func(options DoctorOptions) (int, error) {
		got = options
		return 0, nil
	})
	if err != nil {
		t.Fatalf("RunDoctor() error = %v", err)
	}
	if code != 0 {
		t.Fatalf("RunDoctor() code = %d, want 0", code)
	}
	return got
}
