package cmd

import (
	"io"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/reproio/nvault"
)

func YamlConverter(input io.Reader, output io.Writer, cryptor Cryptor) error {
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
}
