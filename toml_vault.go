package toml_vault

import (
	"errors"
	"os"

	"github.com/manifoldco/promptui"
)

type Cryptor interface {
	Encryptor
	Decryptor
}

type Config struct {
	AwsConfig
	GcpConfig

	Cryptor string

	Salt            string
	Cipher          string
	KeyLen          int
	Digest          string
	SignatureKeyLen int

	UseSignPassphrase bool
	Passphrase        string
	SignPassphrase    string
}

func (c *Config) GetPassphrase() error {
	if c.Cryptor == "simple" {
		return nil
	}

	var err error

	passphrase := os.Getenv("YAML_VAULT_PASSPHRASE")
	if passphrase == "" {
		prompt := promptui.Prompt{
			Label: "Enter passphrase",
			Validate: func(input string) error {
				if input == "" {
					return errors.New("Please input passphrase")
				}
				return nil
			},
		}
		passphrase, err = prompt.Run()
		if err != nil {
			return err
		}
	}
	c.Passphrase = passphrase

	signPassphrase := os.Getenv("YAML_VAULT_SIGN_PASSPHRASE")
	if signPassphrase == "" && c.UseSignPassphrase {
		prompt := promptui.Prompt{
			Label: "Enter sign passphrase",
			Validate: func(input string) error {
				if input == "" {
					return errors.New("Please input sign passphrase")
				}
				return nil
			},
		}
		signPassphrase, err = prompt.Run()
		if err != nil {
			return err
		}
	}
	c.SignPassphrase = signPassphrase

	return nil
}
func Encrypt(config Config, p *Placeholder, keys []Key) error {
	_, err := NewEncryptor(config)
	if err != nil {
		return err
	}
	return nil
}

func Decrypt(config Config, p *Placeholder, keys []Key) error {
	_, err := NewEncryptor(config)
	if err != nil {
		return err
	}
	return nil
}
