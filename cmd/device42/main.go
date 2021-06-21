package main

import (
	"errors"
	"log"
	"os"

	"github.com/chopnico/device42"

	"github.com/chopnico/device42/internal/cli/ipam"

	"github.com/urfave/cli/v2"
)

func main(){
	var api device42.Api

	app := cli.NewApp()
	app.Name = "device42"
	app.Usage = "device42 CLI"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name: "username",
			Usage: "account `USERNAME`",
			EnvVars: []string{"DEVICE42_USERNAME"},
		},
		&cli.StringFlag{
			Name: "password",
			Usage: "account `PASSWORD`",
			EnvVars: []string{"DEVICE42_PASSWORD"},
		},
		&cli.StringFlag{
			Name: "base-url",
			Usage: "`BASE-URL`of the device42 endpoint (https://device42.example.local)",
			EnvVars: []string{"DEVICE42_BASEURL"},
		},
		&cli.BoolFlag{
			Name: "ignore-ssl",
			Usage: "ignore ssl errors",
			Value: false,
		},
		&cli.IntFlag{
			Name: "timeout",
			Usage: "http timeout",
			Value: 0,
		},
		&cli.BoolFlag{
			Name: "json",
			Usage: "print as json",
			Value: false,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.String("username") == "" {
			return errors.New(device42.ErrorEmptyUsername)
		} else if c.String("password") == "" {
			return errors.New(device42.ErrorEmptyPassword)
		} else if c.String("base-url") == "" {
			return errors.New(device42.ErrorEmptyBaseUrl)
		}

		a, err := device42.NewApiBasicAuth(
			c.String("username"),
			c.String("password"),
			c.String("base-url"),
			c.Bool("ignore-ssl"),
			c.Int("timeout"),
		)
		if err != nil {
			return err
		}

		api = (*a)

		return nil
	}

	ipam.NewCommand(app, &api)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
