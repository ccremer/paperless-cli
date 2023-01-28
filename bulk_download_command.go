package main

import (
	"github.com/ccremer/clustercode/pkg/paperless"
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type BulkDownloadCommand struct {
	cli.Command

	PaperlessURL   string
	PaperlessToken string
	PaperlessUser  string

	TargetPath string
	Content    string
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
		},
	}
	return c
}

func (c *BulkDownloadCommand) Action(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context)

	log.V(1)
	clt := paperless.NewClient(c.PaperlessURL, c.PaperlessUser, c.PaperlessToken)
	documents, queryErr := clt.QueryDocuments(ctx.Context, paperless.QueryParams{
		TruncateContent: true,
	})
	if queryErr != nil {
		return queryErr
	}
	documentIDs := paperless.MapToDocumentIDs(documents)
	downloadErr := clt.BulkDownload(ctx.Context, c.TargetPath, paperless.BulkDownloadParams{
		FollowFormatting: true,
		Content:          paperless.BulkDownloadContent(c.Content),
		DocumentIDs:      documentIDs,
	})
	if downloadErr != nil {
		return downloadErr
	}
	log.Info("Downloaded zip archive")
	return downloadErr
}
