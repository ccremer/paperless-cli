package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ccremer/plogr"
	"github.com/urfave/cli/v2"
)

var (
	// These will be populated by Goreleaser
	version = "unknown"
	commit  = "-dirty-"
	date    = time.Now().Format("2006-01-02")

	appName     = "paperless-cli"
	appLongName = "CLI tool to interact with paperless-ngx remote API "

	// envPrefix is the global prefix to use for the keys in environment variables
	envPrefix = "PAPERLESS_"
)

func main() {
	app := NewApp()
	err := app.Run(os.Args)
	if err != nil {
		plogr.DefaultErrorPrinter.Println(err.Error())
		os.Exit(1)
	}
}

func NewApp() *cli.App {
	app := &cli.App{
		Name:    appName,
		Usage:   appLongName,
		Version: fmt.Sprintf("%s, revision=%s, date=%s", version, commit, date),

		Before: setupLogging,
		Flags: []cli.Flag{
			newLogLevelFlag(),
		},
		Commands: []*cli.Command{
			&newUploadCommand().Command,
			&newConsumeCommand().Command,
		},
	}
	return app
}

// env combines envPrefix with given suffix delimited by underscore.
func env(suffix string) string {
	return envPrefix + suffix
}

// envVars combines envPrefix with each given suffix delimited by underscore.
func envVars(suffixes ...string) []string {
	arr := make([]string, len(suffixes))
	for i := range suffixes {
		arr[i] = env(suffixes[i])
	}
	return arr
}

func actions(actions ...cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		for _, action := range actions {
			if err := action(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}
