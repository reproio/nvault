package cmd

import (
	"testing"

	"github.com/reproio/toml_vault"
)

func TestParseKeys(t *testing.T) {
	tests := []struct {
		clikey   string
		expected []toml_vault.Path
		err      string
	}{
		{
			"$.test.[0]./test/,$./test/.[0].test",
			[]toml_vault.Path{
				toml_vault.Path{
					toml_vault.PathFragment{"string", "$"},
					toml_vault.PathFragment{"string", "test"},
					toml_vault.PathFragment{"number", "0"},
					toml_vault.PathFragment{"regexp", "test"},
				},
				toml_vault.Path{
					toml_vault.PathFragment{"string", "$"},
					toml_vault.PathFragment{"regexp", "test"},
					toml_vault.PathFragment{"number", "0"},
					toml_vault.PathFragment{"string", "test"},
				},
			},
			"",
		},
		{`"has not root"`, nil, "`$` must be at first"},
	}

	for _, test := range tests {
		paths, err := ParseKeys(test.clikey)

		if test.err != "" {
			if err == nil {
				t.Errorf("unexpected: clikey %v returns %v expeced %v", test.clikey, err, test.err)
			}
			continue
		}

		for i, expected := range test.expected {
			if !expected.Equal(paths[i]) {
				t.Errorf("unexpected: clikey %v returns %v expeced %v", test.clikey, paths[i], expected)
			}
		}
	}
}
