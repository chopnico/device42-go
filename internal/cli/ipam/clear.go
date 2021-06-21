package ipam

import (
	"fmt"

	"github.com/chopnico/device42"

	"github.com/urfave/cli/v2"
)

func clearIp(app *cli.App, api *device42.Api) *cli.Command {
	return &cli.Command{
		Name: "ip",
		Usage: "clear an ip address",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "address",
				Usage: "`ADDRESS` to clear",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			err := api.ClearIp(c.String("address"))
			if err != nil {
				return err
			}

			fmt.Println("cleared")
			return nil
		},
	}
}

func clearCommands(app *cli.App, api *device42.Api) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		clearIp(app, api),
	)

	return commands
}
