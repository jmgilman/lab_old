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

	flagSet := flag.NewFlagSet("", 0)
	flagSet.String(flag_secret_backend, "mock", "")
	_ = flagSet.Parse([]string{expected})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)

	// With no error
	var got string
	b := secretConfig{
		provider: &mocks.MockSecretProvider{
			FnDelete: func(key string) error {
				got = key
				return nil
			},
		},
	}

	result, err := delete(ctx, &b)
	is.NoErr(err)
	is.Equal(expected, got)
	is.Equal(result.Key, expected)

	// With error
	b = secretConfig{
		provider: &mocks.MockSecretProvider{
			FnDelete: func(key string) error {
				return fmt.Errorf("failed")
			},
		},
	}

	_, err = delete(ctx, &b)
	is.Equal(err.Error(), "failed")
}

func TestGenerate(t *testing.T) {
	is := is.New(t)
	expected_key := "key"
	expected_value := "value"

	flagSet := flag.NewFlagSet("", 0)
	flagSet.Int(flag_secret_length, 12, "")
	flagSet.Int(flag_secret_numbers, 2, "")
	flagSet.Int(flag_secret_symbols, 3, "")
	_ = flagSet.Parse([]string{expected_key})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)

	// With no error
	var got_key string
	var got_length int
	var got_numbers int
	var got_symbols int
	s := secretConfig{
		provider: &mocks.MockSecretProvider{
			FnGenerate: func(key string, length int, nums int, symbols int) (string, error) {
				got_key = key
				got_length = length
				got_numbers = nums
				got_symbols = symbols
				return expected_value, nil
			},
		},
	}

	result, err := generate(ctx, &s)
	is.NoErr(err)
	is.Equal(expected_key, got_key)
	is.Equal(12, got_length)
	is.Equal(2, got_numbers)
	is.Equal(3, got_symbols)
	is.Equal(result.Key, expected_key)
	is.Equal(result.Value, expected_value)

	// With error
	s = secretConfig{
		provider: &mocks.MockSecretProvider{
			FnGenerate: func(key string, length int, nums int, symbols int) (string, error) {
				return "", fmt.Errorf("failed")
			},
		},
	}

	_, err = generate(ctx, &s)
	is.Equal(err.Error(), "failed")
}

func TestGet(t *testing.T) {
	is := is.New(t)
	expected_key := "key"
	expected_value := "value"

	flagSet := flag.NewFlagSet("", 0)
	flagSet.String(flag_secret_backend, "mock", "")
	_ = flagSet.Parse([]string{expected_key})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)

	// With no error
	var got_key string
	s := secretConfig{
		provider: &mocks.MockSecretProvider{
			FnGet: func(key string) (string, error) {
				got_key = key
				return expected_value, nil
			},
		},
	}

	result, err := get(ctx, &s)
	is.NoErr(err)
	is.Equal(expected_key, got_key)
	is.Equal(result.Key, expected_key)
	is.Equal(result.Value, expected_value)

	// With error
	s = secretConfig{
		provider: &mocks.MockSecretProvider{
			FnGet: func(key string) (string, error) {
				return "", fmt.Errorf("failed")
			},
		},
	}

	_, err = get(ctx, &s)
	is.Equal(err.Error(), "failed")
}

func TestSet(t *testing.T) {
	is := is.New(t)
	expected_key := "test"
	expected_value := "test"

	flagSet := flag.NewFlagSet("", 0)
	flagSet.String(flag_secret_backend, "mock", "")
	_ = flagSet.Parse([]string{expected_key, expected_value})
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)

	// With no error
	var got_key string
	var got_value string
	s := secretConfig{
		provider: &mocks.MockSecretProvider{
			FnSet: func(key string, value string) error {
				got_key = key
				got_value = value
				return nil
			},
		},
	}

	result, err := set(ctx, &s)
	is.NoErr(err)
	is.Equal(expected_key, got_key)
	is.Equal(expected_value, got_value)
	is.Equal(result.Key, expected_key)
	is.Equal(result.Value, expected_value)

	// With error
	s = secretConfig{
		provider: &mocks.MockSecretProvider{
			FnSet: func(key string, value string) error {
				return fmt.Errorf("failed")
			},
		},
	}

	_, err = set(ctx, &s)
	is.Equal(err.Error(), "failed")
}
