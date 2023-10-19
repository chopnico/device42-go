package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/chopnico/device42-go"
	"github.com/chopnico/output"
	"github.com/urfave/cli/v2"
)

func ipamVLANCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		ipamVLANList(app),
		ipamVLANGet(app),
		ipamVLANSet(app),
		ipamVLANDelete(app),
	)

	return commands
}

func ipamVLANList(app *cli.App) *cli.Command {
	flags := addQuietFlag(
		addDisplayFlags([]cli.Flag{
			&cli.StringFlag{
				Name:     "filter-by-tags",
				Usage:    "allows for filtering of vlans by a list of `TAGS`",
				Required: false,
			},
		},
		))

	return &cli.Command{
		Name:  "list",
		Usage: "list all vlans",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			vlans := &[]device42.VLAN{}
			var err error

			if c.String("filter-by-tags") != "" {
				s := strings.Split(c.String("filter-by-tags"), ",")
				vlans, err = api.GetVLANsByAllTags(s)
			} else {
				vlans, err = api.GetVLANs()
			}
			if err != nil {
				return err
			}
			if c.Bool("quiet") {
				for _, i := range *vlans {
					fmt.Println(i.VlanID)
				}
			} else {
				switch c.String("format") {
				case "json":
					fmt.Print(output.FormatItemsAsJson(vlans))
				case "list":
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemsAsList(vlans, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemsAsList(vlans, p))
					}
				default:
					data := [][]string{}
					for _, i := range *vlans {
						data = append(data,
							[]string{strconv.Itoa(i.VlanID), i.Name, strconv.Itoa(i.Number)},
						)
					}
					headers := []string{"ID", "Name", "Number"}
					fmt.Print(output.FormatTable(data, headers))
				}
			}
			return nil
		},
	}
}

func ipamVLANGet(app *cli.App) *cli.Command {
	flags := addDisplayFlags(nil)

	return &cli.Command{
		Name:      "get",
		Usage:     "get a vlan",
		ArgsUsage: "ID",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "get")
				return errors.New("you must supply a vlan id")
			}

			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			var id int
			_, err := fmt.Sscan(c.Args().First(), &id)
			if err != nil {
				return err
			}
			vlan, err := api.GetVLANByID(id)
			if err != nil {
				return err
			}

			switch c.String("format") {
			case "json":
				fmt.Printf("%s\n", output.FormatItemAsJson(vlan))
			default:
				if c.String("properties") == "" {
					fmt.Print(output.FormatItemAsList(&vlan, nil))
				} else {
					p := strings.Split(c.String("properties"), ",")
					fmt.Print(output.FormatItemAsList(&vlan, p))
				}
			}
			return nil
		},
	}
}

func ipamVLANSet(app *cli.App) *cli.Command {
	flags := addQuietFlag([]cli.Flag{
		&cli.IntFlag{
			Name:     "number",
			Usage:    "the vlan `NUMBER`",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "`NAME` of the vlan",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "`DESCRIPTION` of the vlan",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "notes",
			Usage:    "`NOTES` of the vlan",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "tags",
			Usage:    "`TAGS` of the vlan",
			Required: false,
		},
	})

	return &cli.Command{
		Name:  "set",
		Usage: "add or update a vlan",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			t := strings.Split(c.String("tags"), ",")

			vlan := &device42.VLAN{
				Name:        c.String("name"),
				Number:      c.Int("number"),
				Description: c.String("description"),
				Notes:       c.String("notes"),
				Tags:        t,
			}

			vlan, err := api.SetVLAN(vlan)
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				fmt.Println(vlan.VlanID)
			} else {
				switch c.String("format") {
				case "json":
					fmt.Printf("%s\n", output.FormatItemAsJson(vlan))
				default:
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemAsList(vlan, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemAsList(vlan, p))
					}
				}
			}

			return nil
		},
	}
}

func ipamVLANDelete(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a vlan",
		ArgsUsage: "ID",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "delete")
				return errors.New("you must supply a vlan id")
			}
			for i := 0; i < c.Args().Len(); i++ {
				api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

				id := 0
				_, err := fmt.Sscan(c.Args().Get(i), &id)
				if err != nil {
					return err
				}
				err = api.DeleteVLAN(id)
				if err != nil {
					return err
				}

				fmt.Println("sucessfully deleted vlan with id " + strconv.Itoa(id))
			}
			return nil
		},
	}
}
