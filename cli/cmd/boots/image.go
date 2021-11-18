package main

import (
	"fmt"
	"io"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/HomeOperations/jmgilman/cli/http"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

type imageConfig struct {
	fs       afero.Fs
	provider gcli.ImageProvider
}

func newImageConfig(c *cli.Context) imageConfig {
	return imageConfig{
		fs:       afero.NewOsFs(),
		provider: http.NewImageProvider(),
	}
}

func image() *cli.Command {
	fetch := &cli.Command{
		Name:      "fetch",
		Usage:     "Downloads the specified Container Linux image to the local disk",
		ArgsUsage: "<CHANNEL> <ARCH>",
		Action: func(c *cli.Context) error {
			i := newImageConfig(c)
			return fetch(c, i)
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Usage:   "Disables printing to STDOUT",
			},
		},
	}

	return &cli.Command{
		Name:        "image",
		Usage:       "Provides operations for working with Container Linux images",
		Subcommands: []*cli.Command{fetch},
	}
}

func fetch(c *cli.Context, i imageConfig) error {
	if c.NArg() < 2 {
		fmt.Println(c.NArg())
		return gcli.Exit("must provide a channel and target architecture")
	}

	filename := buildFilename(c.Args().Get(0), c.Args().Get(1))
	out, err := i.fs.Create(filename)
	if err != nil {
		return gcli.Exit(err.Error())
	}
	defer out.Close()

	data, size, err := i.provider.Fetch(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return gcli.Exit(err.Error())
	}
	defer data.Close()

	if c.Bool("quiet") {
		_, err = io.Copy(out, data)
	} else {
		bar := pb.Full.Start64(size)
		barReader := bar.NewProxyReader(data)

		_, err = io.Copy(out, barReader)
		bar.Finish()
	}

	if err != nil {
		return gcli.Exit(err.Error())
	}

	// Validate image
	_, err = out.Seek(io.SeekStart, 0)
	if err != nil {
		return gcli.Exit(err.Error())
	}

	err = i.provider.Validate(out, c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return gcli.Exit(err.Error())
	}

	return nil
}

func buildFilename(channel string, arch string) string {
	return fmt.Sprintf("flatcar_%s_%s.bin.bz2", channel, arch)
}
