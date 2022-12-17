package main

import (
	"github.com/urfave/cli/v2"
)

func newLogLevelFlag() *cli.IntFlag {
	return &cli.IntFlag{
		Name: "log-level", Aliases: []string{"v"}, EnvVars: envVars("LOG_LEVEL"),
		Usage: "number of the log level verbosity",
		Value: 0,
	}
}
