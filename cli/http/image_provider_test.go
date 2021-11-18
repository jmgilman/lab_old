package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/matryer/is"
	"golang.org/x/crypto/openpgp"
)

type MockHTTPClient struct {
	fnDo func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.fnDo(req)
}

type MockPGPClient struct {
	fnReadArmoredKeyRing     func(r io.Reader) (openpgp.EntityList, error)
	fnCheckDetachedSignature func(keyring openpgp.KeyRing, signed io.Reader, signature io.Reader) (signer *openpgp.Entity, err error)
}

func (m *MockPGPClient) ReadArmoredKeyRing(r io.Reader) (openpgp.EntityList, error) {
	return m.fnReadArmoredKeyRing(r)
}

func (m *MockPGPClient) CheckDetachedSignature(keyring openpgp.KeyRing, signed io.Reader, signature io.Reader) (signer *openpgp.Entity, err error) {
	return m.fnCheckDetachedSignature(keyring, signed, signature)
}

func TestBuildURL(t *testing.T) {
	is := is.New(t)
	expected := fmt.Sprintf(baseURL, "alpha", "arm64", imageName)

	fetcher := ImageProvider{
		httpClient: &MockHTTPClient{},
	}
	got := fetcher.buildURL("alpha", "arm64")
	is.Equal(expected, got)
}

func TestDownload(t *testing.T) {
	is := is.New(t)
	expected_data := "test"
	expected_size := 1024
	expected_url := "url"

	// With no error
	var got_url string
	mock := MockHTTPClient{
		fnDo: func(req *http.Request) (*http.Response, error) {
			got_url = req.URL.String()
			return &http.Response{
				Body:          io.NopCloser(strings.NewReader(expected_data)),
				ContentLength: int64(expected_size),
			}, nil
		},
	}

	fetcher := ImageProvider{
		httpClient: &mock,
	}
	res, size, err := fetcher.download(expected_url)
	is.NoErr(err)
	is.Equal(expected_url, got_url)
	is.Equal(int64(expected_size), size)

	res_data, err := io.ReadAll(res)
	is.NoErr(err)
	is.Equal(expected_data, string(res_data))

	// With error
	mock = MockHTTPClient{
		fnDo: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("error")
		},
	}

	fetcher = ImageProvider{
		httpClient: &mock,
	}
	_, _, err = fetcher.download(expected_url)
	is.Equal(err.Error(), "error")
}

func TestFetch(t *testing.T) {
	is := is.New(t)
	expected_url := fmt.Sprintf(baseURL, "alpha", "arm64", imageName)

	var got_url string
	mock := MockHTTPClient{
		fnDo: func(req *http.Request) (*http.Response, error) {
			got_url = req.URL.String()
			return &http.Response{
				Body: io.NopCloser(strings.NewReader("test")),
			}, nil
		},
	}

	fetcher := ImageProvider{
		httpClient: &mock,
	}
	_, _, err := fetcher.Fetch("alpha", "arm64")
	is.NoErr(err)
	is.Equal(expected_url, got_url)
}

func TestValidate(t *testing.T) {
	is := is.New(t)
	expected_url := fmt.Sprintf("%s.sig", fmt.Sprintf(baseURL, "alpha", "arm64", imageName))
	expected_pub_key := publicKey
	expected_data := "test"
	expected_sig_data := "testsignature"

	var got_url string
	mock_http := MockHTTPClient{
		fnDo: func(req *http.Request) (*http.Response, error) {
			got_url = req.URL.String()
			return &http.Response{
				Body: io.NopCloser(strings.NewReader(expected_sig_data)),
			}, nil
		},
	}

	var got_pub_key string
	var got_data string
	var got_sig_data string
	mock_pgp := MockPGPClient{
		fnReadArmoredKeyRing: func(r io.Reader) (openpgp.EntityList, error) {
			data, _ := io.ReadAll(r)
			got_pub_key = string(data)
			return openpgp.EntityList{}, nil
		},
		fnCheckDetachedSignature: func(keyring openpgp.KeyRing, signed, signature io.Reader) (signer *openpgp.Entity, err error) {
			data, _ := io.ReadAll(signed)
			got_data = string(data)

			data, _ = io.ReadAll(signature)
			got_sig_data = string(data)

			return nil, nil
		},
	}

	// With no error
	fetcher := ImageProvider{
		httpClient: &mock_http,
		pgpClient:  &mock_pgp,
	}
	err := fetcher.Validate(io.NopCloser(strings.NewReader(expected_data)), "alpha", "arm64")
	is.NoErr(err)
	is.Equal(expected_url, got_url)
	is.Equal(expected_pub_key, got_pub_key)
	is.Equal(expected_data, got_data)
	is.Equal(expected_sig_data, got_sig_data)

	// With failed validation
	mock_pgp.fnCheckDetachedSignature = func(keyring openpgp.KeyRing, signed, signature io.Reader) (signer *openpgp.Entity, err error) {
		return nil, fmt.Errorf("failed")
	}

	fetcher = ImageProvider{
		httpClient: &mock_http,
		pgpClient:  &mock_pgp,
	}
	err = fetcher.Validate(io.NopCloser(strings.NewReader(expected_data)), "alpha", "arm64")
	is.True(errors.Is(err, gcli.ErrSigCheckFailed))
}
