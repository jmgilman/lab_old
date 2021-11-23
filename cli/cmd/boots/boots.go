package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

type App struct {
	out afero.File
}

func (a *App) Exit(data interface{}, err error) error {
	var result gcli.AppResult
	if err != nil {
		result = gcli.AppResult{
			Error:   err.Error(),
			Success: false,
		}
	} else {
		result = gcli.AppResult{
			Data:    data,
			Success: true,
		}
	}

	text, perr := json.Marshal(result)
	if perr != nil {
		return fmt.Errorf("Error serializing result")
	}

	_, err = a.out.Write(text)
	return err
}

func main() {
	app := App{
		out: os.Stdout,
	}

	image := image(&app)
	secret := secret(&app)

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"V"},
		Usage:   "print version",
	}
	cli_app := &cli.App{
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

	err := cli_app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
