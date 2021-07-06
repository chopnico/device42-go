package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chopnico/device42"

	"github.com/chopnico/output"
	"github.com/urfave/cli/v2"
)

func ipamVrfGroupCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
		ipamVrfGroupShow(app),
	)

	return commands
}

func ipamVrfGroupShow(app *cli.App) *cli.Command {
	flags := globalFlags()

	return &cli.Command{
		Name:  "show",
		Usage: "show a vrf group",
		Flags: flags,
		Action: func(c *cli.Context) error {
			api := c.Context.Value("api").(*device42.Api)
			vrfGroup, err := api.VrfGroups()

			if err != nil {
				return err
			}

			switch c.String("format") {
			case "json":
				fmt.Print(output.FormatItemsAsJson(vrfGroup))
			case "list":
				if c.String("properties") == "" {
					fmt.Print(output.FormatItemsAsList(vrfGroup, nil))
				} else {
					p := strings.Split(c.String("properties"), ",")
					fmt.Print(output.FormatItemsAsList(vrfGroup, p))
				}
			default:
				data := [][]string{}
				for _, i := range *vrfGroup {
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
