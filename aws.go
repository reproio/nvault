package nvault

import "fmt"

type AwsConfig struct {
	AwsKmsKeyID        string
	AwsRegion          string
	AwsAccessKeyID     string
	AwsSecretAccessKey string
}

type AwsCryptor struct {
	AwsConfig
}

func (c *AwsCryptor) Encrypt(value interface{}) (interface{}, error) {
	return c.Decrypt(value)
}

func (c *AwsCryptor) Decrypt(value interface{}) (interface{}, error) {
	fmt.Println("aws cryptor")
	return value, nil
}
