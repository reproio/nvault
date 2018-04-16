package toml_vault

type SimpleCryptor struct {
}

func (c *SimpleCryptor) Encrypt(value interface{}) (interface{}, error) {
	return "Encrypted", nil
}

func (c *SimpleCryptor) Decrypt(value interface{}) (interface{}, error) {
	return "Decrypted", nil
}
