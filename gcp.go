package nvault

import (
	"fmt"
)

type GcpCryptor struct {
	GcpConfig
}

type GcpConfig struct {
	GcpKmsResourceID  string
	GcpCredentialFile string
}

func (c *GcpCryptor) Encrypt(value interface{}) (interface{}, error) {
	return c.Decrypt(value)
}

func (c *GcpCryptor) Decrypt(value interface{}) (interface{}, error) {
	fmt.Println("simple cryptor")
	return value, nil
}
