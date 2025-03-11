package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

var helpTemplate = `NAME:
   {{template "helpNameTemplate" .}}

USAGE:
   {{if .UsageText}}{{wrap .UsageText 3}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{else}}{{if .Args}} [arguments...]{{end}}{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{template "descriptionTemplate" .}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{template "visibleCommandCategoryTemplate" .}}{{end}}
`

func init() {
	cli.AppHelpTemplate = helpTemplate
}

var app = &cli.App{
	Name:      "cqlshvm",
	Usage:     "Version manager for cqlsh",
	UsageText: "cqlshvm <command> [<args>]",
	Commands: []*cli.Command{
		listCommand,
		downloadCommand,
	},
}

func Execute() error {
	if err := app.Run(os.Args); err != nil {
		return err
	}
	return nil
}
