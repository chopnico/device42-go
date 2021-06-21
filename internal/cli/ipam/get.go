package ipam

import (
	"strconv"

	"github.com/chopnico/device42"
	dc "github.com/chopnico/device42/internal/cli"

	"github.com/urfave/cli/v2"
)

func getSubnet(app *cli.App, api *device42.Api) *cli.Command {
	return &cli.Command{
		Name: "subnet",
		Usage: "get a subnet",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "name",
				Usage: "`NAME` name of the subnet",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			subnet, err := api.GetSubnetByName(c.String("name"))
			if err != nil {
				return err
			}

			if c.Bool("json") {
				err := dc.PrintJson(subnet)
				if err != nil {
					return err
				}
				return nil
			}

			data := [][]string{
				[]string{strconv.Itoa(subnet.SubnetId), subnet.Name, subnet.Network, strconv.Itoa(subnet.MaskBits)},
			}
			header := []string{"SubnetId", "Name", "Network", "Mask Bits"}

			dc.PrintTable(data, header)

			return nil
		},
	}
}

func getCommands(app *cli.App, api *device42.Api) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		getSubnet(app, api),
	)

	return commands
}
