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

func ipamIPCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		ipamIPList(app),
		ipamIPGet(app),
		ipamIPSet(app),
		ipamIPDelete(app),
		ipamIPClear(app),
		ipamIPSuggest(app),
	)

	return commands
}

func ipamIPClear(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "address",
			Usage:    "`ADDRESS` to clear",
			Required: true,
		},
	}

	return &cli.Command{
		Name:  "clear",
		Usage: "clear an ip address",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			err := api.ClearIP(c.String("address"))
			if err != nil {
				return err
			}

			fmt.Println("successfully cleared ip address " + c.String("address"))

			return nil
		},
	}
}

func ipamIPSuggest(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.IntFlag{
			Name:     "subnet-id",
			Usage:    "`SUBNET-ID` to chose an ip from",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "mask-bits",
			Usage:    "`MASK-BITS` for ip",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "reserve",
			Usage: "reserve IP address",
		},
	}

	return &cli.Command{
		Name:  "suggest",
		Usage: "suggest an ip address",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			ip, err := api.SuggestIPWithSubnetID(c.Int("subnet-id"), c.Int("mask-bits"), c.Bool("reserve"))
			if err != nil {
				return err
			}

			if c.Bool("reserve") {
				_, err = api.SetIP(ip)
				if err != nil {
					return err
				}
			}

			switch c.String("format") {
			case "json":
				fmt.Printf("%s\n", output.FormatItemAsJson(ip))
			default:
				fmt.Print(output.FormatItemAsList(&ip, []string{"Address"}))
			}
			return nil
		},
	}
}

func ipamIPGet(app *cli.App) *cli.Command {
	flags := addDisplayFlags(nil)

	return &cli.Command{
		Name:      "get",
		Usage:     "get an ip",
		ArgsUsage: "ID",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "get")
				return errors.New("you must supply an ip id")
			}

			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			var id int
			_, err := fmt.Sscan(c.Args().First(), &id)
			if err != nil {
				return err
			}
			ip, err := api.GetIPByID(id)
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
			} else {
				switch c.String("format") {
				case "json":
					fmt.Printf("%s\n", output.FormatItemAsJson(ip))
				default:
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemAsList(&ip, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemAsList(&ip, p))
					}
				}
			}
			return nil
		},
	}
}

func ipamIPSet(app *cli.App) *cli.Command {
	flags := addQuietFlag([]cli.Flag{
		&cli.StringFlag{
			Name:     "address",
			Usage:    "`ADDRESS` of the ip",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "subnet-name",
			Usage:    "`SUBNET-NAME` of the ip",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "label",
			Usage:    "`LABEL` of the ip",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "notes",
			Usage:    "`NOTES` for the ip",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "vrf-group",
			Usage:    "`VRF-GROUP` for the ip",
			Required: false,
		},
	})

	return &cli.Command{
		Name:  "set",
		Usage: "add or update an ip",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			ip := &device42.IP{
				Address:  c.String("address"),
				Subnet:   c.String("subnet-name"),
				Label:    c.String("label"),
				Notes:    c.String("notes"),
				VRFGroup: c.String("vrf-group"),
			}

			ip, err := api.SetIP(ip)
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				fmt.Println(ip.ID)
			} else {
				switch c.String("format") {
				case "json":
					fmt.Printf("%s\n", output.FormatItemAsJson(ip))
				default:
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemAsList(ip, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemAsList(ip, p))
					}
				}
			}

			return nil
		},
	}
}

func ipamIPList(app *cli.App) *cli.Command {
	flags := addQuietFlag(addDisplayFlags(nil))

	return &cli.Command{
		Name:  "list",
		Usage: "list all ips",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

			ips, err := api.GetIPs()
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				for _, i := range *ips {
					fmt.Println(i.ID)
				}
			} else {
				switch c.String("format") {
				case "json":
					fmt.Print(output.FormatItemsAsJson(ips))
				case "list":
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemsAsList(ips, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemsAsList(ips, p))
					}
				default:
					data := [][]string{}
					for _, i := range *ips {
						data = append(data,
							[]string{strconv.Itoa(i.ID), i.Address, i.Subnet, i.Label},
						)
					}
					headers := []string{"ID", "Address", "Subnet", "Label"}
					fmt.Print(output.FormatTable(data, headers))
				}
			}
			return nil
		},
	}
}

func ipamIPDelete(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete an ip",
		ArgsUsage: "ID",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "delete")
				return errors.New("you must supply an ip id")
			}
			for i := 0; i < c.Args().Len(); i++ {
				api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

				var id int
				_, err := fmt.Sscan(c.Args().First(), &id)
				if err != nil {
					return err
				}
				err = api.DeleteIP(id)
				if err != nil {
					return err
				}

				fmt.Println("sucessfully deleted ip with id " + strconv.Itoa(id))
			}
			return nil
		},
	}
}
