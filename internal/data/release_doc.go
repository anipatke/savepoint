package data

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReleaseDocID is a stable identifier for a known supporting document.
type ReleaseDocID string

const (
	// ReleaseDocReleasePRD is the PRD scoped to the selected release.
	ReleaseDocReleasePRD ReleaseDocID = "release-prd"
	// ReleaseDocOverallPRD is the project-wide PRD at the .savepoint root.
	ReleaseDocOverallPRD ReleaseDocID = "overall-prd"
	// ReleaseDocOverallDesign is the project-wide Design at the .savepoint root.
	ReleaseDocOverallDesign ReleaseDocID = "overall-design"
)

// ReleaseDoc is a supporting project document loaded from the .savepoint root.
// Available reports whether the file was found and read; an unavailable doc
// carries an empty Body rather than aborting the load.
type ReleaseDoc struct {
	ID        ReleaseDocID
	Label     string
	Path      string // path relative to the .savepoint root, for display
	Body      string
	Available bool
}

// releaseDocSpec is the static definition of a known supporting document. rel
// returns the document's path relative to the .savepoint root for the given
// release, so release-scoped and root-scoped documents share one loader.
type releaseDocSpec struct {
	id    ReleaseDocID
	label string
	rel   func(release string) string
}

// releaseDocSpecs is the bounded set of documents the Release Docs overlay may
// load: the selected release's PRD, plus the project-wide PRD and Design at the
// root. Loading arbitrary files outside this list is intentionally not
// supported.
var releaseDocSpecs = []releaseDocSpec{
	{
		id:    ReleaseDocReleasePRD,
		label: "Release PRD",
		rel:   func(release string) string { return filepath.Join("releases", release, release+"-PRD.md") },
	},
	{
		id:    ReleaseDocOverallPRD,
		label: "Overall PRD",
		rel:   func(string) string { return "PRD.md" },
	},
	{
		id:    ReleaseDocOverallDesign,
		label: "Overall Design",
		rel:   func(string) string { return "Design.md" },
	},
}

// LoadReleaseDocs reads the supporting documents for the Release Docs overlay
// from the .savepoint root: the selected release's PRD and the project-wide
// PRD/Design. A missing document yields an unavailable entry; only an unexpected
// read error aborts the load and is returned with path context.
func LoadReleaseDocs(root, release string) ([]ReleaseDoc, error) {
	docs := make([]ReleaseDoc, 0, len(releaseDocSpecs))
	for _, spec := range releaseDocSpecs {
		rel := spec.rel(release)
		doc := ReleaseDoc{
			ID:    spec.id,
			Label: spec.label,
			Path:  rel,
		}

		content, err := os.ReadFile(filepath.Join(root, rel))
		switch {
		case err == nil:
			doc.Body = string(content)
			doc.Available = true
		case os.IsNotExist(err):
			doc.Available = false
		default:
			return nil, fmt.Errorf("read release doc %s: %w", rel, err)
		}

		docs = append(docs, doc)
	}
	return docs, nil
}
