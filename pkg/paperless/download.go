package paperless

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ccremer/clustercode/pkg/errors"
	"github.com/go-logr/logr"
)

type BulkDownloadContent string

type BulkDownloadParams struct {
	DocumentIDs      []int
	FollowFormatting bool
	Content          BulkDownloadContent
}

const (
	BulkDownloadBoth     BulkDownloadContent = "both"
	BulkDownloadArchives BulkDownloadContent = "archive"
	BulkDownloadOriginal BulkDownloadContent = "originals"
)

// String implements fmt.Stringer.
func (c BulkDownloadContent) String() string {
	return string(c)
}

// BulkDownload downloads the documents identified by BulkDownloadParams.DocumentIDs and saves to the given targetPath.
// If targetPath is empty, it will use the suggested file name from Paperless in the current working dir.
func (clt *Client) BulkDownload(ctx context.Context, targetFile *os.File, params BulkDownloadParams) error {
	req, err := clt.makeBulkDownloadRequest(ctx, params)
	if err != nil {
		return err
	}

	log := logr.FromContextOrDiscard(ctx)
	log.V(1).Info("Awaiting response")
	resp, err := clt.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s: %s", resp.Status, string(b))
	}

	log.V(1).Info("Writing download content to file", "file", targetFile.Name())
	_, err = io.Copy(targetFile, resp.Body)
	return errors.Wrap(err, "cannot read response body")
}

func (clt *Client) makeBulkDownloadRequest(ctx context.Context, params BulkDownloadParams) (*http.Request, error) {
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
