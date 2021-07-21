package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	device42 "github.com/chopnico/device42-go"

	"github.com/chopnico/output"
	"github.com/urfave/cli/v2"
)

func ipamSubnetCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	// ordered
	commands = append(commands,
		ipamSubnetList(app),
		ipamSubnetGet(app),
		ipamSubnetSet(app),
		ipamSubnetSuggest(app),
		ipamSubnetDelete(app),
	)

	return commands
}

func ipamSubnetGet(app *cli.App) *cli.Command {
	flags := addDisplayFlags(nil)

	return &cli.Command{
		Name:      "get",
		Usage:     "get a subnet",
		ArgsUsage: "ID",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "get")
				return errors.New("you must supply a subnet id")
			}

			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			var id int
			_, err := fmt.Sscan(c.Args().First(), &id)
			if err != nil {
				return err
			}
			subnet, err := api.GetSubnetByID(id)
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
			} else {
				switch c.String("format") {
				case "json":
					fmt.Printf("%s\n", output.FormatItemAsJson(subnet))
				default:
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemAsList(&subnet, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemAsList(&subnet, p))
					}
				}
			}
			return nil
		},
	}
}

func ipamSubnetSuggest(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:     "subnet-id",
			Usage:    "the parent `SUBNET-ID` to suggest a subnet from",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "mask-bits",
			Usage:    "the mask bits for the suggested subnet",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "`NAME` of the subnet",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "create",
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
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			subnet, err := api.SuggestSubnet(
				c.Int("subnet-id"),
				c.Int("mask-bits"),
				c.String("name"),
				c.Bool("create"),
			)
			if err != nil {
				return err
			}

			switch c.String("format") {
			case "json":
				fmt.Printf("%s\n", output.FormatItemAsJson(subnet))
			default:
				fmt.Print(output.FormatItemAsList(subnet, []string{"SubnetID", "Name", "Network", "MaskBits", "VrfGroupName"}))
			}

			return nil
		},
	}
}

func ipamSubnetSet(app *cli.App) *cli.Command {
	flags := addQuietFlag([]cli.Flag{
		&cli.StringFlag{
			Name:     "network",
			Usage:    "`NETWORK` address of the subnet",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "mask-bits",
			Usage:    "`MASK-BITS` of the subnet",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "`NAME` of the subnet",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "vrf-group",
			Usage:    "`VRF-GROUP` of the subnet",
			Required: false,
		},
	})

	return &cli.Command{
		Name:  "set",
		Usage: "add or update a subnet",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			subnet := &device42.Subnet{
				Name:     c.String("name"),
				Network:  c.String("network"),
				MaskBits: c.Int("mask-bits"),
				VrfGroup: c.String("vrf-group"),
			}

			subnet, err := api.SetSubnet(subnet)
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				fmt.Println(subnet.SubnetID)
			} else {
				switch c.String("format") {
				case "json":
					fmt.Printf("%s\n", output.FormatItemAsJson(subnet))
				default:
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemAsList(subnet, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemAsList(subnet, p))
					}
				}
			}

			return nil
		},
	}
}

func ipamSubnetList(app *cli.App) *cli.Command {
	flags := addQuietFlag(
		addDisplayFlags([]cli.Flag{
			&cli.StringFlag{
				Name:     "filter-by-tags",
				Usage:    "allows for filtering of subnets by a list of `TAGS`",
				Required: false,
			},
		},
		))

	return &cli.Command{
		Name:  "list",
		Usage: "list all subnets",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			subnets := &[]device42.Subnet{}
			var err error

			if c.String("filter-by-tags") != "" {
				s := strings.Split(c.String("filter-by-tags"), ",")
				subnets, err = api.GetSubnetsByAllTags(s)
			} else {
				subnets, err = api.GetSubnets()
			}
			if err != nil {
				return err
			}
			if c.Bool("quiet") {
				for _, i := range *subnets {
					fmt.Println(i.SubnetID)
				}
			} else {
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
							[]string{strconv.Itoa(i.SubnetID), i.Name, i.Network, strconv.Itoa(i.MaskBits), strconv.Itoa(i.ParentVlanID), i.VrfGroupName},
						)
					}
					headers := []string{"ID", "Name", "Network", "MaskBits", "VLAN ID", "VRF Group"}
					fmt.Print(output.FormatTable(data, headers))
				}
			}
			return nil
		},
	}
}

func ipamSubnetDelete(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a subnet",
		ArgsUsage: "ID",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "delete")
				return errors.New("you must supply a subnet id")
			}
			for i := 0; i < c.Args().Len(); i++ {
				api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

				var id int
				_, err := fmt.Sscan(c.Args().First(), &id)
				if err != nil {
					return err
				}
				err = api.DeleteSubnet(id)
				if err != nil {
					return err
				}

				fmt.Println("sucessfully deleted subnet with id " + strconv.Itoa(id))
			}
			return nil
		},
	}
}
