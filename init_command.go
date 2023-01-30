package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"gopkg.in/yaml.v3"
)

type InitCommand struct {
	cli.Command
}

func newInitCommand() *InitCommand {
	c := &InitCommand{}
	c.Command = cli.Command{
		Name:  "init",
		Usage: "Initializes a config file",
		Description: `If CONFIG-FILE is "-" the YAML will be printed to stdout.
If empty, it will be written to "config.yaml".`,
		Action:    c.Action,
		ArgsUsage: "[CONFIG-FILE]",
	}
	return c
}

func (c *InitCommand) Action(ctx *cli.Context) error {
	configFilePath := "config.yaml"
	if ctx.NArg() >= 1 {
		configFilePath = ctx.Args().First()
	}

	flags := allAltSrcFlags(ctx)
	values := map[string]any{}
	for _, flag := range flags {
		values[flag.Names()[0]] = getValueFor(flag)
	}
	b, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("cannot serialize flags to yaml: %w", err)
	}

	if configFilePath == "-" {
		fmt.Println(string(b))
		return nil
	}
	if _, statErr := os.Stat(configFilePath); statErr != nil && os.IsNotExist(statErr) {
		return os.WriteFile(configFilePath, b, 0644)
	}
	return fmt.Errorf("target file %q exists already", configFilePath)
}

func getValueFor(flag cli.Flag) any {
	if f, ok := flag.(*altsrc.StringFlag); ok {
		return f.Value
	}
	if f, ok := flag.(*altsrc.BoolFlag); ok {
		return f.Value
	}
	if f, ok := flag.(*altsrc.DurationFlag); ok {
		return f.Value
	}
	if f, ok := flag.(*altsrc.IntFlag); ok {
		return f.Value
	}
	panic(fmt.Errorf("unknown flag type: %v", flag))
}

func allAltSrcFlags(ctx *cli.Context) []cli.Flag {
	flagMap := map[string]cli.Flag{}
	for _, flag := range ctx.App.Flags {
		if f, isAltSrcFlag := flag.(altsrc.FlagInputSourceExtension); isAltSrcFlag {
			flagMap[flag.Names()[0]] = f
		}
	}
	for _, subcommand := range ctx.App.Commands {
		for _, flag := range subcommand.Flags {
			if f, isAltSrcFlag := flag.(altsrc.FlagInputSourceExtension); isAltSrcFlag {
				flagMap[flag.Names()[0]] = f
			}
		}
	}
	flags := make([]cli.Flag, 0)
	for _, flag := range flagMap {
		flags = append(flags, flag)
	}
	return flags
}
