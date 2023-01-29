package main

import (
	"fmt"
	"os"

	"github.com/ccremer/clustercode/pkg/paperless"
	"github.com/ccremer/plogr"
	"github.com/go-logr/logr"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

type UploadCommand struct {
	cli.Command

	PaperlessURL   string
	PaperlessToken string
	PaperlessUser  string

	CreatedAt         cli.Timestamp
	DocumentTitle     string
	DocumentType      string
	Correspondent     string
	DocumentTags      cli.StringSlice
	DeleteAfterUpload bool
}

func newUploadCommand() *UploadCommand {
	c := &UploadCommand{}
	c.Command = cli.Command{
		Name:  "upload",
		Usage: "Uploads local document(s) to Paperless instance",
		Before: before(func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				ctx.Command.Subcommands = nil // required to print usage of subcommand
				_ = cli.ShowCommandHelp(ctx, ctx.Command.Name)
				return fmt.Errorf("At least one file is required")
			}
			return nil
		}, loadConfigFileFn),
		Action: actions(LogMetadata, c.Action),

		Flags: []cli.Flag{
			newURLFlag(&c.PaperlessURL),
			newUsernameFlag(&c.PaperlessUser),
			newTokenFlag(&c.PaperlessToken),
			newCreatedAtFlag(&c.CreatedAt),
			newTitleFlag(&c.DocumentTitle),
			newDocumentTypeFlag(&c.DocumentType),
			newCorrespondentFlag(&c.Correspondent),
			newTagFlag(&c.DocumentTags),
			newDeleteAfterUploadFlag(&c.DeleteAfterUpload),
		},
		ArgsUsage: "[FILES...]",
	}
	return c
}

func (c *UploadCommand) Action(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context)

	params := paperless.UploadParams{}

	if created := c.CreatedAt.Value(); created != nil {
		params.Created = *created
		log = log.WithValues("created", created.Format("2006-02-03"))
	}
	params.DocumentType, params.Title, params.Correspondent = c.DocumentType, c.DocumentTitle, c.Correspondent
	params.Tags = c.DocumentTags.Value()
	log = log.WithValues("title", params.Title, "type", params.DocumentType, "tags", params.Tags)

	clt := paperless.NewClient(c.PaperlessURL, c.PaperlessUser, c.PaperlessToken)
	for _, arg := range ctx.Args().Slice() {
		log.Info("Uploading file", "file", arg)
		err := clt.Upload(ctx.Context, arg, params)
		if err != nil {
			log.Error(err, "Could not upload file")
			continue
		}
		pterm.Success.Println(plogr.DefaultFormatter("File uploaded", map[string]interface{}{
			"file": arg,
		}))
		if c.DeleteAfterUpload {
			c.deleteAfterUpload(arg)
		}
	}
	return nil
}

func (c *UploadCommand) deleteAfterUpload(arg string) {
	err := os.Remove(arg)
	if err != nil {
		pterm.Warning.Println(plogr.DefaultFormatter("File could not be deleted", map[string]interface{}{
			"file":  arg,
			"error": err,
		}))
	}
}
