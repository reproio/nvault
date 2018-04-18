package main

import (
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/reproio/nvault"
	"github.com/reproio/nvault/cmd"
)

func main() {
	cmd.Run(func(input io.Reader, output io.Writer, cryptor cmd.Cryptor) error {
		p := nvault.Placeholder{}

		data, err := ioutil.ReadAll(input)
		if err != nil {
			return err
		}

		if err := yaml.Unmarshal(data, &p); err != nil {
			return err
		}

		if err := cryptor(&p); err != nil {
			return err
		}

		data, err = yaml.Marshal(&p)
		if err != nil {
			return err
		}

		_, err = output.Write(data)
		return err
	})
}
