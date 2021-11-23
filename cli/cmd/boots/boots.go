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

const (
	flag_quiet   = "quiet"
	flag_verbose = "verbose"
)

// App represents the boots CLI application.
type App struct {
	out afero.File
}

// Exit converts the given data and error into a gcli.AppResult and then writes
// the marshalled JSON output to the configured output.
func (a *App) Exit(c *cli.Context, data interface{}, err error) error {
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

	if !c.Bool(flag_quiet) {
		_, err = a.out.Write(text)
	}
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
				Name:    flag_quiet,
				Usage:   "disable output to STDOUT",
				Aliases: []string{"q"},
			},
			&cli.BoolFlag{
				Name:    flag_verbose,
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
