package data

import "errors"

var (
	ErrNoFrontmatter             = errors.New("no frontmatter found")
	ErrNoClosingFrontmatter      = errors.New("no closing frontmatter delimiter found")
	ErrSavepointDirectoryMissing = errors.New(".savepoint directory not found")
	ErrInvalidStatus             = errors.New("invalid router state")
	ErrMissingFrontmatter        = errors.New("missing or invalid frontmatter")
	ErrConfigNotFound            = errors.New("configuration file not found")
	ErrStructureProblem          = errors.New("project structure problem")
)
