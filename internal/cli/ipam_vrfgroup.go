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

func ipamVRFGroupCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		ipamVRFGroupList(app),
		ipamVRFGroupGet(app),
		ipamVRFGroupSet(app),
		ipamVRFGroupDelete(app),
	)

	return commands
}

func ipamVRFGroupDelete(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a vrf group",
		ArgsUsage: "ID",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "delete")
				return errors.New("you must supply a vrf id")
			}
			for i := 0; i < c.Args().Len(); i++ {
				api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

				id := 0
				_, err := fmt.Sscan(c.Args().Get(i), &id)
				if err != nil {
					return err
				}
				err = api.DeleteVRFGroup(id)
				if err != nil {
					return err
				}

				fmt.Println("sucessfully deleted vrf group with id " + strconv.Itoa(id))
			}
			return nil
		},
	}
}

func ipamVRFGroupSet(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "`NAME` of the vrf group",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Usage:    "`DESCRIPTION` of the vrf group",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "buildings",
			Usage:    "`BUILDINGS` where this vrf group is configured",
			Required: false,
		},
	}

	return &cli.Command{
		Name:  "set",
		Usage: "set a vrf group",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

			buildings := strings.Split(c.String("buildings"), ",")

			vrfGroup := device42.VRFGroup{
				Name:        c.String("name"),
				Description: c.String("description"),
				Buildings:   buildings,
			}

			vg, err := api.SetVRFGroup(&vrfGroup)
			if err != nil {
				return err
			}

			switch c.String("format") {
			case "json":
				fmt.Print(output.FormatItemAsJson(vg))
			default:
				if c.String("properties") == "" {
					fmt.Print(output.FormatItemAsList(vg, nil))
				} else {
					p := strings.Split(c.String("properties"), ",")
					fmt.Print(output.FormatItemsAsList(vg, p))
				}
			}
			return nil
		},
	}
}

func ipamVRFGroupGet(app *cli.App) *cli.Command {
	flags := addQuietFlag(
		addDisplayFlags(
			[]cli.Flag{
				&cli.IntFlag{
					Name:     "id",
					Usage:    "name of the vrf group",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "name",
					Usage:    "name of the vrf group",
					Required: false,
				},
			},
		),
	)

	return &cli.Command{
		Name:  "get",
		Usage: "get a vrf a group",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			var (
				vrfGroup *device42.VRFGroup
				err      error
			)

			if c.Int("id") != 0 {
				vrfGroup, err = api.GetVRFGroupByID(c.Int("id"))
				if err != nil {
					return err
				}
			} else if c.String("name") != "" {
				vrfGroup, err = api.GetVRFGroupByName(c.String("name"))
				if err != nil {
					return err
				}
			} else {
				return errors.New("you must either specifiy an id or the name of a vrf group")
			}

			switch c.String("format") {
			case "json":
				fmt.Print(output.FormatItemAsJson(vrfGroup))
			default:
				if c.String("properties") == "" {
					fmt.Print(output.FormatItemAsList(vrfGroup, nil))
				} else {
					p := strings.Split(c.String("properties"), ",")
					fmt.Print(output.FormatItemsAsList(vrfGroup, p))
				}
			}
			return nil
		},
	}
}

func ipamVRFGroupList(app *cli.App) *cli.Command {
	flags := addQuietFlag(addDisplayFlags(nil))

	return &cli.Command{
		Name:  "list",
		Usage: "list a all vrf groups",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			vrfGroups, err := api.GetVRFGroups()
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				for _, i := range *vrfGroups {
					fmt.Println(i.ID)
				}
			} else {
				switch c.String("format") {
				case "json":
					fmt.Print(output.FormatItemsAsJson(vrfGroups))
				case "list":
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemsAsList(vrfGroups, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemsAsList(vrfGroups, p))
					}
				default:
					data := [][]string{}
					for _, i := range *vrfGroups {
						data = append(data,
							[]string{strconv.Itoa(i.ID), i.Name, strings.Join(i.Buildings, ",")},
						)
					}
					headers := []string{"ID", "Name", "Buildings"}
					fmt.Print(output.FormatTable(data, headers))
				}
			}
			return nil
		},
	}
}
