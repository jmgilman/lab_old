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
	expected_filename := buildFilename(expected_channel, expected_arch)

	set := flag.NewFlagSet("test", 0)
	_ = set.Parse([]string{expected_channel, expected_arch})
	ctx := cli.NewContext(&cli.App{}, set, nil)

	var got_fetch_channel string
	var got_fetch_arch string
	var got_validate_channel string
	var got_validate_arch string
	var got_data []byte
	cfg := imageConfig{
		fs: afero.NewMemMapFs(),
		provider: &mocks.MockImageProvider{
			FnFetch: func(channel, arch string) (io.ReadCloser, int64, error) {
				got_fetch_channel = channel
				got_fetch_arch = arch
				return io.NopCloser(expected_data), int64(expected_data.Len()), nil
			},
			FnValidate: func(data io.ReadCloser, channel, arch string) error {
				got_validate_channel = channel
				got_validate_arch = arch
				got_data, _ = io.ReadAll(data)

				return nil
			},
		},
	}

	err := fetch(ctx, cfg)
	is.NoErr(err)
	is.Equal(got_fetch_channel, expected_channel)
	is.Equal(got_fetch_arch, expected_arch)
	is.Equal(got_validate_channel, expected_channel)
	is.Equal(got_validate_arch, expected_arch)
	is.Equal(string(got_data), "test")

	file, err := cfg.fs.Open(expected_filename)
	is.NoErr(err)

	file_data, err := io.ReadAll(file)
	is.NoErr(err)
	is.Equal(string(file_data), "test")

	// With fetch error
	cfg = imageConfig{
		fs: afero.NewMemMapFs(),
		provider: &mocks.MockImageProvider{
			FnFetch: func(channel, arch string) (io.ReadCloser, int64, error) {
				return nil, 0, fmt.Errorf("failed")
			},
			FnValidate: func(data io.ReadCloser, channel, arch string) error {
				return nil
			},
		},
	}

	err = fetch(ctx, cfg)
	is.Equal(err.Error(), "failed")

	// With validate error
	cfg = imageConfig{
		fs: afero.NewMemMapFs(),
		provider: &mocks.MockImageProvider{
			FnFetch: func(channel, arch string) (io.ReadCloser, int64, error) {
				return io.NopCloser(expected_data), int64(expected_data.Len()), nil
			},
			FnValidate: func(data io.ReadCloser, channel, arch string) error {
				return fmt.Errorf("failed")
			},
		},
	}

	err = fetch(ctx, cfg)
	is.Equal(err.Error(), "failed")
}

func TestBuildFilename(t *testing.T) {
	is := is.New(t)
	expected := "flatcar_stable_aarch64.bin.bz2"
	got := buildFilename("stable", "aarch64")
	is.Equal(got, expected)
}
