package nvault_test

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
				{
					nvault.PathFragment{Type: "string", Fragment: "$"},
					nvault.PathFragment{Type: "string", Fragment: "test"},
					nvault.PathFragment{Type: "number", Fragment: "0"},
					nvault.PathFragment{Type: "regexp", Fragment: "test"},
				},
				{
					nvault.PathFragment{Type: "string", Fragment: "$"},
					nvault.PathFragment{Type: "regexp", Fragment: "test"},
					nvault.PathFragment{Type: "number", Fragment: "0"},
					nvault.PathFragment{Type: "string", Fragment: "test"},
				},
			},
			"",
		},
		{`"has not root"`, nil, "`$` must be at first"},
	}

	for _, test := range tests {
		paths, err := nvault.ParseKeys(test.clikey)

		if test.err != "" {
			if err == nil {
				t.Errorf("unexpected: clikey %v returns %v expected %v", test.clikey, err, test.err)
			}
			continue
		}

		for i, expected := range test.expected {
			if !expected.Equal(paths[i]) {
				t.Errorf("unexpected: clikey %v returns %v expected %v", test.clikey, paths[i], expected)
			}
		}
	}
}
