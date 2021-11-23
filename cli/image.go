package cli

import (
	"errors"
	"io"
)

var ErrSigCheckFailed = errors.New("signature check failed")

type ImageProvider interface {
	// Fetch returns a network stream containing the contents of the production
	// Container Linux image at the given channel for the given architecture.
	Fetch(channel, arch, filename string) (io.ReadCloser, int64, error)

	// Validate takes a stream containing a Container Linux image and validates it
	// against the remote PGP signature for the given channel and architecture.
	Validate(data io.ReadCloser, channel, arch, filename string) error
}
