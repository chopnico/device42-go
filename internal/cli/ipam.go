package cli

import (
	"github.com/urfave/cli/v2"
)

func ipamCommands(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:  "ipam",
		Usage: "ip address management",
		Subcommands: []*cli.Command{
			{
				Name:        "ip",
				Usage:       "ip address management",
				Subcommands: ipamIPCommands(app),
			},
			{
				Name:        "subnet",
				Usage:       "subnet management",
				Subcommands: ipamSubnetCommands(app),
			},
			{
				Name:        "vrf-group",
				Usage:       "vrf group management",
				Subcommands: ipamVRFGroupCommands(app),
			},
		},
	}
}
