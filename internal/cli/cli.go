package cli

import (
	"github.com/urfave/cli/v2"
)

func NewCommands(app *cli.App) {
	app.Commands = append(app.Commands,
		ipamCommands(app),
		buildingCommands(app),
	)
}

func addQuietFlag(flags []cli.Flag) []cli.Flag {
	flags = append(flags,
		&cli.BoolFlag{
			Name:     "quiet",
			Usage:    "prints only ids",
			Required: false,
		},
	)

	return flags
}

func addDisplayFlags(flags []cli.Flag) []cli.Flag {
	flags = append(flags,
		&cli.StringFlag{
			Name:     "properties",
			Usage:    "`PROPERTIES` to print (only relevant to list format)",
			Required: false,
		},
	)

	return flags
}
