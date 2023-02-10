package localdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/ccremer/clustercode/pkg/errors"
	"github.com/ccremer/clustercode/pkg/paperless"
)

type metadataContainer struct {
	Documents []paperless.Document `json:"documents,omitempty"`
}

var fileName = ".metadata.json"

// Database is a simple wrapper around a JSON-based file.
type Database struct {
	documents map[int]paperless.Document
	filePath  string
}

// Open reads the database file from the given directory.
// An error is returned if the file doesn't exist or cannot be read.
// There can only be 1 database per directory.
func Open(documentDir string) (*Database, error) {
	filePath := filepath.Join(documentDir, fileName)
	container := metadataContainer{}

	raw, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			raw = []byte("{}")
		} else {
			return nil, fmt.Errorf("cannot open metadata file: %w", err)
		}
	}
	parseErr := json.Unmarshal(raw, &container)
	if parseErr != nil {
		return nil, fmt.Errorf("cannot parse metadata file %s: %w", filePath, err)
	}
	docs := paperless.MapToDocumentMap(container.Documents)
	return &Database{
		filePath:  filePath,
		documents: docs,
	}, nil
}

// FindByID returns the document by the given ID, or nil if not existing.
func (d *Database) FindByID(id int) *paperless.Document {
	if doc, found := d.documents[id]; found {
		return &doc
	}
	return nil
}

// GetAll returns all documents sorted by ID.
func (d *Database) GetAll() []paperless.Document {
	docs := make([]paperless.Document, len(d.documents))
	i := 0
	for _, document := range d.documents {
		docs[i] = document
		i++
	}
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].ID < docs[j].ID
	})
	return docs
}

// Put adds or updates a document.
func (d *Database) Put(doc paperless.Document) {
	d.documents[doc.ID] = doc
}

// Remove deletes the given document
func (d *Database) Remove(doc paperless.Document) {
	delete(d.documents, doc.ID)
}

// Close saves the database.
func (d *Database) Close() error {
	container := metadataContainer{Documents: d.GetAll()}
	b, err := json.Marshal(container)
	if err != nil {
		return fmt.Errorf("cannot save database: %w", err)
	}
	return errors.Wrap(os.WriteFile(d.filePath, b, 0644), "cannot save database")
}
