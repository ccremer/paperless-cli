package main

import (
	"github.com/urfave/cli/v2"
)

type NotifyPushoverCommand struct {
	cli.Command
}

func newNotifyPushoverCommand() *NotifyPushoverCommand {
	c := &NotifyPushoverCommand{}
	c.Command = cli.Command{
		Name:  "pushover",
		Usage: "Notifies a recipient about a consumed document over pushover.net",
		Flags: []cli.Flag{},
	}
	return c
}
