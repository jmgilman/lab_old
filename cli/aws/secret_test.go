package aws

import (
	"errors"
	"flag"
	"fmt"
	"testing"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/matryer/is"
	"github.com/urfave/cli/v2"
)

type mockSSM struct {
	ssmiface.SSMAPI
	fnDelete func(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error)
	fnGet    func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error)
	fnPut    func(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error)
}

func (m *mockSSM) DeleteParameter(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
	return m.fnDelete(input)
}

func (m *mockSSM) GetParameter(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return m.fnGet(input)
}

func (m *mockSSM) PutParameter(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
	return m.fnPut(input)
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	expected := "test"

	// With no error
	var got string
	provider := SecretProvider{
		ssm: &mockSSM{
			fnDelete: func(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
				got = *input.Name
				return nil, nil
			},
		},
	}

	err := provider.Delete(expected)
	is.NoErr(err)
	is.Equal(got, expected)

	// With key error
	provider = SecretProvider{
		ssm: &mockSSM{
			fnDelete: func(input *ssm.DeleteParameterInput) (*ssm.DeleteParameterOutput, error) {
				return nil, awserr.New(ssm.ErrCodeParameterNotFound, "", fmt.Errorf(""))
			},
		},
	}

	err = provider.Delete(expected)
	is.True(errors.Is(err, gcli.ErrSecretNotFound))
}

func TestGenerate(t *testing.T) {
	is := is.New(t)
	expected_key := "test"
	expected_value := "test"

	// Password generation works
	var got_key string
	var got_value string
	provider := SecretProvider{
		generator: func(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
			return expected_value, nil
		},
		ssm: &mockSSM{
			fnPut: func(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
				got_key = *input.Name
				got_value = *input.Value
				return nil, nil
			},
		},
	}

	result, err := provider.Generate(expected_key, 0, 0, 0)
	is.NoErr(err)
	is.Equal(got_key, expected_key)
	is.Equal(got_value, expected_value)
	is.Equal(result, expected_value)

	// With password generation error
	provider = SecretProvider{
		generator: func(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
			return "", fmt.Errorf("failed")
		},
		ssm: &mockSSM{
			fnPut: func(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
				return nil, nil
			},
		},
	}

	_, err = provider.Generate(expected_key, 0, 0, 0)
	is.Equal(err.Error(), "failed generating random password")

	// With SSM error
	provider = SecretProvider{
		generator: func(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
			return "", nil
		},
		ssm: &mockSSM{
			fnPut: func(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
				return nil, fmt.Errorf("failed")
			},
		},
	}

	_, err = provider.Generate(expected_key, 0, 0, 0)
	is.Equal(err.Error(), "error querying AWS: failed")
}

func TestGet(t *testing.T) {
	is := is.New(t)
	expected_key := "test"
	expected_value := "test"

	// With no error
	var got_key string
	provider := SecretProvider{
		ssm: &mockSSM{
			fnGet: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				got_key = *input.Name
				return &ssm.GetParameterOutput{
					Parameter: &ssm.Parameter{
						Value: aws.String(expected_value),
					},
				}, nil
			},
		},
	}

	got_value, err := provider.Get(expected_key)
	is.NoErr(err)
	is.Equal(expected_key, got_key)
	is.Equal(expected_value, got_value)

	// With key error
	provider = SecretProvider{
		ssm: &mockSSM{
			fnGet: func(input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
				return nil, awserr.New(ssm.ErrCodeParameterNotFound, "", fmt.Errorf(""))
			},
		},
	}

	_, err = provider.Get(expected_key)
	is.True(errors.Is(err, gcli.ErrSecretNotFound))
}

func TestPut(t *testing.T) {
	is := is.New(t)
	expected_key := "test"
	expected_value := "test"

	// With no error
	var got_key string
	var got_value string
	provider := SecretProvider{
		ssm: &mockSSM{
			fnPut: func(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
				got_key = *input.Name
				got_value = *input.Value
				return nil, nil
			},
		},
	}

	err := provider.Set(expected_key, expected_value)
	is.NoErr(err)
	is.Equal(expected_key, got_key)
	is.Equal(expected_value, got_value)

	// With SSM error
	provider = SecretProvider{
		generator: func(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
			return "", nil
		},
		ssm: &mockSSM{
			fnPut: func(input *ssm.PutParameterInput) (*ssm.PutParameterOutput, error) {
				return nil, fmt.Errorf("failed")
			},
		},
	}

	err = provider.Set(expected_key, expected_value)
	is.Equal(err.Error(), "error querying AWS: failed")
}

func TestNewSecretProviderConfig(t *testing.T) {
	is := is.New(t)

	// With invalid flag configuration
	set := flag.NewFlagSet("test", 0)
	set.String(flag_access_key, "test", "test")
	_ = set.Parse([]string{fmt.Sprintf("--%s", flag_access_key), "test"})

	ctx := cli.NewContext(&cli.App{}, set, nil)
	_, err := NewSecretProviderConfig(ctx)
	fmt.Println("Error: ", err)
	is.Equal(err.Error(), "must supply both access and secret keys")

	// With all flags
	set = flag.NewFlagSet("test", 0)
	set.String(flag_access_key, "test", "test")
	set.String(flag_secret_key, "test", "test")
	set.String(flag_region, "test", "test")
	set.Bool("verbose", true, "test")
	_ = set.Parse([]string{fmt.Sprintf("--%s", flag_access_key), "test", fmt.Sprintf("--%s", flag_secret_key), "test", fmt.Sprintf("--%s", flag_region), "test"})

	ctx = cli.NewContext(&cli.App{}, set, nil)
	result, err := NewSecretProviderConfig(ctx)
	is.NoErr(err)

	creds, err := result.config.Credentials.Get()
	is.Equal(creds.AccessKeyID, "test")
	is.Equal(creds.SecretAccessKey, "test")
	is.Equal(*result.config.Region, "test")
	is.Equal(*result.config.LogLevel, aws.LogDebug)
}
