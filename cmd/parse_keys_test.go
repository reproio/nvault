package cmd

import (
	"testing"

	"github.com/reproio/nvault"
)

func TestParseKeys(t *testing.T) {
	tests := []struct {
		clikey   string
		expected []nvault.Path
		err      string
	}{
		{
			"$.test.[0]./test/,$./test/.[0].test",
			[]nvault.Path{
				nvault.Path{
					nvault.PathFragment{"string", "$"},
					nvault.PathFragment{"string", "test"},
					nvault.PathFragment{"number", "0"},
					nvault.PathFragment{"regexp", "test"},
				},
				nvault.Path{
					nvault.PathFragment{"string", "$"},
					nvault.PathFragment{"regexp", "test"},
					nvault.PathFragment{"number", "0"},
					nvault.PathFragment{"string", "test"},
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
