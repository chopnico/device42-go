package cli

import (
	"fmt"

	device42 "github.com/chopnico/device42-go"

	"github.com/chopnico/output"
	"github.com/urfave/cli/v2"
)

func ipamIPCommands(app *cli.App) []*cli.Command {
	var commands []*cli.Command

	commands = append(commands,
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
		&cli.StringFlag{
			Name:     "subnet-id",
			Usage:    "`SUBNET-ID` to chose an IP from",
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
			ip, err := api.SuggestIP(c.String("subnet-id"), c.Bool("reserve"))
			if err != nil {
				return err
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
