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

func ipamVrfGroupCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		ipamVrfGroupList(app),
		ipamVrfGroupGet(app),
		ipamVrfGroupSet(app),
		ipamVrfGroupDelete(app),
	)

	return commands
}

func ipamVrfGroupDelete(app *cli.App) *cli.Command {
	flags := globalFlags(
		[]cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "id of the vrf group",
				Required: true,
			},
		},
	)

	return &cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delet a vrf group",
		Flags:   flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)

			err := api.DeleteVrfGroup(c.Int("id"))
			if err != nil {
				return err
			}

			fmt.Println("sucessfully deleted vrf group with id " + strconv.Itoa(c.Int("id")))

			return nil
		},
	}
}

func ipamVrfGroupSet(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Usage:    "`NAME` of the vrf group",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"d"},
			Usage:    "`DESCRIPTION` of the vrf group",
			Required: false,
		},
	}

	return &cli.Command{
		Name:    "set",
		Aliases: []string{"s"},
		Usage:   "set a vrf group",
		Flags:   flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)

			vrfGroup := device42.VrfGroup{
				Name:        c.String("name"),
				Description: c.String("description"),
			}

			vg, err := api.SetVrfGroup(&vrfGroup)
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

func ipamVrfGroupGet(app *cli.App) *cli.Command {
	flags := globalFlags(
		[]cli.Flag{
			&cli.IntFlag{
				Name:     "id",
				Aliases:  []string{"i"},
				Usage:    "name of the vrf group",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "name of the vrf group",
				Required: false,
			},
		},
	)

	return &cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "get a vrf a group",
		Flags:   flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)
			var (
				vrfGroup *device42.VrfGroup
				err      error
			)

			if c.String("id") != "" {
				vrfGroup, err = api.GetVrfGroupById(c.Int("id"))
				if err != nil {
					return err
				}
			} else if c.String("name") != "" {
				vrfGroup, err = api.GetVrfGroupByName(c.String("name"))
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

func ipamVrfGroupList(app *cli.App) *cli.Command {
	flags := globalFlags(nil)

	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "list a all vrf groups",
		Flags:   flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)
			vrfGroups, err := api.GetVrfGroups()
			if err != nil {
				return err
			}

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
			return nil
		},
	}
}
