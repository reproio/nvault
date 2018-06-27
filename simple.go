package nvault

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"os"
// )

// SimpleCryptor ...
type SimpleCryptor struct {
	SimpleConfig
}

// SimpleConfig ...
type SimpleConfig struct {
	Salt            string
	Cipher          string
	KeyLen          int
	Digest          string
	SignatureKeyLen int

	Passphrase        string
	UseSignPassphrase bool
	SignPassphrase    string
}

// Encrypt ...
func (c *SimpleCryptor) Encrypt(value interface{}) (interface{}, error) {
	return "Encrypted", nil
}

// Decrypt ...
func (c *SimpleCryptor) Decrypt(value interface{}) (interface{}, error) {
	return "Decrypted", nil
}
