package main

import (
	"github.com/go-logr/logr"
	"github.com/urfave/cli/v2"
)

type PostConsumeCommand struct {
	cli.Command

	PostConsumeDocumentOptions
}

func newPostConsumeCommand() *PostConsumeCommand {
	c := &PostConsumeCommand{}
	c.Command = cli.Command{
		Name:        "post-consume",
		Usage:       "Triggers an action after consumption of a document.",
		Description: `All subcommands expect the following environment variables (as defined in https://docs.paperless-ngx.com/advanced_usage/#post-consume-script)`,
		Before:      loadConfigFileFn,

		Subcommands: []*cli.Command{
			&newNotifyCommand(&c.PostConsumeDocumentOptions).Command,
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "document-id",
				EnvVars:     []string{"DOCUMENT_ID"},
				Destination: &c.PostConsumeDocumentOptions.Id,
			},
		},
	}
	return c
}

func (c *PostConsumeCommand) Action(ctx *cli.Context) error {
	log := logr.FromContextOrDiscard(ctx.Context)
	log.Info("hello")
	return nil
}
