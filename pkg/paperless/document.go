package paperless

type Document struct {
	// ID of the document, read-only.
	ID int `json:"id"`
	// OriginalFileName of the original document, read-only.
	OriginalFileName string `json:"original_file_name,omitempty"`
	// ArchivedFileName of the archived document, read-only.
	// May be empty if no archived document is available.
	ArchivedFileName string `json:"archived_file_name,omitempty"`
}

func MapToDocumentIDs(docs []Document) []int {
	ids := make([]int, len(docs))
	for i := 0; i < len(docs); i++ {
		ids[i] = docs[i].ID
	}
	return ids
}

func MapToDocumentMap(docs []Document) map[int]Document {
	docM := make(map[int]Document, len(docs))
	for _, doc := range docs {
		docM[doc.ID] = doc
	}
	return docM
}
