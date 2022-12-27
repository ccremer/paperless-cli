package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func newLogLevelFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name: "log-level", Aliases: []string{"v"}, EnvVars: envVars("LOG_LEVEL"),
		Usage: "number of the log level verbosity",
		Value: 0,
	}
}

func newURLFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name: "url", EnvVars: envVars("URL"),
		Usage:       "URL endpoint of the paperless instance.",
		Required:    true,
		Action:      checkEmptyString("url"),
		Destination: dest,
	}
}

func newTokenFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name: "token", EnvVars: envVars("TOKEN"),
		Usage:       "password or token of the paperless instance.",
		Required:    true,
		Action:      checkEmptyString("token"),
		Destination: dest,
	}
}

func newUsernameFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name: "username", EnvVars: envVars("USERNAME"),
		Usage:       "username for BasicAuth of the paperless instance. Leave empty to use token authentication.",
		Destination: dest,
	}
}

func newCreatedAtFlag(dest *cli.Timestamp) *cli.TimestampFlag {
	return &cli.TimestampFlag{
		Name:        "created-at",
		Usage:       `set the "created" date for all given files.`,
		Layout:      "2006-01-02",
		Destination: dest,
	}
}

func newTitleFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "title",
		Usage:       "set the document title for all given files.",
		Destination: dest,
	}
}
func newCorrespondentFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "correspondent",
		Usage:       "set the correspondent for all given files.",
		Destination: dest,
	}
}
func newDocumentTypeFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "type",
		Usage:       "set the document type for all given files.",
		Destination: dest,
	}
}
func newTagFlag(dest *cli.StringSlice) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "tag",
		Usage:       "set the document tag(s) for all given files.",
		Destination: dest,
	}
}

func newDeleteAfterUploadFlag(dest *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name: "delete-after-upload", EnvVars: envVars("DELETE_AFTER_UPLOAD"),
		Usage:       "deletes the file(s) after upload",
		Destination: dest,
	}
}

func newConsumeDirFlag(dest *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name: "consume-dir", EnvVars: []string{"CONSUME_DIR"},
		Usage:       "the directory name which to consume files.",
		Required:    true,
		Destination: dest,
		Action:      checkEmptyString("consume-dir"),
	}
}

func checkEmptyString(flagName string) func(*cli.Context, string) error {
	return func(ctx *cli.Context, s string) error {
		if s == "" {
			subcommands := ctx.Command.Subcommands
			ctx.Command.Subcommands = nil // required to print usage of subcommand
			_ = cli.ShowCommandHelp(ctx, ctx.Command.Name)
			ctx.Command.Subcommands = subcommands
			return fmt.Errorf(`Required flag %q not set`, flagName)
		}
		return nil
	}
}
