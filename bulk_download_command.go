package main

import (
	"fmt"
	"os"

	"github.com/ccremer/clustercode/pkg/archive"
	"github.com/ccremer/clustercode/pkg/paperless"
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type BulkDownloadCommand struct {
	cli.Command

	PaperlessURL   string
	PaperlessToken string
	PaperlessUser  string

	TargetPath   string
	Content      string
	UnzipEnabled bool
}

func newBulkDownloadCommand() *BulkDownloadCommand {
	c := &BulkDownloadCommand{}
	c.Command = cli.Command{
		Name:   "bulk-download",
		Usage:  "Downloads multiple documents at once",
		Action: actions(LogMetadata, c.Action),
		Flags: []cli.Flag{
			newURLFlag(&c.PaperlessURL),
			newUsernameFlag(&c.PaperlessUser),
			newTokenFlag(&c.PaperlessToken),
			newTargetPathFlag(&c.TargetPath),
			newDownloadContentFlag(&c.Content),
			newUnzipFlag(&c.UnzipEnabled),
		},
	}
	return c
}

func (c *BulkDownloadCommand) Action(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context)

	log.V(1)
	clt := paperless.NewClient(c.PaperlessURL, c.PaperlessUser, c.PaperlessToken)

	log.Info("Getting list of documents")
	documents, queryErr := clt.QueryDocuments(ctx.Context, paperless.QueryParams{
		TruncateContent: true,
	})
	if queryErr != nil {
		return queryErr
	}
	documentIDs := paperless.MapToDocumentIDs(documents)

	tmpFile, createTempErr := os.CreateTemp(os.TempDir(), "paperless-bulk-download-")
	if createTempErr != nil {
		return fmt.Errorf("cannot open temporary file: %w", createTempErr)
	}
	defer os.Remove(tmpFile.Name()) // cleanup if not renamed

	log.Info("Downloading documents")
	downloadErr := clt.BulkDownload(ctx.Context, tmpFile, paperless.BulkDownloadParams{
		FollowFormatting: true,
		Content:          paperless.BulkDownloadContent(c.Content),
		DocumentIDs:      documentIDs,
	})
	if downloadErr != nil {
		return downloadErr
	}

	if c.UnzipEnabled {
		return c.unzip(ctx, tmpFile)
	}
	return c.move(ctx, tmpFile)
}

func (c *BulkDownloadCommand) unzip(ctx *cli.Context, tmpFile *os.File) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	downloadFilePath := c.TargetPath
	if c.TargetPath == "" {
		downloadFilePath = "documents"
	}
	if unzipErr := archive.Unzip(ctx.Context, tmpFile.Name(), downloadFilePath); unzipErr != nil {
		return fmt.Errorf("cannot unzip file %q to %q: %w", tmpFile.Name(), downloadFilePath, unzipErr)
	}
	log.Info("Unzipped archive to dir", "dir", downloadFilePath)
	return nil
}

func (c *BulkDownloadCommand) move(ctx *cli.Context, tmpFile *os.File) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	downloadFilePath := c.TargetPath
	if c.TargetPath == "" {
		downloadFilePath = "documents.zip"
	}
	if renameErr := os.Rename(tmpFile.Name(), downloadFilePath); renameErr != nil {
		return fmt.Errorf("cannot move temp file: %w", renameErr)
	}
	log.Info("Downloaded zip archive", "file", downloadFilePath)
	return nil
}
