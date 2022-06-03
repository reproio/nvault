package nvault

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// AwsCryptor ...
type AwsCryptor struct {
	AwsConfig
}

// AwsConfig ...
type AwsConfig struct {
	AwsKmsKeyID        string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

// Encrypt ...
func (c *AwsCryptor) Encrypt(value interface{}) (interface{}, error) {
	if c.AwsKmsKeyID == "" {
		return nil, errors.New("missing Aws KMS Key ID")
	}
	strvalue := fmt.Sprintf("%v", value)

	output, err := serviceAws(&c.AwsConfig).Encrypt(&kms.EncryptInput{
		KeyId:     aws.String(c.AwsKmsKeyID),
		Plaintext: []byte(strvalue),
	})
	if err != nil {
		return value, nil
	}

	encoded := base64.StdEncoding.EncodeToString(output.CiphertextBlob)
	return encoded, nil
}

// Decrypt ...
func (c *AwsCryptor) Decrypt(value interface{}) (interface{}, error) {
	strvalue := fmt.Sprintf("%v", value)

	decoded, err := base64.StdEncoding.DecodeString(strvalue)
	if err != nil {
		return value, err
	}

	output, err := serviceAws(&c.AwsConfig).Decrypt(&kms.DecryptInput{
		CiphertextBlob: decoded,
	})
	if err != nil {
		return value, err
	}

	return string(output.Plaintext), nil
}

func serviceAws(c *AwsConfig) *kms.KMS {
	config := &aws.Config{}

	if c.AwsRegion != "" {
		config.Region = &c.AwsRegion
	}

	if c.AwsAccessKeyID != "" && c.AwsSecretAccessKey != "" {
		config.Credentials = createAwsCredentials(c)
	}
	return kms.New(session.New(config))
}

func createAwsCredentials(c *AwsConfig) *credentials.Credentials {
	defaultProvider := defaults.RemoteCredProvider(
		aws.Config{Region: &c.AwsRegion},
		defaults.Handlers(),
	)
	envProvider := &credentials.EnvProvider{}

	providers := []credentials.Provider{
		defaultProvider,
		envProvider,
	}

	if c.AwsAccessKeyID != "" && c.AwsSecretAccessKey != "" {
		providers = append(providers, &credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     c.AwsAccessKeyID,
				SecretAccessKey: c.AwsSecretAccessKey,
				SessionToken:    "",
				ProviderName:    "",
			},
		})
	}

	return credentials.NewChainCredentials(providers)
}

// WithAwsCredential ...
func WithAwsCredential(awsAccessKeyID, awsSecretAccessKey string) Option {
	return func(c *Config) error {
		c.AwsAccessKeyID = awsAccessKeyID
		c.AwsSecretAccessKey = awsSecretAccessKey
		return nil
	}
}

// WithAwsRegion ...
func WithAwsRegion(awsRegion string) Option {
	return func(c *Config) error {
		c.AwsRegion = awsRegion
		return nil
	}
}
