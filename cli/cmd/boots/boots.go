package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	image := image()
	secret := secret()

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print version",
	}
	app := &cli.App{
		Name:     "boots",
		Version:  "v0.1.1",
		HelpName: "boots",
		Usage:    "A CLI tool for bootstrapping the GLab stack",
		Commands: []*cli.Command{image, secret},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "enable verbose mode",
				Aliases: []string{"v"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
