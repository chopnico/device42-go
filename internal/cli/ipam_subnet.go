package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chopnico/device42"

	"github.com/chopnico/output"
	"github.com/urfave/cli/v2"
)

func ipamSubnetCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		ipamSubnetList(app),
		ipamSubnetSuggest(app),
		ipamSubnetSet(app),
	)

	return commands
}

func ipamSubnetGet(app *cli.App) *cli.Command {
	flags := globalFlags()
	flags = append(flags,
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "`NAME` of the subnet",
		},
		&cli.StringFlag{
			Name:    "id",
			Aliases: []string{"i"},
			Usage:   "`ID` of the subnet",
		},
		&cli.StringFlag{
			Name:    "vlan-id",
			Aliases: []string{"vi"},
			Usage:   "`VLAN-ID` of the subnet",
		},
	)
	return &cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get a subnet",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}

func ipamSubnetSuggest(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:     "subnet-id",
			Aliases:  []string{"s"},
			Usage:    "the parent `SUBNET-ID` to suggest a subnet from",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "mask-bits",
			Aliases:  []string{"m"},
			Usage:    "the mask bits for the suggested subnet",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "create",
			Aliases:  []string{"c"},
			Usage:    "should we go ahead and `CREATE` the subnet",
			Value:    false,
			Required: false,
		},
	}

	return &cli.Command{
		Name:  "suggest",
		Usage: "suggest an subnet from a parent subnet",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)
			subnet, err := api.SuggestSubnet(c.Int("subnet-id"), c.Int("mask-bits"))
			if err != nil {
				return err
			}

			if c.Bool("create") {
				err = api.CreateChildSubnet(c.Int("subnet-id"), c.Int("mask-bits"))
				if err != nil {
					return err
				}
			}

			switch c.String("format") {
			case "json":
				fmt.Printf("%s\n", output.FormatItemAsJson(subnet))
			default:
				fmt.Print(output.FormatItemAsList(&subnet, []string{"Network", "MaskBits", "Name"}))
			}
			return nil
		},
	}
}

func ipamSubnetSet(app *cli.App) *cli.Command {
	flags := globalFlags()
	flags = append(flags,
		&cli.StringFlag{
			Name:     "network",
			Aliases:  []string{"net"},
			Usage:    "`NETWORK` address of the subnet",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "mask-bits",
			Aliases:  []string{"m", "mask"},
			Usage:    "`MASK-BITS` of the subnet",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "`NAME` of the subnet",
			Required: false,
		},
	)
	return &cli.Command{
		Name:    "set",
		Usage:   "add or update a subnet",
		Aliases: []string{"s"},
		Flags:   flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)
			subnet := device42.Subnet{
				Name:     c.String("name"),
				Network:  c.String("network"),
				MaskBits: c.Int("mask-bits"),
			}

			err := api.CreateSubnet(&subnet)
			if err != nil {
				return err
			}

			output.FormatItemAsList(subnet, []string{"Name", "Network", "MaskBits"})

			return nil
		},
	}
}

func ipamSubnetList(app *cli.App) *cli.Command {
	flags := globalFlags()

	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list all subnets",
		Flags:   flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)
			var (
				err     error
				subnets *[]device42.Subnet
			)

			subnets, err = api.Subnets()

			if err != nil {
				return err
			}

			switch c.String("format") {
			case "json":
				fmt.Print(output.FormatItemsAsJson(subnets))
			case "list":
				if c.String("properties") == "" {
					fmt.Print(output.FormatItemsAsList(subnets, nil))
				} else {
					p := strings.Split(c.String("properties"), ",")
					fmt.Print(output.FormatItemsAsList(subnets, p))
				}
			default:
				data := [][]string{}
				for _, i := range *subnets {
					data = append(data,
						[]string{strconv.Itoa(i.SubnetID), i.Name, i.Network, strconv.Itoa(i.MaskBits), strconv.Itoa(i.ParentVlanID)},
					)
				}
				headers := []string{"ID", "Name", "Network", "MaskBits", "VLAN ID"}
				fmt.Print(output.FormatTable(data, headers))
			}
			return nil
		},
	}
}
