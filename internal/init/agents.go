package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	managedBegin = "<!-- SAVEPOINT:BEGIN -->"
	managedEnd   = "<!-- SAVEPOINT:END -->"
)

var agentGuideVariants = []string{
	"AGENTS.md",
	"agents.md",
	"Agents.md",
	"Agents.MD",
	"AGENTS.MD",
}

// FindAgentGuide returns the path of an existing agent guide in dir,
// checking canonical name and casing variants via directory listing so the
// actual on-disk filename is preserved on case-insensitive filesystems.
// Returns "" if none found.
func FindAgentGuide(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	known := make(map[string]bool, len(agentGuideVariants))
	for _, v := range agentGuideVariants {
		known[strings.ToLower(v)] = true
	}
	for _, e := range entries {
		if !e.IsDir() && known[strings.ToLower(e.Name())] {
			return filepath.Join(dir, e.Name())
		}
	}
	return ""
}

// MergeAgentGuide inserts or refreshes the Savepoint managed block in the
// agent guide at targetPath. If the file does not exist, writes the block
// as the entire file. User content outside the block is preserved.
func MergeAgentGuide(targetPath, rendered string) error {
	block := managedBegin + "\n" + strings.TrimSpace(rendered) + "\n" + managedEnd

	existing, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return AtomicWrite(targetPath, []byte(block+"\n"))
		}
		return fmt.Errorf("read agent guide: %w", err)
	}

	merged, _ := replaceManagedBlock(string(existing), block)
	return AtomicWrite(targetPath, []byte(merged))
}

// replaceManagedBlock swaps the managed block into existing, and reports
// whether a complete marker pair was there to swap. A file with no markers, or
// only one of the pair, has no place to put the block: the merged result then
// appends it, which adopts the whole file as an agent guide. Callers that must
// not adopt silently — upgrade — check the flag; callers that mean to adopt —
// init scaffolding — ignore it.
func replaceManagedBlock(existing, block string) (string, bool) {
	begin := strings.Index(existing, managedBegin)
	end := strings.Index(existing, managedEnd)
	if begin != -1 && end != -1 && end > begin {
		endIdx := end + len(managedEnd)
		return existing[:begin] + block + existing[endIdx:], true
	}
	trimmed := strings.TrimRight(existing, "\n")
	return trimmed + "\n\n" + block + "\n", false
}
