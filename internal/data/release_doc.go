package data

import (
	"fmt"
	"os"
	"path/filepath"
)

// ReleaseDocID is a stable identifier for a known supporting document.
type ReleaseDocID string

const (
	ReleaseDocPRD    ReleaseDocID = "prd"
	ReleaseDocDesign ReleaseDocID = "design"
)

// ReleaseDoc is a supporting project document loaded from the .savepoint root.
// Available reports whether the file was found and read; an unavailable doc
// carries an empty Body rather than aborting the load.
type ReleaseDoc struct {
	ID        ReleaseDocID
	Label     string
	Path      string
	Body      string
	Available bool
}

// releaseDocSpec is the static definition of a known supporting document.
type releaseDocSpec struct {
	id       ReleaseDocID
	label    string
	fileName string
}

// releaseDocSpecs is the bounded set of documents the board may load. Loading
// arbitrary files outside this list is intentionally not supported.
var releaseDocSpecs = []releaseDocSpec{
	{id: ReleaseDocPRD, label: "PRD", fileName: "PRD.md"},
	{id: ReleaseDocDesign, label: "Design", fileName: "Design.md"},
}

// LoadReleaseDocs reads the known supporting documents from the .savepoint root.
// A missing document yields an unavailable entry; only an unexpected read error
// aborts the load and is returned with path context.
func LoadReleaseDocs(root string) ([]ReleaseDoc, error) {
	docs := make([]ReleaseDoc, 0, len(releaseDocSpecs))
	for _, spec := range releaseDocSpecs {
		doc := ReleaseDoc{
			ID:    spec.id,
			Label: spec.label,
			Path:  spec.fileName,
		}

		content, err := os.ReadFile(filepath.Join(root, spec.fileName))
		switch {
		case err == nil:
			doc.Body = string(content)
			doc.Available = true
		case os.IsNotExist(err):
			doc.Available = false
		default:
			return nil, fmt.Errorf("read release doc %s: %w", spec.fileName, err)
		}

		docs = append(docs, doc)
	}
	return docs, nil
}
