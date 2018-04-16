package toml_vault

import "fmt"

type GcpConfig struct {
	GcpKmsResourceID  string
	GcpCredentialFile string
}

type GcpCryptor struct {
	GcpConfig
}

func (c *GcpCryptor) Encrypt(value interface{}) (interface{}, error) {
	return c.Decrypt(value)
}

func (c *GcpCryptor) Decrypt(value interface{}) (interface{}, error) {
	fmt.Println("gcp cryptor")
	return value, nil
}
