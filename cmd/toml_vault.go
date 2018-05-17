package cmd

import (
	"io"

	"github.com/BurntSushi/toml"

	"github.com/reproio/nvault"
)

func TomlConverter(input io.Reader, output io.Writer, cryptor Cryptor) error {
	p := nvault.Placeholder{}

	_, err := toml.DecodeReader(input, &p)
	if err != nil {
		return err
	}

	if err := cryptor(&p); err != nil {
		return err
	}

	encoder := toml.NewEncoder(output)
	if err := encoder.Encode(&p); err != nil {
		return err
	}

	return nil
}
