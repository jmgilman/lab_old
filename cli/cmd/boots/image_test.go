package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/HomeOperations/jmgilman/cli/mocks"
	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func TestFetch(t *testing.T) {
	is := is.New(t)
	expected_channel := "stable"
	expected_arch := "aarch64"
	expected_data := bytes.NewBuffer([]byte("test"))
	expected_size := int64(expected_data.Len())
	expected_filename := "flatcar_production_image.bin.bz2"
	expected_output_file := "image.bin.bz2"

	flagSet := flag.NewFlagSet("", 0)
	flagSet.String(flag_image_architecture, expected_arch, "")
	flagSet.String(flag_image_channel, expected_channel, "")
	flagSet.String(flag_image_name, expected_filename, "")
	flagSet.String(flag_image_output, expected_output_file, "")
	_ = flagSet.Parse([]string{})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)
	ctx.Set(flag_image_output, expected_output_file)

	var got_fetch_channel string
	var got_fetch_arch string
	var got_fetch_filename string
	var got_validate_channel string
	var got_validate_arch string
	var got_validate_filename string
	var got_data []byte
	cfg := imageConfig{
		fs: afero.NewMemMapFs(),
		provider: &mocks.MockImageProvider{
			FnFetch: func(channel, arch string, filename string) (io.ReadCloser, int64, error) {
				got_fetch_channel = channel
				got_fetch_arch = arch
				got_fetch_filename = filename
				return io.NopCloser(expected_data), expected_size, nil
			},
			FnValidate: func(data io.ReadCloser, channel, arch string, filename string) error {
				got_validate_channel = channel
				got_validate_arch = arch
				got_validate_filename = filename
				got_data, _ = io.ReadAll(data)

				return nil
			},
		},
	}

	result, err := fetch(ctx, cfg)
	is.NoErr(err)
	is.Equal(got_fetch_channel, expected_channel)
	is.Equal(got_fetch_arch, expected_arch)
	is.Equal(got_fetch_filename, expected_filename)
	is.Equal(got_validate_channel, expected_channel)
	is.Equal(got_validate_arch, expected_arch)
	is.Equal(got_validate_filename, expected_filename)
	is.Equal(string(got_data), "test")

	is.Equal(result.Path, expected_output_file)
	is.Equal(result.Size, expected_size)

	file, err := cfg.fs.Open(expected_output_file)
	is.NoErr(err)

	file_data, err := io.ReadAll(file)
	is.NoErr(err)
	is.Equal(string(file_data), "test")

	// With fetch error
	cfg = imageConfig{
		fs: afero.NewMemMapFs(),
		provider: &mocks.MockImageProvider{
			FnFetch: func(channel, arch, filename string) (io.ReadCloser, int64, error) {
				return nil, 0, fmt.Errorf("failed")
			},
			FnValidate: func(data io.ReadCloser, channel, arch, filename string) error {
				return nil
			},
		},
	}

	_, err = fetch(ctx, cfg)
	is.Equal(err.Error(), "failed")

	// With validate error
	cfg = imageConfig{
		fs: afero.NewMemMapFs(),
		provider: &mocks.MockImageProvider{
			FnFetch: func(channel, arch, filename string) (io.ReadCloser, int64, error) {
				return io.NopCloser(expected_data), int64(expected_data.Len()), nil
			},
			FnValidate: func(data io.ReadCloser, channel, arch, filename string) error {
				return fmt.Errorf("failed")
			},
		},
	}

	_, err = fetch(ctx, cfg)
	is.Equal(err.Error(), "failed")
}
