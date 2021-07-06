package cli

import (
	"github.com/urfave/cli/v2"
)

func ipamCommands(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:    "ipam",
		Aliases: []string{"i"},
		Usage:   "ip address management",
		Subcommands: []*cli.Command{
			{
				Name:        "ip",
				Usage:       "interact with ips",
				Aliases:     []string{"i"},
				Subcommands: ipamIpCommands(app),
			},
			{
				Name:        "subnet",
				Usage:       "interact with subnets",
				Aliases:     []string{"s"},
				Subcommands: ipamSubnetCommands(app),
			},
			{
				Name:        "vrf-group",
				Usage:       "interact with vrf groups",
				Aliases:     []string{"vg"},
				Subcommands: ipamVrfGroupCommands(app),
			},
		},
	}
}
