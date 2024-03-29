package nvault

import (
	"fmt"
	"regexp"

	"github.com/mitchellh/mapstructure"
)

// PathFragment ...
type PathFragment struct {
	Type     string
	Fragment string
}

func (pf PathFragment) String() string {
	return fmt.Sprintf("%s: %s", pf.Fragment, pf.Type)
}

// Path ...
type Path []PathFragment

// Equal ...
func (p Path) Equal(other Path) bool {
	for i, left := range p {
		if left != other[i] {
			return false
		}
	}
	return true
}

// AddRoot ...
func (p Path) AddRoot(root PathFragment) Path {
	return append(Path{root}, p...)
}

func (p Path) String() (s string) {
	for _, f := range p {
		s = fmt.Sprintf("%s/%s", s, f.Fragment)
	}
	return
}

// Match ...
func (p Path) Match(search Path) (result bool) {
	if search[0].Fragment != "$" {
		return
	}

	for i, fragment := range search[1:] {
		if len(p) <= i {
			break
		}

		if fragment.Fragment == "*" {
			result = true
		} else {
			switch fragment.Type {
			case "string":
				result = p[i].Fragment == fragment.Fragment
			case "regexp":
				result, _ = regexp.MatchString(p[i].Fragment, fragment.Fragment)
			case "number":
				result = fragment.Fragment == p[i].Fragment
			}
		}

		if !result {
			break
		}
	}
	return
}

// Placeholder ...
type Placeholder map[string]interface{}

// Matches ...
func (p Placeholder) Matches(search []Path) (results []Path) {
	for _, path := range p.Paths() {
		for _, s := range search {
			if path.Match(s) {
				results = append(results, path)
			}
		}
	}
	return
}

// Paths ...
func (p Placeholder) Paths() (paths []Path) {
	for f, v := range p {
		root := PathFragment{"string", f}
		switch vi := v.(type) {
		case map[string]interface{}:
			v := Placeholder(vi)
			for _, path := range v.Paths() {
				path := path.AddRoot(root)
				paths = append(paths, path)
			}
		case []map[string]interface{}:
			v := Placeholders(vi)
			for _, path := range v.Paths() {
				path := path.AddRoot(root)
				paths = append(paths, path)
			}
		default:
			paths = append(paths, Path{root})
		}
	}
	return paths
}

// Set ...
func (p Placeholder) Set(path Path, value interface{}) error {
	if len(path) == 0 {
		return nil
	}
	switch pi := p[path[0].Fragment].(type) {
	case map[string]interface{}:
		p := Placeholder(pi)
		return p.Set(Path(path[1:]), value)
	case []map[string]interface{}:
		p := Placeholders(pi)
		return p.Set(Path(path[1:]), value)
	default:
		p[path[0].Fragment] = value
	}
	return nil
}

// Get ...
func (p Placeholder) Get(path Path) (interface{}, error) {
	if len(path) == 0 {
		return p, nil
	}
	switch pi := p[path[0].Fragment].(type) {
	case map[string]interface{}:
		p := Placeholder(pi)
		return p.Get(Path(path[1:]))
	case []map[string]interface{}:
		p := Placeholders(pi)
		return p.Get(Path(path[1:]))
	default:
		return pi, nil
	}
}

// Placeholders ...
type Placeholders []map[string]interface{}

// Matches ...
func (ps Placeholders) Matches(search []Path) (results []Path) {
	if len(search) == 0 {
		return ps.Paths()
	}

	for _, path := range ps.Paths() {
		for _, s := range search {
			if path.Match(s) {
				results = append(results, path)
			}
		}
	}
	return
}

// Paths ...
func (ps Placeholders) Paths() (paths []Path) {
	for i, p := range ps {
		root := PathFragment{"string", fmt.Sprintf("%d", i)}
		p := Placeholder(p)
		for _, path := range p.Paths() {
			path := path.AddRoot(root)
			paths = append(paths, path)
		}
	}
	return paths
}

// Set ...
func (ps Placeholders) Set(path Path, value interface{}) error {
	if len(path) == 0 {
		return nil
	}
	for i, p := range ps {
		if path[0].Fragment != fmt.Sprintf("%d", i) {
			continue
		}
		p := Placeholder(p)
		return p.Set(Path(path[1:]), value)
	}
	return nil
}

// Get ...
func (ps Placeholders) Get(path Path) (interface{}, error) {
	if len(path) == 0 {
		return ps, nil
	}
	for i, p := range ps {
		if path[0].Fragment != fmt.Sprintf("%d", i) {
			continue
		}
		p := Placeholder(p)
		return p.Get(Path(path[1:]))
	}
	return nil, nil
}

// Encrypt ...
func (p *Placeholder) Encrypt(s interface{}, config *Config, paths ...Path) error {
	err := Encrypt(p, config, paths...)
	if err != nil {
		return err
	}

	return mapstructure.Decode(p, s)
}

// Decrypt ...
func (p *Placeholder) Decrypt(s interface{}, config *Config, paths ...Path) error {
	err := Decrypt(p, config, paths...)
	if err != nil {
		return err
	}

	return mapstructure.Decode(p, s)
}
