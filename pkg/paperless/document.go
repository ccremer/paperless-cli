package paperless

type Document struct {
	// ID of the document, read-only.
	ID int `json:"id"`
}

func MapToDocumentIDs(docs []Document) []int {
	ids := make([]int, len(docs))
	for i := 0; i < len(docs); i++ {
		ids[i] = docs[i].ID
	}
	return ids
}
