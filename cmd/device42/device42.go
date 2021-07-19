package main

import (
	"context"
	"errors"
	"log"
	"os"

	device42 "github.com/chopnico/device42-go"

	CLI "github.com/chopnico/device42-go/internal/cli"

	"github.com/urfave/cli/v2"
)

// some application and default variables
var (
	AppName  string = "device42"
	AppUsage string = "a device42 cli/tui tool"
	// ldflags will be used to set this. check Makefile
	AppVersion string

	DefaultLoggingLevel = "info"
	DefaultPrintFormat  = "table"
	DefaultTimeOut      = 60
)

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = AppUsage
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "username",
			Usage:       "account `USERNAME`",
			EnvVars:     []string{"DEVICE42_USERNAME"},
			DefaultText: "none",
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "account `PASSWORD`",
			EnvVars:     []string{"DEVICE42_PASSWORD"},
			DefaultText: "none",
		},
		&cli.StringFlag{
			Name:        "host",
			Usage:       "device42 appliance `HOST`",
			EnvVars:     []string{"DEVICE42_HOST"},
			DefaultText: "none",
		},
		&cli.BoolFlag{
			Name:  "ignore-ssl",
			Usage: "ignore ssl errors",
			Value: false,
		},
		&cli.IntFlag{
			Name:  "timeout",
			Usage: "http timeout",
			Value: 0,
		},
		&cli.StringFlag{
			Name:  "format",
			Usage: "printing format (json, list, table)",
			Value: "table",
		},
		&cli.StringFlag{
			Name:  "logging",
			Usage: "set logging level",
			Value: "info",
		},
		&cli.StringFlag{
			Name:  "proxy",
			Usage: "set http proxy",
		},
	}
	app.Before = func(c *cli.Context) error {
		var err error
		var api *device42.API

		if c.String("username") == "" {
			cli.ShowAppHelp(c)
			return errors.New(device42.ErrorEmptyUsername)
		} else if c.String("password") == "" {
			cli.ShowAppHelp(c)
			return errors.New(device42.ErrorEmptyPassword)
		} else if c.String("host") == "" {
			cli.ShowAppHelp(c)
			return errors.New(device42.ErrorEmptyHost)
		}

		api, err = device42.NewAPIBasicAuth(
			c.String("username"),
			c.String("password"),
			c.String("host"),
		)
		if err != nil {
			return err
		}

		// set options
		api.Timeout(c.Int("timeout")).
			LoggingLevel(c.String("logging")).
			Proxy(c.String("proxy"))

		// should we ignore ssl errors?
		if c.Bool("ignore-ssl") {
			api.IgnoreSSLErrors()
		}

		ctx := context.WithValue(c.Context, device42.APIContextKey("api"), api)
		c.Context = ctx

		return nil
	}

	// create cli commands
	CLI.NewCommands(app)

	// run the app
	err := app.Run(os.Args)
	if err != nil {
		if err.Error() != "debugging" {
			log.Fatal(err)
		}
	}
	os.Exit(0)
}
