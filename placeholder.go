package toml_vault

import (
	"github.com/BurntSushi/toml"
)

type Placeholder map[string]interface{}

type Key toml.Key

func (p *Placeholder) Set(key Key, value interface{}) error {
	if len(key) == 0 {
		return nil
	}
	switch pi := (*p)[key[0]].(type) {
	case map[string]interface{}:
		p := Placeholder(pi)
		return p.Set(key[1:], value)
	default:
		(*p)[key[0]] = value
	}
	return nil
}

func (p *Placeholder) Get(key Key) (interface{}, error) {
	switch pi := (*p)[key[0]].(type) {
	case Placeholder:
		g, err := pi.Get(key[1:])
		return &g, err
	default:
		return &pi, nil
	}
}
