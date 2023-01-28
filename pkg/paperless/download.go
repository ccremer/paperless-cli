package paperless

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
// If BulkDownloadParams.DocumentIDs is empty, all documents will be downloaded.
// If targetPath is empty, it will use the suggested file name from Paperless in the current working dir.
func (clt *Client) BulkDownload(ctx context.Context, targetPath string, params BulkDownloadParams) error {
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

	out, err := os.Create(getTargetPathOrFromHeader(targetPath, resp.Header))
	defer out.Close()

	log.V(1).Info("Writing download content to file", "file", out.Name())
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}
	return nil
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
	log.V(1).Info("Preparing request", "path", path)
	req, err := http.NewRequestWithContext(ctx, "POST", path, body)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request: %w", err)
	}
	clt.setAuth(req)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func getTargetPathOrFromHeader(v string, header http.Header) string {
	if v != "" {
		return v
	}
	raw := header.Get("content-disposition")
	fileName := strings.TrimSuffix(strings.TrimPrefix(raw, `attachment; filename="`), `"`)
	return fileName
}
