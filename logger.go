package main

import (
	"log/slog"
	"os"
	"runtime"

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
	backend := pterm.DefaultLogger.
		WithLevel(mapLevel(c.Int(newLogLevelFlag().Name))).
		WithCaller().
		WithCallerOffset(4)
	handler := pterm.NewSlogHandler(backend)
	slogger := slog.New(handler)
	slog.SetDefault(slogger)
	c.Context = logr.NewContextWithSlogLogger(c.Context, slogger)
	return nil
}

func mapLevel(level int) pterm.LogLevel {
	switch level {
	case -1:
		return pterm.LogLevelDisabled
	case 1:
		return pterm.LogLevelDebug
	}
	if level >= 2 {
		return pterm.LogLevelTrace
	}
	return pterm.LogLevelInfo
}
