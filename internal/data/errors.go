package data

import "errors"

var (
	ErrNoFrontmatter             = errors.New("no frontmatter found")
	ErrNoClosingFrontmatter      = errors.New("no closing frontmatter delimiter found")
	ErrSavepointDirectoryMissing = errors.New(".savepoint directory not found")
)
