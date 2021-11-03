package main

import (
	"fmt"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/HomeOperations/jmgilman/cli/aws"
	"github.com/urfave/cli/v2"
)

type SecretsCommand int

const (
	cmdDelete SecretsCommand = iota
	cmdGenerate
	cmdGet
	cmdSet
)

const (
	flag_backend = "backend"
)

var mockProvider gcli.SecretProvider

func secrets() *cli.Command {
	// Register flags
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

	// Setup subcommands
	delete := &cli.Command{
		Name:      "delete",
		Usage:     "Deletes a secret",
		ArgsUsage: "<KEY>",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			return secretsCommand(c, cmdDelete)
		},
	}
	generate := &cli.Command{
		Name:      "generate",
		Usage:     "Generates a new random secret",
		ArgsUsage: "<KEY>",
		Flags:     gen_flags,
		Action: func(c *cli.Context) error {
			return secretsCommand(c, cmdGenerate)
		},
	}
	get := &cli.Command{
		Name:      "get",
		Usage:     "Fetches a secret",
		ArgsUsage: "<KEY>",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			return secretsCommand(c, cmdGet)
		},
	}
	set := &cli.Command{
		Name:      "set",
		Usage:     "Sets a secret",
		ArgsUsage: "<KEY> <VALUE>",
		Flags:     flags,
		Action: func(c *cli.Context) error {
			return secretsCommand(c, cmdSet)
		},
	}

	return &cli.Command{
		Name:        "secrets",
		Usage:       "Provides CRUD operations for secrets",
		Subcommands: []*cli.Command{delete, generate, get, set},
	}
}

func secretsCommand(c *cli.Context, subcommand SecretsCommand) error {
	// Setup configuration
	var provider gcli.SecretProvider
	switch c.String(flag_backend) {
	case "aws":
		config, err := aws.NewSecretProviderConfig(c)
		if err != nil {
			return gcli.Exit(fmt.Sprintf("error configuring secret backend: %s", err))
		}

		p := aws.NewSecretProvider(config)
		provider = &p
	case "mock":
		provider = mockProvider
	default:
		return gcli.Exit(fmt.Sprintf("invalid backend: %s", c.String(flag_backend)))
	}

	// Handle subcommand
	switch subcommand {
	case cmdDelete:
		return delete(c, provider)
	case cmdGenerate:
		return generate(c, provider)
	case cmdGet:
		return get(c, provider)
	case cmdSet:
		return set(c, provider)
	default:
		return gcli.Exit("unknown subcommand")
	}
}

func delete(c *cli.Context, provider gcli.SecretProvider) error {
	if c.NArg() < 1 {
		return gcli.Exit("must provide a key")
	}

	err := provider.Delete(c.Args().First())
	if err != nil {
		return gcli.Exit(err.Error())
	}

	return nil
}

func generate(c *cli.Context, provider gcli.SecretProvider) error {
	if c.NArg() < 1 {
		return gcli.Exit("must provide a key")
	}

	err := provider.Generate(c.Args().First(), c.Int("length"), c.Int("numbers"), c.Int("symbols"))
	if err != nil {
		return gcli.Exit(err.Error())
	}

	return nil
}

func get(c *cli.Context, provider gcli.SecretProvider) error {
	if c.NArg() < 1 {
		return gcli.Exit("must provide a key")
	}

	value, err := provider.Get(c.Args().First())
	if err != nil {
		return gcli.Exit(err.Error())
	}

	fmt.Print(value)
	return nil
}

func set(c *cli.Context, provider gcli.SecretProvider) error {
	if c.NArg() < 2 {
		return gcli.Exit("must provide a key and value")
	}

	err := provider.Set(c.Args().Get(0), c.Args().Get(1))
	if err != nil {
		return gcli.Exit(err.Error())
	}

	return nil
}
