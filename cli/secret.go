package cli

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var Test string = "test"

var ErrSecretNotFound = errors.New("secret not found")

// SecretProvider represents a backend capable of storing sensitive data using a
// key/value format.
type SecretProvider interface {
	// Deletes the secret with the given key
	Delete(key string) error

	// Generates a new random secret value with the given key. Overwrites any
	// previous value that existed with the key.
	Generate(key string, length int, nums int, symbols int) (string, error)

	// Returns the value of the secret with the given key
	Get(key string) (string, error)

	// Sets the value of the secret with the given key. Overwrites any previous
	// value that existed with the key.
	Set(key string, value string) error
}

type SecretProviderConfig interface {
	Flags() []cli.Flag
}
