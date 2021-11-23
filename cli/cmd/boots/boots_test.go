package main

import (
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v2"
)

func TestExit(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()

	flagSet := flag.NewFlagSet("", 0)
	_ = flagSet.Parse([]string{})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)

	data := struct {
		Field string
	}{
		Field: "test",
	}

	// With no error
	expected_json_data := `{"data":{"Field":"test"},"error":"","success":true}`
	file, err := fs.Create("test")
	is.NoErr(err)

	app := App{
		out: file,
	}

	err = app.Exit(ctx, data, nil)
	is.NoErr(err)

	file.Seek(0, io.SeekStart)
	got_data, err := io.ReadAll(file)
	is.NoErr(err)
	is.Equal(string(got_data), expected_json_data)

	// With error
	expected_json_data = `{"data":null,"error":"failed","success":false}`
	file, err = fs.Create("test")
	is.NoErr(err)

	app = App{
		out: file,
	}
	err = app.Exit(ctx, nil, fmt.Errorf("failed"))

	file.Seek(0, io.SeekStart)
	got_data, err = io.ReadAll(file)
	is.NoErr(err)
	is.Equal(string(got_data), expected_json_data)

	// With quiet
	file, err = fs.Create("test")
	is.NoErr(err)

	flagSet = flag.NewFlagSet("", 0)
	flagSet.Bool(flag_quiet, true, "")
	_ = flagSet.Parse([]string{})
	ctx = cli.NewContext(&cli.App{}, flagSet, nil)

	app = App{
		out: file,
	}
	err = app.Exit(ctx, nil, fmt.Errorf("failed"))

	file.Seek(0, io.SeekStart)
	got_data, err = io.ReadAll(file)
	is.NoErr(err)
	is.Equal(len(got_data), 0)
}
