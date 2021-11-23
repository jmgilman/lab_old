package aws

import (
	"fmt"

	gcli "github.com/HomeOperations/jmgilman/cli"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/sethvargo/go-password/password"
	"github.com/urfave/cli/v2"
)

const (
	flag_access_key = "aws-access-key"
	flag_secret_key = "aws-secret-key"
	flag_region     = "aws-region"
	flag_kms_key    = "aws-kms-key"
)

// SecretProvider implements bootstrap.SecretProvider using the AWS SSM
// parameter store as the backend. All parameters are configured to be encrypted
// using the default KMS key associated with the account executing requests.
// Credentials can be configured through the default AWS environment variables.
type SecretProvider struct {
	generator PasswordGenerator
	ssm       ssmiface.SSMAPI
}

// PasswordGenerator is a function which can be used for generating passwords.
type PasswordGenerator func(length int, numDigits int, numSymbols int, noUpper bool, allowRepeat bool) (string, error)

// Delete deletes the secret with the given key
func (s *SecretProvider) Delete(key string) error {
	in := ssm.DeleteParameterInput{
		Name: &key,
	}

	_, err := s.ssm.DeleteParameter(&in)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == ssm.ErrCodeParameterNotFound {
			return gcli.ErrSecretNotFound
		}

		return fmt.Errorf("error querying AWS: %s", err)
	}

	return nil
}

// Generate generates a new random secret value with the given key. Overwrites
// any previous value that existed with the key.
func (s *SecretProvider) Generate(key string, length int, nums int, symbols int) (string, error) {
	res, err := s.generator(length, nums, symbols, false, false)
	if err != nil {
		return "", fmt.Errorf("failed generating random password")
	}

	in := ssm.PutParameterInput{
		Name:      &key,
		Value:     &res,
		Type:      aws.String("SecureString"),
		Overwrite: aws.Bool(true),
	}

	_, err = s.ssm.PutParameter(&in)
	if err != nil {
		return "", fmt.Errorf("error querying AWS: %s", err)
	}

	return res, nil
}

// Get returns the value of the secret with the given key.
func (s *SecretProvider) Get(key string) (string, error) {
	in := ssm.GetParameterInput{
		Name:           &key,
		WithDecryption: aws.Bool(true),
	}

	out, err := s.ssm.GetParameter(&in)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == ssm.ErrCodeParameterNotFound {
			return "", gcli.ErrSecretNotFound
		}

		return "", fmt.Errorf("error querying AWS: %s", err)
	}

	return *out.Parameter.Value, nil
}

// Set sets the value of the secret with the given key. Overwrites any previous
// value that existed with the key.
func (s *SecretProvider) Set(key string, value string) error {
	in := ssm.PutParameterInput{
		Name:      &key,
		Value:     &value,
		Type:      aws.String("SecureString"),
		Overwrite: aws.Bool(true),
	}

	_, err := s.ssm.PutParameter(&in)
	if err != nil {
		return fmt.Errorf("error querying AWS: %s", err)
	}

	return nil
}

// NewSecretProvider creates a new instance of SecretProvider using the given
// configuration.
func NewSecretProvider(config SecretProviderConfig) SecretProvider {
	sess := session.Must(session.NewSession(config.config))
	ssm := ssm.New(sess)

	return SecretProvider{
		generator: password.Generate,
		ssm:       ssm,
	}
}

// SecretProviderConfig provides the configuration details needed for
// instantiating a new SecretProvider.
type SecretProviderConfig struct {
	config *aws.Config
}

// Flags returns the CLI flags that can be used to configure the AWS secret
// provider.
func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  flag_access_key,
			Usage: "AWS access key ID (default: $AWS_ACCESS_KEY_ID)",
		},
		&cli.StringFlag{
			Name:  flag_secret_key,
			Usage: "AWS secret key ID (default: $AWS_SECRET_ACCESS_KEY)",
		},
		&cli.StringFlag{
			Name:  flag_region,
			Usage: "AWS region (default: $AWS_REGION)",
		},
		&cli.StringFlag{
			Name:  flag_kms_key,
			Usage: "KMS key ID to use for encryption (defaults to account default)",
		},
	}
}

// NewSecretProviderConfig creates a new SecretProviderConfig by parsing CLI
// flags contained within the passed cli.Context.
func NewSecretProviderConfig(c *cli.Context) (SecretProviderConfig, error) {
	config := aws.NewConfig()

	if c.IsSet(flag_access_key) || c.IsSet(flag_secret_key) {
		if !(c.IsSet(flag_access_key) && c.IsSet(flag_secret_key)) {
			return SecretProviderConfig{}, fmt.Errorf("must supply both access and secret keys")
		}

		creds := credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     c.String(flag_access_key),
				SecretAccessKey: c.String(flag_secret_key),
				ProviderName:    "flags",
			},
		})

		config = config.WithCredentials(creds)
	}

	if c.IsSet(flag_region) {
		config.Region = aws.String(c.String(flag_region))
	}

	if c.Bool("verbose") {
		config.LogLevel = aws.LogLevel(aws.LogDebug)
	}

	return SecretProviderConfig{
		config: config,
	}, nil
}
