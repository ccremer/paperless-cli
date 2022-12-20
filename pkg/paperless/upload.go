package paperless

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-logr/logr"
)

type UploadParams struct {
	Title         string
	Created       time.Time
	Correspondent string
	DocumentType  string
	Tags          []string
}

func (clt *Client) Upload(ctx context.Context, filePath string, params UploadParams) error {
	req, err := clt.makeFileUploadRequest(ctx, filePath, params)
	if err != nil {
		return err
	}

	resp, err := clt.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	body, _ := io.ReadAll(resp.Body)
	errMessage := string(body)
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("unauthorized")
	default:
		return fmt.Errorf("request failed with status code %d: %v", resp.StatusCode, errMessage)
	}
}

func (clt *Client) makeFileUploadRequest(ctx context.Context, filePath string, params UploadParams) (*http.Request, error) {
	log := logr.FromContextOrDiscard(ctx).WithValues("filePath", filePath)

	log.V(1).Info("Reading file")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read source file: %w", err)
	}
	defer file.Close()

	log.V(1).Info("Preparing payload for file upload")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("document", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("cannot prepare file for upload: %w", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("cannot copy file to request: %w", err)
	}

	writeUploadFormFields(writer, params)

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("cannot write form body: %w", err)
	}

	log.V(1).Info("Preparing request")
	req, err := http.NewRequestWithContext(ctx, "POST", clt.URL+"/api/documents/post_document/", body)
	if err != nil {
		return nil, fmt.Errorf("cannot prepare request: %w", err)
	}
	clt.setAuth(req)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func writeUploadFormFields(writer *multipart.Writer, params UploadParams) {
	if !params.Created.IsZero() {
		_ = writer.WriteField("created", params.Created.Format("2006-01-02"))
	}
	if v, f := params.Correspondent, "correspondent"; v != "" {
		_ = writer.WriteField(f, v)
	}
	if v, f := params.Title, "title"; v != "" {
		_ = writer.WriteField(f, v)
	}
	if v, f := params.DocumentType, "document_type"; v != "" {
		_ = writer.WriteField(f, v)
	}
	for _, tag := range params.Tags {
		_ = writer.WriteField("tags", tag) // we can specify multiple times
	}
}
