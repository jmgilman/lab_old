package main

import (
	"fmt"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/HomeOperations/jmgilman/cli/aws"
	"github.com/urfave/cli/v2"
)

const (
	flag_backend = "backend"
)

// secretConfig holds dependencies utilized by the secret subcommand.
type secretConfig struct {
	provider gcli.SecretProvider
}

// newSecretConfig returns a secretConfig configured with default dependencies.
func newSecretsConfig(c *cli.Context) (*secretConfig, error) {
	sc := &secretConfig{}
	switch c.String(flag_backend) {
	case "aws":
		pc, err := aws.NewSecretProviderConfig(c)
		if err != nil {
			return nil, err
		}

		p := aws.NewSecretProvider(pc)
		sc.provider = &p
	default:
		return nil, fmt.Errorf("invalid backend: %s", c.String(flag_backend))
	}

	return sc, nil
}

// secret returns the secret subcommand.
func secret(a gcli.App) *cli.Command {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "backend",
			Value: "aws",
			Usage: "secret backend to use",
		},
	}
	flags = append(flags, aws.Flags()...)

	gen_flags := []cli.Flag{
		&cli.IntFlag{
			Name:    "length",
			Aliases: []string{"l"},
			Value:   16,
			Usage:   "length of the generated password",
		},
		&cli.IntFlag{
			Name:    "numbers",
			Aliases: []string{"n"},
			Value:   1,
			Usage:   "The quantity of numbers to include in the password",
		},
		&cli.IntFlag{
			Name:    "symbols",
			Aliases: []string{"s"},
			Value:   1,
			Usage:   "The quantity of symbols to include in the password",
		},
	}
	gen_flags = append(flags, gen_flags...)

	delete := &cli.Command{
		Name:      "delete",
		Usage:     "Deletes a secret",
		ArgsUsage: "<KEY>",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			s, err := newSecretsConfig(c)
			if err != nil {
				return a.Exit(nil, err)
			}
			return a.Exit(delete(c, s))
		},
	}
	generate := &cli.Command{
		Name:      "generate",
		Usage:     "Generates a new random secret",
		ArgsUsage: "<KEY>",
		Flags:     gen_flags,
		Action: func(c *cli.Context) error {
			s, err := newSecretsConfig(c)
			if err != nil {
				return a.Exit(nil, err)
			}
			return a.Exit(generate(c, s))
		},
	}
	get := &cli.Command{
		Name:      "get",
		Usage:     "Fetches a secret",
		ArgsUsage: "<KEY>",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			s, err := newSecretsConfig(c)
			if err != nil {
				return a.Exit(nil, err)
			}
			return a.Exit(get(c, s))
		},
	}
	set := &cli.Command{
		Name:      "set",
		Usage:     "Sets a secret",
		ArgsUsage: "<KEY> <VALUE>",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			s, err := newSecretsConfig(c)
			if err != nil {
				return a.Exit(nil, err)
			}
			return a.Exit(set(c, s))
		},
	}

	return &cli.Command{
		Name:        "secret",
		Usage:       "Provides CRUD operations for secrets",
		Subcommands: []*cli.Command{delete, generate, get, set},
	}
}

// deleteResult is the result from calling delete().
type deleteResult struct {
	Key string `json:"key"`
}

// delete deletes a secret.
func delete(c *cli.Context, s *secretConfig) (deleteResult, error) {
	if c.NArg() < 1 {
		return deleteResult{}, fmt.Errorf("must provide a key")
	}

	err := s.provider.Delete(c.Args().First())
	if err != nil {
		return deleteResult{}, err
	}

	return deleteResult{
		Key: c.Args().First(),
	}, nil
}

// generateResult is the result from calling generate().
type generateResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// generate generates a new random secret.
func generate(c *cli.Context, s *secretConfig) (generateResult, error) {
	if c.NArg() < 1 {
		return generateResult{}, fmt.Errorf("must provide a key")
	}

	value, err := s.provider.Generate(c.Args().First(), c.Int("length"), c.Int("numbers"), c.Int("symbols"))
	if err != nil {
		return generateResult{}, err
	}

	return generateResult{
		Key:   c.Args().First(),
		Value: value,
	}, nil
}

// getResult is the result from calling get().
type getResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// get fetches a secret.
func get(c *cli.Context, s *secretConfig) (getResult, error) {
	if c.NArg() < 1 {
		return getResult{}, fmt.Errorf("must provide a key")
	}

	value, err := s.provider.Get(c.Args().First())
	if err != nil {
		return getResult{}, err
	}

	return getResult{
		Key:   c.Args().First(),
		Value: value,
	}, nil
}

// setResult is the result from calling set().
type setResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// set sets a secret.
func set(c *cli.Context, s *secretConfig) (setResult, error) {
	if c.NArg() < 2 {
		return setResult{}, fmt.Errorf("must provide a key and value")
	}

	err := s.provider.Set(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return setResult{}, err
	}

	return setResult{
		Key:   c.Args().Get(0),
		Value: c.Args().Get(1),
	}, nil
}
