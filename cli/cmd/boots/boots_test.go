package main

import (
	"fmt"
	"io"
	"testing"

	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestExit(t *testing.T) {
	is := is.New(t)
	fs := afero.NewMemMapFs()
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

	err = app.Exit(data, nil)
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
	err = app.Exit(nil, fmt.Errorf("failed"))

	file.Seek(0, io.SeekStart)
	got_data, err = io.ReadAll(file)
	is.NoErr(err)
	is.Equal(string(got_data), expected_json_data)
}
