package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/HomeOperations/jmgilman/cli/mocks"
	"github.com/matryer/is"
	"github.com/urfave/cli/v2"
)

func TestDelete(t *testing.T) {
	is := is.New(t)
	expected := "test"
	cmd := cmdDelete

	set := flag.NewFlagSet("test", 0)
	set.String(flag_backend, "mock", "test")
	_ = set.Parse([]string{expected})
	ctx := cli.NewContext(&cli.App{}, set, nil)

	// With no error
	var got string
	mockProvider = &mocks.MockSecretProvider{
		FnDelete: func(key string) error {
			got = key
			return nil
		},
	}

	err := secretsCommand(ctx, cmd)
	is.NoErr(err)
	is.Equal(expected, got)

	// With error
	mockProvider = &mocks.MockSecretProvider{
		FnDelete: func(key string) error {
			return fmt.Errorf("failed")
		},
	}

	err = secretsCommand(ctx, cmd)
	is.Equal(err.Error(), "failed")
}

func TestGenerate(t *testing.T) {
	is := is.New(t)
	expected := "test"
	cmd := cmdGenerate

	set := flag.NewFlagSet("test", 0)
	set.String(flag_backend, "mock", "test")
	_ = set.Parse([]string{expected})
	ctx := cli.NewContext(&cli.App{}, set, nil)

	// With no error
	var got string
	mockProvider = &mocks.MockSecretProvider{
		FnGenerate: func(key string) error {
			got = key
			return nil
		},
	}

	err := secretsCommand(ctx, cmd)
	is.NoErr(err)
	is.Equal(expected, got)

	// With error
	mockProvider = &mocks.MockSecretProvider{
		FnGenerate: func(key string) error {
			return fmt.Errorf("failed")
		},
	}

	err = secretsCommand(ctx, cmd)
	is.Equal(err.Error(), "failed")
}

func TestGet(t *testing.T) {
	is := is.New(t)
	expected := "test"
	cmd := cmdGet

	set := flag.NewFlagSet("test", 0)
	set.String(flag_backend, "mock", "test")
	_ = set.Parse([]string{expected})
	ctx := cli.NewContext(&cli.App{}, set, nil)

	// With no error
	var got string
	mockProvider = &mocks.MockSecretProvider{
		FnGet: func(key string) (string, error) {
			got = key
			return "", nil
		},
	}

	err := secretsCommand(ctx, cmd)
	is.NoErr(err)
	is.Equal(expected, got)

	// With error
	mockProvider = &mocks.MockSecretProvider{
		FnGet: func(key string) (string, error) {
			return "", fmt.Errorf("failed")
		},
	}

	err = secretsCommand(ctx, cmd)
	is.Equal(err.Error(), "failed")
}

func TestSet(t *testing.T) {
	is := is.New(t)
	expected_key := "test"
	expected_value := "test"
	cmd := cmdSet

	set := flag.NewFlagSet("test", 0)
	set.String(flag_backend, "mock", "test")
	_ = set.Parse([]string{expected_key, expected_value})
	ctx := cli.NewContext(&cli.App{}, set, nil)

	// With no error
	var got_key string
	var got_value string
	mockProvider = &mocks.MockSecretProvider{
		FnSet: func(key string, value string) error {
			got_key = key
			got_value = value
			return nil
		},
	}

	err := secretsCommand(ctx, cmd)
	is.NoErr(err)
	is.Equal(expected_key, got_key)
	is.Equal(expected_value, got_value)

	// With error
	mockProvider = &mocks.MockSecretProvider{
		FnSet: func(key string, value string) error {
			return fmt.Errorf("failed")
		},
	}

	err = secretsCommand(ctx, cmd)
	is.Equal(err.Error(), "failed")
}
