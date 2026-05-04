package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// GateResult holds the outcome of a single quality gate.
type GateResult struct {
	Name     string
	Command  string
	Passed   bool
	ExitCode int
	Output   string
}

// RunQualityGates executes configured quality gates (lint, typecheck, test).
func RunQualityGates(root string, overrides ...DoctorDependencies) []GateResult {
	deps := doctorDependencies(overrides)
	configPath := filepath.Join(root, "config.yml")
	cfg, err := deps.ConfigReader.Read(configPath)
	if err != nil {
		return []GateResult{{
			Name:    "config",
			Command: "",
			Passed:  false,
			Output:  fmt.Sprintf("cannot read config: %v", err),
		}}
	}

	timeout := 60 * time.Second
	if cfg.QualityGates.Timeout != "" {
		if d, err := time.ParseDuration(cfg.QualityGates.Timeout); err == nil {
			timeout = d
		}
	}

	var results []GateResult

	if cfg.QualityGates.Lint != nil && *cfg.QualityGates.Lint != "" {
		results = append(results, runGate("lint", *cfg.QualityGates.Lint, root, timeout))
	}
	if cfg.QualityGates.Typecheck != nil && *cfg.QualityGates.Typecheck != "" {
		results = append(results, runGate("typecheck", *cfg.QualityGates.Typecheck, root, timeout))
	}
	if cfg.QualityGates.Test != nil && *cfg.QualityGates.Test != "" {
		results = append(results, runGate("test", *cfg.QualityGates.Test, root, timeout))
	}

	if len(results) == 0 {
		results = append(results, GateResult{
			Name:    "quality_gates",
			Command: "",
			Passed:  true,
			Output:  "no quality gates configured",
		})
	}

	return results
}

func runGate(name, command string, root string, timeout time.Duration) GateResult {
	parts := splitCommand(command)
	if len(parts) == 0 {
		return GateResult{
			Name:    name,
			Command: command,
			Passed:  false,
			Output:  "empty command",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Dir = root
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(output))

	switch {
	case err != nil && ctx.Err() == context.DeadlineExceeded:
		return GateResult{
			Name:    name,
			Command: command,
			Passed:  false,
			Output:  fmt.Sprintf("timed out after %v", timeout),
		}
	case err != nil:
		if exitErr, ok := err.(*exec.ExitError); ok {
			return GateResult{
				Name:     name,
				Command:  command,
				Passed:   false,
				ExitCode: exitErr.ExitCode(),
				Output:   outStr,
			}
		}
		return GateResult{
			Name:    name,
			Command: command,
			Passed:  false,
			Output:  fmt.Sprintf("failed to execute: %v", err),
		}
	default:
		return GateResult{
			Name:     name,
			Command:  command,
			Passed:   true,
			ExitCode: 0,
			Output:   outStr,
		}
	}
}

func splitCommand(command string) []string {
	var parts []string
	current := strings.Builder{}
	inQuote := false

	for i := 0; i < len(command); i++ {
		c := command[i]
		switch {
		case c == '"':
			inQuote = !inQuote
		case c == ' ' && !inQuote:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}
