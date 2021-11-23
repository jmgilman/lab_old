package http

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"golang.org/x/crypto/openpgp"
)

var baseURL string = "https://%s.release.flatcar-linux.net/%s-usr/current/%s"
var imageName string = "flatcar_production_image.bin.bz2"

// httpClient is an interface for processing HTTP requests and returning HTTP
// responses.
type httpClient interface {
	// Do sends an HTTP request and returns an HTTP response.
	Do(req *http.Request) (*http.Response, error)
}

// pgpClient is an interface for validating blobs of signed data using PGP keys.
type pgpClient interface {
	// ReadArmoredKeyRing reads one or more public/private keys from an armor
	// keyring file.
	ReadArmoredKeyRing(r io.Reader) (openpgp.EntityList, error)

	// CheckDetachedSignature takes a signed file and a detached signature and
	// returns the signer if the signature is valid. If the signer isn't known,
	// ErrUnknownIssuer is returned.
	CheckDetachedSignature(keyring openpgp.KeyRing, signed io.Reader, signature io.Reader) (signer *openpgp.Entity, err error)
}

// ImageProvider implements cli.ImageProvider using an httpClient and pgpClient.
type ImageProvider struct {
	httpClient httpClient
	pgpClient  pgpClient
}

// openpgpClient implements pgpClient using the openpgp package.
type openpgpClient struct{}

func (o *openpgpClient) ReadArmoredKeyRing(r io.Reader) (openpgp.EntityList, error) {
	return openpgp.ReadArmoredKeyRing(r)
}

func (o *openpgpClient) CheckDetachedSignature(keyring openpgp.KeyRing, signed io.Reader, signature io.Reader) (signer *openpgp.Entity, err error) {
	return openpgp.CheckDetachedSignature(keyring, signed, signature)
}

// buildUrl returns the fully qualified URL to the requested Container Linux
// production image file.
func (i *ImageProvider) buildURL(channel string, arch string) string {
	return fmt.Sprintf(baseURL, channel, arch, imageName)
}

// download downloads the remote file at the given URL, returning a stream of
// data and it's expected size.
func (i *ImageProvider) download(url string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	resp, err := i.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.ContentLength, nil
}

func (i *ImageProvider) Fetch(channel string, arch string) (io.ReadCloser, int64, error) {
	return i.download(i.buildURL(channel, arch))
}

func (i *ImageProvider) Validate(data io.ReadCloser, channel string, arch string) error {
	url := fmt.Sprintf("%s.sig", i.buildURL(channel, arch))

	sig, _, err := i.download(url)
	if err != nil {
		return err
	}

	keyring, err := i.pgpClient.ReadArmoredKeyRing(strings.NewReader(publicKey))
	if err != nil {
		return err
	}

	_, err = i.pgpClient.CheckDetachedSignature(keyring, data, sig)
	if err != nil {
		return gcli.ErrSigCheckFailed
	}

	return nil
}

func NewImageProvider() gcli.ImageProvider {
	return &ImageProvider{
		httpClient: &http.Client{},
		pgpClient:  &openpgpClient{},
	}
}
