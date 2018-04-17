package nvault

import (
	"errors"
	"os"

	"github.com/manifoldco/promptui"
)

type Config struct {
	SimpleConfig
	AwsConfig
	GcpConfig

	Cryptor string
}

func (c *Config) GetPassphrase() error {
	if c.Cryptor != "simple" {
		return nil
	}

	var err error

	passphrase := os.Getenv("YAML_VAULT_PASSPHRASE")
	if passphrase == "" {
		prompt := promptui.Prompt{
			Label: "Enter passphrase",
			Mask:  '*',
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
			Mask:  '*',
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
