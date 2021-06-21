package ipam

import (
	"github.com/chopnico/device42"

	"github.com/urfave/cli/v2"
)

func NewCommand(app *cli.App, api *device42.Api) {
	app.Commands = append(app.Commands,
		&cli.Command{
			Name: "ipam",
			Aliases: []string{"i"},
			Usage: "ip address management",
			Subcommands: []*cli.Command{
				{
					Name: "suggest",
					Usage: "suggest something",
					Aliases: []string{"s"},
					Subcommands: suggestCommands(app, api),

				},
				{
					Name: "get",
					Usage: "get something",
					Aliases: []string{"g"},
					Subcommands: getCommands(app, api),
				},
				{
					Name: "clear",
					Usage: "clear something",
					Aliases: []string{"c"},
					Subcommands: clearCommands(app, api),
				},
			},
		},
	)
}
