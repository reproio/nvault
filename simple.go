package toml_vault

import "fmt"

type SimpleCryptor struct {
}

func (c *SimpleCryptor) Encrypt(value interface{}) (interface{}, error) {
	return c.Decrypt(value)
}

func (c *SimpleCryptor) Decrypt(value interface{}) (interface{}, error) {
	fmt.Println("simple cryptor")
	return value, nil
}
