package main

import (
	"os"
	"runtime"

	"github.com/ccremer/plogr"
	"github.com/go-logr/logr"
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
	sink := newSink(c.Int(newLogLevelFlag().Name))
	logger = logr.New(sink)
	c.Context = logr.NewContext(c.Context, logger)
	return nil
}

func newSink(level int) *plogr.PtermSink {
	sink := plogr.NewPtermSink()
	sink.ErrorPrinter.ShowLineNumber = true
	for i := 1; i <= level; i++ {
		sink.SetLevelEnabled(i, true)
	}
	return &sink
}
