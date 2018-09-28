package nvault

// NoneCryptor ...
type NoneCryptor struct {
	NoneConfig
}

// NoneConfig ...
type NoneConfig struct {
}

// Encrypt ...
func (c *NoneCryptor) Encrypt(value interface{}) (interface{}, error) {
	return value, nil
}

// Decrypt ...
func (c *NoneCryptor) Decrypt(value interface{}) (interface{}, error) {
	return value, nil
}
