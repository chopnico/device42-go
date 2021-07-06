package cli

import (
	"github.com/urfave/cli/v2"
)

func NewCommands(app *cli.App) {
	app.Commands = append(app.Commands,
		ipamCommands(app),
	)
}

func globalFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "properties",
			Aliases:  []string{"p"},
			Usage:    "`PROPERTIES` to print (only relevant to list format)",
			Required: false,
		},
	}
}
