package main

import (
	"github.com/urfave/cli/v2"
)

type UploadCommand struct {
	cli.Command
}

func newUploadCommand() *UploadCommand {
	c := &UploadCommand{}
	c.Command = cli.Command{
		Name:        "upload",
		Description: "Uploads a local document to Paperless instance",
		Action:      c.Action,
	}
	return c
}

func (c *UploadCommand) Action(ctx *cli.Context) error {
	return LogMetadata(ctx)
}
