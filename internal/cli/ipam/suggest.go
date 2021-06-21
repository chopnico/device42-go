package ipam

import (
	"fmt"

	"github.com/chopnico/device42"
	dc "github.com/chopnico/device42/internal/cli"

	"github.com/urfave/cli/v2"
)

func suggestIp(app *cli.App, api *device42.Api) *cli.Command {
	return &cli.Command{
		Name: "ip",
		Usage: "suggest an ip address",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "subnet-id",
				Usage: "`SUBNET-ID` to chose an IP from",
				Required: true,
			},
			&cli.BoolFlag{
				Name: "reserve",
				Usage: "reserve IP address",
			},
		},
		Action: func(c *cli.Context) error {
			ip, err := api.SuggestIp(c.String("subnet-id"), c.Bool("reserve"))
			if err != nil {
				return err
			}

			if c.Bool("json") {
				err := dc.PrintJson(ip)
				if err != nil {
					return err
				}
				return nil
			}

			fmt.Println(ip.Address)

			return nil
		},
	}
}

func suggestCommands(app *cli.App, api *device42.Api) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		suggestIp(app, api),
	)

	return commands
}
