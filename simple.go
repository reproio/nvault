package nvault

type SimpleConfig struct {
	Salt            string
	Cipher          string
	KeyLen          int
	Digest          string
	SignatureKeyLen int

	UseSignPassphrase bool
	Passphrase        string
	SignPassphrase    string
}

type SimpleCryptor struct {
	SimpleConfig
}

func (c *SimpleCryptor) Encrypt(value interface{}) (interface{}, error) {
	return "Encrypted", nil
}

func (c *SimpleCryptor) Decrypt(value interface{}) (interface{}, error) {
	return "Decrypted", nil
}
