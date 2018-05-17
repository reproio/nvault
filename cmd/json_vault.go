package cmd

import (
	"encoding/json"
	"io"

	"github.com/reproio/nvault"
)

func JsonConverter(input io.Reader, output io.Writer, cryptor Cryptor) error {
	p := nvault.Placeholder{}

	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&p); err != nil {
		return err
	}

	if err := cryptor(&p); err != nil {
		return err
	}

	encoder := json.NewEncoder(output)
	if err := encoder.Encode(&p); err != nil {
		return err
	}

	return nil
}
