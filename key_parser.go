package toml_vault

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/BurntSushi/toml"
)

func ParseKeys(clikey string, metaKeys []toml.Key) (keys []Key, err error) {
	if clikey != "" {
		for _, key := range strings.Split(clikey, ",") {
			kp := keyParser{key}
			keys = append(keys, kp.Parse())
		}
	} else {
		for _, key := range metaKeys {
			keys = append(keys, Key(key))
		}
	}
	return keys, nil
}

type keyParser struct {
	origin string
}

func (kp keyParser) Parse() Key {
	var s scanner.Scanner
	s.Init(strings.NewReader(kp.origin))
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: %s\n", s.Position, s.TokenText())
	}

	return Key{"hoge", "huga"}
}
