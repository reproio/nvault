package toml_vault

import (
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"

	"github.com/reproio/toml_vault"
)

func TestKeyParserParse(t *testing.T) {
	tests := []struct {
		clikey   string
		metaKeys []toml.Key
		expected []toml_vault.Key
		err string
	}{
		{"$", []toml.Key{toml.Key{"test"}}, []toml_vault.Key{toml_vault.Key{"test"}}, ""},
		{"", []toml.Key{toml.Key{"test"}}, []toml_vault.Key{toml_vault.Key{"test"}}, ""},
	}
	for _, test := range tests {
		parsedKey, err := toml_vault.ParseKeys(test.clikey, test.metaKeys)

		if test.err != "" {
			if err != nil {
				t.Errorf("unexpected: clikey %v metaKeys %v returns %v expeced %v", test.clikey, test.metaKeys, err, test.err)
			} else {
				t.Errorf("unexpected: clikey %v metaKeys %v returns no error expeced %v", test.clikey, test.metaKeys, test.expected)
			}
			continue
		}

		if reflect.DeepEqual(parsedKey, test.expected) {
			t.Errorf("unexpected: clikey %v metaKeys %v returns %v expeced %v", test.clikey, test.metaKeys, parsedKey, test.expected)
		}
	}
}
