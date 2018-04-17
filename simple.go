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

type SimpleCryptor struct {
	SimpleConfig
}

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

func (c *SimpleCryptor) Encrypt(value interface{}) (interface{}, error) {
	return "Encrypted", nil
}

func (c *SimpleCryptor) Decrypt(value interface{}) (interface{}, error) {
	return "Decrypted", nil
}
