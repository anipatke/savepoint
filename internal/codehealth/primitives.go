package codehealth

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// timestampLayout is the only accepted time form: UTC, whole seconds.
const timestampLayout = "2006-01-02T15:04:05Z"

func validateTimestamp(field, value string) error {
	t, err := time.Parse(timestampLayout, value)
	if err != nil || t.Format(timestampLayout) != value {
		return fieldError(ErrMalformedTimestamp, field, "%q is not a UTC RFC 3339 time (YYYY-MM-DDTHH:MM:SSZ)", value)
	}
	return nil
}

const digestPrefix = "sha256:"

// validateDigest accepts "sha256:" plus 64 lowercase hex characters.
func validateDigest(field, value string) error {
	raw, ok := strings.CutPrefix(value, digestPrefix)
	if !ok || len(raw) != 64 || !isLowerHex(raw) {
		return fieldError(ErrMalformedIdentity, field, "%q is not a sha256 digest", value)
	}
	return nil
}

func isLowerHex(s string) bool {
	if _, err := hex.DecodeString(s); err != nil {
		return false
	}
	return s == strings.ToLower(s)
}

// validateCommit accepts a full SHA-1 or SHA-256 object name.
func validateCommit(field, value string) error {
	if (len(value) != 40 && len(value) != 64) || !isLowerHex(value) {
		return fieldError(ErrMalformedIdentity, field, "%q is not a full commit id", value)
	}
	return nil
}

// validateToken checks a short single-line label such as a version string.
func validateToken(field, value string, required bool) error {
	if value == "" {
		if required {
			return fieldError(ErrMalformedRecord, field, "is required")
		}
		return nil
	}
	if len(value) > MaxTokenLen || hasControl(value) || strings.TrimSpace(value) != value {
		return fieldError(ErrUnboundedDetail, field, "must be a trimmed single line of at most %d bytes", MaxTokenLen)
	}
	return nil
}

// validateText checks bounded free text.
func validateText(field, value string, max int) error {
	if len(value) > max || hasControl(value) {
		return fieldError(ErrUnboundedDetail, field, "must be single-line text of at most %d bytes", max)
	}
	return nil
}

func hasControl(s string) bool {
	return strings.IndexFunc(s, unicode.IsControl) >= 0
}

// validateRelativePath requires a repository-relative, forward-slash path with
// no traversal, drive letter, URL, or host-specific form. Glob characters are
// allowed because scope and exclusion entries use this check too.
func validateRelativePath(sentinel error, field, p string) error {
	switch {
	case p == "":
		return fieldError(sentinel, field, "path is empty")
	case len(p) > MaxPathLen:
		return fieldError(sentinel, field, "path exceeds %d bytes", MaxPathLen)
	case hasControl(p):
		return fieldError(sentinel, field, "path contains control characters")
	case strings.ContainsAny(p, `\:`):
		return fieldError(sentinel, field, "path %q must use forward slashes and no drive or URL form", p)
	case strings.HasPrefix(p, "/") || strings.HasPrefix(p, "~"):
		return fieldError(sentinel, field, "path %q must be repository-relative", p)
	}
	for _, seg := range strings.Split(p, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return fieldError(sentinel, field, "path %q has an empty or traversing segment", p)
		}
	}
	return nil
}

// sensitiveSegment reports path segments that commonly hold credentials or
// repository internals. Evidence must never point at them.
func sensitiveSegment(seg string) bool {
	s := strings.ToLower(seg)
	switch {
	case s == ".git", s == ".ssh", s == ".aws", s == ".gnupg", s == ".netrc", s == ".npmrc", s == ".pypirc":
		return true
	case s == ".env" || strings.HasPrefix(s, ".env."):
		return true
	case strings.HasPrefix(s, "id_rsa"), strings.HasPrefix(s, "id_ed25519"), strings.HasPrefix(s, "id_ecdsa"):
		return true
	case strings.HasPrefix(s, "credentials"), strings.HasPrefix(s, "secrets"):
		return true
	}
	for _, ext := range []string{".pem", ".key", ".p12", ".pfx", ".kdbx"} {
		if strings.HasSuffix(s, ext) {
			return true
		}
	}
	return false
}

func validatePatternList(field string, patterns []string) error {
	if len(patterns) > MaxScopeEntries {
		return fieldError(ErrUnboundedDetail, field, "has %d entries; at most %d", len(patterns), MaxScopeEntries)
	}
	seen := make(map[string]bool, len(patterns))
	for i, p := range patterns {
		f := indexed(field, i)
		if err := validateRelativePath(ErrInvalidConfig, f, p); err != nil {
			return err
		}
		if seen[p] {
			return fieldError(ErrInvalidConfig, f, "duplicate entry %q", p)
		}
		seen[p] = true
	}
	return nil
}

func indexed(field string, i int) string {
	return field + "[" + strconv.Itoa(i) + "]"
}
