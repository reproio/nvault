package main

import (
	"encoding/json"
	"io"

	"github.com/reproio/nvault"
	"github.com/reproio/nvault/cmd"
)

func main() {
	cmd.Run(func(input io.Reader, output io.Writer, cryptor cmd.Cryptor) error {
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
	})
}
