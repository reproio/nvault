package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/reproio/nvault"
)

func ParseKeys(clikey string) (paths []nvault.Path, err error) {
	for _, key := range strings.Split(clikey, ",") {
		path, err := parseKey(key)
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	if len(paths) == 0 {
		return paths, errors.New("empty paths")
	}

	return paths, nil
}

type ScanType struct {
	Regexp *regexp.Regexp
	Type   string
}

var regs = []ScanType{
	ScanType{regexp.MustCompile(`^'(.*?)'`), "string"},
	ScanType{regexp.MustCompile(`^"(.*?)"`), "string"},
	ScanType{regexp.MustCompile(`^/(.*?)/`), "regexp"},
	ScanType{regexp.MustCompile(`^\s+`), "space"},
	ScanType{regexp.MustCompile(`^\[(\d*)\]`), "number"},
	ScanType{regexp.MustCompile(`^\.`), "delimiter"},
	ScanType{regexp.MustCompile(`^[^\.]+`), "string"},
}

func parseKey(key string) (nvault.Path, error) {
	s := bufio.NewScanner(strings.NewReader(key))

	s.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	Start:
		for _, st := range regs {
			loc := st.Regexp.FindSubmatchIndex(data)
			if len(loc) > 0 {
				if len(loc) > 3 {
					token = data[loc[2]:loc[3]]
				} else {
					token = data[loc[0]:loc[1]]
				}
				token = append([]byte(fmt.Sprintf("%s:", st.Type)), token...)
				advance = advance + loc[1]

				if st.Type == "delimiter" || st.Type == "space" {
					padding := loc[1] - loc[0]
					if len(data) >= padding {
						data = data[padding:]
						goto Start
					}
					break
				}
				return
			}
		}
		if atEOF && len(data) > advance {
			return len(data), data[advance:], nil
		}
		return advance, nil, nil
	})

	var fragments []string
	for s.Scan() {
		fragments = append(fragments, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("Invalid input: %s", err)
	}

	if len(fragments) == 0 {
		return nil, errors.New("no path fragments")
	}

	var path nvault.Path
	for _, fragment := range fragments {
		f := strings.SplitN(fragment, ":", 2)
		path = append(path, nvault.PathFragment{f[0], f[1]})
	}

	if path[0].Fragment != "$" {
		return nil, errors.New("`$` must be at first")
	}

	return path, nil
}
