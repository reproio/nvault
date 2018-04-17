package mapcryptor

import (
	"github.com/mitchellh/mapstructure"

	"github.com/reproio/nvault"
)

func Encrypt(p *nvault.Placeholder, s interface{}, config *nvault.Config, paths ...nvault.Path) error {
	err := nvault.Encrypt(p, config, paths...)
	if err != nil {
		return err
	}

	return mapstructure.Decode(p, s)
}

func Decrypt(p *nvault.Placeholder, s interface{}, config *nvault.Config, paths ...nvault.Path) error {
	err := nvault.Decrypt(p, config, paths...)
	if err != nil {
		return err
	}

	return mapstructure.Decode(p, s)
}
