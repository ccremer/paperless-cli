package paperless

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
)

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

func (clt *Client) makeUpdateDocumentRequest(ctx context.Context, params BulkDownloadParams) (*http.Request, error) {
	log := logr.FromContextOrDiscard(ctx)

	js := map[string]any{
		"content":           params.Content,
		"follow_formatting": params.FollowFormatting,
		"documents":         params.DocumentIDs,
	}
	marshal, err := json.Marshal(js)
	if err != nil {
		return nil, fmt.Errorf("cannot serialize to JSON: %w", err)
	}
	body := bytes.NewReader(marshal)

	path := clt.URL + "/api/documents/bulk_download/"
	log.V(1).Info("Preparing request", "path", path, "document_ids", params.DocumentIDs)
	req, err := http.NewRequestWithContext(ctx, "POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request: %w", err)
	}
	clt.setAuth(req)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
