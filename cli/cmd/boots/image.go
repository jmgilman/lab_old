package main

import (
	"fmt"
	"io"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/HomeOperations/jmgilman/cli/http"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

const (
	flag_image_architecture = "architecture"
	flag_image_channel      = "channel"
	flag_image_name         = "image"
	flag_image_output       = "output"
)

// imageConfig holds dependencies utilized by the image subcommand.
type imageConfig struct {
	fs       afero.Fs
	provider gcli.ImageProvider
}

// newImageConfig returns an imageConfig configured with default dependencies.
func newImageConfig(c *cli.Context) imageConfig {
	return imageConfig{
		fs:       afero.NewOsFs(),
		provider: http.NewImageProvider(),
	}
}

// image returns the image subcommand.
func image(a gcli.App) *cli.Command {
	fetch := &cli.Command{
		Name:  "fetch",
		Usage: "Downloads the specified Container Linux image to the local disk",
		Action: func(c *cli.Context) error {
			i := newImageConfig(c)
			return a.Exit(fetch(c, i))
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flag_image_architecture,
				Aliases: []string{"a"},
				Usage:   "Target architecture",
				Value:   "amd64",
			},
			&cli.StringFlag{
				Name:    flag_image_channel,
				Aliases: []string{"c"},
				Usage:   "Target channel",
				Value:   "stable",
			},
			&cli.StringFlag{
				Name:    flag_image_name,
				Aliases: []string{"i"},
				Usage:   "Target image filename",
				Value:   "flatcar_production_image.bin.bz2",
			},
			&cli.StringFlag{
				Name:        flag_image_output,
				Aliases:     []string{"o"},
				Usage:       "Output filename",
				DefaultText: "target image filename",
			},
		},
	}

	return &cli.Command{
		Name:        "image",
		Usage:       "Provides operations for working with Container Linux images",
		Subcommands: []*cli.Command{fetch},
	}
}

// fetchResult is the result from calling fetch().
type fetchResult struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// fetch downloads the specified Container Linux image to the local disk.
func fetch(c *cli.Context, i imageConfig) (fetchResult, error) {
	arch := c.String(flag_image_architecture)
	channel := c.String(flag_image_channel)
	filename := c.String(flag_image_name)

	var output_file string
	if c.IsSet(flag_image_output) {
		output_file = c.String(flag_image_output)
	} else {
		fmt.Println("Value:", c.String(flag_image_output))
		output_file = filename
	}

	out, err := i.fs.Create(output_file)
	if err != nil {
		return fetchResult{}, err
	}
	defer out.Close()

	data, size, err := i.provider.Fetch(channel, arch, filename)
	if err != nil {
		return fetchResult{}, err
	}
	defer data.Close()

	_, err = io.Copy(out, data)

	if err != nil {
		return fetchResult{}, err
	}

	_, err = out.Seek(io.SeekStart, 0) // Reset reader for validation
	if err != nil {
		return fetchResult{}, err
	}

	err = i.provider.Validate(out, channel, arch, filename)
	if err != nil {
		return fetchResult{}, err
	}

	return fetchResult{
		Path: output_file,
		Size: size,
	}, nil
}
