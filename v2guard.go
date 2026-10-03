package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// schemaVersionLine matches the top-level schema_version key V2 projects add to
// .savepoint/config.yml. Savepoint 1.x never writes it.
var schemaVersionLine = regexp.MustCompile(`(?m)^schema_version:\s*([0-9]+)`)

// v2ProjectNotice returns the message to show, and true, when the nearest
// .savepoint directory at or above start belongs to a V2 (or newer) project.
// This release cannot read those projects, and without the check it would stop
// with a Go panic about a missing epics directory instead of saying so.
func v2ProjectNotice(start string, args []string) (string, bool) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", false
	}
	for {
		config := filepath.Join(dir, ".savepoint", "config.yml")
		if info, err := os.Stat(filepath.Join(dir, ".savepoint")); err == nil && info.IsDir() {
			raw, err := os.ReadFile(config)
			if err != nil {
				return "", false
			}
			m := schemaVersionLine.FindSubmatch(raw)
			if m == nil {
				return "", false
			}
			if n, err := strconv.Atoi(string(m[1])); err != nil || n < 2 {
				return "", false
			}
			return fmt.Sprintf(`This project uses Savepoint V2 (schema_version: %s in %s), but this is Savepoint 1.x, which cannot read it.

Run the current version instead:
  npx savepoint@latest %s

If package.json lists savepoint ^1.x, update it so npx stops running this copy:
  npm install -D savepoint@latest
`, m[1], config, strings.Join(args, " ")), true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
