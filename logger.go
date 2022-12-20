package main

import (
	"os"
	"runtime"

	"github.com/ccremer/plogr"
	"github.com/go-logr/logr"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var logger logr.Logger

func init() {
	// Remove `-v` short option from --version flag
	cli.VersionFlag.(*cli.BoolFlag).Aliases = nil
}

// LogMetadata prints various metadata to the root logger.
// It prints version, architecture and current user ID and returns nil.
func LogMetadata(c *cli.Context) error {
	log := logr.FromContextOrDiscard(c.Context)
	log.WithValues(
		"version", version,
		"date", date,
		"commit", commit,
		"go_os", runtime.GOOS,
		"go_arch", runtime.GOARCH,
		"go_version", runtime.Version(),
		"uid", os.Getuid(),
		"gid", os.Getgid(),
	).Info("Starting up " + appName)
	return nil
}

func setupLogging(c *cli.Context) error {
	logger = newPlogger()
	c.Context = logr.NewContext(c.Context, logger)
	return nil
}

func newPlogger() logr.Logger {
	sink := plogr.NewPtermSink()
	sink.FallbackPrinter = &pterm.Debug
	sink.ErrorPrinter.ShowLineNumber = false
	return logr.New(sink)
}
