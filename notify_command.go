package main

import (
	"github.com/urfave/cli/v2"
)

type NotifyCommand struct {
	cli.Command

	PostConsumeDocumentOptions
}

func newNotifyCommand(documentOptions *PostConsumeDocumentOptions) *NotifyCommand {
	c := &NotifyCommand{}
	c.PostConsumeDocumentOptions = *documentOptions
	c.Command = cli.Command{
		Name:  "notify",
		Usage: "Notifies a recipient about a consumed document",
		Subcommands: []*cli.Command{
			&newNotifyPushoverCommand().Command,
		},
	}
	return c
}
