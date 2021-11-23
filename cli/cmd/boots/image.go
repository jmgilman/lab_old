package main

import (
	"fmt"
	"io"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/HomeOperations/jmgilman/cli/http"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
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
		Name:      "fetch",
		Usage:     "Downloads the specified Container Linux image to the local disk",
		ArgsUsage: "<CHANNEL> <ARCH>",
		Action: func(c *cli.Context) error {
			i := newImageConfig(c)
			return a.Exit(fetch(c, i))
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

// fetchResult is the result from calling fetch().
type fetchResult struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// fetch downloads the specified Container Linux image to the local disk.
func fetch(c *cli.Context, i imageConfig) (fetchResult, error) {
	if c.NArg() < 2 {
		return fetchResult{}, fmt.Errorf("Must provide a channel and target architecture")
	}

	filename := buildFilename(c.Args().Get(0), c.Args().Get(1))
	out, err := i.fs.Create(filename)
	if err != nil {
		return fetchResult{}, err
	}
	defer out.Close()

	data, size, err := i.provider.Fetch(c.Args().Get(0), c.Args().Get(1))
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

	err = i.provider.Validate(out, c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return fetchResult{}, err
	}

	return fetchResult{
		Path: filename,
		Size: size,
	}, nil
}

// buildFilename constructs a filename from an image channel and architecture.
func buildFilename(channel string, arch string) string {
	return fmt.Sprintf("flatcar_%s_%s.bin.bz2", channel, arch)
}
