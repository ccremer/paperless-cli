package app

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ccremer/paperless-cli/pkg/paperless"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func newConfigFileFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name: "config", EnvVars: []string{"CONFIG"},
		Aliases: []string{"C"},
		Value:   "config.yaml",
		Usage:   "path to a config file containing additional config.",
	}
}

func newLogLevelFlag() *altsrc.IntFlag {
	return altsrc.NewIntFlag(&cli.IntFlag{
		Name: "log-level", Aliases: []string{"v"}, EnvVars: []string{"LOG_LEVEL"},
		Usage: "number of the log level verbosity",
		Value: 0,
	})
}

func newURLFlag(dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name: "url", EnvVars: envVars("URL"),
		Usage:       "URL endpoint of the paperless instance.",
		Action:      checkEmptyString("url"),
		Destination: dest,
	})
}

func newTokenFlag(dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name: "token", EnvVars: envVars("TOKEN"),
		Usage:       "password or token of the paperless instance.",
		Action:      checkEmptyString("token"),
		Destination: dest,
	})
}

func newUsernameFlag(dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name: "username", EnvVars: envVars("USERNAME"),
		Usage:       "username for BasicAuth of the paperless instance. Leave empty to use token authentication.",
		Destination: dest,
	})
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

func newConsumeDirFlag(dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name: "consume-dir", EnvVars: []string{"CONSUME_DIR"},
		Usage:       "the directory name which to consume files.",
		Required:    true,
		Destination: dest,
		Action:      checkEmptyString("consume-dir"),
	})
}

func newConsumeDelayFlag(dest *time.Duration) *altsrc.DurationFlag {
	return altsrc.NewDurationFlag(&cli.DurationFlag{
		Name: "consume-delay", EnvVars: []string{"CONSUME_DELAY"},
		Usage:       "the delay after detecting the last file write operation before uploading it.",
		Value:       1 * time.Second,
		Destination: dest,
		Action: func(ctx *cli.Context, duration time.Duration) error {
			if duration.Milliseconds() < 100 {
				return showFlagError(ctx, fmt.Errorf("Duration of flag %q must be at least 100ms", "consume-delay"))
			}
			return nil
		},
	})
}

func newTargetPathFlag(dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name: "target-path", EnvVars: []string{"DOWNLOAD_TARGET_PATH"},
		Usage:       "target file path where documents are downloaded.",
		DefaultText: "documents.zip",
		Destination: dest,
	})
}

func newDownloadContentFlag(dest *string) *altsrc.StringFlag {
	return altsrc.NewStringFlag(&cli.StringFlag{
		Name: "content", EnvVars: []string{"DOWNLOAD_CONTENT"},
		Usage:       "selection of document variant.",
		Value:       paperless.BulkDownloadArchives.String(),
		Destination: dest,
		Action: func(ctx *cli.Context, s string) error {
			enum := []string{
				paperless.BulkDownloadArchives.String(),
				paperless.BulkDownloadOriginal.String(),
				paperless.BulkDownloadBoth.String()}
			for _, key := range enum {
				if s == key {
					return nil
				}
			}
			return fmt.Errorf("parameter %q must be one of [%s]", "content", strings.Join(enum, ", "))
		},
	})
}

func newUnzipFlag(dest *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{
		Name: "unzip", EnvVars: []string{"DOWNLOAD_UNZIP"},
		Usage:       "unzip the downloaded file.",
		Destination: dest,
	})
}

func newOverwriteFlag(dest *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{
		Name: "overwrite", EnvVars: []string{"DOWNLOAD_OVERWRITE"},
		Usage:       "deletes existing file(s) before downloading.",
		Destination: dest,
	})
}

func newIncrementalFlag(dest *bool) *altsrc.BoolFlag {
	return altsrc.NewBoolFlag(&cli.BoolFlag{
		Name: "incremental", EnvVars: []string{"DOWNLOAD_INCREMENTAL"},
		Usage: fmt.Sprintf("only download the missing files and remove deleted documents. Implies --%s and --%s",
			newUnzipFlag(nil).Name, newOverwriteFlag(nil).Name),
		Destination: dest,
	})
}

func loadConfigFileFn(ctx *cli.Context) error {
	path := ctx.String(newConfigFileFlag().Name)
	flags := ctx.Command.Flags
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}
	return altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc(newConfigFileFlag().Name))(ctx)
}

func checkEmptyString(flagName string) func(*cli.Context, string) error {
	return func(ctx *cli.Context, s string) error {
		if s == "" {
			return showFlagError(ctx, fmt.Errorf(`Required flag %q not set`, flagName))
		}
		return nil
	}
}

func showFlagError(ctx *cli.Context, err error) error {
	subcommands := ctx.Command.Subcommands
	ctx.Command.Subcommands = nil // required to print usage of subcommand
	_ = cli.ShowCommandHelp(ctx, ctx.Command.Name)
	ctx.Command.Subcommands = subcommands
	return err
}
