package main

import (
	"fmt"

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

	CreatedAt     cli.Timestamp
	DocumentTitle string
	DocumentType  string
	Correspondent string
	DocumentTags  cli.StringSlice
}

func newUploadCommand() *UploadCommand {
	c := &UploadCommand{}
	c.Command = cli.Command{
		Name:        "upload",
		Description: "Uploads local document(s) to Paperless instance",
		Before: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				ctx.Command.Subcommands = nil // required to print usage of subcommand
				_ = cli.ShowCommandHelp(ctx, ctx.Command.Name)
				return fmt.Errorf("At least one file is required")
			}
			return nil
		},
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
	}
	return nil
}
