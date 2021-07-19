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

func buildingCommands(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:  "building",
		Usage: "building management",
		Subcommands: []*cli.Command{
			buildingList(app),
			buildingGet(app),
			buildingSet(app),
			buildingDelete(app),
		},
	}
}

func buildingList(app *cli.App) *cli.Command {
	flags := addQuietFlag(addDisplayFlags(nil))

	return &cli.Command{
		Name:  "list",
		Usage: "lists all buildings",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			buildings, err := api.GetBuildings()
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				for _, i := range *buildings {
					fmt.Println(i.BuildingID)
				}
			} else {
				switch c.String("format") {
				case "json":
					fmt.Print(output.FormatItemsAsJson(buildings))
				case "list":
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemsAsList(buildings, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemsAsList(buildings, p))
					}
				default:
					data := [][]string{}
					for _, i := range *buildings {
						data = append(data,
							[]string{strconv.Itoa(i.BuildingID), i.Name, i.Address},
						)
					}
					headers := []string{"ID", "Name", "Address"}
					fmt.Print(output.FormatTable(data, headers))
				}
			}

			return nil
		},
	}
}

func buildingGet(app *cli.App) *cli.Command {
	flags := addQuietFlag(
		addDisplayFlags(
			[]cli.Flag{
				&cli.StringFlag{
					Name:     "name",
					Usage:    "get building by `NAME`",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "id",
					Usage:    "get building by `ID`",
					Required: false,
				},
			},
		),
	)

	return &cli.Command{
		Name:  "get",
		Usage: "get a building",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)
			var (
				buildings *[]device42.Building
				err       error
			)

			if c.String("name") != "" {
				buildings, err = api.GetBuildingByName(c.String("name"))
			} else if c.Int("id") != 0 {
				buildings, err = api.GetBuildingByID(c.Int("id"))
			} else {
				_ = cli.ShowCommandHelp(c, "get")
				return errors.New("you must supply a name")
			}
			if err != nil {
				return err
			}

			if c.Bool("quiet") {
				for _, i := range *buildings {
					fmt.Println(i.BuildingID)
				}
			} else {
				switch c.String("format") {
				case "json":
					fmt.Print(output.FormatItemsAsJson(buildings))
				case "list":
					if c.String("properties") == "" {
						fmt.Print(output.FormatItemsAsList(buildings, nil))
					} else {
						p := strings.Split(c.String("properties"), ",")
						fmt.Print(output.FormatItemsAsList(buildings, p))
					}
				default:
					data := [][]string{}
					for _, i := range *buildings {
						data = append(data,
							[]string{strconv.Itoa(i.BuildingID), i.Name, i.Address},
						)
					}
					headers := []string{"ID", "Name", "Address"}
					fmt.Print(output.FormatTable(data, headers))
				}
			}

			return nil
		},
	}
}

func buildingSet(app *cli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "the `NAME` of the building",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "address",
			Usage:    "the `ADDRESS` of the building",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "contact-name",
			Usage:    "the building's main `CONTACT-NAME`",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "notes",
			Usage:    "some `NOTES` about the building",
			Required: false,
		},
	}

	return &cli.Command{
		Name:  "set",
		Usage: "create or update a building",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

			building := device42.Building{
				Name:        c.String("name"),
				Address:     c.String("address"),
				ContactName: c.String("contact-name"),
				Notes:       c.String("notes"),
			}

			b, err := api.SetBuilding(&building)
			if err != nil {
				return err
			}

			switch c.String("format") {
			case "json":
				fmt.Print(output.FormatItemsAsJson(b))
			default:
				if c.String("properties") == "" {
					fmt.Print(output.FormatItemsAsList(b, nil))
				} else {
					p := strings.Split(c.String("properties"), ",")
					fmt.Print(output.FormatItemsAsList(b, p))
				}
			}
			return nil
		},
	}
}

func buildingDelete(app *cli.App) *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "delete a building",
		ArgsUsage: "ID",
		Action: func(c *cli.Context) error {
			if c.Args().Len() == 0 {
				_ = cli.ShowCommandHelp(c, "delete")
				return errors.New("you must supply a building id")
			}
			for i := 0; i < c.Args().Len(); i++ {
				api := c.Context.Value(device42.APIContextKey("api")).(*device42.API)

				var id int
				_, err := fmt.Sscan(c.Args().First(), &id)
				if err != nil {
					return err
				}
				err = api.DeleteBuilding(id)
				if err != nil {
					return err
				}

				fmt.Println("sucessfully deleted building with id " + strconv.Itoa(id))
			}
			return nil
		},
	}
}
