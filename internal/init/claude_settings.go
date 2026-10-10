package init

import (
	"os"
	"path/filepath"
	"strings"
)

// claudeSettingsPath is the template-relative Claude Code settings file. It
// holds the user's own hooks, permissions, and status line, so Savepoint only
// ever creates it; an existing file is never rewritten.
const claudeSettingsPath = ".claude/settings.json"

// sessionStartHook and guardHook are the script names that mark the Savepoint
// hooks as already wired into a settings file.
const (
	sessionStartHook = "session-start.js"
	guardHook        = "guard.js"
)

// sessionStartEntry and guardEntry are the exact settings entries to add by
// hand when the project already has its own .claude/settings.json.
const (
	sessionStartEntry = `"hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"node","args":["${CLAUDE_PROJECT_DIR}/.claude/hooks/session-start.js"]}]}]}`
	guardEntry        = `"hooks":{"PreToolUse":[{"matcher":"Edit|Write|MultiEdit|Bash","hooks":[{"type":"command","command":"node","args":["${CLAUDE_PROJECT_DIR}/.claude/hooks/guard.js"]}]}]}`
)

// settingsNote is the report note for a settings file left alone. It names
// only the entries the file lacks, or "" when it has both hooks. A plain text
// search is enough: the file is never parsed or rewritten, so a user's
// comments and key order are safe.
func settingsNote(existing []byte) string {
	text := string(existing)
	var missing []string
	if !strings.Contains(text, sessionStartHook) {
		missing = append(missing, "the session-start hook: "+sessionStartEntry)
	}
	if !strings.Contains(text, guardHook) {
		missing = append(missing, "the owner-rules guard: "+guardEntry)
	}
	if len(missing) == 0 {
		return ""
	}
	return "kept your file; to add " + strings.Join(missing, " and ") + " merge into its top-level object (join with an existing \"hooks\" key)"
}

// ClaudeSettingsAdvice returns the line init prints when the project already
// has a .claude/settings.json without the hook, and "" when nothing is needed.
func ClaudeSettingsAdvice(targetDir string) string {
	existing, err := os.ReadFile(filepath.Join(targetDir, filepath.FromSlash(claudeSettingsPath)))
	if err != nil {
		return ""
	}
	if note := settingsNote(existing); note != "" {
		return claudeSettingsPath + " " + note
	}
	return ""
}

// upgradeClaudeSettings installs the settings file when absent and otherwise
// leaves it byte-identical, reporting the entry to add when the hook is missing.
func upgradeClaudeSettings(absTarget string, content []byte, dryRun bool, write assetWriter) (UpgradeEntry, error) {
	entry := UpgradeEntry{Path: claudeSettingsPath}
	dest := filepath.Join(absTarget, filepath.FromSlash(claudeSettingsPath))

	existing, err := os.ReadFile(dest)
	switch {
	case err == nil:
		entry.Action, entry.Note = ActionUnchanged, ""
		if note := settingsNote(existing); note != "" {
			entry.Action, entry.Note = ActionInfo, note
		}
		return entry, nil
	case !os.IsNotExist(err):
		return failedEntry(entry, ""), err
	}

	entry.Action = ActionInstalled
	if dryRun {
		return entry, nil
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return failedEntry(entry, ""), err
	}
	if err := write(dest, content); err != nil {
		return failedEntry(entry, ""), err
	}
	return entry, nil
}
